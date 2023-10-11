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
import Vuetify from 'vuetify';
import moment from 'moment';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';
import App from './App.vue';
import router from './router';
import store from './store';

const capitalize = function capitalize(value) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
};

const formatDate = function formatDate(value) {
  if (!value) return '';
  return moment(String(value)).format('YYYY/MM/DD HH:mm:ss:SSS');
};

const opts = {};
const vuetify = new Vuetify(opts);

Vue.filter('formatDate', formatDate);
Vue.filter('capitalize', capitalize);

Vue.use(Vuetify);

new Vue({
  vuetify,
  router,
  store,
  components: { App },
  render: h => h(App),
}).$mount('#app');
