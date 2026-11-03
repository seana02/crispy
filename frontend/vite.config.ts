import { defineConfig } from 'vite';
import solidPlugin from 'vite-plugin-solid';

export default defineConfig({
  plugins: [solidPlugin()],
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    target: 'esnext',
  },
  resolve: {
    alias: {
      '~': '/src',
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    transformMode: {
      web: [/\.[jt]sx?$/],
    },
    deps: {
      optimizer: {
        web: {
          include: [],
        },
      },
    },
  },
});
