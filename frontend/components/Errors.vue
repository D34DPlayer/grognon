<script setup lang="ts">
import { usePage } from '@inertiajs/vue3'
import { computed } from 'vue'

const page = usePage()

const errors = computed(() => {
  return page.props.errors ?? {}
})
const hasErrors = computed(() => {
  return Object.keys(errors.value).length > 0
})
</script>

<template>
  <v-alert v-if="hasErrors" closable type="error" variant="tonal" class="ma-2">
    <ul>
      <li v-for="error, code in errors" :key="code">
        {{ error }}
      </li>
    </ul>
  </v-alert>
</template>
