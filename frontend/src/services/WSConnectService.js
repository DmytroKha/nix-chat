class WSConnectService {

  ws = null;
  currentReconnectDelay = null;
  maxReconnectDelay = null;
  token = null;
  serverUrl = "ws://" + location.host + "/app/ws"
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
    if (!this.ws) {
        this.ws = new WebSocket(this.serverUrl + "?bearer=" + this.user.token);
    }
    this.ws.addEventListener('open', (event) => {
      this.onWebsocketOpen(event)
    });
    this.ws.addEventListener('close', (event) => {
      this.onWebsocketClose(event)
    });
  }

  onWebsocketOpen() {
    console.log("connected to WS!");
    console.log("users",this.users);
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
    console.log("all-rooms out")
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