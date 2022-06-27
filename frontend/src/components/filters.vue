<template>
  <div id="filters">
    <v-row>
      <v-col
      class="d-flex ml-4"
      cols="12"
      sm="2">
        <v-select
        v-model="filters.sources"
        chips
        multiple
        :items="sources"
        label="Sources"
        >
          <template v-slot:prepend-item>
            <v-list-item @click="unselect('sources')">
              <v-list-item-action>
                <b>Unselect All</b>
              </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>
          </template>
          <template v-slot:selection="{ item }">
            <v-chip dark :color="stringToColor(item)" small
            >{{ item }}</v-chip>
          </template>
          <template v-slot:item="{ item, attrs, on }">
            <v-list-item v-on="on" v-bind="attrs" #default="{ active }">
              <v-list-item-action>
                <v-checkbox :input-value="active"></v-checkbox>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>
                  <v-chip dark :color="stringToColor(item)" small
                  >{{ item }}</v-chip>
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-select>
      </v-col>
      <v-col
      class="d-flex ml-4"
      cols="12"
      sm="2">
        <v-select
        v-model="filters.priorities"
        chips
        multiple
        :items="priorities"
        label="Priorities"
        >
          <template v-slot:prepend-item>
            <v-list-item @click="unselect('priorities')">
              <v-list-item-action>
                <b>Unselect All</b>
              </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>
          </template>
          <template v-slot:selection="{ item }">
            <v-chip dark :color="priorityToColor(item)" small
            >{{ item }}</v-chip>
          </template>
          <template v-slot:item="{ item, attrs, on }">
            <v-list-item v-on="on" v-bind="attrs" #default="{ active }">
              <v-list-item-action>
                <v-checkbox :input-value="active"></v-checkbox>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>
                  <v-chip dark :color="priorityToColor(item)" small
                  >{{ item }}</v-chip>
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-select>
      </v-col>
      <v-col
      class="d-flex ml-4"
      cols="12"
      sm="3">
        <v-select
        v-model="filters.rule"
        chips
        :items="rules"
        label="Rules"
        >
          <template v-slot:prepend-item>
            <v-list-item @click="unselect('rule')">
              <v-list-item-action>
                <b>Unselect All</b>
              </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>
          </template>
          <template v-slot:selection="{ item }">
            {{ item }}
          </template>
          <template v-slot:item="{ item, attrs, on }">
            <v-list-item v-on="on" v-bind="attrs">
              <v-list-item-content>
                  {{item}}
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-select>
      </v-col>
      <v-col
      class="d-flex ml-4"
      cols="12"
      sm="2">
        <v-select
        v-model="filters.tags"
        chips
        multiple
        :items="tags"
        label="Tags"
        >
          <template v-slot:prepend-item>
            <v-list-item @click="unselect('tags')">
              <v-list-item-action>
                <b>Unselect All</b>
              </v-list-item-action>
            </v-list-item>
            <v-divider></v-divider>
          </template>
          <template v-slot:selection="{ item }">
            <v-chip dark :color="stringToColor(item)" small
            >{{ item }}</v-chip>
          </template>
          <template v-slot:item="{ item, attrs, on }">
            <v-list-item v-on="on" v-bind="attrs" #default="{ active }">
              <v-list-item-action>
                <v-checkbox :input-value="active"></v-checkbox>
              </v-list-item-action>
              <v-list-item-content>
                <v-list-item-title>
                  <v-chip dark :color="stringToColor(item)" small
                  >{{ item }}</v-chip>
                </v-list-item-title>
              </v-list-item-content>
            </v-list-item>
          </template>
        </v-select>
      </v-col>
      <v-spacer>
      </v-spacer>
      <v-col
      class="d-flex ml-4 pr-10"
      cols="12"
      sm="2">
        <v-select
        v-model="filters.since"
        chips
        :items="since"
        label="Since"
        >
        </v-select>
      </v-col>
    </v-row>
    <v-row>
      <v-col
      class="d-flex ml-4 pr-10"
      cols="12"
      sm="12">
        <v-text-field
          v-model="filters.search"
          label="Search"
          single-line
          hide-details
          append-icon="search"
        ></v-text-field>
        <v-spacer></v-spacer>
        <Counters
          :filters="filters"
          :ticer="ticer"
          @add-item-to-filters="addItemToFilters('priorities', ...arguments)"
        ></Counters>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { requests } from '../http';
import { utils } from '../utils';
import Counters from '../components/counters.vue';

export default {
  name: 'Filters',
  components: {
    Counters,
  },
  props: {
    addItem: { required: false, type: Object },
  },
  data() {
    return {
      priorities: [],
      rules: [],
      sources: [],
      tags: [],
      since: ['5min', '15min', '30min', '1h', '2h', '5h', '12h', '24h', '48h', '1w', '2w', '1M', '3M', '6M', '1y', '2y'],
      filters: {
        priorities: [],
        rule: '',
        tags: [],
        sources: [],
        search: '',
        since: '24h',
      },
      debounce: null,
    };
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
  },
  watch: {
    addItem: {
      handler() {
        this.addItemToFilters(this.addItem.list, this.addItem.item);
      },
      deep: true,
    },
    filters: {
      handler() {
        if (this.$route.query.since !== this.filters.since || this.$route.query.since === '') {
          this.listItems('source');
          this.listItems('priority');
          this.listItems('rule');
          this.listItems('tags');
          this.$router.push({ query: { ...this.$route.query, since: this.filters.since } });
        }
        if (this.filters.priorities !== []) {
          this.$router.push({ query: { ...this.$route.query, priority: this.filters.priorities } });
        }
        if (this.filters.sources !== []) {
          this.$router.push({ query: { ...this.$route.query, source: this.filters.sources } });
        }
        if (this.filters.rule !== '') {
          this.$router.push({ query: { ...this.$route.query, rule: this.filters.rule } });
        }
        if (this.filters.tags !== '') {
          this.$router.push({ query: { ...this.$route.query, tags: this.filters.tags } });
        }
        this.$emit('send-filters', this.filters);
      },
      deep: true,
    },
    ticer: {
      handler() {
        this.listItems('priority');
        this.listItems('rule');
        this.listItems('source');
        this.listItems('tags');
      },
    },
  },
  methods: {
    async listItems(list) {
      await requests.countByEvents(list, '', '', '', '', '', this.filters.since)
        .then((response) => {
          let values = {};
          values = response.data.results;
          switch (list) {
            case 'source':
            case 'sources':
              Object.keys(values).forEach((value) => {
                this.sources.push(value);
              });
              this.sources = [...new Set(this.sources)];
              break;
            case 'priority':
            case 'priorities':
              Object.keys(values).forEach((value) => {
                this.priorities.push(value);
              });
              this.priorities = [...new Set(this.priorities)];
              break;
            case 'rule':
              Object.keys(values).forEach((value) => {
                this.rules.push(value);
              });
              this.rules = [...new Set(this.rules)];
              break;
            case 'tags':
              Object.keys(values).forEach((value) => {
                this.tags.push(value);
              });
              this.tags = [...new Set(this.tags)];
              break;
            default:
              break;
          }
        });
    },
    priorityToColor(prio) {
      return utils.priorityToColor(prio);
    },
    stringToColor(str) {
      return utils.stringToColor(str);
    },
    unselect(list) {
      this.$nextTick(() => {
        switch (list) {
          case 'priority':
          case 'priorities':
            this.filters.priorities = [];
            break;
          case 'rule':
            this.filters.rule = '';
            break;
          case 'source':
          case 'sources':
            this.filters.sources = [];
            break;
          case 'tags':
            this.filters.tags = [];
            break;
          default:
            break;
        }
      });
    },
    addItemToFilters(list, item) {
      switch (list) {
        case 'priority':
        case 'priorities':
          this.filters.priorities.push(item);
          this.filters.priorities = [...new Set(this.filters.priorities)];
          break;
        case 'rule':
          this.filters.rule = item;
          break;
        case 'source':
        case 'sources':
          this.filters.sources.push(item);
          this.filters.sources = [...new Set(this.filters.sources)];
          break;
        case 'search':
        case 'filter':
          this.filters.search = item;
          break;
        case 'tags':
          this.filters.tags.push(item);
          this.filters.tags = [...new Set(this.filters.tags)];
          break;
        default:
          break;
      }
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
    this.listItems('priority');
    this.listItems('rule');
    this.listItems('source');
    this.listItems('tags');
  },
};
</script>
