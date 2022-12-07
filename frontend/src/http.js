import axios from 'axios';
import store from './store';

const production = process.env.NODE_ENV === 'production';

const api = axios.create({
  baseURL: `${production ? `//${window.location.host}` : process.env.VUE_APP_API}/api/v1`,
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
