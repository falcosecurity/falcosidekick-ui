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
  </v-card>
</template>

<script lang="ts">
import { colorCodes, initStats } from '@/api/mapper'
import { FalcoEvent, Priority, Stats } from '@/api/model'
import Vue, { PropType } from 'vue'
import Wait from './Wait.vue'

type Props = {
  time: number|null;
  events: FalcoEvent[];
}

type Computed = {
  pie: any;
  stats: Stats;
}

export default Vue.extend<{}, {}, Computed, Props>({
  components: { Wait },
  name: 'PieChart',
  props: {
    events: { type: Array as PropType<FalcoEvent[]>, required: true },
    time: { type: Number }
  },
  computed: {
    stats () {
      const stats = initStats()

      this.events.forEach((event) => {
        stats[event.priority] += 1
      })

      return stats
    },
    pie () {
      return {
        series: Object.values(this.stats),
        chartOptions: {
          chart: {
            height: 350,
            type: 'donut'
          },
          dataLabels: {
            dropShadow: {
              enabled: true,
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              blur: 1,
              color: '#000',
              opacity: 1
            }
          },
          plotOptions: {
            pie: {
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
          labels: Object.keys(this.stats),
          colors: Object.keys(this.stats).map(priority => colorCodes[priority as Priority])
        }
      }
    }
  }
})
</script>
