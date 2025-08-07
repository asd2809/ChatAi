import { createStore } from 'vuex';

export default createStore({
  state() {
    return {
     isLoggedIn: false, // 用户登录状态
   };
  },
  mutations: {
    setLoginStatus(state, status) {
     state.isLoggedIn = status;
   },
 },
  actions: {
    login({ commit }) {
      commit('setLoginStatus', true);
    },
    logout({ commit }) {
      commit('setLoginStatus', false);
    },
  },
});