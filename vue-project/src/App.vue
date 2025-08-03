<template>
  <div id="app" class="layout">
    <Sidebar />
    <ChatBox />
  </div>
</template>

<script>
import Sidebar from './components/Sidebar.vue';
import ChatBox from './components/ChatBox.vue';

export default {
  components: {
    Sidebar,
    ChatBox
  }
};
</script>

<style>
.layout {
  display: flex;
  height: 100vh;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  margin-top: 60px;
}
</style>