// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/* eslint-disable no-param-reassign */
import Vue from 'vue';
import Vuex from 'vuex';
import createPersistedState from 'vuex-persistedstate';

Vue.use(Vuex);

export const emptyFilters = () => ({
  priorities: [],
  sources: [],
  hostnames: [],
  tags: [],
  rule: '',
  search: '',
  since: '24h',
});

const LISTS = {
  priority: 'priorities',
  priorities: 'priorities',
  source: 'sources',
  sources: 'sources',
  hostname: 'hostnames',
  hostnames: 'hostnames',
  tag: 'tags',
  tags: 'tags',
};

function loadTheme() {
  try {
    return window.localStorage.getItem('fsui-theme') || 'falco';
  } catch (e) {
    return 'falco';
  }
}

export default new Vuex.Store({
  plugins: [createPersistedState({
    storage: window.sessionStorage,
    paths: ['refreshInterval', 'username', 'password'],
  })],
  state: {
    ticer: 0,
    refreshInterval: 10000,
    refreshIntervals: ['off', '10s', '20s', '30s', '1min', '2min'],
    username: '',
    password: '',
    filters: emptyFilters(),
    theme: loadTheme(),
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
    setCredentials(state, payload) {
      state.username = payload.username;
      state.password = payload.password;
    },
    emptyCredentials(state) {
      state.username = '';
      state.password = '';
    },
    setFilters(state, payload) {
      state.filters = { ...state.filters, ...payload };
    },
    toggleFilter(state, { list, item }) {
      if (list === 'rule') {
        state.filters.rule = state.filters.rule === item ? '' : item;
        return;
      }
      if (list === 'search' || list === 'filter') {
        state.filters.search = item;
        return;
      }
      const key = LISTS[list];
      if (!key) return;
      const current = state.filters[key];
      state.filters[key] = current.includes(item)
        ? current.filter(v => v !== item)
        : [...current, item];
    },
    addFilter(state, { list, item }) {
      if (list === 'rule') {
        state.filters.rule = item;
        return;
      }
      if (list === 'search' || list === 'filter') {
        state.filters.search = String(item);
        return;
      }
      const key = LISTS[list];
      if (!key || state.filters[key].includes(item)) return;
      state.filters[key] = [...state.filters[key], item];
    },
    clearFilters(state) {
      state.filters = { ...emptyFilters(), since: state.filters.since };
    },
    setTheme(state, theme) {
      state.theme = theme;
      try {
        window.localStorage.setItem('fsui-theme', theme);
      } catch (e) {
        // storage unavailable, keep in memory only
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
    setCredentials(context, payload) {
      context.commit('setCredentials', payload);
    },
    emptyCredentials(context) {
      context.commit('emptyCredentials');
    },
    setFilters(context, payload) {
      context.commit('setFilters', payload);
    },
    toggleFilter(context, payload) {
      context.commit('toggleFilter', payload);
    },
    addFilter(context, payload) {
      context.commit('addFilter', payload);
    },
    clearFilters(context) {
      context.commit('clearFilters');
    },
    setTheme(context, theme) {
      context.commit('setTheme', theme);
    },
  },
});
/* eslint-enable no-param-reassign */
