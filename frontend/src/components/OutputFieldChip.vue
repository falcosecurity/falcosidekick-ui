<template>
    <v-chip dark
            label
            small
            class="my-1 mr-2 pr-0 py-0"
            color="primary"
            style="border: 1px solid #000"
            v-on="$listeners"
    >
    {{ label }}
    <div class="white black--text rounded ml-2 px-2 py-2 text-body-2"
            style="font-size: 0.7rem!important; line-height: 0.9rem!important;"
    >
        <span v-html="value" />
    </div>
    </v-chip>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  name: 'OutputFieldChip',
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true }
  }
})
</script>
