<template>
    <v-chip dark
            class="my-0 mr-2 pr-0 py-0"
            :color="color(priority)"
            label
            style="border: 1px solid #000"
            v-on="$listeners"
    >
    {{ priority }} <div :class="`bg-${displayMode}-color rounded ml-2 px-2 py-2 text-body-2`">{{ count }}</div>
    </v-chip>
</template>

<script lang="ts">
import { mapPriorityToColor } from '@/api/mapper'
import { Priority } from '@/api/model'
import Vue, { PropType } from 'vue'
import { mapState } from 'vuex'

export default Vue.extend({
  name: 'PriorityCountChip',
  props: {
    priority: { type: String as PropType<Priority>, required: true },
    count: { type: Number, required: true }
  },
  computed: mapState(['displayMode']),
  methods: {
    color (priority: string): string {
      return mapPriorityToColor(priority as Priority)
    }
  }
})
</script>

<style scoped>
.bg-light-color {
  background-color: var(--v-white)!important;
  color: var(--v-black)!important;
}
.bg-dark-color {
  background-color: var(--v-secondary-base)!important;
  color: var(--v-white)!important;
}
</style>
