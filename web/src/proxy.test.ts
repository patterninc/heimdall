import { NextRequest } from 'next/server'
import { describe, expect, it } from 'vitest'

import { config, proxy } from './proxy'

const request = (headers: Record<string, string>) =>
  new NextRequest('http://localhost:4000/jobs', { headers })

describe('proxy', () => {
  it('copies the Heimdall user and version headers into readable cookies', () => {
    const response = proxy(
      request({
        'X-Heimdall-User': 'alice@example.com',
        'X-Heimdall-Version': '1.4.2',
      }),
    )

    const user = response.cookies.get('heimdall-user')
    expect(user?.value).toBe('alice@example.com')
    expect(user?.httpOnly).toBe(false)
    expect(user?.sameSite).toBe('strict')
    expect(response.cookies.get('heimdall-version')?.value).toBe('1.4.2')
  })

  it('sets no cookies without the headers', () => {
    const response = proxy(request({}))

    expect(response.cookies.getAll()).toEqual([])
  })

  it('skips API and static routes', () => {
    const matcher = new RegExp(`^${config.matcher[0]}$`)

    expect(matcher.test('/jobs/job-0001')).toBe(true)
    expect(matcher.test('/api/v1/jobs')).toBe(false)
    expect(matcher.test('/_next/static/chunk.js')).toBe(false)
  })
})
