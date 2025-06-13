<script setup lang="ts">
import type { Column, Connection, CronCreate } from '@/types'
import { getExtensions } from '@/codemirror'
import { useLink } from '@/composables'
import { SCHEDULES } from '@/types'
import { useForm } from '@inertiajs/vue3'
import { computed, ref } from 'vue'
import { Codemirror } from 'vue-codemirror'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const props = defineProps<{
  connectionId?: number
  connections?: Connection[]
  columns?: Column[]
}>()

const form = useForm({
  ConnectionId: props.connectionId,
  Name: '',
  Command: '',
  Schedule: undefined,
} as Partial<CronCreate>)
const isValid = ref(false)

const connectionsOptions = computed(() => {
  if (!props.connections) {
    return []
  }
  return props.connections.map(c => ({
    title: c.ConnectionUrl,
    value: c.ConnectionId,
  }))
})

const invalidConnectionId = computed(() => {
  return !props.connections?.some(c => c.ConnectionId === props.connectionId)
})

function onSubmit() {
  console.log('onSubmit', form, isValid.value)
  form.post('/crons')
}
</script>

<template>
  <v-form v-model="isValid" @submit.prevent="onSubmit">
    <v-card>
      <v-card-title>Create a new cron</v-card-title>
      <v-card-text class="pb-0">
        <div class="d-flex flex-column ga-3">
          <v-select
            v-model="form.ConnectionId"
            label="Connection"
            placeholder="Select a connection"
            :readonly="!invalidConnectionId"
            :items="connectionsOptions"
            :rules="[v => !!v || 'Connection is required']"
          />
          <v-text-field
            v-model="form.Name"
            label="Cron name"
            :rules="[v => !!v || 'Cron name is required']"
          />

          <v-select
            v-model="form.Schedule"
            label="Schedule"
            :items="SCHEDULES"
            :rules="[v => !!v || 'Schedule is required']"
          />
          <h3>SQL Command</h3>
          <Codemirror
            v-model="form.Command"
            placeholder="Type your SQL command here..."
            :indent-with-tab="true"
            :tab-size="2"
            :extensions="getExtensions(columns)"
          />
        </div>
      </v-card-text>

      <v-card-actions class="justify-space-between">
        <v-btn v-bind="useLink('/crons')">
          Cancel
        </v-btn>
        <v-btn :disabled="!isValid" type="submit" color="primary">
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>
