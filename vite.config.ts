import vue from '@vitejs/plugin-vue'
import laravel from 'laravel-vite-plugin'
import { defineConfig } from 'vite'

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
  ],
  build: {
    manifest: true, // Generate manifest.json file
    outDir: 'public/build',
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
  server: {
    hmr: {
      host: 'localhost',
    },
  },
})
