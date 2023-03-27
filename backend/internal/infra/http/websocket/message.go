package websocket

import (
	"encoding/json"
	"github.com/DmytroKha/nix-chat/internal/domain"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const AddFriendAction = "add-friend"
const AddFoeAction = "add-foe"
const GetAllRooms = "all-rooms"
const ChangeNameAction = "change-username"

//const GetOnlineUsers = "on-line-users"

type Message struct {
	Action  string        `json:"action"`
	Message string        `json:"message"`
	Target  *Room         `json:"target"`
	Sender  domain.User   `json:"sender"`
	Users   []domain.User `json:"users"`
	//SenderId int64 `json:"senderId"`
}

func (message *Message) encode() []byte {
	j, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return j
}

// UnmarshalJSON custom unmarshel to create a Client instance for Sender
func (message *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	msg := &struct {
		Sender Client `json:"sender"`
		//SenderId int64 `json:"senderId"`
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	//message.SenderId = msg.SenderId
	return nil
}
