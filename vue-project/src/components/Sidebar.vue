<template>
  <div class="sidebar">
    <!-- 新建会话按钮 -->
    <div class="new-chat-button" @click="newSession">
      <i class="fas fa-comment chat-icon"></i>
      New Chat
    </div>

    <!-- 历史会话标题 -->
    <div class="sidebar-item" v-if="sessions.length > 0">历史会话</div>

    <!-- 会话列表 -->
    <div v-for="session in sessions" :key="session.id" class="session-item">
      <i class="fas fa-comment chat-icon"></i>
      {{ session.title }}
    </div>

    <!-- 用户信息 -->
    <div class="user-info" @click="toggleUserInfo">
      <div class="user-details" :class="{ 'highlight': showUserInfo }">
        <div class="avatar" :style="{ backgroundImage: 'url(' + avatarUrl + ')' }"></div>
        <div class="username">{{ username }}</div>
        <div class="arrow">{{ showUserInfo ? '▲' : '▼' }}</div>
      </div>
    </div>

    <!-- 用户详细信息 -->
    <div v-if="showUserInfo" class="user-details-info">
      <p>邮箱: user@example.com</p>
      <p>注册日期: 2024-01-01</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      showUserInfo: false,
      username: 'Re.As.Cx.Co',
      avatarUrl: '', // 示例头像 URL
      sessions: [] // 会话数组
    };
  },
  created() {
    this.loadSessions();
  },
  methods: {
    newSession() {
      const newSessionId = Date.now();
      const newSessionTitle = '新会话 ' + (this.sessions.length + 1);
      this.sessions.push({ id: newSessionId, title: newSessionTitle });
      console.log('新建会话:', newSessionTitle);
    },
    toggleUserInfo() {
      this.showUserInfo = !this.showUserInfo;
    },
    loadSessions() {
      this.sessions = [
        { id: 1, title: '会话1' },
        { id: 2, title: '会话2' },
        { id: 3, title: '会话3' },
      ];
    }
  }
};
</script>

<style scoped>
/* 引入 Font Awesome 图标所需样式（外部引入 HTML 文件中加上） */
/* <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.0/css/all.min.css"> */

.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 250px;
  background-color: #f5f5f5;
  padding: 10px;
  height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  z-index: 10;
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
  color: #666;
  flex-grow: 1;
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
</style>
