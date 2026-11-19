<template>
  <div>
    <PageToolbar title="Events">
      <template v-slot:actions>
        <div class="segmented" role="group" aria-label="Row density">
          <button :class="{ active: !dense }" title="Comfortable" @click="setDense(false)">
            <v-icon small>mdi-view-agenda-outline</v-icon>
          </button>
          <button :class="{ active: dense }" title="Compact" @click="setDense(true)">
            <v-icon small>mdi-view-headline</v-icon>
          </button>
        </div>
        <button class="btn" @click="saveFile()">
          <v-icon small>mdi-download-outline</v-icon> Export
        </button>
      </template>
    </PageToolbar>
    <div class="page">
      <section class="panel panel--flush">
        <v-data-table
          :class="['events-table', { 'events-table--dense': dense }]"
          :page='page'
          :headers='headers'
          :items='events'
          item-key="uuid"
          :options.sync='options'
          :server-items-length='totalEvents'
          :items-per-page="itemsPerPage"
          :loading='loading'
          loader-height="2"
          :footer-props="{
            'items-per-page-options': itemsPerPageInterval,
          }"
        >
          <template v-slot:item="{item}">
            <tr @click="showDialog(item)">
              <td class="col-time mono">
                <div>{{ item.time | formatTime }}</div>
                <div class="sub">{{ item.time | formatDay }}</div>
              </td>
              <td class="col-priority">
                <span class="prio" @click.stop="addToFilters('priorities', item.priority)">
                  <span class="dot" :style="{ background: priorityToColor(item.priority) }"></span>
                  {{ item.priority }}
                </span>
              </td>
              <td class="col-rule">
                <div class="rule" @click.stop="addToFilters('rule', item.rule)">
                  {{ item.rule }}</div>
                <div class="output">{{ item.output }}</div>
                <div class="fields">
                  <span
                    v-for="(value, key) in item.output_fields" :key="key"
                    class="field"
                    @click.stop="addToFilters('search', value)">
                    <span class="field-key">{{ key }}</span>
                    <span class="field-value">{{ value }}</span>
                  </span>
                </div>
              </td>
              <td class="col-origin">
                <div class="link-ish" @click.stop="addToFilters('sources', item.source)">
                  {{ item.source }}</div>
                <div class="sub link-ish" v-if="item.hostname"
                  @click.stop="addToFilters('hostnames', item.hostname)">
                  {{ item.hostname }}</div>
              </td>
              <td class="col-tags">
                <span
                  v-for="(tag, index) in item.tags" :key="index"
                  class="tag"
                  @click.stop="addToFilters('tags', tag)">
                  <span class="tag-dot" :style="{ background: stringToColor(tag) }"></span>{{ tag }}
                </span>
              </td>
              <td class="col-action">
                <v-icon small>mdi-chevron-right</v-icon>
              </td>
            </tr>
          </template>
        </v-data-table>
      </section>
    </div>

    <v-dialog v-model="dialog" max-width="880" content-class="event-dialog">
      <div class="modal" v-if="json">
        <header class="modal-header">
          <div class="modal-title">{{ json.rule }}</div>
          <div class="modal-meta">
            <span class="prio">
              <span class="dot"
                :style="{ background: priorityToColor(String(json.priority)) }"></span>
              {{ json.priority }}
            </span>
            <span><v-icon x-small>mdi-clock-outline</v-icon> {{ json.time | formatDate }}</span>
            <span><v-icon x-small>mdi-server</v-icon> {{ json.hostname || '—' }}</span>
            <span><v-icon x-small>mdi-source-branch</v-icon> {{ json.source }}</span>
          </div>
          <v-icon class="modal-close" @click="dialog = false">mdi-close</v-icon>
        </header>
        <div class="modal-body">
          <nav class="modal-nav">
            <a :class="{ active: tab === 'details' }" @click="tab = 'details'">Details</a>
            <a :class="{ active: tab === 'fields' }" @click="tab = 'fields'">Fields</a>
            <a :class="{ active: tab === 'json' }" @click="tab = 'json'">JSON</a>
          </nav>
          <div class="modal-content">
            <div v-show="tab === 'details'">
              <h5>Output</h5>
              <pre class="output-block mono">{{ json.output }}</pre>
              <h5>Event</h5>
              <dl>
                <dt>Time</dt><dd class="mono">{{ json.time }}</dd>
                <dt>Source</dt>
                <dd><a @click="addToFilters('sources', json.source)">{{ json.source }}</a></dd>
                <dt>Hostname</dt>
                <dd>
                  <a @click="addToFilters('hostnames', json.hostname)">
                    {{ json.hostname || '—' }}</a>
                </dd>
                <dt>Priority</dt>
                <dd>
                  <a @click="addToFilters('priorities', json.priority)">{{ json.priority }}</a>
                </dd>
                <dt>Rule</dt>
                <dd><a @click="addToFilters('rule', json.rule)">{{ json.rule }}</a></dd>
                <dt>Tags</dt>
                <dd>
                  <span
                    v-for="(tag, index) in json.tags" :key="index"
                    class="tag" @click="addToFilters('tags', tag)">
                    <span class="tag-dot"
                      :style="{ background: stringToColor(String(tag)) }"></span>{{ tag }}
                  </span>
                </dd>
              </dl>
            </div>
            <div v-show="tab === 'fields'">
              <h5>Output fields</h5>
              <dl>
                <template v-for="(value, key) in json.output_fields">
                  <dt :key="`k-${key}`" class="mono">{{ key }}</dt>
                  <dd :key="`v-${key}`" class="mono">
                    <a @click="addToFilters('search', value)">{{ value }}</a>
                  </dd>
                </template>
              </dl>
            </div>
            <div v-show="tab === 'json'" class="json-view">
              <vue-json-pretty :data="json" :collapsedOnClickBrackets="false" />
            </div>
          </div>
        </div>
        <footer class="modal-footer">
          <button class="btn" @click="copyToClipBoard()">
            <v-icon small>{{ copied ? 'mdi-check' : 'mdi-content-copy' }}</v-icon>
            {{ copied ? 'Copied' : 'Copy JSON' }}
          </button>
          <button class="btn btn--primary" @click="dialog = false">Close</button>
        </footer>
      </div>
    </v-dialog>
  </div>
</template>

<script>
import 'vue-json-pretty/lib/styles.css';
import VueJsonPretty from 'vue-json-pretty';
import moment from 'moment';
import { mapActions, mapState } from 'vuex';
import { requests } from '../http';
import { utils } from '../utils';
import PageToolbar from '../components/PageToolbar.vue';

export default {
  name: 'DatatableComponent',
  components: {
    PageToolbar,
    VueJsonPretty,
  },
  filters: {
    formatTime(v) {
      return v ? moment(String(v)).format('HH:mm:ss.SSS') : '';
    },
    formatDay(v) {
      return v ? moment(String(v)).format('YYYY-MM-DD') : '';
    },
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
      loading: true,
      options: {
        page: 1,
        itemsPerPage: 10,
      },
      headers: [
        { text: 'Time', value: 'time', sortable: false },
        { text: 'Priority', value: 'priority', sortable: false },
        { text: 'Rule / output', value: 'rule', sortable: false },
        { text: 'Source / host', value: 'source', sortable: false },
        { text: 'Tags', value: 'tags', sortable: false },
        { text: '', value: 'json', sortable: false },
      ],
      debounce: null,
      dialog: false,
      tab: 'details',
      dense: (() => {
        try { return window.localStorage.getItem('fsui-dense') === '1'; } catch (e) { return false; }
      })(),
      json: '',
      copied: false,
    };
  },
  computed: {
    ...mapState(['filters', 'ticer']),
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
    ...mapActions(['addFilter']),
    addToFilters(list, item) {
      this.addFilter({ list, item });
    },
    setDense(v) {
      this.dense = v;
      try { window.localStorage.setItem('fsui-dense', v ? '1' : '0'); } catch (e) { /* ignore */ }
    },
    showDialog(i) {
      this.json = i;
      this.tab = 'details';
      this.copied = false;
      this.dialog = true;
    },
    copyToClipBoard() {
      navigator.clipboard.writeText(JSON.stringify(this.json, null, 2));
      this.copied = true;
      setTimeout(() => { this.copied = false; }, 1500);
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
.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 32px;
  padding: 0 12px;
  border: 1px solid var(--border);
  border-radius: 6px;
  background: var(--surface);
  color: var(--text);
  font-size: 13px;
  font-weight: 500;
  transition: border-color 0.15s, color 0.15s;
}

.btn .v-icon {
  color: inherit !important;
}

.btn:hover {
  border-color: var(--accent);
  color: var(--accent);
}

.btn--primary {
  background: var(--accent);
  border-color: var(--accent);
  color: #fff;
}

.btn--primary:hover {
  color: #fff;
  filter: brightness(1.08);
}

.panel--flush {
  padding: 0 !important;
  overflow: hidden;
}

/* ---------- events table ---------- */
.events-table.v-data-table {
  background: var(--surface) !important;
  color: var(--text) !important;
}

.events-table th {
  font-size: 10px !important;
  font-weight: 600 !important;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted) !important;
  background: var(--surface-2) !important;
  border-bottom: 1px solid var(--border) !important;
  height: 36px !important;
  white-space: nowrap;
}

.events-table tbody tr {
  cursor: pointer;
}

.events-table tbody tr:hover {
  background: var(--row-hover) !important;
}

.events-table td {
  vertical-align: top;
  padding-top: 10px !important;
  padding-bottom: 10px !important;
  height: auto !important;
  font-size: 13px !important;
  border-bottom: 1px solid var(--border) !important;
}

.events-table .sub {
  color: var(--text-muted);
  font-size: 11px;
  margin-top: 2px;
}

.events-table .col-time {
  width: 120px;
  white-space: nowrap;
  font-size: 12px !important;
}

.events-table .col-priority {
  width: 130px;
  white-space: nowrap;
}

.prio {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-weight: 500;
}

.prio .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex: none;
}

.events-table .rule {
  font-weight: 600;
  color: var(--text);
  display: inline-block;
}

.events-table .rule:hover,
.events-table .link-ish:hover {
  color: var(--accent);
}

.events-table .output {
  color: var(--text-muted);
  font-size: 12px;
  margin-top: 3px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-word;
}

.fields {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
}

.field {
  display: inline-flex;
  font-size: 11px;
  border: 1px solid var(--border);
  border-radius: 4px;
  overflow: hidden;
  max-width: 100%;
}

.field:hover {
  border-color: var(--accent);
}

.field-key {
  padding: 1px 6px;
  background: var(--surface-2);
  color: var(--text-muted);
  border-right: 1px solid var(--border);
}

.field-value {
  padding: 1px 6px;
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.events-table .col-origin {
  width: 160px;
  white-space: nowrap;
}

.events-table .col-tags {
  width: 220px;
}

.tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  padding: 1px 7px;
  margin: 0 4px 4px 0;
  border-radius: 10px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  color: var(--text);
  cursor: pointer;
  white-space: nowrap;
}

.tag:hover {
  border-color: var(--accent);
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.events-table .col-action {
  width: 32px;
  vertical-align: middle !important;
}

.events-table .col-action .v-icon {
  color: var(--text-faint) !important;
}

.events-table .v-data-footer {
  border-top: none !important;
  font-size: 12px;
  color: var(--text-muted);
}

.segmented {
  display: inline-flex;
  height: 32px;
  border: 1px solid var(--border);
  border-radius: 6px;
  overflow: hidden;
  background: var(--surface);
}

.segmented button {
  width: 34px;
  color: var(--text-muted);
}

.segmented button + button {
  border-left: 1px solid var(--border);
}

.segmented button .v-icon {
  color: inherit !important;
}

.segmented button.active {
  background: var(--accent-soft);
  color: var(--accent);
}

.events-table--dense td {
  padding-top: 6px !important;
  padding-bottom: 6px !important;
  vertical-align: middle;
}

.events-table--dense .col-time .sub,
.events-table--dense .col-origin .sub,
.events-table--dense .fields {
  display: none;
}

.events-table--dense .col-rule {
  max-width: 0;
}

.events-table--dense .rule {
  display: inline;
  margin-right: 10px;
}

.events-table--dense .output {
  display: inline;
  white-space: nowrap;
}

.events-table--dense .col-rule {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.events-table--dense .col-tags {
  white-space: nowrap;
  overflow: hidden;
}

.events-table--dense .tag {
  margin-bottom: 0;
}

/* ---------- event modal (Flood-style) ---------- */
.event-dialog {
  border-radius: 8px !important;
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.35) !important;
}

.modal {
  background: var(--sidebar-bg);
  color: var(--sidebar-text-strong);
  font-size: 13px;
}

.modal-header {
  position: relative;
  padding: 20px 24px 16px;
  border-bottom: 1px solid var(--sidebar-border);
}

.modal-title {
  font-size: 18px;
  font-weight: 500;
  padding-right: 32px;
}

.modal-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 18px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--sidebar-text);
}

.modal-meta .v-icon {
  color: var(--sidebar-heading) !important;
}

.modal-close {
  position: absolute !important;
  top: 18px;
  right: 18px;
  color: var(--sidebar-text) !important;
}

.modal-body {
  display: flex;
  min-height: 360px;
  max-height: 62vh;
}

.modal-nav {
  flex: none;
  width: 140px;
  padding: 16px 0;
  border-right: 1px solid var(--sidebar-border);
  display: flex;
  flex-direction: column;
}

.modal-nav a {
  padding: 6px 20px;
  color: var(--sidebar-text);
  border-right: 2px solid transparent;
  margin-right: -1px;
}

.modal-nav a:hover {
  color: var(--sidebar-text-strong);
}

.modal-nav a.active {
  color: var(--accent);
  border-right-color: var(--accent);
}

.modal-content {
  flex: 1;
  min-width: 0;
  padding: 16px 24px 20px;
  overflow-y: auto;
}

.modal-content h5 {
  font-size: 12px;
  font-weight: 600;
  color: var(--sidebar-text-strong);
  margin: 6px 0 8px;
}

.modal-content h5 + * {
  margin-bottom: 14px;
}

.output-block {
  white-space: pre-wrap;
  word-break: break-word;
  font-size: 12px;
  background: var(--sidebar-bg-2);
  border-radius: 6px;
  padding: 10px 12px;
  color: var(--sidebar-text-strong);
}

.modal-content dl {
  display: grid;
  grid-template-columns: max-content 1fr;
  column-gap: 28px;
  row-gap: 6px;
}

.modal-content dt {
  color: var(--sidebar-text);
}

.modal-content dd {
  color: var(--sidebar-text-strong);
  word-break: break-word;
}

.modal-content dd a:hover {
  color: var(--accent);
  text-decoration: underline;
}

.modal-content .tag {
  background: var(--sidebar-bg-2);
  border-color: var(--sidebar-border);
  color: var(--sidebar-text-strong);
}

.json-view {
  font-size: 12px;
}

.json-view .vjs-tree {
  font-family: 'JetBrains Mono', ui-monospace, Menlo, monospace;
  font-size: 12px;
  color: var(--sidebar-text);
}

.json-view .vjs-key {
  color: #9cc8f5;
}

.json-view .vjs-value-string {
  color: #a5d6a7;
}

.json-view .vjs-value-number,
.json-view .vjs-value-boolean {
  color: #f5c38b;
}

.json-view .vjs-tree-node:hover {
  background: var(--sidebar-bg-2);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 24px;
  border-top: 1px solid var(--sidebar-border);
}

.modal-footer .btn:not(.btn--primary) {
  background: transparent;
  border-color: var(--sidebar-border);
  color: var(--sidebar-text-strong);
}
</style>
