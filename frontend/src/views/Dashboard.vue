<template>
  <v-container fluid class="px-3 pt-2">
    <v-row class="pb-0">
      <v-col cols="12">
        <v-card elevation="2">
          <v-toolbar flat>
            <v-spacer />
            <rule-autocomplete class="ml-3" style="max-width: 500px;" v-model="rules" :events="events" />
            <priority-autocomplete class="ml-3" style="max-width: 350px;" v-model="priorities" />
            <time-range-select class="ml-3" style="max-width: 300px;" v-model="time" />
          </v-toolbar>
        </v-card>
      </v-col>
    </v-row>
    <v-row class="pb-0 mt-0">
      <v-col cols="4">
        <priority-chart :events="filtered" />
      </v-col>
      <v-col cols="8">
        <rule-chart :events="filtered" />
      </v-col>
    </v-row>
    <v-row class="pb-0 mt-0">
      <v-col cols="12">
        <rule-time-chart :events="filtered" :time="time" :key="time" />
      </v-col>
    </v-row>
    <v-row class="pb-0 mt-0">
      <v-col cols="12">
        <priority-time-chart :events="filtered" :time="time" :key="time" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import PriorityAutocomplete from '@/components/PriorityAutocomplete.vue'
import PriorityChart from '@/components/PriorityChart.vue'
import PriorityTimeChart from '@/components/PriorityTimeChart.vue'
import RuleAutocomplete from '@/components/RuleAutocomplete.vue'
import RuleChart from '@/components/RuleChart.vue'
import RuleTimeChart from '@/components/RuleTimeChart.vue'
import TimeRangeSelect from '@/components/TimeRangeSelect.vue'
import Vue from 'vue'
import { mapState } from 'vuex'

type Data = {}

type Computed = {
  events: FalcoEvent[];
  filtered: FalcoEvent[];
  time: number|null;
  priorities: string[];
  rules: string[];
}

export default Vue.extend<Data, {}, Computed, {}>({
  components: { TimeRangeSelect, PriorityChart, RuleChart, RuleTimeChart, PriorityTimeChart, PriorityAutocomplete, RuleAutocomplete },
  name: 'Dashboard',
  computed: {
    ...mapState(['events']),
    time: {
      get () {
        const time = this.$route.query.time

        if (time) return +time

        return null
      },
      set (time: number | null): void {
        if (!time) {
          this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, time: undefined } })

          return
        }

        this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, time: time.toString() } })
      }
    },
    priorities: {
      get () {
        const priority: string|string[] = this.$route.query.priority as string|string[]

        if (!priority) return []

        return Array.isArray(priority) ? priority : [priority]
      },
      set (priority: string[]): void {
        if (!priority) {
          this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, priority: undefined } })

          return
        }

        this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, priority } })
      }
    },
    rules: {
      get () {
        const rules: string|string[] = this.$route.query.rules as string|string[]

        if (!rules) return []

        return Array.isArray(rules) ? rules : [rules]
      },
      set (rules: string[]): void {
        if (!rules) {
          this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, rules: undefined } })

          return
        }

        this.$router.push({ name: this.$route.name || '', query: { ...this.$route.query, rules } })
      }
    },
    filtered () {
      const time = this.time
      let timeFilter: null|number = null

      if (typeof time === 'number' && time > 0) {
        timeFilter = (new Date()).getTime() - time
      }

      if (this.rules.length === 0 && this.priorities.length === 0 && !timeFilter) {
        return this.events
      }

      return this.events.filter(event => {
        return (!timeFilter || event.filterTime >= timeFilter) &&
          (!this.priorities.length || this.priorities.includes(event.priority)) &&
          (!this.rules.length || this.rules.includes(event.rule))
      })
    }
  }
})
</script>

<style>
.apexcharts-svg {
  background-color: transparent!important;
}
</style>
