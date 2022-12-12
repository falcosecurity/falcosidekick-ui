<template>
  <v-card class="elevation-0">
    <v-row>
      <v-col
      cols="12"
      sm="12">
        <Filters
        @send-filters="setFilters(...arguments)"
        :addItem="newItem"
        >
        </Filters>
      </v-col>
    </v-row>
    <v-row>
      <v-col
      cols="12"
      sm="4">
        <v-card class="elevation-0">
          <v-card-title>Sources</v-card-title>
          <PieChart
          id="countBySource"
          list="source"
          :filters="filters"
          @picked-item="addToFilters(...arguments)">
          </PieChart>
          </v-card>
      </v-col>
      <v-col
      cols="12"
      sm="4">
        <v-card class="elevation-0">
          <v-card-title>Priorities</v-card-title>
          <PieChart
          id="countByPriority"
          list="priority"
          :filters="filters"
          @picked-item="addToFilters(...arguments)">
          </PieChart>
          </v-card>
      </v-col>
      <v-col
      cols="12"
      sm="4">
        <v-card class="elevation-0">
          <v-card-title>Tags</v-card-title>
          <PieChart
          id="countByTag"
          list="tags"
          :filters="filters"
          @picked-item="addToFilters(...arguments)">
          </PieChart>
          </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col
      cols="12"
      sm="11">
        <v-card class="elevation-0">
          <v-card-title>Rules</v-card-title>
            <BarChart
            id="countByRules"
            list="rule"
            :filters="filters"
            @picked-item="addToFilters(...arguments)">
            </BarChart>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col
      cols="12"
      sm="12">
        <v-card class="elevation-0">
          <v-card-title>Timeline by Priority</v-card-title>
            <TimelineChart
            id="timeline-priority"
            :filters="filters"
            groupby="priority"
            @picked-item="addToFilters(...arguments)">
            </TimelineChart>
        </v-card>
      </v-col>
    </v-row>
    <v-row>
      <v-col
      cols="12"
      sm="12">
        <v-card class="elevation-0">
          <v-card-title>Timeline by Source</v-card-title>
            <TimelineChart
            id="timeline-source"
            :filters="filters"
            groupby="source"
            @picked-item="addToFilters(...arguments)">
            </TimelineChart>
        </v-card>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import BarChart from '../components/charts/bar.vue';
import PieChart from '../components/charts/pie.vue';
import TimelineChart from '../components/charts/timeline.vue';
import Filters from '../components/filters.vue';

export default {
  name: 'Dashboard',
  components: {
    PieChart,
    BarChart,
    TimelineChart,
    Filters,
  },
  data() {
    return {
      newItem: {
        list: '',
        item: '',
      },
      filters: {
        priorities: [],
        sources: [],
        hostnames: [],
        tags: [],
        rule: '',
        since: '24h',
        search: '',
      },
    };
  },
  methods: {
    setFilters(f) {
      this.filters = f;
    },
    addToFilters(l, i) {
      this.newItem = {
        list: l,
        item: i,
      };
    },
  },
};
</script>

