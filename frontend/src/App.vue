<template>
  <v-app app>
    <v-app-bar app class="blue darken-2" dark>
      <v-img
        src="https://github.com/falcosecurity/falcosidekick/raw/master/imgs/falcosidekick_color.png"
        max-height="40"
        max-width="40"
        style="margin-bottom: 5px; margin-right: 15px;">
      </v-img>
      <v-toolbar-title class="text-no-wrap">
        Falcosidekick UI
      </v-toolbar-title>
      <template v-slot:extension>
        <v-tabs>
          <v-tab v-for="(page) in pages" :key="page.title"
          class="blue darken-2" :to="{ path: page.href, query: $route.query }">
            {{page.title}}
          </v-tab>
        </v-tabs>
        refresh
        <v-select
          style="margin-top: 15px; margin-left: 15px; max-width: 80px;"
          v-model="refresh"
          :items="refreshIntervals"
          dense
        ></v-select>
      </template>
      <v-spacer/>
      <Counters v-if="$store.state.username && $store.state.password"></Counters>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-footer app absolute
      class="blue darken-2 text-no-wrap" dark
    >
      <span>
        2022 - <a href="https://github.com/falcosecurity/falcosidekick-ui">Falco Authors</a>
      </span>
      <v-spacer/>
      <span v-if="$store.state.username && $store.state.password">
        logged as <b>{{$store.state.username}}</b>
        <v-btn text x-small class="ml-3" @click="logout">
          Logout
        </v-btn>
      </span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapActions } from 'vuex';
import Counters from './components/counters.vue';

export default {
  name: 'App',
  components: {
    Counters,
  },
  data() {
    return {
      pages: [
        {
          href: '/dashboard',
          router: true,
          title: 'Dashboard',
        }, {
          href: '/events',
          router: true,
          title: 'Events',
        },
        {
          href: '/info',
          router: true,
          title: 'Info',
        },
        // {
        //   href: 'test',
        //   router: true,
        //   title: 'Test',
        // },
      ],
      timer: '',
      refresh: '10s',
    };
  },
  computed: {
    refreshInterval() {
      return this.$store.state.refreshInterval;
    },
    refreshIntervals() {
      return this.$store.state.refreshIntervals;
    },
  },
  watch: {
    refresh: {
      handler() {
        this.setRefreshInterval(this.refresh);
        this.setTimer();
        if (this.$route.query.refresh !== this.refresh || this.$route.query.refresh === '') {
          this.$router.push({ query: { ...this.$route.query, refresh: this.refresh } });
        }
      },
    },
  },
  methods: {
    ...mapActions([
      'increment',
      'setRefreshInterval',
    ]),
    logout() {
      this.$store.state.username = '';
      this.$store.state.password = '';
      this.$router.push('/login');
    },
    cancelAutoUpdate() {
      clearInterval(this.timer);
    },
    setTimer() {
      clearInterval(this.timer);
      this.timer = setInterval(() => {
        if (this.refreshInterval !== 0) {
          this.increment();
        }
      }, this.$store.state.refreshInterval);
    },
  },
  created() {
    if (typeof this.$route.query.refresh !== 'undefined') {
      this.refresh = this.$route.query.refresh;
      this.setRefreshInterval(this.refresh);
    }
    this.setTimer();
  },
  beforeDestroy() {
    this.cancelAutoUpdate();
  },
};

</script>

<style scoped>
a:link {
  text-decoration: none;
  color: white;
}
a:visited {
  text-decoration: none;
  color: white;
}
a:hover {
  text-decoration: underline;
}
</style>
