<template>
  <v-card flat class="primary mt-2" style="border: 1px solid;" v-bind="$attrs" v-on="$listeners">
    <v-system-bar color="primary white--text">{{ label }}</v-system-bar>
    <v-card-text class="pa-2 white black--text" style="font-size: 0.8rem; line-height: 1.1rem;">
        <span v-html="value" />
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  name: 'OutputFieldCard',
  inheritAttrs: false,
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true }
  }
})
</script>
