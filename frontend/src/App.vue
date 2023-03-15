<template>
  <router-view />
</template>

<script>
// import Dashboard from "./components/pages/Dashboard.vue";

export default {
  name: "App",
  components: {
    // Dashboard,
  },
};
</script>

<style>
*{
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}
button{
  padding: 3px 5px;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100%;
}

body,
html {
  height: 100%;
  margin: 0;
  background: #827c87;
  background: -webkit-linear-gradient(to right, #827c87, #5a5464, #7f7585);
  background: linear-gradient(to right, #827c87, #5a5464, #7f7585)

  
}
.header {
  height: 60px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #83b6ef;
  padding: 0 10px;
}
</style>
