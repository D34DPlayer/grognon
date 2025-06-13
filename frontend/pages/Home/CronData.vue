<script setup lang="ts">
import type { Cron, CronOutput } from '@/types'
import { displayTime } from '@/utils'
import { useLink } from '@/composables'
import { computed } from 'vue'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const props = defineProps<{
  cron?: Cron
  cronOutputs?: CronOutput[]
  data?: Record<string, any>[]
}>()

const columns = computed(() => {
  const c = ['timestamp']
  if (!props.cronOutputs || props.cronOutputs.length === 0) {
    return c
  }

  for (const output of props.cronOutputs) {
    c.push(output.Name)
  }

  return c
})
</script>

<template>
  <div class="d-flex flex-column ga-2">
    <v-card>
      <v-card-title>
        {{ props.cron?.Name }}
        <v-chip :text="props.cron?.Schedule" size="small" />
      </v-card-title>
      <v-card-actions>
        <v-btn v-bind="useLink(`/crons/${props.cron?.CronId}`)">
          View cron
        </v-btn>
        <v-btn v-bind="useLink(`/connections/${props.cron?.ConnectionId}`)">
          View connection
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card>
      <v-card-title>
        Data
      </v-card-title>
      <v-card-text>
        <v-table>
          <thead>
            <tr>
              <th v-for="(column, index) in columns" :key="index">
                {{ column }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, rowIndex) in props.data" :key="rowIndex">
              <td v-for="(column, colIndex) in columns" :key="colIndex">
                {{ column === "timestamp" ? displayTime(row[column]) : row[column] }}
              </td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </div>
</template>
