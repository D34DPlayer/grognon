<script setup lang="ts">
import { router } from '@inertiajs/vue3'
import type { Connection, Cron, CronOutput } from '@/types'
import { getExtensions } from '@/codemirror'
import { useLink } from '@/composables'
import { Codemirror } from 'vue-codemirror'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const props = defineProps<{
  connection?: Connection
  cron?: Cron
  cronOutputs?: CronOutput[]
}>()

const deleteCron = () => {
  if (confirm('Are you sure you want to delete this cron?')) {
    router.delete(`/crons/${props.cron?.CronId}`)
  }
}
</script>

<template>
  <div class="d-flex flex-column ga-2">
    <v-card v-if="props.cron">
      <v-card-title>
        {{ props.cron.Name }} <v-chip
          :text="props.cron.Schedule"
          size="small"
        />
      </v-card-title>
      <v-card-text>
        <v-col>
          <v-row>
            Connection: {{ props.connection?.ConnectionUrl }}
          </v-row>
          <v-row>
            Created at: {{ props.cron.CreatedAt }}
          </v-row>
          <v-row>
            Last run at: {{ props.cron.LastRunAt }}
          </v-row>
        </v-col>
      </v-card-text>

      <v-card-actions>
        <v-btn
          v-bind="useLink(`/connections/${props.cron.ConnectionId}`)"
        >
          View Connection
        </v-btn>
        <v-btn
          @click="deleteCron"
          color="error"
          v-if="props.cron.CronId"
        >
          Delete Cron
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card v-if="props.cron">
      <v-card-title>
        Cron Command
      </v-card-title>
      <v-card-text>
        <Codemirror
          :model-value="props.cron.Command"
          placeholder="Type your SQL command here..."
          :indent-with-tab="true"
          :tab-size="2"
          :extensions="getExtensions()"
          disabled
        />
      </v-card-text>
      <v-card-actions>
        <v-btn
          v-bind="useLink(`/crons/${props.cron.CronId}/data`)"
        >
          View data
        </v-btn>
      </v-card-actions>
    </v-card>

    <v-card v-if="props.cronOutputs">
      <v-card-title>
        Cron Outputs
      </v-card-title>
      <v-card-text>
        <v-table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="output in props.cronOutputs" :key="output.Name">
              <td>{{ output.Name }}</td>
              <td>{{ output.Type }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </div>
</template>
