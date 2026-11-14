<template>
  <div>
    <PageToolbar title="Dashboard"></PageToolbar>
    <div class="page">
      <div class="grid grid-3">
        <section class="panel">
          <header><h3>Sources</h3></header>
          <PieChart
            chart-id="countBySource" list="source" caption="events"
            :filters="filters" @picked-item="addToFilters"></PieChart>
        </section>
        <section class="panel">
          <header><h3>Priorities</h3></header>
          <PieChart
            chart-id="countByPriority" list="priority" caption="events"
            :filters="filters" @picked-item="addToFilters"></PieChart>
        </section>
        <section class="panel">
          <header><h3>Tags</h3></header>
          <PieChart
            chart-id="countByTag" list="tags" caption="tag hits"
            :filters="filters" @picked-item="addToFilters"></PieChart>
        </section>
      </div>

      <section class="panel">
        <header><h3>Timeline by priority</h3></header>
        <TimelineChart
          id="timeline-priority" groupby="priority"
          :filters="filters" @picked-item="addToFilters"></TimelineChart>
      </section>

      <section class="panel">
        <header><h3>Timeline by source</h3></header>
        <TimelineChart
          id="timeline-source" groupby="source"
          :filters="filters" @picked-item="addToFilters"></TimelineChart>
      </section>

      <section class="panel">
        <header><h3>Top rules</h3></header>
        <BarChart list="rule" :filters="filters" @picked-item="addToFilters"></BarChart>
      </section>
    </div>
  </div>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import BarChart from '../components/charts/bar.vue';
import PieChart from '../components/charts/pie.vue';
import TimelineChart from '../components/charts/timeline.vue';
import PageToolbar from '../components/PageToolbar.vue';

export default {
  name: 'Dashboard',
  components: {
    PieChart,
    BarChart,
    TimelineChart,
    PageToolbar,
  },
  computed: {
    ...mapState(['filters']),
  },
  methods: {
    ...mapActions(['addFilter']),
    addToFilters(list, item) {
      this.addFilter({ list, item });
    },
  },
};
</script>

<style>
.page {
  padding: 20px 24px 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.grid {
  display: grid;
  gap: 16px;
}

.grid-3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.grid-2 {
  grid-template-columns: repeat(2, minmax(0, 1fr));
  align-items: start;
}

@media (max-width: 1280px) {
  .grid-3 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .grid-2 { grid-template-columns: minmax(0, 1fr); }
}

@media (max-width: 900px) {
  .grid-3 { grid-template-columns: minmax(0, 1fr); }
}

.panel {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 14px 18px 16px;
  box-shadow: var(--shadow);
  min-width: 0;
}

.panel > header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.panel h3 {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
}
</style>
