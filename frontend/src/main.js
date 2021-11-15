import Vue from 'vue';
import Vuetify from 'vuetify';
import moment from 'moment';
import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';
import router from './router';
import App from './App.vue';

const capitalize = function capitalize(value) {
  if (!value) return '';
  return value.charAt(0).toUpperCase() + value.slice(1);
};

const formatDate = function formatDate(value) {
  if (!value) return '';
  return moment(String(value)).format('YYYY/MM/DD hh:mm:ss:SSS');
};

const opts = {};
const vuetify = new Vuetify(opts);

Vue.filter('formatDate', formatDate);
Vue.filter('capitalize', capitalize);

Vue.use(Vuetify);

new Vue({
  vuetify,
  router,
  components: { App },
  render: h => h(App),
}).$mount('#app');
