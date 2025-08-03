<template>
  <div id="app">
    <div class="layout">
      <Sidebar @new-session="handleNewSession" @select-session="handleSelectSession" />
      <!-- 这里仅传currentSessionId，或干脆不传 -->
      <ChatBox :sessionId="currentSession ? currentSession.id : null" />
    </div>
  </div>
</template>

<script>
import Sidebar from './components/Sidebar.vue';
import ChatBox from './components/ChatBox.vue';

export default {
  components: { Sidebar, ChatBox },
  data() {
    return {
      sessions: [],
      currentSession: null,
    };
  },
  methods: {
    handleNewSession(newSession) {
      this.sessions.push(newSession);
      this.currentSession = newSession;
    },
    handleSelectSession(session) {
      this.currentSession = session;
    },
    loadSessions() {
      const historicalSessions = [
        { id: 1, title: '会话1' },
        { id: 2, title: '会话2' },
        { id: 3, title: '会话3' },
      ];
      this.sessions = historicalSessions;
      this.currentSession = this.sessions[0];
    }
  },
  created() {
    this.loadSessions();
  }
};
</script>

<style>
.layout {
  display: flex;
  height: calc(100vh - 60px);
  margin-top: 60px;
}
</style>
