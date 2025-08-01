<template>
  <div class="chat-container">
    <div id="chat-box" class="chat-box">
      <div v-for="msg in messages" :key="msg.id" class="message">
        <span :class="msg.isSelf ? 'user-message' : 'bot-message'">
          {{ msg.text }}
        </span>
      </div>
    </div>
    <input type="text" v-model="inputMessage" placeholder="Type a message...">
    <button @click="sendMessage">Send</button>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
import moment from 'moment';

export default {
  setup() {
    const messages = ref([]);
    const inputMessage = ref('');
    const ws = new WebSocket('ws://localhost:8080/ws');

    onMounted(() => {
      ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        messages.value.push({
          text: message.text,
          isSelf: false,
          time: moment().format('HH:mm:ss'),
          id: Date.now(),
        });
      };
    });

    onUnmounted(() => {
      ws.close();
    });

    const sendMessage = () => {
      const message = inputMessage.value.trim();
      if (message) {
        ws.send(JSON.stringify({ text: message, isSelf: true }));
        messages.value.push({
          text: message,
          isSelf: true,
          time: moment().format('HH:mm:ss'),
          id: Date.now(),
        });
        inputMessage.value = '';
      }
    };

    return {
      messages,
      inputMessage,
      sendMessage,
    };
  },
};
</script>

<style scoped>
.chat-container {
  display: flex;
  flex-direction: column;
  height: 500px;
  border: 1px solid #ccc;
  padding: 10px;
}

.chat-box {
  flex: 1;
  overflow-y: auto;
  padding: 5px;
  margin-bottom: 10px;
}

.message {
  display: flex;
  align-items: flex-start;
  margin-bottom: 5px;
}

.user-message {
  margin-left: auto;
  background-color: #007bff;
  color: white;
  padding: 5px 10px;
  border-radius: 5px;
}

.bot-message {
  background-color: #f8f9fa;
  padding: 5px 10px;
  border-radius: 5px;
}
</style>