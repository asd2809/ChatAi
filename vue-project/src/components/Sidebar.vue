<template>
  <div class="sidebar">
    <!-- 新建会话按钮 -->
    <div class="new-chat-button" @click="newSession">
      <i class="fas fa-comment chat-icon"></i>
      New Chat
    </div>

    <!-- 历史会话标题 -->
    <div class="sidebar-item" v-if="sessions.length > 0">History</div>

    <!-- 会话列表 -->
    <div v-for="session in sessions" :key="session.id" class="session-item">
      <i class="fas fa-comment chat-icon"></i>
      {{ session.title }}
    </div>

    <!-- 用户信息 -->
<div class="user-info" @click="toggleUserInfo">
  <div class="user-details" :class="{ 'highlight': showUserInfo }">
    <div class="avatar" :style="{ backgroundImage: 'url(' + avatarUrl + ')' }"></div>
    <div class="username">
      <template v-if="isLoggedIn">
        {{ username }}
      </template>
      <template v-else>
        请登录
      </template>
    </div>
    <div class="arrow">{{ showUserInfo ? '▲' : '▼' }}</div>
  </div>
</div>


    <!-- 用户详细信息 -->
    <div v-if="showUserInfo" class="user-details-info">
      <p>E-mail: user@example.com</p>
      <p>Registered: 2024-01-01</p>
      <button v-if="!isLoggedIn" @click="toggleRegisterModal">Register</button>
    </div>

    <!-- 注册模态框 -->
    <div v-if="showRegisterModal" class="modal">
      <div class="modal-content">
        <span class="close" @click="toggleRegisterModal">&times;</span>
        <h2>Register</h2>
        <form @submit.prevent="handleRegister">
          <input type="text" placeholder="Username" v-model="newUsername" required />
          <input type="email" placeholder="Email" v-model="newEmail" required />
          <input type="password" placeholder="Password" v-model="newPassword" required />
          <input type="password" placeholder="Confirm Password" v-model="confirmPassword" required />
          <button type="submit">Register</button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';

const router = useRouter();
const store = useStore();

const isLoggedIn = computed(() => store.state.isLoggedIn);

const newUsername = ref('');
const newEmail = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const showRegisterModal = ref(false);
const showUserInfo = ref(false);
const sessions = ref([]);

const toggleRegisterModal = () => {
  showRegisterModal.value = !showRegisterModal.value;
};

const handleRegister = async () => {
  if (newPassword.value !== confirmPassword.value) {
    alert('Passwords do not match');
    return;
  }
  try {
    const response = await fetch('/api/register', {
       method: 'POST',
       headers: { 'Content-Type': 'application/json' },
       body: JSON.stringify({
         username: newUsername.value,
         email: newEmail.value,
         password: newPassword.value
       })
    });
    const data = await response.json();
    console.log('Registration successful:', data);
    router.push('/login');
  } catch (error) {
    console.error('Registration failed:', error);
    alert('An error occurred during registration');
  }
};

const toggleUserInfo = () => {
  showUserInfo.value = !showUserInfo.value;
};

const loadSessions = () => {
  sessions.value = [
    { id: 1, title: 'Session 1' },
    { id: 2, title: 'Session 2' },
    { id: 3, title: 'Session 3' },
  ];
};

const newSession = () => {
  const newSessionId = Date.now();
  const newSessionTitle = 'New Chat ' + (sessions.value.length + 1);
  sessions.value.push({ id: newSessionId, title: newSessionTitle });
  console.log('New session:', newSessionTitle);
};

onMounted(() => {
  loadSessions();
});
</script>

<style scoped>
.sidebar {
  width: 250px;
  background-color: #f5f5f5;
  padding: 10px;
  height: 100%;
  display: flex;
  flex-direction: column;
  position: relative;
  box-sizing: border-box;
}

.new-chat-button {
  background-color: #4a6ef5;
  color: white;
  border-radius: 15px;
  padding: 15px;
  margin: 20px 0;
  text-align: center;
  font-size: 16px;
  font-weight: bold;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80%;
  cursor: pointer;
  margin-left: auto;
  margin-right: auto;
}

.new-chat-button:hover {
  background-color: #3b5bd8;
}

.chat-icon {
  font-size: 18px;
  color: white;
  margin-right: 8px;
}

.sidebar-item {
  font-weight: bold;
  color: #444;
  margin: 10px 0 5px 10px;
  font-size: 14px;
}

.session-item {
  background-color: #fff;
  color: #333;
  border-radius: 5px;
  padding: 10px;
  margin: 10px 0;
  font-size: 14px;
  display: flex;
  align-items: center;
  cursor: pointer;
}

.session-item:hover {
  background-color: #f0f0f0;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-top: 1px solid #ccc;
  cursor: pointer;
  margin-top: auto;
}

.user-details {
  display: flex;
  align-items: center;
  width: 100%;
}

.user-details.highlight {
  background-color: #e0e0e0;
}

.avatar {
  width: 50px;
  height: 50px;
  background-color: #ccc;
  background-size: cover;
  background-position: center;
  border-radius: 50%;
  margin-right: 10px;
}

.username {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}
.username::before {
  content: '';
  display: inline-block;
  margin-right: 4px;
}


.arrow {
  font-size: 12px;
  color: #666;
}

.user-details-info {
  margin-top: 10px;
  padding: 10px;
  background-color: #fff;
  border: 1px solid #ccc;
  border-radius: 5px;
  position: absolute;
  bottom: 0;
  width: 100%;
}

.modal {
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0, 0.5);
}

.modal-content {
  background-color: #fff;
  padding: 20px;
  border-radius: 5px;
  width: 300px;
}

.close {
  float: right;
  font-size: 28px;
  font-weight: bold;
  cursor: pointer;
}
</style>