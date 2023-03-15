<template>
  <div class="hello">
    <div id="app">
      <div class="container h-100">
        <div class="header">
          <h1>Profile</h1>
          <button @click="navigate">dasboard</button>
        </div>
        <div class="row justify-content-center h-100">

          <div>
            <h2>Change avatar</h2>
            <input
              type="file"
              id="avatar"
              name="avatar"
              accept="image/png, image/jpeg"
              @change="changeAvatar"
            />
          </div>

          <div>
            <h2>Change user name</h2>

            <div class="input-group">
              <input
                v-model="user.name"
                class="form-control username"
                placeholder="name"
              />

              <div class="input-group-append">
                <span class="input-group-text send_btn" @click="changeName">
                  >
                </span>
              </div>
            </div>
          </div>

          <div>
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
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import router from "../router";
export default {
  name: "ProfilePage",
  data() {
    return {
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
      profileError: "",
    };
  },
  beforeMount() {
    this.getUserData();
  },
  mounted: function () {
    this.user.name = localStorage.getItem('name');
    this.user.token = localStorage.getItem('token');
    this.user.uid = localStorage.getItem('uid');
  },
  methods: {
    navigate() {
      router.push({ path: "/dashboard" });
      //localStorage.removeItem('name');
      //localStorage.removeItem('token');
      //localStorage.removeItem('uid');
    },
    async getUserData() {
      try {
        const result = await this.axios.get(
          "http://localhost:8080/api/v1/users/me"
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          console.log("Change password failed 1");
        } else {
          this.user = result.data;
        }
      } catch (e) {
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
          this.profileError = "Change password failed 1";
        } else {
          this.user.password = result.data;
        }
      } catch (e) {
        this.profileError = "Change password failed 2";
        console.log(e);
      }
    },
    async changeName() {
      try {
        const result = await this.axios.put(
          "http://localhost:8080/api/v1/users/change-name",
          this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.profileError = "Change name failed";
        } else {
          this.user.name = result.data;
        }
      } catch (e) {
        this.profileError = "Change name failed";
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
          this.profileError = "Change avatar failed";
        } else {
          this.user.name = result.data;
        }
      } catch (e) {
        this.profileError = "Change avatar failed";
        console.log(e);
      }
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
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
</style>
