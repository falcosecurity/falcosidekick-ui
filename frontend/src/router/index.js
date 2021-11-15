import Vue from 'vue';
import VueRouter from 'vue-router';
import TestPage from '../views/TestPage.vue';
import Dashboard from '../views//DashboardPage.vue';
import EventsPage from '../views/EventsPage.vue';
import InfoPage from '../views/InfoPage.vue';

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
    path: '*',
    redirect: {
      name: 'dashboard',
    },
  },
];

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes,
});

export default router;
