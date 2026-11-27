import { screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { CLUSTER_COLUMNS, toClusterRows } from '@/modules/Clusters/Helper'
import { clusters } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import { withStableIds } from '@/test/snapshot'
import ListEmptyState from './ListEmptyState'
import ListPage from './ListPage'

const filters = <button type='button'>Filters</button>

describe('ListPage', () => {
  it('matches the snapshot with rows', () => {
    const { container } = renderWithProviders(
      <ListPage
        title='Clusters'
        resultCount='2'
        qaTestId='clusters-page'
        filters={filters}
        tableProps={{
          columns: CLUSTER_COLUMNS,
          rows: toClusterRows(clusters),
          sortedColumn: 'name',
          sortDirection: 'asc',
          onSort: () => {},
          qaTestId: 'clusters-table',
        }}
      />,
    )

    expect(screen.getByText('2 results')).toBeInTheDocument()
    expect(withStableIds(container)).toMatchSnapshot()
  })

  it('matches the snapshot for the empty state', async () => {
    const onClearFilters = vi.fn()
    const { container, user } = renderWithProviders(
      <ListPage
        title='Clusters'
        resultCount='0'
        qaTestId='clusters-page'
        filters={filters}
        tableProps={{
          columns: CLUSTER_COLUMNS,
          rows: [],
          emptyState: <ListEmptyState onClearFilters={onClearFilters} />,
          qaTestId: 'clusters-table',
        }}
      />,
    )

    expect(withStableIds(container)).toMatchSnapshot()

    await user.click(screen.getByRole('button', { name: 'Clear filters' }))
    expect(onClearFilters).toHaveBeenCalledOnce()
  })

  it('omits the result count while the first page loads', () => {
    renderWithProviders(
      <ListPage
        title='Clusters'
        qaTestId='clusters-page'
        filters={filters}
        tableProps={{
          columns: CLUSTER_COLUMNS,
          rows: [],
          isLoading: true,
          qaTestId: 'clusters-table',
        }}
      />,
    )

    expect(screen.getByText('Clusters')).toBeInTheDocument()
    expect(screen.queryByText(/results/)).not.toBeInTheDocument()
  })
})
