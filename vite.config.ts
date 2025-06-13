import path from 'node:path'
import vue from '@vitejs/plugin-vue'
import laravel from 'laravel-vite-plugin'
import Components from 'unplugin-vue-components/vite'
import { defineConfig } from 'vite'
import vuetify from 'vite-plugin-vuetify'

export default defineConfig({
  plugins: [
    laravel({
      input: 'frontend/app.ts',
      publicDirectory: 'public',
      buildDirectory: 'build',
      refresh: true,
    }),
    vue({
      include: [/\.vue$/],
    }),
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
    manifest: true, // Generate manifest.json file
    outDir: 'public/build',
    sourcemap: true,
    rollupOptions: {
      input: 'frontend/app.ts',
      output: {
        entryFileNames: 'assets/[name].js',
        chunkFileNames: 'assets/[name].js',
        assetFileNames: 'assets/[name].[ext]',
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
  server: {
    hmr: {
      host: 'localhost',
    },
  },
})
