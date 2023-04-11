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
            <img :src="user.photo" width="100"/>
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

          <div class="alert alert-danger" role="alert" v-show="profileError">
            {{ profileError }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import router from "../router";
import { wsConnect } from "@/services/WSConnectService";
export default {
  name: "ProfilePage",
  data() {
    return {
      user: {
        name: "",
        oldPassword: "",
        newPassword: "",
        photo: "",
      },
      profileError: "",
    };
  },
  mounted: function () {
    this.user.name = wsConnect.user.name;
    this.user.photo = wsConnect.user.photo
  },
  methods: {
    navigate() {
      router.push({ path: "/dashboard" });
    },
    async changePass() {
      try {
        const result = await this.axios.put(
          "http://localhost:8080/api/v1/users/change-pwd?bearer=" +
            wsConnect.user.token,
            this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.profileError = "Change password failed 1";
        } else {
          this.profileError = "Password was changed successfully";
        }
      } catch (e) {
        this.profileError = "Change password failed 2";
        console.log(e);
      }
    },
    async changeName() {
      try {
        const result = await this.axios.put(
          "http://localhost:8080/api/v1/users/change-name?bearer=" +
            wsConnect.user.token,
            this.user
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.profileError = "Change name failed";
        } else {
          wsConnect.user.name = result.data;
          for (let i = 0; i < wsConnect.users.length; i++) {
            if (wsConnect.users[i].id == wsConnect.user.id) {
              wsConnect.ws.send(
                  JSON.stringify({ action: "user-left", sender: wsConnect.users[i] })
              );
              wsConnect.users[i].name = wsConnect.user.name;
              wsConnect.ws.send(
                  JSON.stringify({ action: "user-join", sender: wsConnect.users[i]})
              );
              break;
            }
          }
          this.profileError = "Name was changed successfully";
        }
      } catch (e) {
        this.profileError = "Change name failed";
        console.log(e);
      }
    },
    async changeAvatar(e) {
      var formData = new FormData();
      // //Attach file
      formData.append("image", e.target.files[0]);
      console.log(formData)
      try {
        const result = await this.axios.put(
            "http://localhost:8080/api/v1/users/change_avtr?bearer="+
            wsConnect.user.token,
            e.target.files[0]
        );
        if (
          result.data.status !== "undefined" &&
          result.data.status == "error"
        ) {
          this.profileError = "Change avatar failed";
        } else {
          wsConnect.user.photo = result.data
          this.user.photo = result.data
          for (let i = 0; i < wsConnect.users.length; i++) {
            if (wsConnect.users[i].id == wsConnect.user.id) {
              wsConnect.ws.send(
                  JSON.stringify({ action: "user-left", sender: wsConnect.users[i] })
              );
              wsConnect.users[i].photo = wsConnect.user.photo;
              wsConnect.ws.send(
                  JSON.stringify({ action: "user-join", sender: wsConnect.users[i]})
              );
              break;
            }}
          this.profileError = "Avatar was changed successfully";
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
