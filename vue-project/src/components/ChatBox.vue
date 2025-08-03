<template>
  <div class="chat-box">
    <div class="messages" ref="messagesContainer">
      <div v-for="(message, index) in messages" :key="index" class="message" :class="{ 'self': message.self, 'other': !message.self }">
        <div :class="['message-bubble', { 'self-bubble': message.self }]">
          <div class="message-text">{{ message.text }}</div>
          <div class="message-timestamp">{{ message.timestamp }}</div>
        </div>
      </div>
    </div>
    <div class="input-area">
      <input type="text" v-model="inputMessage" @keyup.enter="sendMessage" placeholder="输入消息..." />
      <button @click="sendMessage">发送</button>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      messages: [],
      inputMessage: '',
      socket: null
    };
  },
  methods: {
    sendMessage() {
      if (this.inputMessage.trim() === '') {
        console.log('输入不能为空');
        return;
      }
      const message = {
        text: this.inputMessage,
        self: true, // 发送消息为 self
        timestamp: new Date().toLocaleTimeString()
      };
      this.messages.push(message);
      this.inputMessage = '';
      this.scrollToBottom();
      this.socket.send(JSON.stringify(message));
      console.log('发送消息:', message);
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const messagesContainer = this.$refs.messagesContainer;
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
      });
    },
    connectWebSocket() {
      this.socket = new WebSocket('ws://localhost:8080/ws');
      this.socket.onopen = () => {
        console.log('WebSocket 连接已建立');
      };
      this.socket.onmessage = (event) => {
        const message = JSON.parse(event.data);
        // 确保接收到的消息是 `self: false`
        message.self = false; // 设置接收到的消息为非自己发送
        this.messages.push(message);
        this.scrollToBottom();
      };
      this.socket.onclose = () => {
        console.log('WebSocket 连接已关闭');
      };
      this.socket.onerror = (error) => {
        console.log('WebSocket 连接发生错误:', error);
      };
    },
    disconnectWebSocket() {
      if (this.socket) {
        this.socket.close();
      }
    }
  },
  mounted() {
    // 在页面加载时显示系统的欢迎消息
    const welcomeMessage = {
      text: '我是小智，有什么可以帮助你的吗？',
      self: false,  // 系统消息，显示在左边
      timestamp: new Date().toLocaleTimeString()
    };
    this.messages.push(welcomeMessage);
    this.scrollToBottom();

    this.connectWebSocket();
  },
  beforeUnmount() {
    this.disconnectWebSocket();
  }
};
</script>

<style scoped>
.chat-box {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px;
  margin-left: 250px; /* Adjust this value to match the sidebar width */
  max-width: calc(100vw - 250px); /* Ensure chat box does not exceed sidebar width */
}

.messages {
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 10px;
  width: 100%; /* Ensure messages take full width */
}

.message {
  display: flex;
  align-items: flex-start;
  padding: 5px;
  margin: 5px 0;
  justify-content: flex-end; /* 默认右对齐 */
}

.message.self {
  justify-content: flex-end; /* 自己发送的消息右对齐 */
}

.message.other {
  justify-content: flex-start; /* 其他消息左对齐 */
}

.message-bubble {
  max-width: 60%;
  padding: 8px 12px;
  border-radius: 20px;
  background-color: #b3e5fc; /* Light blue for self messages */
}

.message.other .message-bubble {
  background-color: #e0e0e0; /* Light grey for other messages */
}

.message-text {
  margin-bottom: 2px;
}

.message-timestamp {
  font-size: 12px;
  color: #999;
}

.input-area {
  display: flex;
  align-items: center;
}

input {
  flex-grow: 1;
  padding: 5px;
  margin-right: 5px;
}

button {
  padding: 5px 10px;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #ddd;
}
</style>
