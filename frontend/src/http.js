import axios from 'axios';

// const axiosParams = function axiosParams() {
//   const params = new URLSearchParams();
//   return params;
// };

const api = axios.create({
  baseURL: 'http://localhost:2802/api/v1',
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
    });
  },
  getConfiguration() {
    return api.request({
      url: '/configuration',
      method: 'get',
      params: {},
    });
  },
  getVersion() {
    return api.request({
      url: '/version',
      method: 'get',
      params: {},
    });
  },
  countEvents(source, priority, rule, filter, tags, since) {
    return api.request({
      url: '/events/count',
      method: 'get',
      params: {
        source: `${source}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
      },
    });
  },
  countByEvents(group, source, priority, rule, filter, tags, since) {
    return api.request({
      url: `/events/count/${group}`,
      method: 'get',
      params: {
        source: `${source}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
      },
    });
  },
  searchEvents(source, priority, rule, filter, tags, since, page, limit) {
    return api.request({
      url: '/events/search',
      method: 'get',
      params: {
        source: `${source}`,
        priority: `${priority}`,
        rule: `${rule}`,
        filter: `${filter}`,
        tags: `${tags}`,
        since: `${since}`,
        page: `${page}`,
        limit: `${limit}`,
      },
    });
  },
};

export default {
};
