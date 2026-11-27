import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { useUser } from './useUser'

describe('useUser', () => {
  it('returns the decoded heimdall-user cookie', () => {
    document.cookie = 'theme=dark'
    document.cookie = `heimdall-user=${encodeURIComponent('alice@example.com')}`

    const { result } = renderHook(() => useUser())

    expect(result.current).toBe('alice@example.com')
  })

  it('returns null when the cookie is missing', () => {
    const { result } = renderHook(() => useUser())

    expect(result.current).toBeNull()
  })
})
