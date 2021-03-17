<template>
    <v-autocomplete dense
                    multiple
                    :items="options"
                    outlined
                    hide-details
                    clearable
                    name="Priorities"
                    label="Priorities"
                    :value="value"
                    v-bind="$attrs"
                    @input="input"
    >
      <template v-slot:selection="{ item, index }">
          <v-chip small v-if="index <= 1" label outlined>
            <span>{{ item }}</span>
          </v-chip>
          <span
            v-if="index === 2"
            class="grey--text caption"
          >
            (+{{ selected.length - 2 }} others)
          </span>
        </template>
    </v-autocomplete>
</template>

<script lang="ts">
import Vue from 'vue'
import { initStats } from '../api/mapper'
import debounce from 'lodash.debounce'

const debounced = debounce((emit: () => void) => { emit() }, 800)

type Data = {
  options: string[];
  selected: string[];
}

type Methods = { input(priorities: string[]): void }

export default Vue.extend<Data, Methods, {}, { value: string[] }>({
  props: {
    value: { type: Array, default: () => [] }
  },
  data () {
    return {
      options: Object.keys(initStats()),
      selected: this.value
    }
  },
  methods: {
    input (priorities: string[]): void {
      this.selected = priorities

      debounced(() => { this.$emit('input', priorities) })
    }
  }
})
</script>
