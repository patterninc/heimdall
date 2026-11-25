import { vi } from 'vitest'

/** Shared `next/navigation` router mock; assert on it to check navigation. */
export const mockRouter = {
  push: vi.fn(),
  replace: vi.fn(),
  back: vi.fn(),
  forward: vi.fn(),
  refresh: vi.fn(),
  prefetch: vi.fn(),
}

/** Path the mocked `usePathname()` returns; tests set it, the setup file resets it to `/jobs`. */
export const mockLocation = { pathname: '/jobs' }
