package websocket

import (
	"context"
	"encoding/json"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"github.com/google/uuid"
	"log"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
	//users          []database.User
	users          []domain.User
	roomRepository domain.RoomRepository
	userRepository database.UserRepository
}

// NewWebsocketServer creates a new WsServer type
func NewWebsocketServer(roomRepository domain.RoomRepository, userRepository database.UserRepository) *WsServer {
	wsServer := &WsServer{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		rooms:          make(map[*Room]bool),
		roomRepository: roomRepository,
		userRepository: userRepository,
	}

	// Add users from database to server
	wsServer.users, _ = userRepository.FindAll()

	return wsServer
}

// Run our websocket server, accepting various requests
func (server *WsServer) Run(ctx context.Context) {
	go server.listenPubSubChannel(ctx)
	for {
		select {

		case client := <-server.register:
			server.registerClient(client, ctx)

		case client := <-server.unregister:
			server.unregisterClient(client, ctx)

		}

	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string, ctx context.Context) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	if foundRoom == nil {
		// Try to run the room from the repository, if it is found.
		foundRoom = server.runRoomFromRepository(name, ctx)
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string, private bool, ctx context.Context) *Room {
	room := NewRoom(name, private)
	server.roomRepository.Save(room)

	go room.RunRoom(ctx)
	server.rooms[room] = true

	return room
}

func (server *WsServer) listOnlineClients(client *Client) {
	var uniqueUsers = make(map[string]bool)
	for _, user := range server.users {
		if ok := uniqueUsers[user.GetId()]; !ok {
			//if ok := uniqueUsers[user.Id]; !ok {
			message := &Message{
				Action: UserJoinedAction,
				Sender: user,
				//SenderId: user.Id,
			}
			uniqueUsers[user.GetId()] = true
			//uniqueUsers[user.Id] = true
			client.send <- message.encode()
		}
	}
}

func (server *WsServer) registerClient(client *Client, ctx context.Context) {
	uid := client.ID.String()
	if user := server.findUserByID(uid); user == nil {
		// Add user to the repo
		var userRepo database.User
		userRepo.Uid = client.GetId()
		userRepo.Name = client.GetName()
		server.userRepository.Save(userRepo)
	}

	// Publish user in PubSub
	server.publishClientJoined(client, ctx)

	server.listOnlineClients(client)
	server.clients[client] = true

}

func (server *WsServer) unregisterClient(client *Client, ctx context.Context) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		// Publish user left in PubSub
		server.publishClientLeft(client, ctx)
	}
}

func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.GetId() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

// NEW: Try to find a room in the repo, if found Run it.
func (server *WsServer) runRoomFromRepository(name string, ctx context.Context) *Room {
	var room *Room
	dbRoom, _ := server.roomRepository.FindByName(name)
	if dbRoom != nil {
		room = NewRoom(dbRoom.GetName(), dbRoom.GetPrivate())
		room.ID, _ = uuid.Parse(dbRoom.GetId())

		go room.RunRoom(ctx)
		server.rooms[room] = true
	}

	return room
}

// Publish userJoined message in pub/sub
func (server *WsServer) publishClientJoined(client *Client, ctx context.Context) {

	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
		//SenderId: client.ID,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Publish userleft message in pub/sub
func (server *WsServer) publishClientLeft(client *Client, ctx context.Context) {

	message := &Message{
		Action: UserLeftAction,
		Sender: client,
		//SenderId: client.ID,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Listen to pub/sub general channels
func (server *WsServer) listenPubSubChannel(ctx context.Context) {

	pubsub := config.Redis.Subscribe(ctx, PubSubGeneralChannel)

	ch := pubsub.Channel()

	for msg := range ch {

		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}

		switch message.Action {
		case UserJoinedAction:
			server.handleUserJoined(message)
		case UserLeftAction:
			server.handleUserLeft(message)
		case JoinRoomPrivateAction:
			server.handleUserJoinPrivate(message, ctx)
		}

	}
}

func (server *WsServer) handleUserJoined(message Message) {
	// Add the user to the slice
	server.users = append(server.users, message.Sender)
	//user, err := server.userRepository.Find(message.SenderId)
	//if err != nil {
	//	log.Printf("handleUserJoined: %s", err)
	//}
	//server.users = append(server.users, user)
	server.broadcastToClients(message.encode())

}

func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		if user.GetId() == message.Sender.GetId() {
			//if user.Id == message.SenderId {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
			break // added this break to only remove the first occurrence
		}
	}
	server.broadcastToClients(message.encode())
}

func (server *WsServer) findUserByID(ID string) domain.User {
	var foundUser domain.User
	//var foundUser database.User
	for _, client := range server.users {
		uid := client.GetId()
		if uid == ID {
			//if client.Id == ID {
			foundUser = client
			break
		}
	}

	return foundUser
}

func (server *WsServer) handleUserJoinPrivate(message Message, ctx context.Context) {
	//clientId, _ := strconv.Atoi(message.Message)
	targetClients := server.findClientsByID(message.Message)
	//targetClients := server.findClientsByID(int64(clientId))
	for _, targetClient := range targetClients {
		//targetClient.joinRoom(message.Target.GetName(), message.SenderId, ctx)
		targetClient.joinRoom(message.Target.GetName(), message.Sender, ctx)
	}
}

func (server *WsServer) findClientsByID(ID string) []*Client {
	var foundClients []*Client
	for client := range server.clients {
		if client.GetId() == ID {
			foundClients = append(foundClients, client)
		}
	}

	return foundClients

}

func (server *WsServer) AppendUser(user domain.User) {
	server.users = append(server.users, user)
}
