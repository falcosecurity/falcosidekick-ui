<template>
<div>
    <slot v-if="ready" />
</div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend<{ ready: boolean }, {}, {}>({
  data: () => ({
    ready: false
  }),
  created () {
    setTimeout(() => { this.ready = true }, 100)
  }
})
</script>
