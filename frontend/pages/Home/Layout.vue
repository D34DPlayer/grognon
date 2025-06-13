<script setup lang="ts">
import Sidebar from './Sidebar.vue'
</script>

<template>
  <v-row>
    <v-col cols="3">
      <Sidebar />
    </v-col>

    <v-col>
      <slot />
    </v-col>
  </v-row>
</template>
