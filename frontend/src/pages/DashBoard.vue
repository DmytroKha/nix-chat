<template>
  <div class="dashboard">
    <div class="header">
      <h1>Dashboard</h1>
      <div>
        <button @click="logout">logout</button>
        <button @click="navigate">profile</button>
      </div>
    </div>
    <div class="content">
      <div class="sidebar">
        <div>
          <h2 v-on:click="showUsersList = 1">on-line</h2>
          <div v-if="showUsersList == 1">on-line users
            <div class="row" v-if="users.length">
              <div class="col-2 card profile"  v-for="user in users" :key="user.id">
                <div class="card-header">{{ user.name }}</div>
                <div class="card-body">
                  <button class="btn btn-primary" @click="joinPrivateRoom(user)">
                    Send msg
                  </button>
                  <button class="btn btn-primary" @click="addFriend(user)">
                    Add to FL
                  </button>
                  <button class="btn btn-primary" @click="addToBlackList(user)">
                    Add to BL
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div>
          <h2 v-on:click="showUsersList = 2">multi-rooms</h2>
          <div v-if="showUsersList == 2">multi-users rooms
            <div class="input-group">
              <input
                  v-model="roomInput"
                  class="form-control name"
                  placeholder="Type the room you want to join"
                  @keyup.enter.exact="joinRoom"
              />
              <div class="input-group-append">
                <span class="input-group-text send_btn" @click="joinRoom(roomInput)"> > </span>
              </div>
              <div class="input-group-append">
<!--            <span class="input-group-text send_btn" @click="getAllRooms">-->
<!--              >-->
<!--            </span>-->
              </div>
            </div>
            <div class="row" v-if="users.length">
              <div class="col-2 card profile"  v-for="room in rooms" :key="room.id">
                <div class="card-header">{{ room.name }}</div>
                <div class="card-body">
                  <button class="btn btn-primary" @click="joinRoom(room.name)">
                    Send msg
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div>
          <h2 v-on:click="showUsersList = 3">friends</h2>
          <div v-if="showUsersList == 3">users from the friend list
            <div class="row" v-if="friends.length">
              <div class="col-2 card profile"  v-for="user in friends" :key="user.id">
                <div class="card-header">{{ user.name }}</div>
                <div class="card-body">
                  <button class="btn btn-primary" @click="joinPrivateRoom(user)">
                    Send msg
                  </button>
                  <button class="btn btn-primary" @click="removeFromFriendList(user)">
                    Remove from friends
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div>
          <h2 v-on:click="showUsersList = 4">black list</h2>
          <div v-if="showUsersList == 4">users from black list
            <div class="row" v-if="blackList.length">
              <div class="col-2 card profile"  v-for="user in blackList" :key="user.id">
                <div class="card-header">{{ user.name }}</div>
                <div class="card-body">
                  <button class="btn btn-primary" @click="removeFromBlackList(user)">
                    Remove from blacklist
                  </button>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
      <div class="main">
<!--        <div class="row" v-if="users.length">-->
<!--          <div class="col-2 card profile"  v-for="user in users" :key="user.id">-->
<!--            <div class="card-header">{{ user.name }}</div>-->
<!--            <div class="card-body">-->
<!--              <button class="btn btn-primary" @click="joinPrivateRoom(user)">-->
<!--                Send Message-->
<!--              </button>-->
<!--              <button class="btn btn-primary" @click="addFriend(user)">-->
<!--                Add to friend list-->
<!--              </button>-->
<!--              <button class="btn btn-primary" @click="addFoe(user)">-->
<!--                Add to black list-->
<!--              </button>-->
<!--            </div>-->
<!--          </div>-->
<!--        </div>-->

        <div class="chat" v-for="(room, key) in chatRooms" :key="key">
          <div class="card">
            <div class="card-header msg_head">
              <div class="d-flex bd-highlight justify-content-center">
                {{ room.name }}
                <span class="card-close" @click="leaveRoom(room)">leave</span>
              </div>
            </div>
            <div class="card-body msg_card_body">
              <div
                v-for="(message, key) in room.messages"
                :key="key"
                class="d-flex justify-content-start mb-4"
              >
                <div class="msg_cotainer">
                  {{ message.message }}
                  <span class="msg_name" v-if="message.sender">
                    <img :src="message.sender.photo" width="30"/>
                    {{ message.sender.name }}
                  </span>
                </div>
              </div>
            </div>
            <div class="card-footer">
              <div class="input-group">
                <textarea
                  v-model="room.newMessage"
                  name=""
                  class="form-control type_msg"
                  placeholder="Type your message..."
                  @keyup.enter.exact="sendMessage(room)"
                ></textarea>
                <div class="input-group-append">
                  <span
                    class="input-group-text send_btn"
                    @click="sendMessage(room)"
                    >></span
                  >
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import router from "../router";
import { wsConnect } from "@/services/WSConnectService";
// import wsConnect from "../services/WSConnectService";
export default {
  name: "DashBoard",
  props: {
    msg: String,
  },
  data() {
    return {
      //showWhiteListUsers: true,
      //showOnlineListUsers: true,
      showUsersList: 1,
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      roomInput: null,
      rooms: [],
      chatRooms: [],
      // user: {
      //   id: "",
      //   name: "",
      //   token: "",
      //   friends: [],
      //   blackList: [],
      // },
      friends: [],
      blackList: [],
      users: [],
      initialReconnectDelay: 1000,
      currentReconnectDelay: 0,
      maxReconnectDelay: 16000,
      loginError: "",
      newRoom: true,
    };
  },
  //beforeMount() {
  //  this.connect();
  //},
  mounted: function () {
    // this.user.name = localStorage.getItem("name");
    // this.user.token = localStorage.getItem("token");
      if(
          !wsConnect.ws
    )
     {
    this.connect();
    } else {
        //console.log("1 users", this.users)
        this.users = wsConnect.users;
        this.rooms = wsConnect.rooms;
        this.chatRooms = wsConnect.chatRooms;
        this.blackList = wsConnect.user.blackList;
        this.friends = wsConnect.user.friends;
      }
  },
  methods: {
    navigate() {
      router.push({ path: "/profile" });
    },
    logout() {
      // localStorage.removeItem("name");
      // localStorage.removeItem("token");
      // localStorage.removeItem("id");
      // localStorage.removeItem("ws");
      // console.log(this.user.name);
      for (let i = 0; i < wsConnect.users.length; i++) {
        if (this.users[i].id == wsConnect.user.id) {
          //console.log("user-left", wsConnect.users[i]);
          wsConnect.ws.send(
            JSON.stringify({ action: "user-left", sender: wsConnect.users[i] })
          );
          wsConnect.users.splice(i, 1);
          break;
        }
      }
      //this.handleUserLeft({ action: "user-left", sender: this.user})
      //console.log("logout", wsConnect.user);
      // this.ws.send(
      //     JSON.stringify({ action: "user-left", message: wsConnect.user })
      // );
      wsConnect.ws = null;
      router.push({ path: "/" });
    },
    connect() {
      wsConnect.connectToWebsocket();
      this.connectToWebsocket();
      //this.getAllRooms();
    },
    connectToWebsocket() {
      wsConnect.ws.addEventListener("message", (event) => {
        this.handleNewMessage(event);
      });
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
      this.currentReconnectDelay = 1000;
    },
    onWebsocketClose() {
      wsConnect.ws = null;

      setTimeout(() => {
        this.reconnectToWebsocket();
      }, this.currentReconnectDelay);
    },
    reconnectToWebsocket() {
      if (wsConnect.ws === null) {
        if (this.currentReconnectDelay < this.maxReconnectDelay) {
          this.currentReconnectDelay *= 2;
        }
        this.connectToWebsocket();
      }
    },
    handleNewMessage(event) {
      let data = event.data;
      data = data.split(/\r?\n/);
      this.users = wsConnect.users;
      this.rooms = wsConnect.rooms;
      this.chatRooms = wsConnect.chatRooms;
      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        switch (msg.action) {
          case "send-message":
            this.handleChatMessage(msg);
            wsConnect.chatRooms = this.chatRooms;
            break;
          case "user-join":
            this.handleUserJoined(msg);
            wsConnect.users = this.users;
            break;
          case "user-left":
            this.handleUserLeft(msg);
            wsConnect.users = this.users;
            break;
          case "room-joined":
            this.handleRoomJoined(msg);
            wsConnect.rooms = this.rooms;
            break;
          case "all-rooms":
            this.handleAllRoomsJoined(msg);
            wsConnect.rooms = this.rooms;
            break;
          case "add-friend":
            this.friends = wsConnect.user.friends;
            this.handleFriendsJoined(msg);
            wsConnect.user.friends = this.friends;
            break;
          case "get-friends":
            this.handleFriends(msg);
            break;
          case "add-to-black-list":
            this.blackList = wsConnect.user.blackList;
            this.handleBlackListJoined(msg);
            wsConnect.user.blackList = this.blackList;
            break;
          case "get-black-list":
            this.handleBlackList(msg);
            break;
          default:
            break;
        }
      }
    },
    handleChatMessage(msg) {
      const room = this.findRoom(msg.target.id);
      if (typeof room !== "undefined") {
        room.messages.push(msg);
      }
      console.log("->send-message", msg)
      console.log("->room", room)
    },
    handleUserJoined(msg) {
      console.log('++++',msg.sender)
      console.log('++++',this.users)
      if (!this.userExists(msg.sender)) {
        this.users.push(msg.sender);
        // console.log('++++',this.users)
        // console.log('++++',msg.sender)
      }
    },
    handleUserLeft(msg) {
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id == msg.sender.id) {
          this.users.splice(i, 1);
          // console.log('---',this.users)
          // console.log('---',msg.sender)
          return;
        }
      }
    },
    handleRoomJoined(msg) {
      this.room = msg.target;
      this.room.name = this.room.private ? msg.sender.name : this.room.name;
      this.room["messages"] = [];
      this.chatRooms.push(this.room);
      console.log("0. rooms", this.chatRooms);
      console.log("1. rooms", this.rooms);
      console.log("2. room", this.room.id);
      console.log("2.1. room", this.room.private);
      this.newRoom = !this.room.private
      for (let i = 0; i < this.rooms.length; i++) {
        console.log("4. room", this.rooms[i].id);
        if (this.rooms[i].id === this.room.id) {
          this.newRoom = false
          console.log("5. newRoom",  this.newRoom);
        }
      }
      if (this.newRoom == true) {
        this.rooms.push(this.room);
        console.log("6. rooms",  this.rooms);
      }
    },
    handleAllRoomsJoined(msg) {
      this.room = msg.target;
      this.room.name = this.room.private ? msg.sender.name : this.room.name;
      this.room["messages"] = [];
      this.rooms.push(this.room);
      //console.log("rooms", this.rooms);
    },
    // handleOnlineUsers(msg) {
    //   this.users = [];
    //   for (let i = 0; i < msg.users.length; i++) {
    //      this.users.push(msg.users[i])
    //   }
    // },
    handleFriendsJoined(msg) {
      const usr = wsConnect.user;
      if (typeof usr !== "undefined") {
        var inList = false
        for (let i = 0; i < this.friends.length; i++) {
          if (this.friends[i].id == msg.sender.id) {
            inList = true;
            break;
          }
        }
        if (!inList) {
          usr.friends.push(msg.sender);
        }
      }
    },
    handleBlackListJoined(msg) {
      const usr = wsConnect.user;
      if (typeof usr !== "undefined") {
      var inList = false
      for (let i = 0; i < this.blackList.length; i++) {
        if (this.blackList[i].id == msg.sender.id) {
          inList = true;
          break;
        }
      }
      if (!inList) {
          usr.blackList.push(msg.sender);
      }
      }
    },
    handleFriends(msg) {
      console.log("1. friend", msg);
      console.log("2. friend", msg.users);
      var friends = msg.users;
      if (typeof friends !== "undefined") {
        for (let i = 0; i < friends.length; i++) {
          this.friends.push(msg.users[i]);
        }
      }
    },
    handleBlackList(msg) {
      console.log("1. bl", msg);
      console.log("2. bl", msg.users);
      var blackList = msg.users;
      if (typeof blackList !== "undefined") {
        for (let i = 0; i < blackList.length; i++) {
          this.blackList.push(msg.users[i]);
        }
      }
    },
    sendMessage(room) {
      if (room.newMessage !== "") {
        console.log("<-send-message", room.newMessage)
        console.log("<-room", room)
        wsConnect.ws.send(
          JSON.stringify({
            action: "send-message",
            message: room.newMessage,
            target: {
              id: room.id,
              name: room.name,
            },
          })
        );
        room.newMessage = "";
      }
    },
    findRoom(roomId) {
      console.log("<->room", roomId)
      for (let i = 0; i < this.chatRooms.length; i++) {
        if (this.chatRooms[i].id === roomId) {
          return this.chatRooms[i];
        }
      }
    },
    joinRoom(roomName) {
      wsConnect.ws.send(
        //JSON.stringify({ action: "join-room", message: this.roomInput })
          JSON.stringify({ action: "join-room", message: roomName })
      );
      this.roomInput = "";
    },
    getAllRooms() {
      console.log("--------rooms", this.rooms)
      console.log("--------chatRooms", this.chatRooms)
      wsConnect.ws.send(
          JSON.stringify({ action: "all-rooms"})
      );
      console.log("+++++++++rooms", this.rooms)
      console.log("+++++++++chatRooms", this.chatRooms)
    },
    // getOnlineUsers() {
    //   wsConnect.ws.send(
    //       JSON.stringify({ action: "on-line-users", message: "get on-line users" })
    //   );
    //},
    leaveRoom(room) {
      wsConnect.ws.send(JSON.stringify({ action: "leave-room", message: room.id.toString() }));

      for (let i = 0; i < this.chatRooms.length; i++) {
        if (this.chatRooms[i].id === room.id) {
          this.chatRooms.splice(i, 1);
          break;
        }
      }
    },
    joinPrivateRoom(room) {
      console.log("room.id", room.id)
      console.log("0. room.id", this.chatRooms);
      console.log("1. room.id", this.rooms);
      wsConnect.ws.send(
        JSON.stringify({ action: "join-room-private", message: room.id.toString() })
      );
    },
    addFriend(friend) {
      wsConnect.ws.send(
        JSON.stringify({ action: "add-friend", sender: friend })
      );
    },
    removeFromFriendList(friend) {
      console.log("removeFromFriendList", friend);
      wsConnect.ws.send(JSON.stringify({ action: "remove-from-friend-list", sender: friend }));

      for (let i = 0; i < this.friends.length; i++) {
        if (this.friends[i].id === friend.id) {
          this.friends.splice(i, 1);
          break;
        }
      }
    },
    addToBlackList(bl) {
      wsConnect.ws.send(JSON.stringify({ action: "add-to-black-list", sender: bl }));
    },
    removeFromBlackList(bl) {
      console.log("removeFromBlackList", bl);
      wsConnect.ws.send(JSON.stringify({ action: "remove-from-black-list", sender: bl }));

      for (let i = 0; i < this.blackList.length; i++) {
        if (this.blackList[i].id === bl.id) {
          this.blackList.splice(i, 1);
          break;
        }
      }
    },
    userExists(user) {
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id == user.id) {
          return true;
        }
      }
      return false;
    },
  },
};
</script>

<style scoped>
.dashboard {
  height: 100%;
}
.content {
  display: flex;
  justify-content: space-between;
  height: calc(100% - 60px);
}
.sidebar {
  width: 30%;
  height: 100%;
  background-color: beige;
  padding: 10px;
}
.sidebar h2 {
  cursor: pointer;
}
.main {
  width: 80%;
  padding: 10px;
  background-color: rgb(190, 190, 144);
  height: 100%;
}
.chat {
  margin: 15px;
}

.room,
.form {
  margin-top: auto;
  margin-bottom: auto;
}

.card {
  height: 500px;
  border-radius: 10px;
  background-color: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card.profile {
  height: 150px;
  margin: 15px;
}

.sidebar .card.profile {
  height: initial;
  margin: 1px;
}

.card.profile .card-header {
  color: #fff;
}

.msg_head {
  color: #fff;
  text-align: center;
  font-size: 26px;
}

.msg_card_body {
  overflow-y: auto;
}
.card-header {
  border-radius: 15px 15px 0 0;
  border-bottom: 0;
}

.card-close {
  font-size: 0.5em;
  text-decoration: underline;
  float: right;
  position: absolute;
  right: 15px;
  cursor: pointer;
}

.card-footer {
  border-radius: 0 0 15px 15px;
  border-top: 0;
}
.container {
  align-content: center;
}

.type_msg {
  background-color: rgba(86, 10, 134, 0.6);
  border: 0;
  color: white;
  height: 60px;
  width: 100%;
  overflow-y: auto;
}
.type_msg:focus {
  box-shadow: none;
  outline: 0px;
  background-color: rgba(255, 255, 255, 0.6);
}

.send_btn {
  border-radius: 0 15px 15px 0;
  background-color: rgba(0, 0, 0, 0.3);
  border: 0;
  color: white;
  cursor: pointer;
}

.msg_cotainer {
  margin-top: auto;
  margin-bottom: auto;
  margin-left: 10px;
  border-radius: 25px;
  background-color: #82ccdd;
  padding: 10px 15px;
  position: relative;
  min-width: 100px;
  line-height: 1.2rem;
}
.msg_cotainer_send {
  margin-top: auto;
  margin-bottom: auto;
  margin-right: 10px;
  border-radius: 25px;
  background-color: #75d5fd;
  padding: 10px;
  position: relative;
}

.msg_name {
  display: block;
  font-size: 0.8em;
  font-style: italic;
  color: #545454;
}

.msg_head {
  position: relative;
}
</style>
