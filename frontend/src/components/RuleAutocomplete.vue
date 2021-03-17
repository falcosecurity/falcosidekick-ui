<template>
    <v-autocomplete dense
                    multiple
                    :items="options"
                    outlined
                    hide-details
                    name="Rules"
                    label="Rules"
                    clearable
                    :value="value"
                    v-bind="$attrs"
                    @input="input"
    >
    <template v-slot:selection="{ item, index }">
        <v-chip small v-if="index === 0" label outlined>
          <span>{{ item }}</span>
        </v-chip>
        <span
          v-if="index === 1"
          class="grey--text caption"
        >
          (+{{ value.length - 1 }} others)
        </span>
      </template>
    </v-autocomplete>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import Vue from 'vue'
import debounce from 'lodash.debounce'

const debounced = debounce((emit: () => void) => { emit() }, 800)

type Data = {
  selected: string[];
}

type Computed = {
  options: string[];
}

type Props = {
  events: FalcoEvent[];
  value: string[];
}

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, Computed, Props>({
  props: {
    events: { type: Array, default: () => [] },
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      selected: []
    }
  },
  computed: {
    options () {
      return Object.keys(this.events.reduce<{ [rule: string]: null }>((acc, event) => {
        acc[event.rule] = null

        return acc
      }, {})).sort()
    }
  },
  methods: {
    input (rules: string[]): void {
      this.selected = rules

      debounced(() => { this.$emit('input', rules) })
    }
  }
})
</script>
