<template>
  <div id="app">
    <!-- 顶部导航栏 -->
    <NavBar />
    <!-- 页面主体，左边侧边栏 + 右边聊天框 -->
    <div class="layout">
      <Sidebar @new-session="handleNewSession" @select-session="handleSelectSession" />
      <ChatBox :sessionId="currentSession ? currentSession.id : null" />
    </div>
  </div>
</template>


<script>
import Sidebar from './components/Sidebar.vue';
import ChatBox from './components/ChatBox.vue';
import NavBar from './components/NavBar.vue';
export default {
  components: { Sidebar, ChatBox,NavBar  },
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
  margin-top:40px;
}
</style>
