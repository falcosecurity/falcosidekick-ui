<template>
  <v-card style="height: 100%">
    <v-card-title>
      Priority Timeline
    </v-card-title>
    <v-card-text>
      <wait>
        <apexchart type="bar" :height="height" :options="chart" :series="series"></apexchart>
      </wait>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { HistoryCategories, Serie } from '@/types'
import { colorCodes, initStats } from '@/api/mapper'
import { FalcoEvent, Priority } from '@/api/model'
import { calculateStart, timeline } from '@/helper/time'
import Vue, { PropType } from 'vue'
import Wait from './Wait.vue'

type Props = {
  time: number|null;
  events: FalcoEvent[];
}

type Computed = {
  chart: any;
  series: Serie[];
  categories: HistoryCategories;
}

type Data = {
  start: number;
  stepSize: number;
  steps: number;
  height: number;
  range: number;
}

export default Vue.extend<Data, {}, Computed, Props>({
  components: { Wait },
  name: 'PriorityTimeChart',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, required: true },
    time: { type: Number }
  },
  data () {
    const { range, steps, start, stepSize } = timeline(this.time, 15)

    return {
      range,
      steps,
      start,
      stepSize,
      height: 400
    }
  },
  watch: {
    events () {
      this.start = calculateStart(this.range, this.stepSize)
    }
  },
  computed: {
    categories () {
      return Array.from(Array(this.steps).keys()).reduce<HistoryCategories>((acc, step) => {
        const time = this.start + ((step + 1) * this.stepSize)
        let label = (new Date(time)).toISOString().replace('T', ' ').slice(0, -5)

        if (this.range >= 60 * 60000) {
          label = `${label.slice(0, -2)}00`
        }

        acc.push({ time, label })

        return acc
      }, [])
    },
    series () {
      const init = Object.keys(initStats()).reduce<{ [prioriry: string ]: Serie }>((acc, priority) => {
        acc[priority] = { name: priority, data: [...Array(this.steps).keys()].map(() => 0) }

        return acc
      }, {})

      const result: { [rule: string ]: Serie } = this.events.reduce<{ [priority: string ]: Serie }>((acc, event) => {
        if (event.filterTime < this.start) {
          return acc
        }

        for (const index in this.categories) {
          if (this.categories[(+index) + 1] && event.filterTime < this.categories[(+index) + 1].time) {
            acc[event.priority].data[index] += 1
            return acc
          }
        }

        acc[event.priority].data[this.steps - 1] += 1

        return acc
      }, init)

      return Object.values(result).reverse()
    },
    chart () {
      return {
        chart: {
          type: 'bar',
          height: this.height,
          stacked: true,
          toolbar: {
            show: false
          }
        },
        colors: this.series.map(serie => colorCodes[serie.name as Priority]),
        plotOptions: {
          bar: {
            horizontal: false
          }
        },
        xaxis: {
          categories: this.categories.map(({ label }) => label.split(' '))
        },
        legend: {
          position: 'right',
          offsetY: 40
        },
        fill: {
          opacity: 1
        }
      }
    }
  }
})
</script>
