<template>
  <div class="chat-interface">
    <div class="sidebar">
      <div class="sidebar-header">
        <h2>功能菜单</h2>
      </div>
      <ul class="sidebar-list">
        <li @click="handleLogin">登录</li>
        <li>Kimi+</li>
        <li>PPT 助手</li>
        <li>
          历史会话
          <ul class="submenu">
            <li v-for="(session, index) in sessions" :key="index">{{ session }}</li>
          </ul>
        </li>
      </ul>
    </div>
    <div class="chat-box">
      <div class="chat-header">
        <h2>对话</h2>
      </div>
      <div class="chat-messages" v-for="(message, index) in messages" :key="index">
        <div :class="{'self': message.self, 'other': !message.self}">
          {{ message.content }}
        </div>
      </div>
      <div class="chat-input">
        <input type="text" v-model="inputMessage" @keyup.enter="sendMessage" placeholder="输入消息..." />
        <button @click="sendMessage">发送</button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      sessions: ['会话1', '会话2', '会话3'],
      messages: [
        { content: '你好！', self: false },
        { content: '你好，有什么可以帮助你的吗？', self: true }
      ],
      inputMessage: ''
    };
  },
  methods: {
    handleLogin() {
      // 登录逻辑
      console.log('登录');
    },
    sendMessage() {
      if (this.inputMessage.trim()) {
        this.messages.push({ content: this.inputMessage, self: true });
        this.inputMessage = '';
        // 发送消息到后端的逻辑
        console.log('发送消息:', this.inputMessage);
      }
    }
  }
};
</script>

<style scoped>
.chat-interface {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 250px;
  background-color: #f5f5f5;
  padding: 10px;
}

.sidebar-header {
  font-size: 18px;
  margin-bottom: 10px;
}

.sidebar-list {
  list-style: none;
  padding: 0;
}

.sidebar-list li {
  margin: 10px 0;
  cursor: pointer;
}

.submenu {
  list-style: none;
  padding-left: 20px;
}

.chat-box {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  padding: 10px;
}

.chat-header {
  font-size: 18px;
  margin-bottom: 10px;
}

.chat-messages {
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 10px;
}

.chat-input {
  display: flex;
  align-items: center;
}

.chat-input input {
  flex-grow: 1;
  padding: 5px;
  margin-right: 5px;
}

.chat-input button {
  padding: 5px 10px;
}

.self {
  align-self: flex-end;
  background-color: #e0e0e0;
  padding: 5px;
  border-radius: 5px;
}

.other {
  align-self: flex-start;
  background-color: #d1e7dd;
  padding: 5px;
  border-radius: 5px;
}
</style>