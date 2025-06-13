<script setup lang="ts">
import { useLink } from '@/composables'

const links = [
  { name: 'Connections', path: '/connections' },
  { name: 'Crons', path: '/crons' },
]
</script>

<template>
  <v-sheet rounded="lg">
    <v-list rounded="lg">
      <v-list-item
        v-for="link in links"
        :key="link.name"
        :title="link.name"
        v-bind="useLink(link.path)"
      />
    </v-list>
  </v-sheet>
</template>
