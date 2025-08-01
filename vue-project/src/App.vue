<template>
  <div id="app">
    <WebSocketChat />
  </div>
</template>

<script>
import WebSocketChat from './components/WebSocketChat.vue';

export default {
  components: {
    WebSocketChat
  }
};
</script>

<style>
/* 添加一些全局样式 */
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>