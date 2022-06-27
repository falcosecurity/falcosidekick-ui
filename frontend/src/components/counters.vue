<template>
  <div>
    <v-chip
    dark
    class="rounded-0"
    color="blue lighten-3">Total
    </v-chip>
    <v-chip
    dark
    class="rounded-0"
    style="margin-right: 5px;">{{countByPriority.statistics.all}}
    </v-chip>
    <span
      v-for="(value,key) in countByPriority.results" :key="key">
      <v-chip
      dark
      @click="addItemToList(key)"
      class="rounded-0"
      :color="priorityToColor(key)">
        {{key}}
      </v-chip>
      <v-chip
      dark
      class="rounded-0"
      style="margin-right: 5px;">
        {{value}}
      </v-chip>
    </span>
  </div>
</template>

<script>
import { requests } from '../http';
import { utils } from '../utils';

export default {
  name: 'Counters',
  props: {
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
  data() {
    return {
      count: 0,
      countByPriority: {
        statistics: {
          all: 0,
        },
      },
    };
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
  methods: {
    priorityToColor(prio) {
      return utils.priorityToColor(prio);
    },
    updateChart() {
      // this.countByPriority = {
      //   statistics: {
      //     all: 0,
      //   },
      // };
      requests.countByEvents(
        'priority',
        this.filters.sources,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
      )
        .then((response) => {
          this.countByPriority = response.data;
        });
    },
    addItemToList(item) {
      this.$emit('add-item-to-filters', item);
    },
  },
  mounted() {
    this.updateChart();
  },
};
</script>
