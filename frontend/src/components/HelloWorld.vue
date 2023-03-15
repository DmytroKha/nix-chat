<template>
  <div class="hello">
    <div id="app">
      <div class="container h-100">
        <div class="row justify-content-center h-100">
          <div class="col-12 form" v-if="!ws">
            <!--<h2>Anonymous log-in</h2>
            <div class="input-group">
              <input v-model="user.name" class="form-control name" placeholder="Please fill in your (nick)name"
                     @keyup.enter.exact="connect"/>
              <div class="input-group-append">
              <span class="input-group-text send_btn" @click="connect">
                >
              </span>
              </div>
            </div>
            -->

            <h2>New user</h2>

            <div class="input-group">
              <input
                v-model="user.username"
                class="form-control username"
                placeholder="username"
              />

              <input
                v-model="user.password"
                type="password"
                class="form-control password"
                placeholder="password"
              />
              <input
                v-model="user.confirmPassword"
                type="password"
                class="form-control password"
                placeholder="confirm password"
              />

              <div class="input-group-append">
                <span class="input-group-text send_btn" @click="register">
                  >
                </span>
              </div>
            </div>

            <h2>Registered users</h2>

            <div class="input-group">
              <input
                v-model="user.username"
                class="form-control username"
                placeholder="username"
              />

              <input
                v-model="user.password"
                type="password"
                class="form-control password"
                placeholder="password"
              />

              <div class="input-group-append">
                <span class="input-group-text send_btn" @click="login">
                  >
                </span>
              </div>
            </div>

            <div class="alert alert-danger" role="alert" v-show="loginError">
              {{ loginError }}
            </div>
          </div>
          <div class="col-12">
            <div class="row">
              <div
                class="col-2 card profile"
                v-for="user in users"
                :key="user.id"
              >
                <div class="card-header">{{ user.name }}</div>
                <div class="card-body">
                  <button
                    class="btn btn-primary"
                    @click="joinPrivateRoom(user)">
                    Send Message
                  </button>
                  <button
                      class="btn btn-primary"
                      @click="addFriend(user)">
                    Add to friend list
                  </button>
                  <button
                      class="btn btn-primary"
                      @click="addFoe(user)">
                    Add to black list
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div class="col-12 room" v-if="ws != null">
            <h2>Change user password</h2>

            <div class="input-group">
              <input
                v-model="user.oldPassword"
                type="password"
                class="form-control password"
                placeholder="old password"
              />
              <input
                v-model="user.newPassword"
                type="password"
                class="form-control password"
                placeholder="new password"
              />
              <div class="input-group-append">
                <span class="input-group-text send_btn" @click="changePass">
                  >
                </span>
              </div>
            </div>

            <div class="alert alert-danger" role="alert" v-show="loginError">
              {{ loginError }}
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
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
  data() {
    return {
      ws: null,
      serverUrl: "ws://localhost:8080/ws",
      roomInput: null,
      rooms: [],
      user: {
        name: "",
        username: "",
        password: "",
        confirmPassword: "",
        oldPassword: "",
        newPassword: "",
        token: "",
        friends: [],
        foes: [],
      },
      users: [],
      initialReconnectDelay: 1000,
      currentReconnectDelay: 0,
      maxReconnectDelay: 16000,
      loginError: "",
    };
  },
  mounted: function () {},
  methods: {
    connect() {
      this.connectToWebsocket();
    },
    async register() {
      try {
        //const result = await this.axios.post("http://" + location.host + '/api/v1/auth/login', this.user);
        const result = await this.axios.post(
          "http://localhost:8080/api/v1/auth/register",
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.loginError = "Register failed";
        } else {
          this.user.token = result.data;
          this.connectToWebsocket();
        }
      } catch (e) {
        this.loginError = "Register failed";
        console.log(e);
      }
    },
    async login() {
      try {
        //const result = await this.axios.post("http://" + location.host + '/api/v1/auth/login', this.user);
        const result = await this.axios.post(
          "http://localhost:8080/api/v1/auth/login",
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.loginError = "Login failed";
        } else {
          this.user.token = result.data;
          this.connectToWebsocket();
        }
      } catch (e) {
        this.loginError = "Login failed";
        console.log(e);
      }
    },
    async changePass() {
      try {
        const result = await this.axios.put(
          "http://localhost:8080/api/v1/users/change-pwd?bearer=" +
            this.user.token,
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.loginError = "Change password failed 1";
        } else {
          this.user.password = result.data;
        }
      } catch (e) {
        this.loginError = "Change password failed 2";
        console.log(e);
      }
    },
    async changeName() {
      try {
        const result = await this.axios.post(
          "http://localhost:8080/api/v1/users/change_name",
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.loginError = "Change name failed";
        } else {
          this.user.name = result.data;
        }
      } catch (e) {
        this.loginError = "Change name failed";
        console.log(e);
      }
    },
    async changeAvatar() {
      try {
        const result = await this.axios.post(
          "http://localhost:8080/api/v1/users/change_avtr",
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.loginError = "Change avatar failed";
        } else {
          this.user.name = result.data;
        }
      } catch (e) {
        this.loginError = "Change avatar failed";
        console.log(e);
      }
    },
    connectToWebsocket() {
      if (this.user.token != "") {
        this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);
      } else {
        this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
      }
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
      if (this.currentReconnectDelay < this.maxReconnectDelay) {
        this.currentReconnectDelay *= 2;
      }
      this.connectToWebsocket();
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

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
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
