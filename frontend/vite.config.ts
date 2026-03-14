/// <reference types="vitest/config" />

import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  base: '/',
  plugins: [
    tsconfigPaths(),
    react({
      babel: {
        plugins: [['babel-plugin-react-compiler']],
      },
    }),
  ],
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/tests/vitest.setup.ts'],
    globals: true,
    watch: false,
    include: ['./src/tests/*.test.tsx'],
  },
  preview: {
    port: 3000,
    strictPort: true,
  },
  server: {
    port: 3000,
    strictPort: true,
    host: '0.0.0.0',
    origin: 'http://0.0.0.0:3000',
    proxy: {
      '/ws': {
        target: 'http://backend-dev:8080',
        ws: true,
        changeOrigin: true,
      },
      '/api': {
        target: 'http://backend-dev:8080',
        changeOrigin: true,
      },
    },
  },
});
