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
          <h2 v-on:click="showWhiteListUsers = true">white list</h2>
          <div v-if="showWhiteListUsers">users from white list</div>
        </div>
        <div>
          <h2 v-on:click="showWhiteListUsers = false">black list</h2>
          <div v-if="!showWhiteListUsers">users from black list</div>
        </div>
      </div>
      <div class="main">
        <div class="row">
          <div class="col-2 card profile" v-for="user in users" :key="user.id">
            <div class="card-header">{{ user.name }}</div>
            <div class="card-body">
              <button class="btn btn-primary" @click="joinPrivateRoom(user)">
                Send Message
              </button>
              <button class="btn btn-primary" @click="addFriend(user)">
                Add to friend list
              </button>
              <button class="btn btn-primary" @click="addFoe(user)">
                Add to black list
              </button>
            </div>
          </div>
        </div>

        <div class="input-group">
          <input
              v-model="roomInput"
              class="form-control name"
              placeholder="Type the room you want to join"
              @keyup.enter.exact="joinRoom"
          />
          <div class="input-group-append">
                <span class="input-group-text send_btn" @click="joinRoom">
                  >
                </span>
          </div>
        </div>

        <div class="chat" v-for="(room, key) in rooms" :key="key">
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
                  <span class="msg_name" v-if="message.sender">{{
                    message.sender.name
                  }}</span>
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
export default {
  name: "DashBoard",
  props: {
    msg: String,
  },
  data() {
    return {
      showWhiteListUsers: true,
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      roomInput: null,
      rooms: [],
      user: {
        uid: "",
        name: "",
        username: "",
        password: "",
        confirmPassword: "",
        oldPassword: "",
        newPassword: "",
        token: "",
        friends: [],
        foes: [],
        ws: false,
      },
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
      this.user.name = localStorage.getItem('name');
      this.user.token = localStorage.getItem('token');
      this.user.ws = localStorage.getItem('ws');
      console.log(this.ws)
      console.log(this.user.ws)
      console.log(6,this.users)
      if(
         this.user.ws == false || this.user.ws == null
     )
      {
        this.connect()
      }
  },

  methods: {
    navigate() {
      //localStorage.setItem('name', this.user.username);
      //localStorage.setItem('token', this.user.token);
      //localStorage.setItem('uid', this.user.uid);
      router.push({ path: "/profile" });
    },
    logout(){
      localStorage.removeItem('name');
      localStorage.removeItem('token');
      localStorage.removeItem('uid');
      localStorage.removeItem('ws');
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id == this.user.uid) {
          this.ws.send(JSON.stringify({ action: "user-left", sender: this.users[i]}));
        }
      }
      this.ws = null;
      router.push({ path: "/" });
    },

     connect() {
       this.connectToWebsocket();
     },
     connectToWebsocket() {
       if (this.user.token != "") {
         this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);
         localStorage.setItem('ws', true);
       }
       //else {
       //  this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
      //}
      this.ws.addEventListener("open", (event) => {
        this.onWebsocketOpen(event);
      });
      this.ws.addEventListener("message", (event) => {
        this.handleNewMessage(event);
      });
      this.ws.addEventListener("close", (event) => {
        this.onWebsocketClose(event);
      });
    },
    onWebsocketOpen() {
      console.log("connected to WS!");
      this.currentReconnectDelay = 1000;
    },
    onWebsocketClose() {
      this.ws = null;

      setTimeout(() => {
        this.reconnectToWebsocket();
      }, this.currentReconnectDelay);
    },
    reconnectToWebsocket() {
      if(this.ws === null){
        if (this.currentReconnectDelay < this.maxReconnectDelay) {
          this.currentReconnectDelay *= 2;
        }
        this.connectToWebsocket();
      }
    },
    handleNewMessage(event) {
      let data = event.data;
      data = data.split(/\r?\n/);

      for (let i = 0; i < data.length; i++) {
        let msg = JSON.parse(data[i]);
        switch (msg.action) {
          case "send-message":
            this.handleChatMessage(msg);
            break;
          case "user-join":
            this.handleUserJoined(msg);
            break;
          case "user-left":
            this.handleUserLeft(msg);
            break;
          case "room-joined":
            this.handleRoomJoined(msg);
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
    },
    handleUserJoined(msg) {
      if (!this.userExists(msg.sender)) {
        this.users.push(msg.sender);
        this.user.uid = msg.sender.id
        localStorage.setItem('uid', this.user.uid);
        console.log(5,this.users)
      }
    },
    handleUserLeft(msg) {
      for (let i = 0; i < this.users.length; i++) {
        if (this.users[i].id == msg.sender.id) {
          this.users.splice(i, 1);
          return;
        }
      }
    },
    handleRoomJoined(msg) {
      this.room = msg.target;
      this.room.name = this.room.private ? msg.sender.name : this.room.name;
      this.room["messages"] = [];
      this.rooms.push(this.room);
    },
    sendMessage(room) {
      if (room.newMessage !== "") {
        this.ws.send(
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
      for (let i = 0; i < this.rooms.length; i++) {
        if (this.rooms[i].id === roomId) {
          return this.rooms[i];
        }
      }
    },
    joinRoom() {
      this.ws.send(
        JSON.stringify({ action: "join-room", message: this.roomInput })
      );
      this.roomInput = "";
    },
    leaveRoom(room) {
      this.ws.send(JSON.stringify({ action: "leave-room", message: room.id }));

      for (let i = 0; i < this.rooms.length; i++) {
        if (this.rooms[i].id === room.id) {
          this.rooms.splice(i, 1);
          break;
        }
      }
    },
    joinPrivateRoom(room) {
      this.ws.send(
        JSON.stringify({ action: "join-room-private", message: room.id })
      );
    },
    addFriend(friend) {
      this.ws.send(
        JSON.stringify({ action: "add-friend", message: friend.id })
      );
    },
    addFoe(foe) {
      this.ws.send(JSON.stringify({ action: "add-foe", message: foe.id }));
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
  width: 20%;
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
  border-radius: 15px;
  background-color: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.card.profile {
  height: 150px;
  margin: 15px;
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
