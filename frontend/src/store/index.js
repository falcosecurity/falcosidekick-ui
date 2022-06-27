/* eslint-disable no-param-reassign */
import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    ticer: 0,
    refreshInterval: 10000,
    refreshIntervals: ['off', '10s', '20s', '30s', '1min', '2min'],
  },
  mutations: {
    increment(state) {
      state.ticer += 1;
    },
    setRefreshInterval(state, payload) {
      switch (payload) {
        case '10s':
          state.refreshInterval = 10000;
          break;
        case '20s':
          state.refreshInterval = 20000;
          break;
        case '30s':
          state.refreshInterval = 30000;
          break;
        case '1min':
          state.refreshInterval = 60000;
          break;
        case '2min':
          state.refreshInterval = 120000;
          break;
        case 'off':
          state.refreshInterval = 0;
          break;
        default:
          break;
      }
    },
  },
  actions: {
    increment(context) {
      context.commit('increment');
    },
    setRefreshInterval(context, payload) {
      context.commit('setRefreshInterval', payload);
    },
  },
});
/* eslint-enable no-param-reassign */
