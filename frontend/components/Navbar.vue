<script setup lang="ts">
import { useLink } from '@/composables'

const links = [
  { href: '/', text: 'Home' },
]
</script>

<template>
  <v-app-bar flat>
    <v-container class="mx-auto d-flex align-center justify-center">
      <v-app-bar-title>Grognon</v-app-bar-title>

      <v-btn
        v-for="link in links"
        :key="link.text"
        :text="link.text"
        v-bind="useLink(link.href)"
        variant="text"
      />
    </v-container>
  </v-app-bar>
</template>
