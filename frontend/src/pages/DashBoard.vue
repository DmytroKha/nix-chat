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
                  <button class="btn btn-primary" @click="addFoe(user)">
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
            <span class="input-group-text send_btn" @click="getAllRooms">
              >
            </span>
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
          <h2 v-on:click="showUsersList = 3">white list</h2>
          <div v-if="showUsersList == 3">users from white list</div>
        </div>
        <div>
          <h2 v-on:click="showUsersList = 4">black list</h2>
          <div v-if="showUsersList == 4">users from black list</div>
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
                    <img :src="message.sender.photo" width="30" height="30"/>
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
      //   uid: "",
      //   name: "",
      //   username: "",
      //   password: "",
      //   confirmPassword: "",
      //   oldPassword: "",
      //   newPassword: "",
      //   token: "",
      //   friends: [],
      //   foes: [],
      // },
      users: [],
      initialReconnectDelay: 1000,
      currentReconnectDelay: 0,
      maxReconnectDelay: 16000,
      loginError: "",
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
      }
  },
  methods: {
    navigate() {
      router.push({ path: "/profile" });
    },
    logout() {
      // localStorage.removeItem("name");
      // localStorage.removeItem("token");
      // localStorage.removeItem("uid");
      // localStorage.removeItem("ws");
      // console.log(this.user.name);
      for (let i = 0; i < wsConnect.users.length; i++) {
        if (this.users[i].id == wsConnect.user.uid) {
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
      // if (this.user.token != "") {
      //   this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);
      //   localStorage.setItem('ws', true);
      // }
      //else {
      //  this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
      //}
      // this.ws.addEventListener("open", (event) => {
      //   this.onWebsocketOpen(event);
      // });
      // this.ws.addEventListener("message", (event) => {
      //   this.handleNewMessage(event);
      // });
      wsConnect.ws.addEventListener("message", (event) => {
        //console.log("event connectToWebsocket", event);
        this.handleNewMessage(event);
      });
      // this.ws.addEventListener("close", (event) => {
      //   this.onWebsocketClose(event);
      // });
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
      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        switch (msg.action) {
          case "send-message":
            this.handleChatMessage(msg);
            break;
          case "user-join":
            // console.log("users user-join users", this.users);
            // console.log("users user-join ws", wsConnect.users);
            // console.log("users user-join msg", msg);
            this.handleUserJoined(msg);
            wsConnect.users = this.users;
            break;
          case "user-left":
            // console.log("users log-out users", this.users);
            // console.log("users log-out ws", wsConnect.users);
            // console.log("users log-out msg", msg);
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
          // case "on-line-users":
          //   //console.log('!!! on-line-users', msg)
          //   this.handleOnlineUsers(msg);
          //   break;
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
    },
    handleUserJoined(msg) {
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
      //console.log("rooms", this.rooms);
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
    sendMessage(room) {
      if (room.newMessage !== "") {
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
      wsConnect.ws.send(JSON.stringify({ action: "leave-room", message: room.id }));

      for (let i = 0; i < this.chatRooms.length; i++) {
        if (this.chatRooms[i].id === room.id) {
          this.chatRooms.splice(i, 1);
          break;
        }
      }
    },
    joinPrivateRoom(room) {
      wsConnect.ws.send(
        JSON.stringify({ action: "join-room-private", message: room.id })
      );
    },
    addFriend(friend) {
      wsConnect.ws.send(
        JSON.stringify({ action: "add-friend", message: friend.id })
      );
    },
    addFoe(foe) {
      wsConnect.ws.send(JSON.stringify({ action: "add-foe", message: foe.id }));
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
