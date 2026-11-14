<template>
  <aside class="sidebar">
    <div class="brand">
      <img
        src="https://github.com/falcosecurity/falcosidekick/raw/master/imgs/falcosidekick_color.png"
        alt="" class="brand-logo">
      <span class="brand-name">Falcosidekick</span>
    </div>

    <nav class="nav">
      <router-link
        v-for="page in pages" :key="page.href"
        :to="{ path: page.href, query: $route.query }"
        class="nav-item"
        active-class="nav-item--active">
        <v-icon small>{{ page.icon }}</v-icon>
        <span>{{ page.title }}</span>
      </router-link>
    </nav>

    <div class="summary">
      <div class="summary-label">Events · last {{ filters.since }}</div>
      <div class="summary-total">{{ formatNumber(total) }}</div>
      <div class="severity-bar" v-if="total > 0">
        <span
          v-for="p in priorityItems" :key="p.name"
          :style="{ flexGrow: p.count, background: priorityToColor(p.name) }"
          :title="`${p.name}: ${p.count}`"></span>
      </div>
    </div>

    <div class="search">
      <v-icon small class="search-icon">mdi-magnify</v-icon>
      <input
        v-model="search"
        type="text"
        placeholder="Search events"
        spellcheck="false">
      <v-icon v-if="search" small class="search-clear" @click="search = ''">mdi-close</v-icon>
    </div>

    <div class="groups">
      <section class="group">
        <h4>Filter by priority</h4>
        <ul>
          <li
            v-for="p in priorityItems" :key="p.name"
            :class="{ active: filters.priorities.includes(p.name) }"
            @click="toggle('priorities', p.name)">
            <span class="dot" :style="{ background: priorityToColor(p.name) }"></span>
            <span class="label">{{ p.name }}</span>
            <span class="badge">{{ formatNumber(p.count) }}</span>
          </li>
        </ul>
      </section>

      <section class="group" v-for="g in groups" :key="g.key">
        <h4>Filter by {{ g.title }}</h4>
        <ul>
          <li
            v-for="item in visible(g)" :key="item.name"
            :class="{ active: filters[g.key].includes(item.name) }"
            @click="toggle(g.key, item.name)">
            <span class="dot dot--ring" :style="{ borderColor: stringToColor(item.name) }"></span>
            <span class="label" :title="item.name">{{ item.name }}</span>
            <span class="badge">{{ formatNumber(item.count) }}</span>
          </li>
        </ul>
        <button
          v-if="lists[g.key].length > g.limit"
          class="more"
          @click="$set(expanded, g.key, !expanded[g.key])">
          {{ expanded[g.key] ? 'Show less' : `Show all ${lists[g.key].length}` }}
        </button>
      </section>
    </div>

    <footer class="sidebar-footer">
      <div class="footer-row">
        <label>Refresh</label>
        <select v-model="refresh">
          <option v-for="r in refreshIntervals" :key="r" :value="r">{{ r }}</option>
        </select>
      </div>
      <div class="footer-row">
        <label>Theme</label>
        <div class="themes">
          <button
            v-for="t in themes" :key="t.id"
            :class="['theme-swatch', { active: theme === t.id }]"
            :title="t.name"
            :style="{ background: t.swatch }"
            @click="setTheme(t.id)"></button>
        </div>
      </div>
      <div class="footer-row user" v-if="$store.state.username">
        <span><v-icon x-small>mdi-account-circle-outline</v-icon>
          {{ $store.state.username }}</span>
        <button class="link" @click="$emit('logout')">Log out</button>
      </div>
      <router-link
        :to="{ path: '/info', query: $route.query }"
        class="about"
        active-class="about--active"
        title="Version, configuration and outputs">
        <v-icon x-small>mdi-information-outline</v-icon>
        <span>About</span>
        <span class="version mono" v-if="version">{{ version }}</span>
      </router-link>
    </footer>
  </aside>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import numeral from 'numeral';
import { requests } from '../http';
import { utils } from '../utils';

export const THEMES = [
  { id: 'falco', name: 'Falco', swatch: 'linear-gradient(135deg, #08303c 50%, #00a3bb 50%)' },
  { id: 'flood', name: 'Flood', swatch: 'linear-gradient(135deg, #1f2a38 50%, #ffffff 50%)' },
  { id: 'midnight', name: 'Midnight', swatch: 'linear-gradient(135deg, #141b24 50%, #202a37 50%)' },
];

export default {
  name: 'Sidebar',
  props: {
    pages: { type: Array, required: true },
    value: { type: String, default: '10s' },
  },
  data() {
    return {
      lists: {
        priorities: [],
        sources: [],
        hostnames: [],
        tags: [],
      },
      groups: [
        {
          key: 'sources', group: 'source', title: 'source', limit: 6,
        },
        {
          key: 'hostnames', group: 'hostname', title: 'hostname', limit: 6,
        },
        {
          key: 'tags', group: 'tags', title: 'tag', limit: 8,
        },
      ],
      expanded: {},
      total: 0,
      search: '',
      debounce: null,
      themes: THEMES,
      version: '',
    };
  },
  computed: {
    ...mapState(['filters', 'ticer', 'theme']),
    refreshIntervals() {
      return this.$store.state.refreshIntervals;
    },
    refresh: {
      get() { return this.value; },
      set(v) { this.$emit('input', v); },
    },
    priorityItems() {
      const sorted = utils.sortPriorities(this.lists.priorities.map(p => p.name));
      return sorted.map(name => this.lists.priorities.find(p => p.name === name));
    },
  },
  watch: {
    'filters.since': function since() { this.load(); },
    'filters.search': {
      handler(v) { if (v !== this.search) this.search = v; },
      immediate: true,
    },
    search(v) {
      clearTimeout(this.debounce);
      this.debounce = setTimeout(() => {
        if (v !== this.filters.search) this.setFilters({ search: v });
      }, 350);
    },
    ticer() { this.load(); },
  },
  methods: {
    ...mapActions(['toggleFilter', 'setFilters', 'setTheme']),
    toggle(list, item) {
      this.toggleFilter({ list, item });
    },
    visible(g) {
      const items = this.lists[g.key];
      return this.expanded[g.key] ? items : items.slice(0, g.limit);
    },
    formatNumber(n) {
      return n >= 10000 ? numeral(n).format('0.0a') : numeral(n).format('0,0');
    },
    priorityToColor(p) { return utils.priorityToColor(p); },
    stringToColor(s) { return utils.stringToColor(s); },
    fetch(group) {
      return requests.countByEvents(group, '', '', '', '', '', '', this.filters.since)
        .then(r => Object.entries(r.data.results || {})
          .map(([name, count]) => ({ name, count }))
          .sort((a, b) => b.count - a.count));
    },
    load() {
      this.fetch('priority').then((l) => {
        this.lists.priorities = l;
        this.total = l.reduce((s, p) => s + p.count, 0);
      });
      this.groups.forEach((g) => {
        this.fetch(g.group).then((l) => { this.lists[g.key] = l; });
      });
    },
  },
  mounted() {
    this.load();
    requests.getVersion()
      .then((r) => { this.version = (r.data && r.data.GitVersion) || ''; })
      .catch(() => {});
  },
};
</script>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: 248px;
  display: flex;
  flex-direction: column;
  background: var(--sidebar-bg);
  color: var(--sidebar-text);
  font-size: 13px;
  z-index: 5;
}

.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 18px 18px 14px;
}

.brand-logo {
  width: 26px;
  height: 26px;
}

.brand-name {
  color: var(--sidebar-text-strong);
  font-weight: 600;
  font-size: 15px;
  letter-spacing: 0.01em;
}

.nav {
  display: flex;
  gap: 4px;
  padding: 0 12px 14px;
  border-bottom: 1px solid var(--sidebar-border);
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 3px;
  padding: 7px 4px 6px;
  border-radius: 6px;
  color: var(--sidebar-text);
  text-decoration: none;
  font-size: 11px;
  font-weight: 500;
  transition: background 0.15s, color 0.15s;
}

.nav-item .v-icon {
  color: inherit !important;
  font-size: 18px !important;
}

.nav-item:hover {
  background: var(--sidebar-bg-2);
  color: var(--sidebar-text-strong);
}

.nav-item--active {
  background: var(--sidebar-bg-2);
  color: var(--sidebar-text-strong);
  box-shadow: inset 0 -2px 0 var(--accent);
}

.summary {
  padding: 16px 18px 14px;
}

.summary-label {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--sidebar-heading);
  font-weight: 600;
}

.summary-total {
  color: var(--sidebar-text-strong);
  font-size: 30px;
  font-weight: 300;
  line-height: 1.2;
  margin: 2px 0 10px;
  font-variant-numeric: tabular-nums;
}

.severity-bar {
  display: flex;
  height: 4px;
  border-radius: 2px;
  overflow: hidden;
  gap: 2px;
}

.severity-bar span {
  min-width: 3px;
}

.search {
  position: relative;
  margin: 4px 12px 8px;
}

.search input {
  width: 100%;
  background: var(--sidebar-bg-2);
  border: 1px solid transparent;
  border-radius: 6px;
  padding: 7px 28px 7px 30px;
  color: var(--sidebar-text-strong);
  font-size: 13px;
  outline: none;
  transition: border-color 0.15s;
}

.search input::placeholder {
  color: var(--sidebar-heading);
}

.search input:focus {
  border-color: var(--accent);
}

.search-icon,
.search-clear {
  position: absolute !important;
  top: 50%;
  transform: translateY(-50%);
  color: var(--sidebar-heading) !important;
}

.search-icon {
  left: 9px;
}

.search-clear {
  right: 8px;
  cursor: pointer;
}

.groups {
  flex: 1;
  overflow-y: auto;
  padding-bottom: 12px;
  scrollbar-color: var(--sidebar-border) transparent;
}

.group {
  padding: 12px 0 2px;
}

.group h4 {
  padding: 0 18px 6px;
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--sidebar-heading);
  font-weight: 600;
}

.group ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.group li {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 4px 18px;
  cursor: pointer;
  user-select: none;
  border-left: 2px solid transparent;
  transition: background 0.1s, color 0.1s;
}

.group li:hover {
  color: var(--sidebar-text-strong);
  background: rgba(255, 255, 255, 0.03);
}

.group li.active {
  color: var(--sidebar-text-strong);
  background: var(--sidebar-bg-2);
  border-left-color: var(--accent);
}

.group li.active .badge {
  background: var(--accent);
  color: #fff;
}

.dot {
  flex: none;
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dot--ring {
  background: transparent;
  border: 2px solid;
}

.label {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.badge {
  flex: none;
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 8px;
  background: var(--sidebar-badge);
  color: var(--sidebar-text);
  font-variant-numeric: tabular-nums;
}

.more {
  margin: 2px 18px 0 29px;
  font-size: 11px;
  color: var(--sidebar-heading);
}

.more:hover {
  color: var(--accent);
}

.sidebar-footer {
  border-top: 1px solid var(--sidebar-border);
  padding: 10px 18px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
}

.footer-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.footer-row label {
  color: var(--sidebar-heading);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-weight: 600;
}

.footer-row select {
  background: var(--sidebar-bg-2);
  color: var(--sidebar-text-strong);
  border-radius: 4px;
  padding: 2px 8px;
  font-size: 12px;
  outline: none;
  appearance: auto;
}

.themes {
  display: flex;
  gap: 6px;
}

.theme-swatch {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid var(--sidebar-border);
  transition: transform 0.1s;
}

.theme-swatch:hover {
  transform: scale(1.1);
}

.theme-swatch.active {
  border-color: var(--accent);
}

.user .v-icon {
  color: inherit !important;
}

.about {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 2px -6px 0;
  padding: 3px 6px;
  border-radius: 4px;
  font-size: 11px;
  color: var(--sidebar-heading);
  text-decoration: none;
}

.about .v-icon {
  color: inherit !important;
}

.about:hover,
.about--active {
  color: var(--sidebar-text-strong);
  background: var(--sidebar-bg-2);
}

.about .version {
  margin-left: auto;
  font-size: 10px;
  opacity: 0.8;
}

.link {
  color: var(--sidebar-heading);
  font-size: 11px;
}

.link:hover {
  color: var(--accent);
}
</style>
