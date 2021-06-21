<template>
  <v-card flat :class="`bg-${displayMode}-color mt-2`" style="border: 1px solid;" v-bind="$attrs" v-on="$listeners">
    <v-system-bar :color="`bg-${displayMode}-color white--text`">{{ label }}</v-system-bar>
    <v-card-text class="pa-2 text" style="font-size: 0.8rem; line-height: 1.1rem;">
        <span v-html="value" />
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  name: 'OutputFieldCard',
  inheritAttrs: false,
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true }
  },
  computed: mapState(['displayMode'])
})
</script>

<style scoped>
.bg-light-color {
  background-color: var(--v-primary-base)!important;
  border-color: var(--v-primary-base)!important;
}
.bg-dark-color {
  background-color: var(--v-primary-darken3)!important;
  border-color: var(--v-primary-darken3)!important;
}

.bg-light-color .text {
  background-color: var(--v-white)!important;
  color: var(--v-black)!important;
}
.bg-dark-color .text {
  background-color: var(--v-secondary-base)!important;
  color: var(--v-white)!important;
}
</style>
