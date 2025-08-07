<!-- src/views/Register.vue -->
<template>
  <div class="register">
    <h2>注册</h2>
    <form @submit.prevent="handleRegister">
      <input type="text" placeholder="用户名" v-model="newUsername" required />
      <input type="email" placeholder="邮箱" v-model="newEmail" required />
      <input type="password" placeholder="密码" v-model="newPassword" required />
      <input type="password" placeholder="确认密码" v-model="confirmPassword" required />
      <button type="submit">注册</button>
    </form>
  </div>
</template>

<script>
import { useStore } from 'vuex';

export default {
  setup() {
    const store = useStore();
    const handleRegister = () => {
      // 调用API进行注册
    };
    return { handleRegister };
  },
  data() {
    return {
      newUsername: '',
      newEmail: '',
      newPassword: '',
      confirmPassword: ''
    };
  }
};
</script>

<style scoped>
.register {
  max-width: 300px;
 margin: 50px auto;
}
</style>