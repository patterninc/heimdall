import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { resourceStatusVariant } from '@/common/Services/status'
import { buildCluster } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import { withStableIds } from '@/test/snapshot'
import { mockRouter } from '@/test/router'
import DetailPage from './DetailPage'

const cluster = buildCluster()

const renderDetailPage = (props: { isError?: boolean } = {}) =>
  renderWithProviders(
    <DetailPage
      backHref='/clusters'
      backLabel='Back to clusters'
      resourceLabel='cluster'
      title={cluster.name}
      subtitle={cluster.id}
      status={cluster.status}
      statusVariant={resourceStatusVariant(cluster.status)}
      right={<button type='button'>Header action</button>}
      qaTestId='cluster-details'
      {...props}
    >
      <p>Detail content</p>
    </DetailPage>,
  )

describe('DetailPage', () => {
  it('matches the snapshot with header, status and content', () => {
    const { container } = renderDetailPage()

    expect(screen.getByTestId('cluster-details-status')).toHaveTextContent(
      cluster.status,
    )
    expect(withStableIds(container)).toMatchSnapshot()
  })

  it('matches the snapshot in the error state', () => {
    const { container } = renderDetailPage({ isError: true })

    expect(screen.queryByText('Header action')).not.toBeInTheDocument()
    expect(withStableIds(container)).toMatchSnapshot()
  })

  it('matches the snapshot with only the required header props', () => {
    const { container } = renderWithProviders(
      <DetailPage
        backHref='/clusters'
        backLabel='Back to clusters'
        resourceLabel='cluster'
        title='cluster-emr'
        qaTestId='cluster-details'
      >
        <p>Detail content</p>
      </DetailPage>,
    )

    expect(
      screen.queryByTestId('cluster-details-status'),
    ).not.toBeInTheDocument()
    expect(withStableIds(container)).toMatchSnapshot()
  })

  it('navigates to the list from the back button', async () => {
    const { user } = renderDetailPage()

    await user.click(screen.getByTestId('cluster-details-back'))

    expect(mockRouter.push).toHaveBeenCalledWith('/clusters')
  })
})
