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


import axios from 'axios';
import store from './store';

const production = process.env.NODE_ENV === 'production';

const api = axios.create({
  baseURL: `${production ? `//${window.location.host}${window.location.pathname}` : process.env.VUE_APP_API}api/v1`,
  headers: {
    'Content-type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': '*',
    'X-Requested-With': 'XMLHttpRequest',
  },
  params: new URLSearchParams(),
});

// Response interceptor for handling 401 in OIDC mode
api.interceptors.response.use(
  response => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      if (store.state.authMode === 'oidc') {
        store.commit('clearSSO');
        window.location.href = `${window.location.pathname}#/login`;
      }
    }
    return Promise.reject(error);
  },
);

const getAuthConfig = () => {
  if (store.state.authMode === 'oidc' || store.state.authMode === 'none') {
    return undefined;
  }
  return {
    username: store.state.username,
    password: store.state.password,
  };
};

export const requests = {
  authMe() {
    return api.request({
      url: '/auth/me',
      method: 'get',
      params: {},
    });
  },
  listOutputs() {
    return api.request({
      url: '/outputs',
      method: 'get',
      params: {},
      auth: getAuthConfig(),
    });
  },
  getConfiguration() {
    return api.request({
      url: '/configuration',
      method: 'get',
      params: {},
      auth: getAuthConfig(),
    });
  },
  getVersion() {
    return api.request({
      url: '/version',
      method: 'get',
      params: {},
      auth: getAuthConfig(),
    });
  },
  countEvents(source, hostname, priority, rule, filter, tags, since) {
    return api.request({
      url: '/events/count',
      method: 'get',
      params: {
        source: `${source}`,
        hostname: `${hostname}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
      },
      auth: getAuthConfig(),
    });
  },
  countByEvents(group, source, hostname, priority, rule, filter, tags, since) {
    return api.request({
      url: `/events/count/${group}`,
      method: 'get',
      params: {
        source: `${source}`,
        hostname: `${hostname}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
      },
      auth: getAuthConfig(),
    });
  },
  searchEvents(source, hostname, priority, rule, filter, tags, since, page, limit) {
    return api.request({
      url: '/events/search',
      method: 'get',
      params: {
        source: `${source}`,
        hostname: `${hostname}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
        page: `${page}`,
        limit: `${limit}`,
      },
      auth: getAuthConfig(),
    });
  },
  authenticate(username, password) {
    return api.request({
      url: '/auth',
      method: 'post',
      auth: {
        username: `${username}`,
        password: `${password}`,
      },
    });
  },
  logout() {
    return api.request({
      url: '/auth/logout',
      method: 'post',
      headers: {
        'X-Requested-With': 'XMLHttpRequest',
      },
    });
  },
};

export default {
};
