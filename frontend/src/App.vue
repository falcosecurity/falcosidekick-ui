<template>
  <v-app>
    <template v-if="showChrome">
      <Sidebar :pages="pages" v-model="refresh" @logout="logout"></Sidebar>
      <main class="main">
        <router-view></router-view>
      </main>
    </template>
    <router-view v-else></router-view>
  </v-app>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import { Chart as ChartJS } from 'chart.js';
import Sidebar from './components/Sidebar.vue';
import { utils } from './utils';

const QUERY_KEYS = {
  priorities: 'priority',
  sources: 'source',
  hostnames: 'hostname',
  tags: 'tags',
  rule: 'rule',
  search: 'filter',
  since: 'since',
};

const asArray = (v) => {
  if (v === undefined || v === null || v === '') return [];
  return Array.isArray(v) ? v : [v];
};

export default {
  name: 'App',
  components: {
    Sidebar,
  },
  data() {
    return {
      pages: [
        { href: '/dashboard', title: 'Dashboard', icon: 'mdi-view-dashboard-outline' },
        { href: '/events', title: 'Events', icon: 'mdi-format-list-bulleted' },
      ],
      timer: '',
      refresh: '10s',
    };
  },
  computed: {
    ...mapState(['filters', 'theme', 'refreshInterval']),
    showChrome() {
      return this.$route.name !== 'login'
        && this.$store.state.username && this.$store.state.password;
    },
  },
  watch: {
    refresh: {
      handler() {
        this.setRefreshInterval(this.refresh);
        this.setTimer();
        this.syncQuery({ refresh: this.refresh });
      },
    },
    filters: {
      handler() {
        const q = {};
        Object.entries(QUERY_KEYS).forEach(([key, param]) => {
          const v = this.filters[key];
          q[param] = Array.isArray(v) ? v : (v || undefined);
        });
        this.syncQuery(q);
      },
      deep: true,
    },
    theme: {
      handler() { this.applyTheme(); },
      immediate: true,
    },
  },
  methods: {
    ...mapActions(['increment', 'setRefreshInterval', 'setFilters', 'emptyCredentials']),
    syncQuery(patch) {
      const query = { ...this.$route.query, ...patch };
      Object.keys(query).forEach((k) => {
        const v = query[k];
        if (v === undefined || v === '' || (Array.isArray(v) && v.length === 0)) delete query[k];
      });
      if (JSON.stringify(query) === JSON.stringify(this.$route.query)) return;
      this.$router.replace({ query }).catch(() => {});
    },
    applyTheme() {
      document.documentElement.setAttribute('data-theme', this.theme);
      this.$vuetify.theme.dark = this.theme === 'midnight';
      this.$nextTick(() => {
        ChartJS.defaults.color = utils.cssVar('--text-muted');
        ChartJS.defaults.borderColor = utils.cssVar('--grid');
        const accent = utils.cssVar('--accent');
        this.$vuetify.theme.themes.light.primary = accent;
        this.$vuetify.theme.themes.dark.primary = accent;
        this.increment();
      });
    },
    logout() {
      this.emptyCredentials();
      this.$router.push('/login');
    },
    setTimer() {
      clearInterval(this.timer);
      this.timer = setInterval(() => {
        if (this.refreshInterval !== 0) {
          this.increment();
        }
      }, this.refreshInterval || 60000);
    },
    loadFiltersFromQuery() {
      const q = this.$route.query;
      const f = {};
      Object.entries(QUERY_KEYS).forEach(([key, param]) => {
        if (q[param] === undefined) return;
        f[key] = Array.isArray(this.filters[key]) ? asArray(q[param]) : String(q[param]);
      });
      this.setFilters(f);
    },
  },
  created() {
    ChartJS.defaults.font.family = "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif";
    ChartJS.defaults.font.size = 11;
    ChartJS.defaults.plugins.legend.labels.boxWidth = 8;
    ChartJS.defaults.plugins.legend.labels.boxHeight = 8;
    ChartJS.defaults.plugins.legend.labels.usePointStyle = true;
    ChartJS.defaults.plugins.tooltip.backgroundColor = 'rgba(20, 27, 36, 0.92)';
    ChartJS.defaults.plugins.tooltip.padding = 8;
    ChartJS.defaults.plugins.tooltip.cornerRadius = 4;
    this.setTimer();
    // the initial navigation resolves asynchronously (route guards),
    // so read the URL only once the router is ready
    this.$router.onReady(() => {
      this.loadFiltersFromQuery();
      if (typeof this.$route.query.refresh !== 'undefined') {
        this.refresh = this.$route.query.refresh;
        this.setRefreshInterval(this.refresh);
      }
    });
  },
  beforeDestroy() {
    clearInterval(this.timer);
  },
};
</script>

<style scoped>
.main {
  margin-left: 248px;
  min-height: 100vh;
  background: var(--bg);
}
</style>
