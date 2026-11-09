<template>
  <div class="donut">
    <div class="donut-chart">
      <canvas ref="canvas" :id="chartId" width="150" height="150"></canvas>
      <div class="donut-center">
        <div class="donut-total">{{ formatNumber(total) }}</div>
        <div class="donut-caption">{{ caption }}</div>
      </div>
    </div>
    <ul class="donut-legend">
      <li
        v-for="item in legend" :key="item.name"
        @click="$emit('picked-item', list, item.name)">
        <span class="dot" :style="{ background: item.color }"></span>
        <span class="name" :title="item.name">{{ item.name }}</span>
        <span class="count">{{ formatNumber(item.count) }}</span>
        <span class="pct">{{ item.pct }}%</span>
      </li>
      <li v-if="hidden > 0" class="more">+{{ hidden }} more</li>
      <li v-if="!items.length && loaded" class="empty">No events</li>
    </ul>
  </div>
</template>

<script>
import {
  Chart as ChartJS,
  DoughnutController,
  Tooltip,
  ArcElement,
} from 'chart.js';
import numeral from 'numeral';
import { requests } from '../../http';
import { utils } from '../../utils';

ChartJS.register(DoughnutController, Tooltip, ArcElement);

const MAX_LEGEND = 6;

export default {
  name: 'PieChart',
  props: {
    list: { type: String },
    caption: { type: String, default: 'events' },
    chartId: { type: String, default: 'pie-chart' },
    filters: { type: Object, required: true },
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
    total() {
      return this.items.reduce((s, i) => s + i.count, 0);
    },
    legend() {
      return this.items.slice(0, MAX_LEGEND).map(i => ({
        ...i,
        pct: this.total ? Math.round((i.count / this.total) * 100) : 0,
      }));
    },
    hidden() {
      return Math.max(0, this.items.length - MAX_LEGEND);
    },
  },
  watch: {
    filters: {
      handler() {
        this.updateChart();
      },
      deep: true,
    },
    ticer: {
      handler() {
        this.updateChart();
      },
    },
  },
  data() {
    return {
      loaded: false,
      items: [],
    };
  },
  methods: {
    formatNumber(n) {
      return n >= 10000 ? numeral(n).format('0.0a') : numeral(n).format('0,0');
    },
    pickItem(event, array) {
      if (array.length > 0) {
        this.$emit('picked-item', this.list, this.items[array[0].index].name);
      }
    },
    colorFor(key) {
      return this.list === 'priority' ? utils.priorityToColor(key) : utils.stringToColor(key);
    },
    async updateChart() {
      const response = await requests.countByEvents(
        this.list,
        this.filters.sources,
        this.filters.hostnames,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
      );
      this.items = Object.entries(response.data.results || {})
        .sort(([, v1], [, v2]) => v2 - v1)
        .map(([name, count]) => ({ name, count, color: this.colorFor(name) }));
      const data = {
        labels: this.items.map(i => i.name),
        datasets: [{
          label: this.list,
          backgroundColor: this.items.map(i => i.color),
          borderColor: utils.cssVar('--surface'),
          borderWidth: 2,
          hoverOffset: 4,
          data: this.items.map(i => i.count),
        }],
      };
      if (this.chart) {
        this.chart.data = data;
        this.chart.update();
      } else {
        this.chart = new ChartJS(this.$refs.canvas, {
          type: 'doughnut',
          data,
          options: {
            animation: false,
            responsive: true,
            maintainAspectRatio: false,
            cutout: '72%',
            onClick: this.pickItem,
            plugins: { legend: { display: false } },
          },
        });
      }
      this.loaded = true;
    },
  },
  mounted() {
    this.updateChart();
  },
  beforeDestroy() {
    if (this.chart) this.chart.destroy();
  },
};
</script>

<style scoped>
.donut {
  display: flex;
  align-items: center;
  gap: 22px;
  padding: 4px 0;
}

.donut-chart {
  position: relative;
  flex: none;
  width: 150px;
  height: 150px;
}

.donut-center {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.donut-total {
  font-size: 22px;
  font-weight: 300;
  color: var(--text);
  font-variant-numeric: tabular-nums;
}

.donut-caption {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}

.donut-legend {
  flex: 1;
  min-width: 0;
  list-style: none;
  padding: 0;
  margin: 0;
  font-size: 12px;
}

.donut-legend li {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 3px 6px;
  margin: 0 -6px;
  border-radius: 4px;
  cursor: pointer;
  color: var(--text);
}

.donut-legend li:hover {
  background: var(--row-hover);
}

.dot {
  flex: none;
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

.name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.count {
  font-variant-numeric: tabular-nums;
  color: var(--text);
}

.pct {
  width: 34px;
  text-align: right;
  color: var(--text-muted);
  font-variant-numeric: tabular-nums;
}

.donut-legend li.more,
.donut-legend li.empty {
  color: var(--text-muted);
  cursor: default;
}

.donut-legend li.more:hover {
  background: none;
}
</style>
