<!-- src/views/Login.vue -->
<template>
  <div class="login">
    <h2>登录</h2>
    <form @submit.prevent="handleLogin">
      <input type="text" placeholder="用户名" v-model="username" required />
      <input type="password" placeholder="密码" v-model="password" required />
      <button type="submit">登录</button>
    </form>
  </div>
</template>

<script>
import { useStore } from 'vuex';

export default {
  setup() {
    const store = useStore();
    const handleLogin = () => {
     // 调用API进行登录
    };
    return { handleLogin };
  },
  data() {
    return {
      username: '',
      password: ''
    };
  }
};
</script>

<style scoped>
.login {
  max-width: 300px;
  margin: 50px auto;
}
</style>