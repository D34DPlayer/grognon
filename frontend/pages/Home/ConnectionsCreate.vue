<script setup lang="ts">
import type { ConnectionCreate } from '@/types'
import { useLink } from '@/composables'
import { DB_TYPES } from '@/types'
import { useForm } from '@inertiajs/vue3'
import { ref } from 'vue'
import Layout from '../Layout.vue'
import HomeLayout from './Layout.vue'

defineOptions({
  layout: [Layout, HomeLayout],
})

const form = useForm({
  DbType: undefined,
  ConnectionUrl: '',
} as Partial<ConnectionCreate>)
const isValid = ref(false)

function onSubmit() {
  form.post('/connections')
}
</script>

<template>
  <v-form v-model="isValid" @submit.prevent="onSubmit">
    <v-card>
      <v-card-title>
        Create a new connection
      </v-card-title>
      <v-card-text class="pb-0">
        <div class="d-flex flex-column ga-3">
          <v-select
            v-model="form.DbType"
            label="DB type"
            placeholder="Select a DB type"
            :items="DB_TYPES"
            :rules="[
              (v) => !!v || 'DB type is required',
            ]"
          />
          <v-text-field
            v-model="form.ConnectionUrl"
            label="Connection url"
            :rules="[
              (v) => !!v || 'Connection url is required',
              (v) => v.length > 0 || 'Connection url must be at least 1 character',
            ]"
          />
        </div>
      </v-card-text>

      <v-card-actions class="justify-space-between">
        <v-btn v-bind="useLink('/connections')">
          Cancel
        </v-btn>
        <v-btn :disabled="!isValid" type="submit" color="primary">
          Create
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-form>
</template>
