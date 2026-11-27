<template>
  <header class="toolbar">
    <div class="toolbar-main">
      <h1>{{ title }}</h1>
      <div class="toolbar-controls" v-if="filterable">
        <v-autocomplete
          class="rule-select"
          :value="filters.rule || null"
          :items="rules"
          placeholder="All rules"
          prepend-inner-icon="mdi-filter-variant"
          outlined dense hide-details clearable
          :menu-props="{ maxHeight: 360, offsetY: true }"
          @change="setFilters({ rule: $event || '' })"
        ></v-autocomplete>
        <v-select
          class="since-select"
          :value="filters.since"
          :items="since"
          prepend-inner-icon="mdi-clock-outline"
          outlined dense hide-details
          :menu-props="{ offsetY: true }"
          @change="setFilters({ since: $event })"
        ></v-select>
        <slot name="actions"></slot>
      </div>
      <slot v-else name="actions"></slot>
    </div>
    <div class="active-filters" v-if="filterable && active.length">
      <span
        v-for="f in active" :key="`${f.list}:${f.value}`"
        class="pill">
        <span class="pill-key">{{ f.label }}</span>
        <span class="pill-dot" v-if="f.color" :style="{ background: f.color }"></span>
        <span class="pill-value">{{ f.value }}</span>
        <v-icon x-small @click="remove(f)">mdi-close</v-icon>
      </span>
      <button class="clear" @click="clearFilters">Clear all</button>
    </div>
  </header>
</template>

<script>
import { mapActions, mapState } from 'vuex';
import { requests } from '../http';
import { utils } from '../utils';

export default {
  name: 'PageToolbar',
  props: {
    title: { type: String, required: true },
    filterable: { type: Boolean, default: true },
  },
  data() {
    return {
      rules: [],
      since: ['5min', '15min', '30min', '1h', '2h', '5h', '12h', '24h', '48h', '1w', '2w', '1M', '3M', '6M', '1y', '2y'],
    };
  },
  computed: {
    ...mapState(['filters', 'ticer']),
    active() {
      const f = this.filters;
      const out = [];
      f.priorities.forEach(v => out.push({
        list: 'priorities', label: 'priority', value: v, color: utils.priorityToColor(v),
      }));
      f.sources.forEach(v => out.push({ list: 'sources', label: 'source', value: v }));
      f.hostnames.forEach(v => out.push({ list: 'hostnames', label: 'host', value: v }));
      f.tags.forEach(v => out.push({ list: 'tags', label: 'tag', value: v }));
      if (f.rule) out.push({ list: 'rule', label: 'rule', value: f.rule });
      if (f.search) out.push({ list: 'search', label: 'search', value: f.search });
      return out;
    },
  },
  watch: {
    'filters.since': function since() { this.loadRules(); },
    ticer() { this.loadRules(); },
  },
  methods: {
    ...mapActions(['setFilters', 'toggleFilter', 'clearFilters']),
    remove(f) {
      if (f.list === 'rule') this.setFilters({ rule: '' });
      else if (f.list === 'search') this.setFilters({ search: '' });
      else this.toggleFilter({ list: f.list, item: f.value });
    },
    loadRules() {
      if (!this.filterable) return;
      requests.countByEvents('rule', '', '', '', '', '', '', this.filters.since)
        .then((r) => {
          this.rules = Object.keys(r.data.results || {}).sort();
        });
    },
  },
  mounted() {
    this.loadRules();
  },
};
</script>

<style scoped>
.toolbar {
  position: sticky;
  top: 0;
  z-index: 4;
  background: var(--bg);
  border-bottom: 1px solid var(--border);
  padding: 12px 24px;
}

.toolbar-main {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 34px;
}

h1 {
  font-size: 17px;
  font-weight: 600;
  color: var(--text);
  margin-right: auto;
  letter-spacing: -0.01em;
}

.toolbar-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rule-select {
  width: 320px;
}

.since-select {
  width: 116px;
  flex: none;
}

.active-filters {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin-top: 10px;
}

.pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  padding: 3px 6px 3px 9px;
  border-radius: 12px;
  background: var(--accent-soft);
  color: var(--text);
  max-width: 420px;
}

.pill-key {
  color: var(--text-muted);
}

.pill-value {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pill-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex: none;
}

.pill .v-icon {
  color: var(--text-muted) !important;
  cursor: pointer;
}

.pill .v-icon:hover {
  color: var(--text) !important;
}

.clear {
  font-size: 12px;
  color: var(--text-muted);
  margin-left: 4px;
}

.clear:hover {
  color: var(--accent);
}
</style>
