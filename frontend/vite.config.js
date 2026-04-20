import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // Proxy to your Go backend
        changeOrigin: true,
      },
      // You can add other backend endpoints here (e.g., '/api')
    }
  },
  build: {
    outDir: '../web', // Build the app into the Go root 'web' folder
    emptyOutDir: true,
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        login: resolve(__dirname, 'login.html'),
        // game: resolve(__dirname, 'game.html'), // Make sure to add all your MPA pages here!
      }
    }
  }
});