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
      </template>
      <v-spacer/>
      <Counters></Counters>
    </v-app-bar>
    <v-main>
      <router-view></router-view>
    </v-main>
    <v-footer app absolute
    class="blue darken-2" dark>
      <span>
        2022 - <a href="https://github.com/falcosecurity/falcosidekick-ui">Falco Authors</a>
        </span>
    </v-footer>
  </v-app>
</template>

<script>
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
    };
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
