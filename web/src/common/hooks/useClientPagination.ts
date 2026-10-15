import { useState } from 'react'

export const CLIENT_PAGE_SIZE = 20

/**
 * Pages an in-memory list. The page snaps back to 1 whenever `resetKey` changes (filters, sort),
 * so a narrower result set never lands on an out-of-range page.
 */
export function useClientPagination<T>(
  rows: T[],
  resetKey: string,
  pageSize = CLIENT_PAGE_SIZE,
) {
  const [state, setState] = useState({ key: resetKey, page: 1 })
  // Commit the reset (not just derive it), so returning to an earlier key can't restore its old page.
  if (state.key !== resetKey) setState({ key: resetKey, page: 1 })
  const totalPages = Math.max(1, Math.ceil(rows.length / pageSize))
  const requestedPage = state.key === resetKey ? state.page : 1
  const currentPage = Math.min(requestedPage, totalPages)

  return {
    pageRows: rows.slice((currentPage - 1) * pageSize, currentPage * pageSize),
    paginationProps: {
      currentPage,
      totalPages,
      onPageChange: (page: number) => setState({ key: resetKey, page }),
    },
  }
}
