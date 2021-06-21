<template>
  <v-chip label :class="`bg-${displayMode}-color`" v-html="value" v-on="$listeners" />
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  name: 'OutputFieldChip',
  props: {
    value: { type: String, required: true }
  },
  computed: mapState(['displayMode'])
})
</script>

<style scoped>
.bg-light-color {
  background-color: var(--v-info-lighten3)!important;
  color: var(--v-black)!important;
}
.bg-dark-color {
  background-color: var(--v-info-darken3)!important;
  color: var(--v-white)!important;
}
</style>
