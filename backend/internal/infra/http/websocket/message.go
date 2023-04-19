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
const AddToBlackListAction = "add-to-black-list"
const GetAllRooms = "all-rooms"
const GetBlackList = "get-black-list"
const RemoveFromBlackListAction = "remove-from-black-list"
const GetFriends = "get-friends"
const RemoveFromFriendsAction = "remove-from-friends"

type Message struct {
	Action  string        `json:"action"`
	Message string        `json:"message"`
	Target  *Room         `json:"target"`
	Sender  domain.User   `json:"sender"`
	Users   []domain.User `json:"users"`
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
		*Alias
	}{
		Alias: (*Alias)(message),
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}
	message.Sender = &msg.Sender
	return nil
}
