<template>
  <div class="sidebar">
    <div class="sidebar-item" @click="newSession">新建会话</div>
    <div class="sidebar-item">历史会话</div>
    <div class="user-info" @click="toggleUserInfo">
      <div class="user-details" :class="{ 'highlight': showUserInfo }">
        <div class="avatar" :style="{ backgroundImage: 'url(' + avatarUrl + ')' }"></div>
        <div class="username">{{ username }}</div>
        <div class="arrow">{{ showUserInfo ? '▲' : '▼' }}</div>
      </div>
    </div>
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
      avatarUrl: '',
    };
  },
  methods: {
    newSession() {
     console.log('新建会话');
   },
    toggleUserInfo() {
     this.showUserInfo = !this.showUserInfo;
   }
  }
};
</script>

<style scoped>
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
  justify-content: space-between;
  z-index: 10;
}

.sidebar-item {
  margin: 10px 0;
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  padding: 10px 0;
  border-top: 1px solid #ccc;
  cursor: pointer;
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
}
</style>