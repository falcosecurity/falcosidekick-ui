<template>
  <div>
    <Bar
      :chart-options="chartOptions"
      :chart-data="chartData"
      :chart-id="chartId"
      :dataset-id-key="datasetIdKey"
      :width="width"
      :height="height"
      v-if="loaded"
    />
  </div>
</template>

<script>
import { Bar } from 'vue-chartjs/legacy';
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  BarElement,
  CategoryScale,
  LinearScale,
} from 'chart.js';
import { requests } from '../../http';

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);

export default {
  name: 'BarChart',
  components: {
    Bar,
  },
  props: {
    list: { type: String },
    chartId: { type: String, default: 'pie-chart' },
    datasetIdKey: { type: String, default: 'label' },
    width: { type: Number, default: 400 },
    height: { type: Number, default: 400 },
    filters: {
      type: Object,
      default() {
        return {
          priorities: [],
          rule: '',
          tags: [],
          sources: [],
          search: '',
          since: '',
        };
      },
    },
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
  },
  watch: {
    filters: {
      handler() {
        this.updateChart();
      },
      deep: true,
      immediate: true,
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
      values: {},
      bgColor: '#1976d2',
      chartData: {
        labels: [],
        datasets: [],
      },
      chartOptions: {
        animation: false,
        responsive: true,
        maintainAspectRatio: false,
        onClick: this.pickItem,
      },
    };
  },
  methods: {
    pickItem(event, array) {
      this.$emit('picked-item', this.list, this.chartData.labels[array[0].index]);
    },
    async updateChart() {
      this.loaded = false;
      await requests.countByEvents(
        this.list,
        this.filters.sources,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
      )
        .then((response) => {
          this.values = Object.entries(response.data.results)
            .sort(([, v1], [, v2]) => v2 - v1)
            .reduce((obj, [k, v]) => ({ ...obj, [k]: v }), {});
        });
      let lb = []; // eslint-disable-line
      let dt = []; // eslint-disable-line
      Object.entries(this.values).forEach(([key, value]) => {
        lb.push(key);
        dt.push(value);
      });
      this.chartData.labels = lb;
      this.chartData.datasets = [{
        label: this.list,
        backgroundColor: this.bgColor,
        data: dt,
      }];
      this.loaded = true;
    },
  },
  created() {
    if (typeof this.$route.query.source !== 'undefined') {
      this.filters.sources = this.$route.query.source;
    }
    if (typeof this.$route.query.priority !== 'undefined') {
      this.filters.priorities = this.$route.query.priority;
    }
    if (typeof this.$route.query.rule !== 'undefined') {
      this.filters.rule = this.$route.query.rule;
    }
    if (typeof this.$route.query.tags !== 'undefined') {
      this.filters.tags = this.$route.query.tags;
    }
    if (typeof this.$route.query.filter !== 'undefined') {
      this.filters.search = this.$route.query.filter;
    }
    if (typeof this.$route.query.since !== 'undefined') {
      this.filters.since = this.$route.query.since;
    }
  },
  mounted() {
    this.updateChart();
  },
};
</script>
