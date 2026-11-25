import '@testing-library/jest-dom/vitest'
import { cleanup, configure } from '@testing-library/react'
import { afterAll, afterEach, beforeAll, vi } from 'vitest'

import { server } from './fixtures/server'
import { createMediaQueryList } from './matchMedia'
import { mockLocation, mockRouter } from './router'

// The app and pattern-ui tag elements with `qa-test-id`, so `*ByTestId` queries read that attribute.
configure({ testIdAttribute: 'qa-test-id', asyncUtilTimeout: 3000 })

vi.mock('next/navigation', () => ({
  useRouter: () => mockRouter,
  usePathname: () => mockLocation.pathname,
  useSearchParams: () => new URLSearchParams(),
}))

// jsdom has no layout APIs; pattern-ui and usehooks-ts read them on mount.
class NoopObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
  takeRecords() {
    return []
  }
}
Object.assign(globalThis, {
  ResizeObserver: NoopObserver,
  IntersectionObserver: NoopObserver,
})
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: (query: string) => createMediaQueryList(query, false),
})
Element.prototype.scrollIntoView = () => {}

beforeAll(() => server.listen({ onUnhandledFrame: 'error' }))
afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
  server.resetHandlers()
  Object.values(mockRouter).forEach((fn) => fn.mockReset())
  mockLocation.pathname = '/jobs'
  localStorage.clear()
  document.cookie
    .split(';')
    .map((cookie) => cookie.split('=')[0].trim())
    .filter(Boolean)
    .forEach((name) => {
      document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`
    })
})
afterAll(() => server.close())
