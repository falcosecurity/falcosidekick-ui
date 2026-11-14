<template>
  <div>
    <PageToolbar title="About" :filterable="false">
      <template v-slot:actions>
        <a class="btn" href="docs" target="_blank" rel="noopener">
          <v-icon small>mdi-api</v-icon> API docs
        </a>
      </template>
    </PageToolbar>
    <div class="page">
      <div class="grid grid-3 info-grid">
        <section class="panel">
          <header><h3>Version</h3></header>
          <dl class="kv">
            <template v-for="(value, key) in version">
              <dt :key="`k-${key}`">{{ key }}</dt>
              <dd :key="`v-${key}`" class="mono">{{ value }}</dd>
            </template>
          </dl>
        </section>
        <section class="panel">
          <header><h3>Configuration</h3></header>
          <dl class="kv">
            <template v-for="(value, key) in configuration">
              <dt :key="`k-${key}`">{{ key }}</dt>
              <dd :key="`v-${key}`" class="mono">{{ value }}</dd>
            </template>
          </dl>
        </section>
        <section class="panel">
          <header>
            <h3>Outputs</h3>
            <span class="count-badge">{{ outputs.length }}</span>
          </header>
          <ul class="outputs">
            <li v-for="value in outputs" :key="value">
              <span class="status"></span>{{ value }}
            </li>
          </ul>
        </section>
      </div>
    </div>
  </div>
</template>

<script>
import { requests } from '../http';
import PageToolbar from '../components/PageToolbar.vue';

export default {
  name: 'Info',
  components: {
    PageToolbar,
  },
  data() {
    return {
      outputs: [],
      version: {},
      configuration: {},
    };
  },
  computed: {
    ticer() {
      return this.$store.state.ticer;
    },
  },
  watch: {
    ticer: {
      handler() {
        this.listOutputs();
      },
    },
  },
  methods: {
    listOutputs() {
      requests.listOutputs()
        .then((response) => {
          this.outputs = response.data;
        });
    },
    getConfiguration() {
      requests.getConfiguration()
        .then((response) => {
          this.configuration = response.data;
        });
    },
    getVersion() {
      requests.getVersion()
        .then((response) => {
          this.version = response.data;
        });
    },
  },
  mounted() {
    this.listOutputs();
    this.getConfiguration();
    this.getVersion();
  },
};
</script>


<style scoped>
.info-grid {
  align-items: start;
}

.kv {
  display: grid;
  grid-template-columns: max-content 1fr;
  column-gap: 24px;
  row-gap: 0;
  font-size: 13px;
}

.kv dt,
.kv dd {
  padding: 7px 0;
  border-bottom: 1px solid var(--border);
}

.kv dt {
  color: var(--text-muted);
}

.kv dd {
  color: var(--text);
  font-size: 12px;
  word-break: break-all;
}

.kv dt:nth-last-of-type(1),
.kv dd:last-of-type {
  border-bottom: none;
}

.count-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: 10px;
  background: var(--accent-soft);
  color: var(--accent);
}

.outputs {
  list-style: none;
  padding: 0;
  margin: 0;
}

.outputs li {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 0;
  font-size: 13px;
  color: var(--text);
  border-bottom: 1px solid var(--border);
}

.outputs li:last-child {
  border-bottom: none;
}

.status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #34c38f;
  box-shadow: 0 0 0 3px rgba(52, 195, 143, 0.15);
}

a.btn {
  text-decoration: none;
}
</style>
