<template>
  <v-card style="height: 100%">
    <v-card-title>
      Rule Timeline
    </v-card-title>
    <v-card-text>
      <wait>
        <apexchart type="bar" :height="height" :options="chart" :series="series"></apexchart>
      </wait>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import Vue, { PropType } from 'vue'
import Wait from './Wait.vue'
import { timeline, calculateStart } from '@/helper/time'
import { HistoryCategories, Serie } from '@/types'

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
  name: 'RuleTimeChart',
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
      return Array.from(Array(this.steps).keys()).reduce<Array<{ time: number; label: string }>>((acc, step) => {
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
      const unsorted: { [rule: string ]: Serie } = this.events.reduce<{ [rule: string ]: Serie }>((acc, event) => {
        if (Object.prototype.hasOwnProperty.call(acc, event.rule) === false) {
          acc[event.rule] = { name: event.rule, data: [...Array(this.steps).keys()].map(() => 0) }
        }

        if (event.filterTime < this.start) {
          return acc
        }

        for (const index in this.categories) {
          if (this.categories[(+index) + 1] && event.filterTime < this.categories[(+index) + 1].time) {
            acc[event.rule].data[index] += 1
            return acc
          }
        }

        acc[event.rule].data[this.steps - 1] += 1

        return acc
      }, {})

      const sorted = Object.keys(unsorted).sort().reduce<{ [rule: string]: Serie }>((acc, rule) => {
        acc[rule] = unsorted[rule]

        return acc
      }, {})

      return Object
        .values(sorted)
        .filter(serie => {
          const sum = serie.data.reduce<number>((sum, count) => sum + count, 0)

          return sum > 0
        })
    },
    chart () {
      return {
        chart: {
          type: 'bar',
          height: this.height,
          stacked: true,
          toolbar: {
            show: false
          },
          zoom: {
            enabled: false
          }
        },
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
