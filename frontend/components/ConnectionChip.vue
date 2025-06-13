<script setup lang="ts">
const props = defineProps<{
  connected: boolean
}>()
</script>

<template>
  <v-chip
    :text="props.connected ? 'Connected' : 'Disconnected'"
    :color="props.connected ? 'success' : 'error'"
    size="small"
  />
</template>
