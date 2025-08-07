<template>
  <nav class="navbar">
    <h1 class="navbar-title">AI 智能助手</h1>
  </nav>
</template>

<script>
export default {
  name: 'NavBar'
};
</script>
<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  height: 60px;
  width: 100%;
  background-color: #4f46e5;
  color: white;
  display: flex;
  align-items: center;
  padding: 0 20px;
  z-index: 100;
  box-sizing: border-box;
}


.navbar-title {
  font-size: 20px;
  font-weight: bold;
}
</style>
