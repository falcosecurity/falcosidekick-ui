<template>
  <v-container fluid class="px-3 py-2">
    <v-row class="pb-0">
      <v-col cols="12">
        <v-card elevation="2">
          <v-card-text>
              <priority-counter-chip dark
                      v-for="(count, priority) in stats"
                      :key="priority"
                      :priority="priority"
                      :count="count"
                      @click="search = priority"
                >
                {{ priority }} <div class="white black--text rounded ml-2 px-2 py-2 text-body-2">{{ count }}</div>
              </priority-counter-chip>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="pt-0 mt-0">
      <v-col cols="12">
        <v-card elevation="2">
          <v-toolbar flat>
            <v-spacer />
            <time-range-select class="ml-3" style="max-width: 300px;" v-model="time" />
            <v-text-field outlined hide-details dense class="ml-3" style="max-width: 300px;" label="Search" v-model="search" clearable />
          </v-toolbar>
          <v-divider />
          <event-table :events="filtered" :search.sync="search" />
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { initStats, mapPriorityToColor } from '@/api/mapper'
import { FalcoEvent, Priority, Stats } from '@/api/model'
import EventTable from '@/components/EventTable.vue'
import PriorityCounterChip from '@/components/PriorityCounterChip.vue'
import TimeRangeSelect from '@/components/TimeRangeSelect.vue'
import Vue from 'vue'
import { mapState } from 'vuex'

type Data = {
  search: string|null;
  time: number|null;
}

type Methods = {
  color (priority: string): string;
}

type Computed = {
  stats: Stats & { total: number };
  events: FalcoEvent[];
  filtered: FalcoEvent[];
}

export default Vue.extend<Data, Methods, Computed, {}>({
  components: { EventTable, PriorityCounterChip, TimeRangeSelect },
  name: 'Dashboard',
  data: () => ({
    search: '',
    time: null
  }),
  computed: {
    ...mapState(['events']),
    filtered () {
      const time = this.time

      if (typeof time === 'number' && time > 0) {
        const timeFilter = (new Date()).getTime() - time

        return this.events.filter(event => event.filterTime >= timeFilter)
      }

      return this.events
    },
    stats () {
      const stats: Stats & { total: number } = { total: 0, ...initStats() }

      this.events.forEach(event => {
        stats[event.priority] += 1
        stats.total += 1
      })

      return stats
    }
  },
  methods: {
    color (priority: string): string {
      return mapPriorityToColor(priority as Priority)
    }
  }
})
</script>
