<template>
  <div class="register">
    <div class="col-12 form" v-if="!ws">
      <h2>Register New user</h2>

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
          <button
            class="input-group-text send_btn"
            :disabled=isDisabled
            @click="register"
          >
            &gt;
          </button>
        </div>
      </div>
      <div class="alert alert-danger" role="alert" v-show="registerError">
        {{ registerError }}
      </div>

      <div class="navigation">
        <button><router-link to="/">Login Page</router-link></button>
      </div>
    </div>
  </div>
</template>

<script>
import router from "../router";
export default {
  name: "RegisterPage",
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
      registerError: "",
    };
  },
  computed: {
    isDisabled() {
      return !(
          this.user.confirmPassword === this.user.password &&
          this.user.username !== "" &&
          this.user.password !== ""
      );
    },
  },
  mounted: function () {},
  methods: {
    navigate() {
      router.push({ path: "/" });
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
          this.registerError = "Register failed";
        } else {
          this.user.token = result.data;
        }
      } catch (e) {
        this.registerError = "Register failed";
        console.log(e);
      }
    },
  },
};
</script>

<style scoped>
.register {
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
.input-group {
  display: flex;
  justify-content: center;
}
.container {
  align-content: center;
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
