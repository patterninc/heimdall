import { screen, waitFor, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'

import { CommandDetails } from '@/modules/Commands/CommandDetails/CommandDetails'
import Commands from '@/modules/Commands/Commands'
import { api, buildCommand, commands, server } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import { recordRequests } from '@/test/requests'

const rowNames = () =>
  screen
    .getAllByRole('row')
    .slice(1)
    .map((row) => within(row).getAllByRole('cell')[0].textContent)

describe('Commands list page', () => {
  it('renders commands sorted by name with status and details links', async () => {
    renderWithProviders(<Commands />)

    expect(await screen.findByText('3 results')).toBeInTheDocument()
    expect(rowNames()).toEqual(['ping', 'shell', 'spark-sql'])
    for (const command of commands) {
      const row = screen.getByTestId(`command-row-${command.id}`)
      expect(within(row).getByText(command.status)).toBeInTheDocument()
      expect(within(row).getByRole('button', { name: 'View' })).toHaveAttribute(
        'href',
        `/commands/${command.id}`,
      )
    }
  })

  it('re-sorts on the client when a column header is clicked', async () => {
    const { user } = renderWithProviders(<Commands />)
    await screen.findByText('3 results')

    const nameHeader = screen.getByText('Name').closest('th')!
    await user.click(within(nameHeader).getByRole('button'))

    await waitFor(() =>
      expect(rowNames()).toEqual(['spark-sql', 'shell', 'ping']),
    )
  })

  it('applies a filter typed in the filter menu, then removes it', async () => {
    const requests = recordRequests('/api/v1/commands')
    const { user } = renderWithProviders(<Commands />)
    await screen.findByText(/results/)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(
      within(await screen.findByTestId('drawer-panel')).getByText('Plugin'),
    )
    const input = await screen.findByPlaceholderText('Enter plugin')
    await user.type(input, 'spark')

    await waitFor(() => expect(requests.last()!.get('plugin')).toBe('spark'))
    await user.clear(input)
    await waitFor(() => expect(requests.last()!.has('plugin')).toBe(false))
    requests.stop()
  })

  it('sends URL filters to the API', async () => {
    const requests = recordRequests('/api/v1/commands')
    renderWithProviders(<Commands />, {
      searchParams: '?plugin=spark&status=ACTIVE&user=alice',
    })
    await screen.findByText('3 results')

    expect(Object.fromEntries(requests.last()!)).toEqual({
      plugin: 'spark',
      status: 'ACTIVE',
      username: 'alice',
    })
    requests.stop()
  })

  it('pages large result sets on the client', async () => {
    server.use(
      http.get(api('/commands'), () =>
        HttpResponse.json({
          data: Array.from({ length: 25 }, (_, i) =>
            buildCommand({
              id: `c${i}`,
              name: `cmd-${String(i).padStart(2, '0')}`,
            }),
          ),
        }),
      ),
    )
    const { user } = renderWithProviders(<Commands />)

    expect(await screen.findByText('25 results')).toBeInTheDocument()
    expect(screen.getByText('1 of 2')).toBeInTheDocument()
    expect(rowNames()).toHaveLength(20)

    await user.click(screen.getByRole('button', { name: 'Next page' }))

    expect(screen.getByText('2 of 2')).toBeInTheDocument()
    expect(rowNames()).toEqual([
      'cmd-20',
      'cmd-21',
      'cmd-22',
      'cmd-23',
      'cmd-24',
    ])
  })

  it('shows the empty state without a clear action when no filters are set', async () => {
    server.use(
      http.get(api('/commands'), () => HttpResponse.json({ data: [] })),
    )
    renderWithProviders(<Commands />)

    expect(await screen.findByText('0 results')).toBeInTheDocument()
    expect(screen.getByText('No Data Available')).toBeInTheDocument()
    expect(
      screen.queryByRole('button', { name: 'Clear filters' }),
    ).not.toBeInTheDocument()
  })
})

describe('Command details page', () => {
  it('renders the overview, links, properties and tags', async () => {
    renderWithProviders(<CommandDetails id='cmd-spark' />)

    expect(await screen.findByText('spark-sql')).toBeInTheDocument()
    expect(screen.getByText('cmd-spark', { selector: 'p' })).toBeInTheDocument()
    expect(screen.getByText('ACTIVE')).toBeInTheDocument()
    expect(screen.getByText('Runs Spark SQL on EMR')).toBeInTheDocument()
    expect(screen.getByText('Async')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /Logs/ })).toHaveAttribute(
      'href',
      's3://logs/spark',
    )
    expect(screen.getByText('spark.executor.memory')).toBeInTheDocument()
    expect(screen.getByText('type:emr')).toBeInTheDocument()
  })

  it('hides empty sections for a command without context', async () => {
    renderWithProviders(<CommandDetails id='cmd-shell' />)

    expect(await screen.findByText('Sync')).toBeInTheDocument()
    expect(screen.queryByText('Links')).not.toBeInTheDocument()
    expect(screen.queryByText('Context properties')).not.toBeInTheDocument()
    expect(screen.queryByText('Cluster tags')).not.toBeInTheDocument()
  })

  it('shows the error state when the command does not exist', async () => {
    renderWithProviders(<CommandDetails id='missing' />)

    expect(
      await screen.findByText("We couldn't load this command"),
    ).toBeInTheDocument()
  })

  it('shows the error state when the API fails', async () => {
    server.use(http.get(api('/commands'), () => HttpResponse.error()))
    renderWithProviders(<CommandDetails id='cmd-spark' />)

    expect(
      await screen.findByText("We couldn't load this command"),
    ).toBeInTheDocument()
  })
})
