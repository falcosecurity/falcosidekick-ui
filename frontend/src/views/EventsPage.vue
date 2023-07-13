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
      <v-col
      align="right"
      cols="12"
      sm="12"
      class="pr-7"
      >
        <v-btn
        color="info"
        small
        @click="saveFile()"
        >
          Export
        </v-btn>
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
            <td>
              <v-btn
                x-small
                :icon="true"
                @click="showDialog(item);"
              >
                <v-icon>mdi-code-json</v-icon>
              </v-btn>
            </td>
          </tr>
        </template>
      </v-data-table>
    </v-row>
    <v-dialog
      v-model="dialog"
      width="auto"
      >
      <v-tabs
      text
      >
        <v-tab href="#details" @click="dialogDetails = true; dialogJson = false;">
          Details
        </v-tab>
        <v-tab href="#json" @click="dialogDetails = false; dialogJson = true;">
          Json
        </v-tab>
      </v-tabs>
      <v-card>
        <v-card-text v-show="dialogDetails">
          <v-list-item style="padding-top: 10px;">
            <v-list-item-content>
              <v-list-item-title>Time</v-list-item-title>
                <v-list-item-subtitle>{{json.time}}</v-list-item-subtitle>
              <v-list-item-title>Source</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip dark
                  @click="addToFilters('sources', json.source)"
                  :color="stringToColor(String(json.source))">{{ json.source }}</v-chip>
                </v-list-item-subtitle>
              <v-list-item-title>Hostname</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip dark
                  @click="addToFilters('hostnames', json.hostname)"
                  :color="stringToColor(String(json.hostname))">{{ json.hostname }}</v-chip>
                </v-list-item-subtitle>
              <v-list-item-title>Priority</v-list-item-title>
                <v-list-item-subtitle>
                  <v-chip dark
                  @click="addToFilters('priorities', json.priority)"
                  :color="priorityToColor(String(json.priority))">{{ json.priority }}</v-chip>
                </v-list-item-subtitle>
              <v-list-item-title>Rule</v-list-item-title>
              <v-list-item-subtitle>{{json.rule}}</v-list-item-subtitle>
              <v-list-item-title>Output</v-list-item-title>
              <v-list-item-subtitle>
                <pre style="white-space: pre-wrap;">{{json.output}}</pre>
              </v-list-item-subtitle>
              <v-list-item-title>Fields</v-list-item-title>
              <v-list-item-subtitle>
                <v-list-item-subtitle>
                  <span style="white-space: pre-wrap;"
                  v-for="(value,key) in json.output_fields" :key="key" :value="value">
                  <v-chip label
                  @click="addToFilters('search', key)"
                  class="rounded-0 mb-1"
                  color="blue lighten-3">
                    {{key}}
                  </v-chip>
                  <v-chip label
                  @click="addToFilters('search', value)"
                  class="rounded-0 mb-1"
                  style="margin-left: -4px; margin-right: 5px;">
                    {{value}}
                  </v-chip>
                  </span>
                </v-list-item-subtitle>
              </v-list-item-subtitle>
              <v-list-item-title>Tags</v-list-item-title>
              <v-list-item-subtitle>
                <v-list-item-subtitle>
                  <v-chip
                  v-for="(tag, index) in json.tags" :key="index"
                  dark
                  @click="addToFilters('tags', tag)"
                  class="mb-1 mr-1 mt-1"
                  :color="stringToColor(String(tag))">{{ tag }}</v-chip>
                </v-list-item-subtitle>
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-card-text>
        <v-card-text v-show="dialogJson">
          <pre style="white-space: pre-wrap;">
            <vue-json-pretty
              :data="json"
              :collapsedOnClickBrackets="false"
            />
          </pre>
        </v-card-text>
        <v-card-actions class="justify-center">
          <v-btn
            style="margin-right: 5px;"
            color="primary" large @click="dialog = false"
          >
            Close
          </v-btn>
          <v-tooltip
          v-model="copied"
          top
          >
            <template v-slot:activator="{ attrs }">
              <v-btn
                color="normal"
                large
                @click="copied = true; copyToClipBoard()"
                v-bind="attrs"
              >
              COPY JSON
              </v-btn>
            </template>
            <span>Copied</span>
          </v-tooltip>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-card>
</template>

<script>
import 'vue-json-pretty/lib/styles.css';
import VueJsonPretty from 'vue-json-pretty';
import { requests } from '../http';
import { utils } from '../utils';
import Counters from '../components/counters.vue';
import Filters from '../components/filters.vue';

export default {
  name: 'DatatableComponent',
  components: {
    Counters,
    Filters,
    VueJsonPretty,
  },
  data() {
    return {
      search: '',
      page: 1,
      itemsPerPage: 10,
      itemsPerPageInterval: [10, 50, 100, 200, 500, 1000],
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
        { text: '', value: 'json' },
      ],
      debounce: null,
      dialog: false,
      dialogDetails: true,
      dialogJson: false,
      json: '',
      copied: false,
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
    showDialog(i) {
      this.dialog = true;
      this.json = i;
    },
    copyToClipBoard() {
      navigator.clipboard.writeText(JSON.stringify(this.json));
    },
    saveFile() {
      let text;
      requests.searchEvents(
        this.filters.sources,
        this.filters.hostnames,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
        0,
        10000,
      )
        .then((response) => {
          text = JSON.stringify(response.data.results);
          const filename = 'falco_events.json';
          const element = document.createElement('a');
          const href = `data:application/json;charset=utf-8,${encodeURIComponent(text)}`;
          element.setAttribute('href', href);
          element.setAttribute('download', filename);

          element.style.display = 'none';
          document.body.appendChild(element);

          element.click();
          document.body.removeChild(element);
        });
    },
  },
  mounted() {
    this.searchEvents();
  },
};
</script>

<style>
.vjs-key {
  color: rgb(43, 51, 42);
}
.vjs-value-string {
  color: rgb(73, 105, 247);
}
.v-list-item__subtitle {
  padding-bottom: 15px;
  font-size: 125%;
}
.v-list-item__title {
  padding-bottom: 2px;
}
</style>
