<template>
  <v-card class="elevation-0">
    <v-row>
      <v-col
      class="ml-5 mt-5"
      cols="12"
      sm="12">
        <v-card class="elevation-0">
          <v-card-title>Outputs</v-card-title>
            <span
            v-for="(value) in outputs" :key="value">
            <v-chip
            dark
            color="#1976d2"
            class="rounded-0 mr-1">
              {{value}}
            </v-chip>
          </span>
        </v-card>
        <v-card class="elevation-0">
          <v-card-title>Configuration</v-card-title>
            <span
              v-for="(value,key) in configuration" :key="key">
                {{key}}:
              <v-chip
              dark small
              color="#1976d2"
              class="rounded-0"
              style="margin-right: 5px;">
                {{value}}
              </v-chip>
              <br>
          </span>
        </v-card>
        <v-card class="elevation-0">
          <v-card-title>Version</v-card-title>
            <span
              v-for="(value,key) in version" :key="key">
                {{key}}:
              <v-chip
              dark small
              color="#1976d2"
              class="rounded-0"
              style="margin-right: 5px;">
                {{value}}
              </v-chip>
              <br>
          </span>
        </v-card>
        <v-card class="elevation-0">
          <v-card-title>
            API<br>
            <v-btn
            href="docs"
            elevation="0"
            class="ml-2"
            color="#1976d2"
            dark icon>
            <v-icon>mdi-link</v-icon>
            </v-btn>
          </v-card-title>
        </v-card>
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import { requests } from '../http';

export default {
  Name: 'Info',
  data() {
    return {
      outputs: [],
      version: {},
      configuration: {},
    };
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
    this.getVersion();
  },
};
</script>

