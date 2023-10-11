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
  },
  params: new URLSearchParams(),
});

export const requests = {
  listOutputs() {
    return api.request({
      url: '/outputs',
      method: 'get',
      params: {},
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
    });
  },
  getConfiguration() {
    return api.request({
      url: '/configuration',
      method: 'get',
      params: {},
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
    });
  },
  getVersion() {
    return api.request({
      url: '/version',
      method: 'get',
      params: {},
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
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
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
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
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
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
      auth: {
        username: store.state.username,
        password: store.state.password,
      },
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
};

export default {
};
