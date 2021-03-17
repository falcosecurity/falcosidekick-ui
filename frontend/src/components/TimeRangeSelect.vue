<template>
    <v-select outlined
              hide-details
              dense
              label="Timerange"
              item-value="value"
              item-text="label"
              :items="options"
              clearable
              v-bind="$attrs"
              v-on="$listeners"
    />
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  inheritAttrs: false,
  data () {
    return {
      options: [
        {
          label: 'last 15 minutes',
          value: 15 * 60000
        },
        {
          label: 'last 30 minutes',
          value: 30 * 60000
        },
        {
          label: 'last hour',
          value: 60 * 60000
        },
        {
          label: 'last 6 hours',
          value: 6 * 60 * 60000
        },
        {
          label: 'last 12 hours',
          value: 12 * 60 * 60000
        },
        {
          label: 'last 24 hours',
          value: 24 * 60 * 60000
        },
        {
          label: 'last 2 days',
          value: 2 * 24 * 60 * 60000
        },
        {
          label: 'last 7 days',
          value: 7 * 24 * 60 * 60000
        }
      ]
    }
  }
})
</script>
