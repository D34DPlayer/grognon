<script setup lang="ts">
import { Head } from '@inertiajs/vue3'
</script>

<template>
  <Head>
    <title>Grognon - The database sniffer</title>
  </Head>
  <v-app>
    <Navbar />
    <v-main class="bg-grey-lighten-3">
      <Errors />
      <v-container>
        <slot />
      </v-container>
    </v-main>
  </v-app>
</template>
