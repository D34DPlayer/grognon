<script setup lang="ts">
import type { Column, Connection } from '@/types'
import ConnectionChip from '@/components/ConnectionChip.vue'
import { useLink } from '@/composables'
import { computed, defineProps } from 'vue'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const props = defineProps<{
  connectionId: number
  connection: Connection
  columns: Column[] | null
}>()

const colsPerTable = computed(() => {
  if (!props.columns) {
    return {}
  }

  const cols: Record<string, Column[]> = {}
  for (const column of props.columns) {
    if (!cols[column.TableName]) {
      cols[column.TableName] = []
    }
    cols[column.TableName].push(column)
  }
  return cols
})
</script>

<template>
  <div class="d-flex flex-column ga-2">
    <v-card>
      <v-card-title>
        {{ props.connection.ConnectionUrl }}
        <ConnectionChip :connected="props.connection.Connected" />
      </v-card-title>
      <v-card-text>
        <v-col>
          <template v-if="props.connection.Connected === false">
            <v-row class="mb-2">
              <v-alert color="error">
                Error during connection: {{ props.connection.LastError }}
              </v-alert>
            </v-row>
          </template>
          <v-row>
            DB type: {{ props.connection.DbType }}
          </v-row>
          <v-row>
            Created at: {{ props.connection.CreatedAt }}
          </v-row>
          <v-row>
            Last connected at: {{ props.connection.LastConnectedAt }}
          </v-row>
        </v-col>
      </v-card-text>

      <v-card-actions>
        <v-btn
          v-bind="useLink(`/connections/${props.connectionId}/crons`)"
        >
          View Crons
        </v-btn>
        <v-btn
          v-bind="useLink(`/connections/${props.connectionId}/crons/create`)"
        >
          Create Cron
        </v-btn>
      </v-card-actions>
    </v-card>
    <v-card v-if="props.columns">
      <v-card-title>Schema</v-card-title>
      <v-card-text>
        <v-table>
          <thead>
            <tr>
              <td>Table</td>
              <td>Column</td>
              <td>Type</td>
              <td>Nullable</td>
              <td>Default</td>
              <td>Primary Key</td>
            </tr>
          </thead>
          <tbody>
            <template v-for="cols, table in colsPerTable" :key="table">
              <tr v-for="column, i in cols" :key="i">
                <td v-if="i === 0" :rowspan="cols.length">
                  {{ column.TableName }}
                </td>
                <td>{{ column.Name }}</td>
                <td>{{ column.Type }}</td>
                <td>{{ !column.Notnull }}</td>
                <td>{{ column.DfltValue }}</td>
                <td>{{ !!column.PK }}</td>
              </tr>
            </template>
          </tbody>
        </v-table>
      </v-card-text>
    </v-card>
  </div>
</template>
