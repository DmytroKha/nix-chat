package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents the websocket client at the server
type Client struct {
	// The actual websocket connection.
	conn     *websocket.Conn
	wsServer *WsServer
	send     chan []byte
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	rooms    map[*Room]bool
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string, ID string, photo string) *Client {
	client := &Client{
		Name:     name,
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		rooms:    make(map[*Room]bool),
	}

	if ID != "" {
		client.ID, _ = uuid.Parse(ID)
	}

	if photo != "" {
		client.Photo = "../././file_storage/" + photo
	}

	return client

}

func (client *Client) readPump() {
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		client.disconnect()
		cancel()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage, ctx)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.wsServer.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	client.wsServer.unregister <- client
	close(client.send)
	client.conn.Close()
}

// ServeWs handles websocket requests from clients requests.
func ServeWs(wsServer *WsServer, ctx echo.Context) error {

	c := ctx.Request().Context()
	userCtxValue := c.Value("user")
	if userCtxValue == nil {
		err := fmt.Errorf("Not authenticated")
		log.Println(err)
		return err
	}

	user := userCtxValue.(domain.User)

	conn, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	client := newClient(conn, wsServer, user.GetName(), user.GetUid(), user.GetPhoto())

	go client.writePump()
	go client.readPump()
	//go client.readPump(context.Background())
	wsServer.register <- client

	return nil

}

// Refactored method
// Use the ID of the target room instead of the name to find it.
// Added case for joining private room
func (client *Client) handleNewMessage(jsonMessage []byte, ctx context.Context) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	if message.Sender == nil || message.Sender.GetName() == "" {
		message.Sender = client
	}

	switch message.Action {
	case SendMessageAction:
		roomID := message.Target.GetUid()
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.broadcast <- &message
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(message, ctx)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessage(message, ctx)

	case GetAllRooms:
		client.handleJoinAllRoomsMessage(message, ctx)

		//case GetOnlineUsers:
		//	client.handleJoinOnlineUsersMessage(ctx)
	case UserJoinedAction:
		client.handleUserJoinMessage(message)

	case UserLeftAction:
		client.handleUserLeaveMessage(message)

	case AddToBlackListAction:
		client.handleJoinToBlackListMessage(message)

	case GetBlackList:
		client.handleBlackList(message)

	}

}

// Refactored method
// Use new joinRoom method
func (client *Client) handleJoinRoomMessage(message Message, ctx context.Context) {
	roomName := message.Message

	client.joinRoom(roomName, nil, ctx)
}

func (client *Client) handleJoinAllRoomsMessage(message Message, ctx context.Context) {
	sender := message.Sender

	client.joinToAllRooms(sender, ctx)
}

func (client *Client) handleJoinToBlackListMessage(message Message) {
	sender := message.Sender

	client.joinToBlackList(sender)
}

func (client *Client) handleBlackList(message Message) {
	sender := message.Sender
	users, _ := client.wsServer.userRepository.GetUserBlackList(sender)

	msg := Message{
		Action: GetBlackList,
		//Target: room,
		Sender: sender,
		Users:  users,
	}

	client.send <- msg.encode()
}

//func (client *Client) handleJoinOnlineUsersMessage(ctx context.Context) {
//	//client.wsServer.publishClientJoined(client, ctx)
//	//client.wsServer.listOnlineClients(client)
//	//client.wsServer.checkOnlineClients(ctx)
//	message := &Message{
//		Action: GetOnlineUsers,
//		Sender: client,
//		Users:  client.wsServer.users,
//	}
//	client.send <- message.encode()
//}

// Refactored method
// Added nil check
func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByID(message.Message)
	if room == nil {
		return
	}
	if _, ok := client.rooms[room]; ok {
		delete(client.rooms, room)
	}

	room.unregister <- client
}

func (client *Client) handleUserLeaveMessage(message Message) {
	//client.wsServer.publishClientLeft(client, ctx)
	//message := &Message{
	//	Action: UserLeftAction,
	//	Sender: client,
	//}
	//client.send <- message.encode()
	//client.wsServer.handleUserLeft(message)
	for i, user := range client.wsServer.users {
		if user.GetUid() == message.Sender.GetUid() {
			//if user.Id == message.SenderId {
			client.wsServer.users[i] = client.wsServer.users[len(client.wsServer.users)-1]
			client.wsServer.users = client.wsServer.users[:len(client.wsServer.users)-1]
			break // added this break to only remove the first occurrence
		}
	}
	client.wsServer.unregister <- client
}

func (client *Client) handleUserJoinMessage(message Message) {
	//client.wsServer.publishClientJoined(client, ctx)
	//message := &Message{
	//	Action: UserJoinedAction,
	//	Sender: client,
	//}
	//client.send <- message.encode()
	//client.wsServer.handleUserJoined(message)

	if client.ID.String() == message.Sender.GetUid() && client.Name != message.Sender.GetName() {
		client.Name = message.Sender.GetName()
	}
	if client.ID.String() == message.Sender.GetUid() && client.Photo != message.Sender.GetPhoto() {
		client.Photo = message.Sender.GetPhoto()
	}
	client.wsServer.register <- client
}

// New method
// When joining a private room we will combine the IDs of the users
// Then we will bothe join the client and the target.
func (client *Client) handleJoinRoomPrivateMessage(message Message, ctx context.Context) {

	//userId, _ := strconv.Atoi(message.Message)
	//target := client.wsServer.findUserByID(int64(userId))
	target := client.wsServer.findUserByID(message.Message)

	if target == nil {
		return
	}

	// create unique room name combined to the two IDs
	//roomName := message.Message + client.ID.String()
	roomName := ""
	if message.Message < client.ID.String() {
		roomName = message.Message + client.ID.String()
	} else {
		roomName = client.ID.String() + message.Message
	}

	// Join room
	joinedRoom := client.joinRoom(roomName, target, ctx)

	// Invite target user
	if joinedRoom != nil {
		client.inviteTargetUser(target, joinedRoom, ctx)
	}

}

// New method
// Joining a room both for public and private roooms
// When joiing a private room a sender is passed as the opposing party
func (client *Client) joinRoom(roomName string, sender domain.User, ctx context.Context) *Room {
	//func (client *Client) joinRoom(roomName string, senderId int64, ctx context.Context) *Room {

	room := client.wsServer.findRoomByName(roomName, ctx)
	if room == nil {
		room = client.wsServer.createRoom(roomName, sender != nil, ctx)
	}

	// Don't allow to join private rooms through public room message
	if sender == nil && room.Private {
		return nil
	}

	if !client.isInRoom(room) {

		client.rooms[room] = true
		room.register <- client

		client.notifyRoomJoined(room, sender)
	}

	return room

}

func (client *Client) joinToAllRooms(sender domain.User, ctx context.Context) {
	//func (client *Client) joinRoom(roomName string, senderId int64, ctx context.Context) *Room {

	rooms := client.wsServer.findAllRooms(ctx)

	// Don't allow to join private rooms through public room message
	//if sender == nil && room.Private {
	//	return nil
	//}
	for _, room := range rooms {
		//client.rooms[room] = true
		//room.register <- client
		var r Room
		r.ID, _ = uuid.Parse(room.GetUid())
		r.Name = room.GetName()
		r.Private = room.GetPrivate()

		client.notifyGetAllRooms(&r, sender)
	}

}

func (client *Client) joinToBlackList(sender domain.User) {

	//rooms := client.wsServer.findAllRooms(ctx)

	// Don't allow to join private rooms through public room message
	//if sender == nil && room.Private {
	//	return nil
	//}
	//for _, room := range rooms {
	//	//client.rooms[room] = true
	//	//room.register <- client
	//	var r Room
	//	r.ID, _ = uuid.Parse(room.GetUid())
	//	r.Name = room.GetName()
	//	r.Private = room.GetPrivate()

	client.notifyBlackList(sender)
	//}

}

// New method
// Check if the client is not yet in the room
func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}

// New method
// Notify the client of the new room he/she joined
func (client *Client) notifyRoomJoined(room *Room, sender domain.User) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
		//SenderId: senderId,
	}

	client.send <- message.encode()
}

func (client *Client) notifyGetAllRooms(room *Room, sender domain.User) {
	message := Message{
		Action: GetAllRooms,
		Target: room,
		Sender: sender,
		//SenderId: senderId,
	}

	client.send <- message.encode()
}

func (client *Client) notifyBlackList(sender domain.User) {
	//var users []domain.User
	//users = append(users, user)
	message := Message{
		Action: AddToBlackListAction,
		//Target: room,
		Sender: sender,
		//Users:  users,
	}

	client.send <- message.encode()
}

func (client *Client) GetUid() string {
	return client.ID.String()
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) GetPhoto() string {
	return client.Photo
}

func (client *Client) GetId() int64 {
	return 0
}

// Send out invite message over pub/sub in the general channel.
func (client *Client) inviteTargetUser(target domain.User, room *Room, ctx context.Context) {
	inviteMessage := &Message{
		Action:  JoinRoomPrivateAction,
		Message: target.GetUid(),
		Target:  room,
		Sender:  client,
		//SenderId: client.ID,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, inviteMessage.encode()).Err(); err != nil {
		log.Println(err)
	}
}
