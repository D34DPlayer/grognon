<script setup lang="ts">
import { useLink } from '@/composables'
import Layout from '@/pages/Layout.vue'

defineOptions({
  layout: Layout,
})

function goBack() {
  window.history.back()
}
</script>

<template>
  <v-card max-width="max-content" class="mx-auto">
    <v-card-title>
      Page not found
    </v-card-title>
    <v-card-actions>
      <v-btn color="primary" @click="goBack">
        Previous page
      </v-btn>
      <v-btn v-bind="useLink('/')">
        Home page
      </v-btn>
    </v-card-actions>
  </v-card>
</template>
