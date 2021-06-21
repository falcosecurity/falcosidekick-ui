<template>
  <v-card style="height: 100%">
    <v-card-title>
      Priorities
    </v-card-title>
    <v-card-text>
      <wait>
        <apexchart type="donut" height="350" :options="pie.chartOptions" :series="pie.series"></apexchart>
      </wait>
    </v-card-text>

    <event-overlay v-model="selectedEvents" />
  </v-card>
</template>

<script lang="ts">
import { colorCodes, initStats } from '@/api/mapper'
import { FalcoEvent, Priority, Stats } from '@/api/model'
import { DisplayMode } from '@/types'
import { ApexOptions } from 'apexcharts'
import Vue, { PropType } from 'vue'
import { mapState } from 'vuex'
import EventOverlay from './EventOverlay.vue'
import Wait from './Wait.vue'

type Props = {
  time: number|null;
  events: FalcoEvent[];
}

type Computed = {
  pie: any;
  displayMode: DisplayMode;
  theme: any;
}

type Data = {
  selectedEvents: FalcoEvent[];
  stats: Stats;
  eventsPerPriority: { [key: string]: FalcoEvent[] };
}

export default Vue.extend<Data, {}, Computed, Props>({
  components: { Wait, EventOverlay },
  name: 'PieChart',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, required: true },
    time: { type: Number }
  },
  data: () => ({
    selectedEvents: [],
    stats: initStats(),
    eventsPerPriority: {}
  }),
  watch: {
    events: {
      immediate: true,
      handler () {
        const stats = initStats()
        const eventsPerPriority: { [key: string]: FalcoEvent[] } = {}

        this.events.forEach((event) => {
          stats[event.priority] += 1

          if (Object.prototype.hasOwnProperty.call(eventsPerPriority, event.priority) === false) {
            eventsPerPriority[event.priority] = []
          }

          eventsPerPriority[event.priority].push(event)
        })

        this.stats = stats
        this.eventsPerPriority = eventsPerPriority
      }
    }
  },
  computed: {
    ...mapState(['displayMode']),
    theme () {
      if (this.displayMode === DisplayMode.LIGHT) return { stroke: { colors: ['#FFFFFF'] } }

      return {
        stroke: { colors: ['#1E1E1E'] },
        theme: { mode: 'dark' }
      }
    },
    pie (): { series: number[]; chartOptions: ApexOptions} {
      const labels = Object.keys(this.stats)

      return {
        series: Object.values(this.stats),
        chartOptions: {
          ...this.theme,
          chart: {
            height: 350,
            type: 'donut',
            selection: {
              enabled: false
            },
            events: {
              dataPointSelection: (_: any, __: any, config: { dataPointIndex: number }) => {
                this.selectedEvents = [...this.eventsPerPriority[labels[config.dataPointIndex]]]
              }
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
          dataLabels: {
            dropShadow: {
              enabled: true,
              top: 0,
              left: 0,
              blur: 1,
              color: '#000',
              opacity: 1
            }
          },
          plotOptions: {
            pie: {
              expandOnClick: false,
              donut: {
                labels: {
                  show: true,
                  total: {
                    showAlways: true,
                    show: true
                  }
                }
              }
            }
          },
          labels,
          colors: Object.keys(this.stats).map(priority => colorCodes[priority as Priority])
        }
      }
    }
  }
})
</script>
