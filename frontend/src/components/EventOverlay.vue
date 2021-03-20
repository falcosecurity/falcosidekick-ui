<template>
  <v-dialog v-model="dialog">
    <v-card>
      <v-card-title>
        Events
        <v-spacer />
        <v-text-field outlined hide-details dense class="ml-3" style="max-width: 300px;" label="Search" v-model="search" clearable />
        <v-btn icon class="ml-4" @click="close">
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </v-card-title>

      <v-divider />

      <event-table :events="value" v-if="value" :search.sync="search" />

      <v-divider />

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn
          color="primary"
          text
          @click="close"
        >
          Close
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import Vue from 'vue'
import EventTable from './EventTable.vue'

type Data = { dialog: boolean; search: string }

type Props = { value: FalcoEvent[] }

type Methods = { close (): void }

export default Vue.extend<Data, Methods, {}, Props>({
  components: { EventTable },
  props: {
    value: { type: Array, default: () => [] }
  },
  data: () => ({
    search: '',
    dialog: false
  }),
  watch: {
    value (value: FalcoEvent[]) {
      if (value.length === 0) {
        return
      }

      this.search = ''

      setTimeout(() => { this.dialog = true }, 200)
    }
  },
  methods: {
    close (): void {
      this.dialog = false

      setTimeout(() => { this.$emit('input', []) }, 200)
    }
  }
})
</script>
