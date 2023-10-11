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

import Vue from 'vue';
import VueRouter from 'vue-router';
import TestPage from '../views/TestPage.vue';
import Dashboard from '../views//DashboardPage.vue';
import EventsPage from '../views/EventsPage.vue';
import InfoPage from '../views/InfoPage.vue';
import LoginPage from '../views/LoginPage.vue';
import store from '../store';

Vue.use(VueRouter);

const routes = [
  {
    path: '/test',
    name: 'test',
    component: TestPage,
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: Dashboard,
  },
  {
    path: '/events',
    name: 'events',
    component: EventsPage,
  },
  {
    path: '/info',
    name: 'info',
    component: InfoPage,
  },
  {
    path: '/login',
    name: 'login',
    component: LoginPage,
  },
  {
    path: '*',
    redirect: {
      name: 'dashboard',
    },
  },
];

const router = new VueRouter({
  mode: 'hash',
  base: process.env.BASE_URL,
  routes,
});

router.beforeEach((to, from, next) => {
  if (store.state.username === '' || store.state.password === '') {
    if (to.name !== 'login') {
      router.push('/login');
    }
  }
  next();
});

export default router;
