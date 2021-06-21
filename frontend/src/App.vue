<template>
  <v-app id="inspire">
    <v-navigation-drawer app width="240" style="padding-bottom: 60px;">
      <v-sheet :class="`bg-${displayMode}-color`" class="pa-4">
        <v-img src="https://raw.githubusercontent.com/falcosecurity/falcosidekick/master/imgs/falcosidekick_color.png" contain width="128" />
        <v-divider class="my-3"></v-divider>
        <h1 class="text-h5">Falcosidekick UI</h1>
        <h2 class="text-subtitle-1">Latest Events from Falco</h2>

        <v-chip class="chip my-2 mr-2 pr-0 py-0" label style="border: 1px solid #000">
          Retention <div class="white black--text rounded ml-2 px-2 py-2 text-body-2">{{ retention }}</div>
        </v-chip>

        <v-divider class="my-3" />

        <h2 class="text-subtitle-1">Outputs</h2>

        <v-chip class="chip my-2 mr-2" label v-for="output in outputs" :key="output">
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

      <div style="position: absolute; bottom: 10px;" class="px-2 text-center">
        <v-btn color="black" @click="dark" dark depressed>dark</v-btn>
        <v-btn color="grey lighten-3 black--text" @click="light" class="ml-2" depressed>light</v-btn>
      </div>
    </v-navigation-drawer>

    <v-main :class="`bg-main-${displayMode}`">
      <router-view />
    </v-main>
  </v-app>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'
import { FETCH_EVENTS, INIT_DISPLAY_MODE, SET_DISPLAY_MODE } from './store'
import { DisplayMode } from './types'

type Computed = {
  retention: number;
  outputs: string[];
  displayMode: DisplayMode;
}

export default Vue.extend<{}, { dark(): void; light(): void }, Computed, {}>({
  name: 'App',
  computed: mapState(['retention', 'outputs', 'displayMode']),
  methods: {
    light () {
      this.$store.commit(SET_DISPLAY_MODE, DisplayMode.LIGHT)
    },
    dark () {
      this.$store.commit(SET_DISPLAY_MODE, DisplayMode.DARK)
    }
  },
  watch: {
    displayMode: {
      immediate: true,
      handler (displayMode: DisplayMode) {
        this.$vuetify.theme.dark = displayMode === DisplayMode.DARK
      }
    }
  },
  created () {
    this.$store.dispatch(FETCH_EVENTS)

    if (!sessionStorage.getItem('displayMode') && window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      this.$store.commit(INIT_DISPLAY_MODE, DisplayMode.DARK)
    }

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
      if (!e.matches) {
        this.$store.commit(SET_DISPLAY_MODE, DisplayMode.LIGHT)
        return
      }

      this.$store.commit(SET_DISPLAY_MODE, DisplayMode.DARK)
    })
  }
})
</script>

<style>
:root {
  --v-white: #fff;
  --v-black: #000;
}
</style>

<style scoped>
.bg-main-light {
  background-color: rgb(240, 240, 240)!important;
}

.bg-light-color {
  background-color: var(--v-info-lighten2)!important;
  color: var(--v-black)!important;
}
.bg-dark-color {
  background-color: var(--v-info-darken1)!important;
  color: var(--v-white)!important;
}

.bg-light-color .chip {
  border-color: var(--v-primary-base)!important;
  background-color: var(--v-primary-base)!important;
  color: var(--v-white)!important;
}
.bg-dark-color .chip {
  border-color: var(--v-primary-darken3)!important;
  background-color: var(--v-primary-darken3)!important;
  color: var(--v-white)!important;
}
</style>
