<template>
  <div class="chat-box">
    <div class="messages" ref="messagesContainer">
      <div
        v-for="(message, index) in messages"
        :key="message.id || index"
        class="message"
        :class="{ self: message.self, other: !message.self }"
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
  props: ['sessionId'], // 可以用来标识当前会话，但这里没用到
  data() {
    return {
      messages: [],
      inputMessage: '',
      socket: null,
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
        timestamp: new Date().toLocaleTimeString(),
        id: Date.now() + Math.random(),
      };
      this.messages.push(message);
      this.scrollToBottom();

      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify({ text: this.inputMessage }));
      } else {
        console.warn('WebSocket 未连接或未打开，发送失败');
      }
      this.inputMessage = '';
    },
    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) container.scrollTop = container.scrollHeight;
      });
    },
    connectWebSocket() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
      this.socket = new WebSocket('ws://localhost:8080/ws');

      this.socket.onopen = () => {
        console.log('✅ WebSocket 连接已建立');
      };

      this.socket.onmessage = (event) => {
        console.log('[收到后端消息]：', event.data);
        try {
          const data = JSON.parse(event.data);
          const text = data.text ? data.text.trim() : '[无内容]';
          const timestamp = data.timestamp
            ? new Date(data.timestamp).toLocaleTimeString()
            : new Date().toLocaleTimeString();
          this.messages.push({
            text,
            self: false,
            timestamp,
            id: Date.now() + Math.random(),
          });
          this.scrollToBottom();
        } catch (error) {
          console.error('❌ 接收消息解析失败:', error, event.data);
        }
      };

      this.socket.onerror = (error) => {
        console.error('❌ WebSocket 错误:', error);
      };

      this.socket.onclose = (event) => {
        console.warn('⚠️ WebSocket 连接关闭', event);
      };
    },
    disconnectWebSocket() {
      if (this.socket) {
        this.socket.close();
        this.socket = null;
      }
    }
  },
  mounted() {
    this.messages.push({
      text: '我是小智，有什么可以帮助你的吗？',
      self: false,
      timestamp: new Date().toLocaleTimeString(),
      id: Date.now(),
    });
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
  height: 100vh;
  box-sizing: border-box;
}

.messages {
  flex-grow: 1;
  overflow-y: auto;
  margin-bottom: 10px;
  width: 100%;
  background: #f9f9f9;
  border-radius: 8px;
  padding: 10px;
  box-sizing: border-box;
}

.message {
  display: flex;
  align-items: flex-start;
  padding: 5px;
  margin: 5px 0;
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
  word-break: break-word;
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
  text-align: right;
}

.input-area {
  display: flex;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
}

input {
  flex-grow: 1;
  padding: 8px;
  margin-right: 5px;
  border: 1px solid #ccc;
  border-radius: 4px;
  font-size: 14px;
}

button {
  padding: 8px 16px;
  cursor: pointer;
  background-color: #1976d2;
  border: none;
  border-radius: 4px;
  color: white;
  font-weight: bold;
  transition: background-color 0.3s ease;
}

button:hover {
  background-color: #1565c0;
}
</style>
