<script setup lang="ts">
import type { Connection, Cron } from '@/types'
import { useLink } from '@/composables'
import { computed } from 'vue'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [
    Layout,
    HomeLayout,
  ],
})

const props = defineProps<{
  crons?: Cron[]
  connectionId?: number
  connection?: Connection
}>()

const createLink = computed(() => {
  return props.connection ? `/connections/${props.connectionId}/crons/create` : '/crons/create'
})
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <v-card v-for="cron in props.crons" :key="cron.CronId" v-bind="useLink(`/crons/${cron.CronId}`)">
      <v-card-title>
        {{ cron.Name }} <v-chip :text="cron.Schedule" size="small" />
      </v-card-title>
    </v-card>

    <v-card v-bind="useLink(createLink)">
      <v-card-title class="text-center">
        Add new cron
      </v-card-title>
    </v-card>
  </div>
</template>
