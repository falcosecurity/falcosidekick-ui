<template>
    <v-data-table :items="events"
                  :headers="headers"
                  :items-per-page="-1"
                  hide-default-footer
                  :search="search"
                  :custom-filter="defaultFilter"
                  disable-sort
    >
    <template #item="{ item }">
        <tr>
        <td><span v-html="item.time.toISOString().split('T').join(' ').slice(0, -5)" /></td>
        <td>
            <v-chip label :color="color(item.priority)" class="mb-3 white--text" @click="$emit('update:search', item.priority)">
              <span v-html="highlightMatches(item.priority)" />
            </v-chip><br>
            <v-chip label
                    @click="$emit('update:search', item.rule)"
                    class="black--text"
                    color="light-blue lighten-5"
                    v-html="highlightMatches(item.rule)"
            /><br>
        </td>
        <td class="py-3" style="font-size: 0.8rem">
            <div v-html="highlightMatches(item.output)" />
            <div class="mt-3">
            <template v-for="(value, label) in item.outputFields">
                <output-field-chip dark
                                   v-if="value.length <= 100"
                                   :key="label"
                                   :value="highlightMatches(value)"
                                   :label="label"
                                   @click="$emit('update:search', value)"
                />
                <output-field-card v-else
                                  :key="label"
                                  :value="highlightMatches(value)"
                                  :label="label"
                                  @click="$emit('update:search', value)"
                />
            </template>
            </div>
        </td>
        </tr>
    </template>
    </v-data-table>
</template>

<script lang="ts">
import { mapPriorityToColor } from '@/api/mapper'
import { FalcoEvent, Priority } from '@/api/model'
import Vue, { PropType } from 'vue'
import { DataTableHeader } from 'vuetify'
import OutputFieldCard from './OutputFieldCard.vue'
import OutputFieldChip from './OutputFieldChip.vue'

export default Vue.extend({
  components: { OutputFieldChip, OutputFieldCard },
  name: 'EventTable',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, default: () => [] },
    search: { type: String }
  },
  data: (): { headers: DataTableHeader[] } => ({
    headers: [
      { text: 'Time', value: 'time' },
      { text: 'Rule / Priority', value: 'rule' },
      { text: 'Output', value: 'output', width: '60%' }
    ]
  }),
  methods: {
    color (priority: string): string {
      return mapPriorityToColor(priority as Priority)
    },
    highlightMatches (text: string): string {
      if (!this.search) return text

      const matchExists = text
        .toLowerCase()
        .includes(this.search.toLowerCase())

      if (!matchExists) return text

      return text.replace(new RegExp(this.search, 'ig'), matchedText => `<strong>${matchedText}</strong>`)
    },
    defaultFilter (value: any, search: string, item: FalcoEvent) {
      search = (search || '').toLocaleLowerCase()

      const existInFields = Object
        .values(item.outputFields)
        .some((value) => {
          return value
            .toLowerCase()
            .includes(search)
        })

      if (existInFields) return true

      const checkPriority = item.priority != null &&
        search != null &&
        item.priority.toLocaleLowerCase().indexOf(search) !== -1

      if (checkPriority) return true

      return value != null &&
        search != null &&
        typeof value !== 'boolean' &&
        value.toString().toLocaleLowerCase().indexOf(search) !== -1
    }
  }
})
</script>

<style scoped>
  >>> strong {
    display: contents;
  }
</style>
