class WSConnectService {

  ws = null;
  currentReconnectDelay = null;
  maxReconnectDelay = null;
  token = null;
  serverUrl = "ws://localhost:8080/ws"
  users = []
  rooms = []
  chatRooms = []
  user = {
    id: "",
    name: "",
    username: "",
    password: "",
    confirmPassword: "",
    oldPassword: "",
    newPassword: "",
    token: "",
    friends: [],
    blackList: [],
    photo: ""
  }

  connectToWebsocket() {
    //this.token = localStorage.getItem('token');
    //if (this.token != "") {
    if (!this.ws) {
        this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);
        // } else {
        //   this.ws = new WebSocket(this.serverUrl + "?name=" + this.user.name);
      //}
    }
    this.ws.addEventListener('open', (event) => {
      this.onWebsocketOpen(event)
    });
    //this.ws.addEventListener('message', (event) => { this.handleNewMessage(event) });
    this.ws.addEventListener('close', (event) => {
      this.onWebsocketClose(event)
    });
  }

  onWebsocketOpen() {
    console.log("connected to WS!");
    this.getAllRooms();
    this.getBlackList();
    this.getFriendList();
    this.currentReconnectDelay = 1000;
  }

  onWebsocketClose() {
    this.ws = null;

    setTimeout(() => {
      this.reconnectToWebsocket();
    }, this.currentReconnectDelay);
  }

  reconnectToWebsocket() {
    if (this.currentReconnectDelay < this.maxReconnectDelay) {
      this.currentReconnectDelay *= 2;
    }
    this.connectToWebsocket();
  }

  getAllRooms() {
    wsConnect.ws.send(
        JSON.stringify({ action: "all-rooms"})
    );
  }

  getBlackList() {
    wsConnect.ws.send(
        JSON.stringify({action: "get-black-list"})
    );
  }

  getFriendList() {
    wsConnect.ws.send(
        JSON.stringify({action: "get-friends"})
    );
  }
  }

export const wsConnect = new WSConnectService();