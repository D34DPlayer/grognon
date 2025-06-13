<script setup lang="ts">
import { sql } from '@codemirror/lang-sql'
import { Codemirror } from 'vue-codemirror'

const extensions = [
  sql(),
]

const model = defineModel<string>({
  type: String,
  default: '',
})
</script>

<template>
  <Codemirror v-model="model" placeholder="Enter your SQL script here" :extensions="extensions" />
</template>
