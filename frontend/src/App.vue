<template>
  <v-app id="inspire">
    <v-navigation-drawer app width="240">
      <v-sheet
        color="light-blue lighten-3"
        class="pa-4"
      >
        <v-img src="https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png" contain width="128" />
        <v-divider class="my-3"></v-divider>
        <h1 class="text-h5">Falcosidekick UI</h1>
        <h2 class="text-subtitle-1">Latest Events from Falco</h2>

        <v-chip class="my-2 mr-2 pr-0 py-0" color="primary" label style="border: 1px solid #000">
          Retention <div class="white black--text rounded ml-2 px-2 py-2 text-body-2">{{ retention }}</div>
        </v-chip>

        <v-divider class="my-3" />

        <h2 class="text-subtitle-1">Outputs</h2>

        <v-chip class="my-2 mr-2" color="primary" label v-for="output in outputs" :key="output">
          {{ output }}
        </v-chip>
      </v-sheet>

      <v-divider />

      <v-list class="py-0">
        <v-list-item to="/" color="light-blue">
          <v-list-item-content>
            <v-list-item-title>Dashboard</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
        <v-divider />
        <v-list-item to="/events" color="light-blue">
          <v-list-item-content>
            <v-list-item-title>Events</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>

      <v-divider />
    </v-navigation-drawer>

    <v-main class="grey lighten-3">
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { FETCH_EVENTS } from './store'

type Computed = {
  retention: number;
  outputs: string[];
}

export default Vue.extend<{}, {}, Computed, {}>({
  name: 'App',
  computed: mapState(['retention', 'outputs']),
  created () {
    this.$store.dispatch(FETCH_EVENTS)
  }
})
</script>
