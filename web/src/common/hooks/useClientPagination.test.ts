import { act, renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { CLIENT_PAGE_SIZE, useClientPagination } from './useClientPagination'

const rows = Array.from({ length: 45 }, (_, i) => i)

describe('useClientPagination', () => {
  it('returns the first page by default', () => {
    const { result } = renderHook(() => useClientPagination(rows, 'a'))

    expect(result.current.pageRows).toHaveLength(CLIENT_PAGE_SIZE)
    expect(result.current.pageRows[0]).toBe(0)
    expect(result.current.paginationProps).toMatchObject({
      currentPage: 1,
      totalPages: 3,
    })
  })

  it('moves to the requested page', () => {
    const { result } = renderHook(() => useClientPagination(rows, 'a'))

    act(() => result.current.paginationProps.onPageChange(3))

    expect(result.current.paginationProps.currentPage).toBe(3)
    expect(result.current.pageRows).toEqual([40, 41, 42, 43, 44])
  })

  it('resets to page 1 when the reset key changes and does not restore the old page', () => {
    const { result, rerender } = renderHook(
      ({ key }) => useClientPagination(rows, key),
      { initialProps: { key: 'a' } },
    )

    act(() => result.current.paginationProps.onPageChange(2))
    rerender({ key: 'b' })
    expect(result.current.paginationProps.currentPage).toBe(1)

    rerender({ key: 'a' })
    expect(result.current.paginationProps.currentPage).toBe(1)
  })

  it('clamps to the last page when the list shrinks', () => {
    const { result, rerender } = renderHook(
      ({ list }) => useClientPagination(list, 'a', 10),
      { initialProps: { list: rows } },
    )

    act(() => result.current.paginationProps.onPageChange(5))
    rerender({ list: rows.slice(0, 15) })

    expect(result.current.paginationProps).toMatchObject({
      currentPage: 2,
      totalPages: 2,
    })
    expect(result.current.pageRows).toEqual([10, 11, 12, 13, 14])
  })

  it('reports one page for an empty list', () => {
    const { result } = renderHook(() => useClientPagination([], 'a'))

    expect(result.current.pageRows).toEqual([])
    expect(result.current.paginationProps.totalPages).toBe(1)
  })
})
