<template>
    <v-chip dark
            label
            small
            :class="`my-1 mr-2 pr-0 py-0 bg-${displayMode}-color`"
            style="border: 1px solid #000"
            v-on="$listeners"
    >
    {{ label }}
    <div class="text rounded ml-2 px-2 py-2 text-body-2"
            style="font-size: 0.7rem!important; line-height: 0.9rem!important;"
    >
        <span v-html="value" />
    </div>
    </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  name: 'OutputFieldChip',
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true }
  },
  computed: mapState(['displayMode'])
})
</script>

<style scoped>
.bg-light-color {
  background-color: var(--v-primary-base)!important;
  border-color: var(--v-primary-base)!important;
}
.bg-dark-color {
  background-color: var(--v-primary-darken3)!important;
  border-color: var(--v-primary-darken3)!important;
}

.bg-light-color .text {
  background-color: var(--v-white)!important;
  color: var(--v-black)!important;
}
.bg-dark-color .text {
  background-color: var(--v-secondary-base)!important;
  color: var(--v-white)!important;
}
</style>
