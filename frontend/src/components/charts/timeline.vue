<template>
  <div>
    <canvas
      :id="id"
      height="52"
    ></canvas>
  </div>
</template>

<script>
import { Chart as ChartJS, registerables } from 'chart.js';
import dayjs from 'dayjs';
import 'chartjs-adapter-moment';
import { requests } from '../../http';
import { utils, PRIORITY_ORDER } from '../../utils';

ChartJS.register(...registerables);

export default {
  name: 'TimelineChart',
  props: {
    id: { type: String, default: 'timeline-chart' },
    // datasetIdKey: { type: String, default: 'label' },
    width: { type: Number, default: 400 },
    height: { type: Number, default: 400 },
    groupby: { type: String, default: '' },
    filters: {
      type: Object,
      default() {
        return {
          priorities: [],
          rule: [],
          tags: [],
          hostnames: [],
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
    },
    ticer: {
      handler() {
        this.updateChart();
      },
    },
  },
  data() {
    return {
      chart: {},
      bgColor: '#f00',
      chartData: {
        type: 'bar',
        data: {
          datasets: [{
            label: 'test',
            backgroundColor: '#3a9bf0',
            data: [],
          }],
        },
        options: {
          responsive: true,
          animation: false,
          onClick: this.pickItem,
          interaction: { mode: 'index', intersect: false },
          plugins: {
            legend: {
              position: 'top',
              align: 'end',
              labels: { padding: 14 },
            },
          },
          scales: {
            x: {
              type: 'time',
              display: true,
              stacked: true,
              grid: { display: false },
              ticks: { maxRotation: 0, autoSkipPadding: 24 },
            },
            y: {
              display: true,
              stacked: true,
              border: { display: false },
              grid: { drawBorder: false },
              ticks: {
                beginAtZero: true,
                maxTicksLimit: 5,
                precision: 0,
              },
            },
          },
        },
      },
      stats: [],
    };
  },
  methods: {
    pickItem(event, array) {
      if (array.length > 0) {
        this.$emit('picked-item', this.groupby, this.chartData.data.datasets[array[0].datasetIndex].label);
      }
    },
    searchEvents(page, last) {
      const limit = 500;
      let results = [];
      let returned = 0;
      let l = last;
      requests.searchEvents(
        this.filters.sources,
        this.filters.hostnames,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
        page,
        limit,
      )
        .then((response) => {
          results = response.data.results; // eslint-disable-line
          returned = response.data.statistics.returned; // eslint-disable-line
          if (results.length === 0) {
            this.chart.update();
            return;
          }
          let unit = '';
          let delta = '';
          let oldest = '';
          switch (this.filters.since) {
            case '5min':
              this.chartData.options.scales.x.time.unit = 'second';
              unit = 'second'; delta = 5; oldest = 300; break;
            case '15min':
              this.chartData.options.scales.x.time.unit = 'second';
              unit = 'second'; delta = 25; oldest = 1500; break;
            case '30min':
              this.chartData.options.scales.x.time.unit = 'minute';
              unit = 'second'; delta = 48; oldest = 3000; break;
            case '1h':
              this.chartData.options.scales.x.time.unit = 'minute';
              unit = 'minute'; delta = 1; oldest = 60; break;
            case '2h':
              this.chartData.options.scales.x.time.unit = 'minute';
              unit = 'minute'; delta = 2; oldest = 120; break;
            case '5h':
              this.chartData.options.scales.x.time.unit = 'minute';
              unit = 'minute'; delta = 5; oldest = 300; break;
            case '12h':
              this.chartData.options.scales.x.time.unit = 'minute';
              unit = 'minute'; delta = 12; oldest = 720; break;
            case '24h':
              this.chartData.options.scales.x.time.unit = 'hour';
              unit = 'minute'; delta = 24; oldest = 1440; break;
            case '48h':
              this.chartData.options.scales.x.time.unit = 'hour';
              unit = 'minute'; delta = 48; oldest = 2880; break;
            case '1w':
              this.chartData.options.scales.x.time.unit = 'day';
              unit = 'hour'; delta = 3; oldest = 168; break;
            case '2w':
              this.chartData.options.scales.x.time.unit = 'day';
              unit = 'hour'; delta = 6; oldest = 336; break;
            case '1M':
              this.chartData.options.scales.x.time.unit = 'day';
              unit = 'hour'; delta = 12; oldest = 724; break;
            case '3M':
              this.chartData.options.scales.x.time.unit = 'day';
              unit = 'hour'; delta = 36; oldest = 2160; break;
            case '6M':
              this.chartData.options.scales.x.time.unit = 'month';
              unit = 'hour'; delta = 72; oldest = 4320; break;
            case '1y':
              this.chartData.options.scales.x.time.unit = 'month';
              unit = 'day'; delta = 15; oldest = 365; break;
            case '2y':
              this.chartData.options.scales.x.time.unit = 'quarter';
              unit = 'day'; delta = 30; oldest = 730; break;
            default:
              break;
          }
          if (l === undefined) {
            l = dayjs().add(delta, unit).toISOString();
          }
          Object.values(results).forEach((value) => {
            let f = '';
            switch (this.groupby) {
              case 'priority':
                f = value.priority;
                break;
              case 'source':
                f = value.source;
                break;
              default:
                f = 'all';
                break;
            }
            if (this.stats[f] === undefined) {
              this.stats[f] = {
                count: 0,
                data: [{ x: l, y: 0 }],
              };
            }
            this.stats[f].count += 1;
            if (dayjs(value.time)
              .isBefore(dayjs(l).subtract(delta, unit))) {
              Object.keys(this.stats).forEach((key) => {
                this.stats[key].data.push({
                  x: l,
                  y: this.stats[key].count,
                });
                this.stats[key].count = 0;
              });
              l = value.time;
              this.stats[f].count = 0;
            }
          });
          if (limit === returned) {
            let d = 1;
            if (page === 0) {
              d = 2;
            }
            this.searchEvents(page + d, l);
          }
          let i = 0;
          Object.keys(this.stats).forEach((key) => {
            this.stats[key].data.push({
              x: dayjs().subtract(oldest - delta, unit).toISOString(),
              y: 0,
            });
            let bgc = '';
            switch (this.groupby) {
              case 'priority':
                bgc = utils.priorityToColor(key);
                break;
              case 'source':
                bgc = utils.stringToColor(key);
                break;
              default:
                bgc = utils.cssVar('--accent');
                break;
            }
            this.chartData.data.datasets[i] = {
              label: key,
              barThickness: this.barWidth(oldest / delta),
              borderRadius: 1,
              backgroundColor: bgc,
              data: this.stats[key].data,
            };
            i += 1;
          });
          this.chartData.data.datasets.sort((data1, data2) => this.compareDatasets(data1, data2));
          this.chart.update();
        });
    },
    barWidth(buckets) {
      const width = (this.chart && this.chart.chartArea)
        ? this.chart.chartArea.right - this.chart.chartArea.left
        : 1000;
      return Math.max(2, Math.min(18, Math.floor((width / buckets) * 0.72)));
    },
    compareDatasets(data1, data2) {
      if (this.groupby === 'priority') {
        return PRIORITY_ORDER.indexOf(data1.label) - PRIORITY_ORDER.indexOf(data2.label);
      }
      const obj1 = data1.label.toUpperCase();
      const obj2 = data2.label.toUpperCase();
      if (obj1 < obj2) { return -1; }
      if (obj1 > obj2) { return 1; }
      return 0;
    },
    async updateChart() {
      this.chartData.data.datasets = [];
      this.stats = {};
      this.searchEvents(0);
    },
  },
  mounted() {
    const ctx = document.getElementById(this.id);
    this.chart = new ChartJS(ctx, this.chartData); // eslint-disable-line
    this.updateChart();
  },
};
</script>
