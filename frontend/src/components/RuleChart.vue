<template>
  <v-card style="height: 100%">
    <v-card-title>
      Rules
    </v-card-title>
    <v-card-text>
      <apexchart type="bar" height="350" :options="chart.chartOptions" :series="chart.series"></apexchart>
    </v-card-text>

    <event-overlay v-model="selectedEvents" />
  </v-card>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import { DisplayMode } from '@/types'
import Vue, { PropType } from 'vue'
import { mapState } from 'vuex'
import EventOverlay from './EventOverlay.vue'

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
  selectedEvents: FalcoEvent[];
  eventsPerRule: { [rule: string]: FalcoEvent[] };
  rules: { [rule: string]: number };
}

// convert a label string into label chunks to trigger multiline labels in apexcharts
const convertLabel = (label: string) => label
  .split(' ')
  .reduce<string[][]>((labelParts, item, index) => {
    const chunkIndex = Math.floor(index / 3)

    if (!labelParts[chunkIndex]) {
      labelParts[chunkIndex] = []
    }

    labelParts[chunkIndex].push(item)

    return labelParts
  }, [])
  .map<string>((labelPrts) => labelPrts.join(' '))

export default Vue.extend<Data, {}, Computed, Props>({
  components: { EventOverlay },
  name: 'RuleChart',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, required: true },
    time: { type: Number }
  },
  data: () => ({
    selectedEvents: [],
    eventsPerRule: {},
    rules: {}
  }),
  watch: {
    events: {
      immediate: true,
      handler (events: FalcoEvent[]) {
        const rules: { [rule: string]: number } = {}
        const eventsPerRule: { [rule: string]: FalcoEvent[] } = {}

        events.forEach((event) => {
          if (Object.prototype.hasOwnProperty.call(rules, event.rule) === false) {
            rules[event.rule] = 0
            eventsPerRule[event.rule] = []
          }

          rules[event.rule] += 1
          eventsPerRule[event.rule].push(event)

          return rules
        })

        this.rules = Object.keys(rules).sort().reduce<{ [rule: string]: number }>((acc, rule) => {
          acc[rule] = rules[rule]

          return acc
        }, {})

        this.eventsPerRule = eventsPerRule
      }
    }
  },
  computed: {
    ...mapState(['displayMode']),
    theme () {
      if (this.displayMode === DisplayMode.LIGHT) return {}

      return {
        theme: { mode: 'dark' }
      }
    },
    chart () {
      const labels = Object.keys(this.rules)

      return {
        series: [{
          name: 'Counter',
          data: Object.values(this.rules)
        }],
        chartOptions: {
          ...this.theme,
          chart: {
            height: 350,
            type: 'bar',
            toolbar: {
              show: false
            },
            selection: {
              enabled: false
            },
            zoom: {
              enabled: false
            },
            events: {
              dataPointSelection: (_: any, __: any, config: { dataPointIndex: number }) => {
                this.selectedEvents = [...this.eventsPerRule[labels[config.dataPointIndex]]]
              }
            }
          },
          dataLabels: {
            enabled: false
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
            labels: {
              rotate: -45,
              rotateAlways: true,
              hideOverlappingLabels: false,
              maxHeight: 200,
              style: {
                colors: [],
                fontSize: '12px'
              }
            },
            categories: Object.keys(this.rules).map<Array<string[]|string>>((label) => convertLabel(label)),
            tickPlacement: 'on'
          }
        }
      }
    }
  }
})
</script>
