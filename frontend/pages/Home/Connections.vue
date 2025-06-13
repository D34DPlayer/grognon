<script setup lang="ts">
import type { Connection } from '@/types'
import ConnectionChip from '@/components/ConnectionChip.vue'
import { useLink } from '@/composables'
import { displayTime } from '@/utils'
import { defineProps } from 'vue'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const props = defineProps<{
  connections: Connection[]
}>()
</script>

<template>
  <div class="d-flex flex-column ga-3">
    <v-card
      v-for="connection in props.connections"
      :key="connection.ConnectionId"
      v-bind="useLink(`/connections/${connection.ConnectionId}`)"
    >
      <v-card-title>
        {{ connection.ConnectionUrl }} <ConnectionChip :connected="connection.Connected" />
      </v-card-title>
      <v-card-text>
        Created at: {{ displayTime(connection.CreatedAt) }}
      </v-card-text>
    </v-card>

    <v-card v-bind="useLink('/connections/create')">
      <v-card-title class="text-center">
        Add new connection
      </v-card-title>
    </v-card>
  </div>
</template>
