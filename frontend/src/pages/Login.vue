<template>
  <div class="login">
    <h2>Login</h2>
    <div class="col-12 form" v-if="!ws">
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
          <button
            class="input-group-text send_btn"
            :disabled="isDisabled()"
            @click="login"
          >
            >
          </button>
        </div>
      </div>

      <div class="alert alert-danger" role="alert" v-show="loginError">
        {{ loginError }}
      </div>
    </div>

    <div class="navigation">
      <button @click="navigate">Register Page</button>
    </div>
  </div>
</template>

<script>
import router from "../router";
import { wsConnect } from "@/services/WSConnectService";
export default {
  name: "LoginPage",
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
    navigate() {
      router.push({ path: "/register" });
    },
    isDisabled() {
      return !(this.user.username !== "" && this.user.password !== "");
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
          this.loginError = "Login failed 1";
        } else {
          wsConnect.user.token = result.data.token;
          wsConnect.user.uid = result.data.uid;
          wsConnect.user.name = this.user.username;
          //localStorage.setItem('name', this.user.username);
          //localStorage.setItem('token', this.user.token);
          //localStorage.setItem('uid', result.data.uid);
          //this.connect();
          router.push({ path: "/dashboard" });
        }
      } catch (e) {
        this.loginError = "Login failed 2";
        console.log(e);
        //router.push({ path: "/dashboard" });
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.login {
  text-align: center;
  position: absolute;
  top: 50%;
  left: 50%;
  margin-right: -50%;
  transform: translate(-50%, -50%);
}

h2 {
  margin-bottom: 40px;
}

.room,
.form {
  margin-top: auto;
  margin-bottom: auto;
}

.container {
  align-content: center;
}
.input-group {
  display: flex;
  justify-content: center;
}

.send_btn {
  border-radius: 0 15px 15px 0;
  background-color: rgba(0, 0, 0, 0.3);
  border: 0;
  color: white;
  cursor: pointer;
}

.navigation {
  margin-top: 40px;
}
</style>
