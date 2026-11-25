import { screen, waitFor, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'

import Clusters from '@/modules/Clusters/Clusters'
import ClustersDetails from '@/modules/Clusters/ClustersDetails/ClustersDetails'
import { api, buildCluster, clusters, server } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import { recordRequests } from '@/test/requests'

const rowNames = () =>
  screen
    .getAllByRole('row')
    .slice(1)
    .map((row) => within(row).getAllByRole('cell')[0].textContent)

describe('Clusters list page', () => {
  it('renders clusters sorted by name with status and details links', async () => {
    renderWithProviders(<Clusters />)

    expect(await screen.findByText('2 results')).toBeInTheDocument()
    expect(rowNames()).toEqual(['emr-prod', 'local'])
    for (const cluster of clusters) {
      const row = screen.getByTestId(`cluster-row-${cluster.id}`)
      expect(within(row).getByText(cluster.status)).toBeInTheDocument()
      expect(within(row).getByRole('button', { name: 'View' })).toHaveAttribute(
        'href',
        `/clusters/${cluster.id}`,
      )
    }
  })

  it('re-sorts on the client when a column header is clicked', async () => {
    const { user } = renderWithProviders(<Clusters />)
    await screen.findByText('2 results')

    const userHeader = screen.getByText('User').closest('th')!
    await user.click(within(userHeader).getByRole('button'))
    await user.click(within(userHeader).getByRole('button'))

    await waitFor(() => expect(rowNames()).toEqual(['local', 'emr-prod']))
  })

  it('applies a filter typed in the filter menu, then removes it', async () => {
    const requests = recordRequests('/api/v1/clusters')
    const { user } = renderWithProviders(<Clusters />)
    await screen.findByText(/results/)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(
      within(await screen.findByTestId('drawer-panel')).getByText('Version'),
    )
    const input = await screen.findByPlaceholderText('Enter version')
    await user.type(input, '7.1.0')

    await waitFor(() => expect(requests.last()!.get('version')).toBe('7.1.0'))
    await user.clear(input)
    await waitFor(() => expect(requests.last()!.has('version')).toBe(false))
    requests.stop()
  })

  it('sends URL filters to the API', async () => {
    const requests = recordRequests('/api/v1/clusters')
    renderWithProviders(<Clusters />, {
      searchParams: '?name=emr&version=7.1.0&status=ACTIVE,INACTIVE',
    })
    await screen.findByText('2 results')

    expect(Object.fromEntries(requests.last()!)).toEqual({
      name: 'emr',
      version: '7.1.0',
      status: 'ACTIVE,INACTIVE',
    })
    requests.stop()
  })

  it('caps the result count at the API limit', async () => {
    server.use(
      http.get(api('/clusters'), () =>
        HttpResponse.json({
          data: Array.from({ length: 101 }, (_, i) =>
            buildCluster({ id: `cl${i}`, name: `cluster-${i}` }),
          ),
        }),
      ),
    )
    renderWithProviders(<Clusters />)

    expect(await screen.findByText('100+ results')).toBeInTheDocument()
    expect(screen.getByText('1 of 6')).toBeInTheDocument()
  })
})

describe('Cluster details page', () => {
  it('renders the overview, runtime, properties and tags', async () => {
    renderWithProviders(<ClustersDetails id='cluster-emr' />)

    expect(await screen.findByText('emr-7.1.0')).toBeInTheDocument()
    expect(screen.getByText('emr-prod')).toBeInTheDocument()
    expect(screen.getByText('ACTIVE')).toBeInTheDocument()
    expect(
      screen.getByText('arn:aws:iam::123456789012:role/emr-exec'),
    ).toBeInTheDocument()
    expect(screen.getByText('Production EMR on EKS')).toBeInTheDocument()
    expect(screen.getByText('spark.executor.memory')).toBeInTheDocument()
    expect(screen.getByText('env:prod')).toBeInTheDocument()
  })

  it('hides the runtime, properties and tags sections when they are empty', async () => {
    renderWithProviders(<ClustersDetails id='cluster-local' />)

    expect(await screen.findByText('INACTIVE')).toBeInTheDocument()
    expect(screen.queryByText('Runtime')).not.toBeInTheDocument()
    expect(screen.queryByText('Context properties')).not.toBeInTheDocument()
    expect(screen.queryByText('Tags')).not.toBeInTheDocument()
  })

  it('shows the error state when the cluster does not exist', async () => {
    renderWithProviders(<ClustersDetails id='missing' />)

    expect(
      await screen.findByText("We couldn't load this cluster"),
    ).toBeInTheDocument()
  })
})
