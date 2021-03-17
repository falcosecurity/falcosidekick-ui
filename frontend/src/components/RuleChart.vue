<template>
  <v-card style="height: 100%">
    <v-card-title>
      Rules
    </v-card-title>
    <v-card-text>
      <apexchart type="bar" height="350" :options="chart.chartOptions" :series="chart.series"></apexchart>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { FalcoEvent } from '@/api/model'
import Vue, { PropType } from 'vue'

type Props = {
  time: number|null;
  events: FalcoEvent[];
}

type Computed = {
  chart: any;
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

export default Vue.extend<{}, {}, Computed, Props>({
  name: 'RuleChart',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, required: true },
    time: { type: Number }
  },
  computed: {
    rules () {
      const unsorted = this.events.reduce<{ [rule: string]: number }>((acc, event) => {
        if (Object.prototype.hasOwnProperty.call(acc, event.rule) === false) {
          acc[event.rule] = 0
        }

        acc[event.rule] += 1

        return acc
      }, {})

      const sorted = Object.keys(unsorted).sort().reduce<{ [rule: string]: number }>((acc, rule) => {
        acc[rule] = unsorted[rule]

        return acc
      }, {})

      return sorted
    },
    chart () {
      return {
        series: [{
          name: 'Counter',
          data: Object.values(this.rules)
        }],
        chartOptions: {
          chart: {
            height: 350,
            type: 'bar',
            toolbar: {
              show: false
            },
            zoom: {
              enabled: false
            }
          },
          dataLabels: {
            enabled: false
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
