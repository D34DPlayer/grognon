import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import laravel from 'laravel-vite-plugin'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import vuetify from 'vite-plugin-vuetify'

export default defineConfig({
  plugins: [
    laravel({
      input: ['frontend/app.ts'],
      ssr: 'frontend/ssr.ts', // Enable SSR
      publicDirectory: 'public',
      buildDirectory: 'bootstrap',
      refresh: true,
    }),
    vue(),
    Components({
      dts: 'frontend/components.d.ts',
      dirs: ['frontend/components'],
      directoryAsNamespace: true,
      collapseSamePrefixes: true,
    }),
    vuetify({
      autoImport: true,
    }),
  ],
  build: {
    ssr: true, // Enable SSR
    outDir: 'bootstrap',
    sourcemap: true,
    rollupOptions: {
      input: 'frontend/ssr.ts',
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name][extname]',
        manualChunks: undefined, // Disable automatic chunk splitting
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './frontend'),
      '$': path.resolve(__dirname, '.'),
    },
  },
  ssr: {
    noExternal: ['vuetify'],
  },
})
