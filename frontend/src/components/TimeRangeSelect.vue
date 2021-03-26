<template>
  <v-select
    outlined
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

interface BuilderConfig {
  seconds: number[];
  minutes: number[];
  hours: number[];
  days: number[];
}

type BuilderConfigKeys = 'minutes' | 'hours' | 'days'

interface Option {
  label: string;
  value: number;
}

export default Vue.extend({
  inheritAttrs: false,
  data () {
    return {
      builderConfig: {
        minutes: [15, 30],
        hours: [1, 6, 12, 24],
        days: [2, 7]
      },
      options: [] as Option[]
    }
  },
  created () {
    this.options = this.buildTimeFilters()
  },
  methods: {
    buildTimeFilters (): Option[] {
      const additional = {
        minutes: 1,
        hours: 60,
        days: 24 * 60
      }

      const config = this.builderConfig as BuilderConfig

      return (Object.keys(config) as BuilderConfigKeys[]).reduce(
        (accumulator: Option[], current) => {
          const timeFilters: Option[] = config[current].map((time: number) => {
            let label = `last ${time === 1 ? '' : `${time} `}${current}`
            if (time === 1) label = label.substr(0, label.length - 1)

            const value = 60000 * time * additional[current]

            return {
              label,
              value
            }
          })

          accumulator = [...accumulator, ...timeFilters]

          return accumulator
        },
        []
      )
    }
  }
})
</script>
