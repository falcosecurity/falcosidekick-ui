<template>
  <div class="ranked">
    <ul>
      <li
        v-for="(item, i) in visible" :key="item.name"
        :class="{ active: filters.rule === item.name }"
        @click="$emit('picked-item', list, item.name)">
        <span class="rank">{{ i + 1 }}</span>
        <span class="name" :title="item.name">{{ item.name }}</span>
        <span class="track">
          <span class="fill" :style="{ width: `${(item.count / max) * 100}%` }"></span>
        </span>
        <span class="count">{{ formatNumber(item.count) }}</span>
      </li>
      <li v-if="loaded && !items.length" class="empty">No events</li>
    </ul>
    <button
      v-if="items.length > limit"
      class="more"
      @click="expanded = !expanded">
      {{ expanded ? 'Show top ' + limit : `Show all ${items.length} rules` }}
    </button>
  </div>
</template>

<script>
import numeral from 'numeral';
import { requests } from '../../http';

export default {
  name: 'BarChart',
  props: {
    list: { type: String },
    limit: { type: Number, default: 10 },
    filters: { type: Object, required: true },
  },
  data() {
    return {
      loaded: false,
      items: [],
      expanded: false,
    };
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
    visible() {
      return this.expanded ? this.items : this.items.slice(0, this.limit);
    },
    max() {
      return this.items.length ? this.items[0].count : 1;
    },
  },
  watch: {
    filters: {
      handler() {
        this.updateChart();
      },
      deep: true,
      immediate: true,
    },
    ticer: {
      handler() {
        this.updateChart();
      },
    },
  },
  methods: {
    formatNumber(n) {
      return numeral(n).format('0,0');
    },
    async updateChart() {
      const response = await requests.countByEvents(
        this.list,
        this.filters.sources,
        this.filters.hostnames,
        this.filters.priorities,
        this.filters.rule,
        this.filters.search,
        this.filters.tags,
        this.filters.since,
      );
      this.items = Object.entries(response.data.results || {})
        .sort(([, v1], [, v2]) => v2 - v1)
        .map(([name, count]) => ({ name, count }));
      this.loaded = true;
    },
  },
};
</script>

<style scoped>
ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

li {
  display: grid;
  grid-template-columns: 22px minmax(180px, 2fr) 3fr 56px;
  align-items: center;
  gap: 12px;
  padding: 6px 8px;
  margin: 0 -8px;
  border-radius: 4px;
  font-size: 13px;
  color: var(--text);
  cursor: pointer;
}

li:hover {
  background: var(--row-hover);
}

li.active .name {
  color: var(--accent);
  font-weight: 500;
}

.rank {
  color: var(--text-faint);
  font-size: 11px;
  text-align: right;
  font-variant-numeric: tabular-nums;
}

.name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.track {
  height: 4px;
  border-radius: 2px;
  background: var(--surface-2);
  box-shadow: inset 0 0 0 1px var(--border);
  overflow: hidden;
}

.fill {
  display: block;
  height: 100%;
  background: var(--accent);
  border-radius: 2px;
}

.count {
  text-align: right;
  font-variant-numeric: tabular-nums;
  color: var(--text-muted);
}

li.empty {
  display: block;
  color: var(--text-muted);
  cursor: default;
}

.more {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-muted);
}

.more:hover {
  color: var(--accent);
}
</style>
