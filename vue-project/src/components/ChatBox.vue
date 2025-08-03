<template>
  <div class="chat-box">
    <div class="messages" ref="messagesContainer">
      <div
        v-for="(message, index) in messages"
        :key="index"
        class="message"
        :class="{ 'self': message.self, 'other': !message.self }"
      >
        <div :class="['message-bubble', { 'self-bubble': message.self }]">
          <div class="message-text">{{ message.text }}</div>
          <div class="message-timestamp">{{ message.timestamp }}</div>
        </div>
      </div>
    </div>
    <div class="input-area">
      <input
        type="text"
        v-model="inputMessage"
        @keyup.enter="sendMessage"
        placeholder="输入消息..."
      />
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
        self: true,
        timestamp: new Date().toLocaleTimeString()
      };
      this.messages.push(message);
      this.scrollToBottom();

      // 仅发送 text 字段给后端
      this.socket.send(JSON.stringify({ text: this.inputMessage }));
      this.inputMessage = '';
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
        try {
          const serverData = JSON.parse(event.data);
          const message = {
            text: serverData.text || '[返回为空]',
            self: false,
            timestamp: new Date().toLocaleTimeString()
          };
          this.messages.push(message);
          this.scrollToBottom();
        } catch (error) {
          console.error('接收到无效 JSON 消息:', event.data);
        }
      };

      this.socket.onclose = () => {
        console.log('WebSocket 连接已关闭');
      };

      this.socket.onerror = (error) => {
        console.log('WebSocket 错误:', error);
      };
    },
    disconnectWebSocket() {
      if (this.socket) {
        this.socket.close();
      }
    }
  },
  mounted() {
    const welcomeMessage = {
      text: '我是小智，有什么可以帮助你的吗？',
      self: false,
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
  margin-left: 250px;
  max-width: calc(100vw - 250px);
}

.messages {
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 10px;
  width: 100%;
}

.message {
  display: flex;
  align-items: flex-start;
  padding: 5px;
  margin: 5px 0;
  justify-content: flex-end;
}

.message.self {
  justify-content: flex-end;
}

.message.other {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 60%;
  padding: 8px 12px;
  border-radius: 20px;
  background-color: #b3e5fc;
}

.message.other .message-bubble {
  background-color: #e0e0e0;
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
  width: 100%;
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
