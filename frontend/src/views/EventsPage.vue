<template>
  <v-card class="elevation-0">
    <v-row>
      <v-col
      cols="12"
      sm="12">
        <Filters
        @send-filters="setFilters(...arguments)"
        :addItem="newItem">
        </Filters>
      </v-col>
    </v-row>
    <v-row>
      <v-data-table
        class="mt-10 ml-5 mr-5"
        :search="search"
        :page='page'
        :headers='headers'
        :items='events'
        item-key="time"
        :options.sync='options'
        :server-items-length='totalEvents'
        :items-per-page="itemsPerPage"
        :loading='loading'
        :footer-props="{
          'items-per-page-options': itemsPerPageInterval,
          firstIcon: 'mdi-minus',
          lastIcon: 'mdi-plus',
          prevIcon: 'mdi-chevron-left',
          nextIcon: 'mdi-chevron-right'
        }"
      >
        <template v-slot:item="{item}">
          <tr>
            <td>{{ item.time | formatDate }}</td>
            <td>
              <v-chip dark
              @click="addToFilters('sources', item.source)"
              :color="stringToColor(item.source)">
              {{ item.source }}
              </v-chip>
            </td>
            <td>
              <v-chip dark v-if="item.hostname"
              @click="addToFilters('hostnames', item.hostname)"
              :color="stringToColor(item.hostname)">
              {{ item.hostname }}
              </v-chip>
            </td>
            <td>
              <v-chip
              @click="addToFilters('priorities', item.priority)"
              :color="priorityToColor(item.priority)"
              dark>
                {{ item.priority }}
              </v-chip>
            </td>
            <td>{{ item.rule }}</td>
            <td>
              <div>{{ item.output }}</div>
              <div>
                <span
                v-for="(value,key) in item.output_fields" :key="key">
                <v-chip small label
                @click="addToFilters('search', key)"
                class="rounded-0 mb-1"
                color="blue lighten-3">
                  {{key}}
                </v-chip>
                <v-chip small label
                @click="addToFilters('search', value)"
                class="rounded-0 mb-1"
                style="margin-left: -4px; margin-right: 5px;">
                  {{value}}
                </v-chip>
                </span>
              </div>
            </td>
            <td>
              <v-chip
              class="mb-1 mr-1 mt-1"
              @click="addToFilters('tags', tag)"
              v-for="(tag, index) in item.tags" :key="index"
              dark small
              :color="stringToColor(tag)">
              {{tag}}
              </v-chip>
            </td>
          </tr>
        </template>
      </v-data-table>
    </v-row>
  </v-card>
</template>

<script>
import { requests } from '../http';
import { utils } from '../utils';
import Counters from '../components/counters.vue';
import Filters from '../components/filters.vue';

export default {
  name: 'DatatableComponent',
  components: {
    Counters,
    Filters,
  },
  data() {
    return {
      search: '',
      page: 1,
      itemsPerPage: 10,
      itemsPerPageInterval: [10, 50, 100, 200],
      totalEvents: 0,
      events: [],
      priorities: [],
      rules: [],
      sources: [],
      hostnames: [],
      tags: [],
      filters: {
        sources: [],
        hostnames: [],
        priorities: [],
        rule: '',
        tags: [],
        since: '24h',
        search: '',
      },
      newItem: {
        list: '',
        item: '',
      },
      loading: true,
      options: {
        page: 1,
        itemsPerPage: 10,
      },
      headers: [
        { text: 'Timestamp', value: 'time' },
        { text: 'Source', value: 'source' },
        { text: 'Hostname', value: 'hostname' },
        { text: 'Priority', value: 'priority' },
        { text: 'Rule', value: 'rule' },
        { text: 'Output', value: 'output' },
        { text: 'Tags', value: 'tags' },
      ],
      debounce: null,
    };
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
  },
  watch: {
    options: {
      handler() {
        this.searchEvents();
      },
    },
    filters: {
      handler() {
        this.resetPage();
        this.searchEvents();
      },
      deep: true,
    },
    search: {
      handler() {
        clearTimeout(this.debounce);
        this.debounce = setTimeout(() => {
          this.resetPage();
          this.searchEvents();
        }, 600);
      },
    },
    ticer: {
      handler() {
        if (this.page === 1) {
          this.searchEvents();
        }
      },
    },
  },
  methods: {
    resetPage() {
      this.page = 1;
      this.options.page = 1;
    },
    searchEvents() {
      this.loading = true;
      const { page, itemsPerPage } = this.options;
      if (this.page !== page) {
        this.page = page;
        this.$router.push({ query: { ...this.$route.query, page: this.page } });
      }
      if (this.itemsPerPage !== itemsPerPage) {
        this.itemsPerPage = itemsPerPage;
        this.$router.push({ query: { ...this.$route.query, limit: this.itemsPerPage } });
      }
      if (this.$route.query.filter !== this.search && this.search !== '') {
        this.$router.push({ query: { ...this.$route.query, filter: this.search } });
      }
      requests.searchEvents(
        this.filters.sources,
        this.filters.hostnames,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
        page,
        itemsPerPage,
      )
        .then((response) => {
          this.loading = false;
          this.events = response.data.results;
          this.totalEvents = response.data.statistics.all;
        });
    },
    priorityToColor(prio) {
      return utils.priorityToColor(prio);
    },
    stringToColor(str) {
      return utils.stringToColor(str);
    },
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
  mounted() {
    this.searchEvents();
  },
};
</script>
