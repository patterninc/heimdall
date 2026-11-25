import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { mockMatchMedia } from '@/test/matchMedia'
import { useIsSmUp } from './useIsSmUp'

describe('useIsSmUp', () => {
  it('is true at the sm breakpoint and up', () => {
    const matchMedia = mockMatchMedia(() => true)

    const { result } = renderHook(() => useIsSmUp())

    expect(result.current).toBe(true)
    expect(matchMedia).toHaveBeenCalledWith('(min-width: 640px)')
  })

  it('is false below the sm breakpoint', () => {
    mockMatchMedia(() => false)

    const { result } = renderHook(() => useIsSmUp())

    expect(result.current).toBe(false)
  })
})
