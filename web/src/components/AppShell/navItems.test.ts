import { describe, expect, it } from 'vitest'

import { APP_ROUTES, isRouteActive } from './navItems'

describe('isRouteActive', () => {
  it('matches the route and its detail pages only', () => {
    expect(isRouteActive('/jobs', '/jobs')).toBe(true)
    expect(isRouteActive('/jobs/job-0001', '/jobs')).toBe(true)
    expect(isRouteActive('/jobsearch', '/jobs')).toBe(false)
    expect(isRouteActive('/clusters', '/jobs')).toBe(false)
  })

  it('lists the three resource pages', () => {
    expect(APP_ROUTES.map((route) => route.href)).toEqual([
      '/jobs',
      '/commands',
      '/clusters',
    ])
  })
})
