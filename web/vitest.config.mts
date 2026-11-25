import react from '@vitejs/plugin-react'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vitest/config'

// Dates render through the runtime time zone; pin it so snapshots and assertions match in every environment.
process.env.TZ = 'UTC'

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: { '@': fileURLToPath(new URL('./src', import.meta.url)) },
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
    css: false,
    server: {
      deps: {
        inline: ['@patterninc/pattern-ui'],
      },
    },
    coverage: {
      provider: 'v8',
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.test.{ts,tsx}', 'src/test/**', 'src/app/**/*.tsx'],
      reporter: ['text', 'json-summary'],
      // CI gate: `pnpm test` fails when coverage drops below these numbers. Raise them as coverage grows.
      thresholds: {
        statements: 85,
        branches: 80,
        functions: 85,
        lines: 85,
      },
    },
  },
})
