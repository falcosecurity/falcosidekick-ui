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

    <event-overlay v-model="selectedEvents" />
  </v-card>
</template>

<script lang="ts">
import { DisplayMode, HistoryCategories, Serie } from '@/types'
import { colorCodes, initStats } from '@/api/mapper'
import { FalcoEvent, Priority } from '@/api/model'
import { calculateStart, timeline } from '@/helper/time'
import Vue, { PropType } from 'vue'
import Wait from './Wait.vue'
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

        Object.keys(initStats()).forEach((priority) => {
          series[priority] = { name: priority, data: [...Array(this.steps).keys()].map(() => 0) }
          eventsPerBucket[priority] = [...Array(this.steps).keys()].map(() => [])
        })

        events.forEach((event) => {
          if (event.filterTime < this.start) {
            return
          }

          for (const index in categories) {
            if (categories[(+index) + 1] && event.filterTime < categories[(+index) + 1].time) {
              series[event.priority].data[index] += 1
              eventsPerBucket[event.priority][index].push(event)
              return
            }
          }

          series[event.priority].data[this.steps - 1] += 1
          eventsPerBucket[event.priority][this.steps - 1].push(event)
        })

        this.series = Object.values(series).reverse()
        this.eventsPerBucket = eventsPerBucket
        this.categories = categories
        this.start = start
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
          events: {
            dataPointSelection: (_: any, __: any, config: { dataPointIndex: number; seriesIndex: number }) => {
              const priority = labels[config.seriesIndex]

              this.selectedEvents = [...this.eventsPerBucket[priority][config.dataPointIndex]]
            }
          }
        },
        colors: this.series.map(serie => colorCodes[serie.name as Priority]),
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
