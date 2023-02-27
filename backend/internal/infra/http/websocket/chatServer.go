package websocket

import (
	"context"
	"encoding/json"
	"github.com/DmytroKha/nix-chat/config"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
	"log"
	"strconv"
)

const PubSubGeneralChannel = "general"

type WsServer struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	users          []database.User
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
	for _, user := range server.users {
		message := &Message{
			Action: UserJoinedAction,
			//Sender: user,
			SenderId: user.Id,
		}
		client.send <- message.encode()
	}
}

func (server *WsServer) registerClient(client *Client, ctx context.Context) {
	// Add user to the repo
	user := database.User{Id: client.ID, Name: client.Name}
	server.userRepository.Save(user)

	// Publish user in PubSub
	server.publishClientJoined(client, ctx)

	server.listOnlineClients(client)
	server.clients[client] = true
}

func (server *WsServer) unregisterClient(client *Client, ctx context.Context) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)

		// Remove user from repo
		server.userRepository.Delete(client.ID)

		// Publish user left in PubSub
		server.publishClientLeft(client, ctx)
	}
}

func (server *WsServer) findRoomByID(ID int64) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findClientByID(ID int64) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID == ID {
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
		room.ID = dbRoom.GetId()

		go room.RunRoom(ctx)
		server.rooms[room] = true
	}

	return room
}

// Publish userJoined message in pub/sub
func (server *WsServer) publishClientJoined(client *Client, ctx context.Context) {

	message := &Message{
		Action: UserJoinedAction,
		//Sender: client,
		SenderId: client.ID,
	}

	if err := config.Redis.Publish(ctx, PubSubGeneralChannel, message.encode()).Err(); err != nil {
		log.Println(err)
	}
}

// Publish userleft message in pub/sub
func (server *WsServer) publishClientLeft(client *Client, ctx context.Context) {

	message := &Message{
		Action: UserLeftAction,
		//Sender: client,
		SenderId: client.ID,
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
	//server.users = append(server.users, message.Sender)
	user, err := server.userRepository.Find(message.SenderId)
	if err != nil {
		log.Printf("handleUserJoined: %s", err)
	}
	server.users = append(server.users, user)
	server.broadcastToClients(message.encode())
}

func (server *WsServer) handleUserLeft(message Message) {
	// Remove the user from the slice
	for i, user := range server.users {
		//if user.GetId() == message.Sender.GetId() {
		//	server.users[i] = server.users[len(server.users)-1]
		//	server.users = server.users[:len(server.users)-1]
		//}
		if user.Id == message.SenderId {
			server.users[i] = server.users[len(server.users)-1]
			server.users = server.users[:len(server.users)-1]
		}
	}
	server.broadcastToClients(message.encode())
}

func (server *WsServer) findUserByID(ID int64) database.User {
	var foundUser database.User
	for _, client := range server.users {
		if client.Id == ID {
			foundUser = client
			break
		}
	}

	return foundUser
}

func (server *WsServer) handleUserJoinPrivate(message Message, ctx context.Context) {
	// Find client for given user, if found add the user to the room.
	clientId, _ := strconv.Atoi(message.Message)
	targetClient := server.findClientByID(int64(clientId))
	if targetClient != nil {
		//targetClient.joinRoom(message.Target.GetName(), message.Sender, ctx)
		targetClient.joinRoom(message.Target.GetName(), message.SenderId, ctx)
	}
}
