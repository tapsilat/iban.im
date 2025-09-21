const { defineConfig } = require('vite');
const vue = require('@vitejs/plugin-vue');
const path = require('path');

// https://vitejs.dev/config/
module.exports = defineConfig({
  plugins: [vue()],
  root: '.',
  publicDir: 'public',
  server: {
    port: 4881,
    proxy: {
      '/graph': {
        target: 'http://localhost:4880',
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  css: {
    preprocessorOptions: {
      scss: {
        // Enable @import of node_modules without ~ prefix
        additionalData: ''
      }
    }
  },
  build: {
    outDir: 'dist',
    emptyOutDir: true,
  },
});
