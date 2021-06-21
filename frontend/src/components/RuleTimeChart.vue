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

    <event-overlay v-model="selectedEvents" />
  </v-card>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import Vue, { PropType } from 'vue'
import Wait from './Wait.vue'
import { timeline, calculateStart } from '@/helper/time'
import { DisplayMode, HistoryCategories, Serie } from '@/types'
import { ApexOptions } from 'apexcharts'
import EventOverlay from './EventOverlay.vue'
import { mapState } from 'vuex'

type Props = {
  time: number|null;
  events: FalcoEvent[];
}

type Computed = {
  chart: any;
  displayMode: DisplayMode;
  theme: any;
}

type Data = {
  start: number;
  stepSize: number;
  steps: number;
  height: number;
  range: number;
  series: Serie[];
  categories: HistoryCategories;
  selectedEvents: FalcoEvent[];
  eventsPerBucket: { [priority: string]: { [atep: number]: FalcoEvent[] } };
}

export default Vue.extend<Data, {}, Computed, Props>({
  components: { Wait, EventOverlay },
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
      height: 400,
      series: [],
      selectedEvents: [],
      eventsPerBucket: {},
      categories: []
    }
  },
  watch: {
    events: {
      immediate: true,
      handler (events: FalcoEvent[]) {
        const start = calculateStart(this.range, this.stepSize)

        const categories = Array.from(Array(this.steps).keys()).reduce<HistoryCategories>((acc, step) => {
          const time = start + ((step + 1) * this.stepSize)
          let label = (new Date(time)).toISOString().replace('T', ' ').slice(0, -5)

          if (this.range >= 60 * 60000) {
            label = `${label.slice(0, -2)}00`
          }

          acc.push({ time, label })

          return acc
        }, [])

        const eventsPerBucket: { [priority: string]: FalcoEvent[][] } = {}
        const series: { [priority: string ]: Serie } = {}

        events.forEach((event) => {
          if (Object.prototype.hasOwnProperty.call(series, event.rule) === false) {
            series[event.rule] = { name: event.rule, data: [...Array(this.steps).keys()].map(() => 0) }
            eventsPerBucket[event.rule] = [...Array(this.steps).keys()].map(() => [])
          }

          if (event.filterTime < start) {
            return
          }

          for (const index in categories) {
            if (categories[(+index) + 1] && event.filterTime < categories[(+index) + 1].time) {
              series[event.rule].data[index] += 1
              eventsPerBucket[event.rule][index].push(event)
              return
            }
          }

          series[event.rule].data[this.steps - 1] += 1
          eventsPerBucket[event.rule][this.steps - 1].push(event)
        })

        const sorted = Object.keys(series).sort().reduce<{ [rule: string]: Serie }>((acc, rule) => {
          acc[rule] = series[rule]

          return acc
        }, {})

        this.eventsPerBucket = eventsPerBucket
        this.categories = categories
        this.start = start
        this.series = Object
          .values(sorted)
          .filter(serie => {
            const sum = serie.data.reduce<number>((sum, count) => sum + count, 0)

            return sum > 0
          })
      }
    }
  },
  computed: {
    ...mapState(['displayMode']),
    theme () {
      if (this.displayMode === DisplayMode.LIGHT) return {}

      return {
        theme: { mode: 'dark', palette: 'palette1' }
      }
    },
    chart (): ApexOptions {
      const labels = this.series.map(serie => serie.name)

      return {
        ...this.theme,
        chart: {
          type: 'bar',
          height: this.height,
          stacked: true,
          toolbar: {
            show: false
          },
          zoom: {
            enabled: false
          },
          events: {
            dataPointSelection: (_: any, __: any, config: { dataPointIndex: number; seriesIndex: number }) => {
              const rule = labels[config.seriesIndex]

              this.selectedEvents = [...this.eventsPerBucket[rule][config.dataPointIndex]]
            }
          }
        },
        plotOptions: {
          bar: {
            horizontal: false
          }
        },
        states: {
          active: {
            allowMultipleDataPointsSelection: false,
            filter: {
              type: 'none'
            }
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
