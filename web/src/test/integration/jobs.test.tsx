import { screen, waitFor, within } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'

import JobDetails from '@/modules/Jobs/JobDetails/JobDetails'
import Jobs from '@/modules/Jobs/Jobs'
import { api, buildJob, jobs, jobsPage, server } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import { recordRequests } from '@/test/requests'
import { mockRouter } from '@/test/router'

describe('Jobs list page', () => {
  it('renders every job from the API with its status and a link to its details', async () => {
    renderWithProviders(<Jobs />)

    expect(await screen.findByText('3 results')).toBeInTheDocument()
    for (const job of jobs) {
      const row = screen.getByTestId(`job-row-${job.id}`)
      expect(within(row).getByText(job.name)).toBeInTheDocument()
      expect(within(row).getByText(job.status)).toBeInTheDocument()
      expect(within(row).getByRole('button', { name: 'View' })).toHaveAttribute(
        'href',
        `/jobs/${job.id}`,
      )
    }
    expect(screen.getAllByText('2026-01-15 10:30:00')).toHaveLength(3)
  })

  it('requests newest jobs first and flips the sort when the active column is clicked', async () => {
    const requests = recordRequests('/api/v1/jobs')
    const { user } = renderWithProviders(<Jobs />)
    await screen.findByText('3 results')

    expect(requests.last()!.get('order_by')).toBe('created_at')
    expect(requests.last()!.get('direction')).toBe('desc')

    const createdAtHeader = screen.getByText('Created At').closest('th')!
    await user.click(within(createdAtHeader).getByRole('button'))

    await waitFor(() => expect(requests.last()!.get('direction')).toBe('asc'))
    requests.stop()
  })

  it('sends URL filters to the API', async () => {
    const requests = recordRequests('/api/v1/jobs')
    renderWithProviders(<Jobs />, {
      searchParams: '?status=FAILED,RUNNING&user=bob&tags=team:data',
    })
    await screen.findByText('3 results')

    expect(Object.fromEntries(requests.last()!)).toMatchObject({
      status: 'FAILED,RUNNING',
      username: 'bob',
      tags: 'team:data',
    })
    requests.stop()
  })

  it('applies a filter typed in the filter menu after the debounce, then removes it', async () => {
    const requests = recordRequests('/api/v1/jobs')
    const { user } = renderWithProviders(<Jobs />)
    await screen.findByText('3 results')

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(
      within(await screen.findByTestId('drawer-panel')).getByText('User'),
    )
    const input = await screen.findByPlaceholderText('Enter user')
    await user.type(input, 'bob')

    await waitFor(() => expect(requests.last()!.get('username')).toBe('bob'))
    await user.clear(input)
    await waitFor(() => expect(requests.last()!.has('username')).toBe(false))
    requests.stop()
  })

  it('pages forward with the keyset cursor and shows an open-ended count', async () => {
    const firstPage = Array.from({ length: 20 }, (_, i) =>
      buildJob({ id: `job-a${i}`, name: `first-${i}` }),
    )
    server.use(
      http.get(api('/jobs'), ({ request }) => {
        const cursor = new URL(request.url).searchParams.get('cursor')
        return HttpResponse.json(
          cursor === 'next-1'
            ? jobsPage([buildJob({ id: 'job-b0', name: 'second-0' })])
            : jobsPage(firstPage, 'next-1'),
        )
      }),
    )
    const { user } = renderWithProviders(<Jobs />)

    expect(await screen.findByText('20+ results')).toBeInTheDocument()
    await user.click(screen.getByRole('button', { name: 'Next page' }))

    expect(await screen.findByText('second-0')).toBeInTheDocument()
    expect(await screen.findByText('21 results')).toBeInTheDocument()
    expect(screen.queryByText('first-0')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'Previous page' }))
    expect(await screen.findByText('first-0')).toBeInTheDocument()
  })

  it('shows the empty state with a clear-filters action when filters match nothing', async () => {
    server.use(http.get(api('/jobs'), () => HttpResponse.json(jobsPage([]))))
    renderWithProviders(<Jobs />, { searchParams: '?name=nothing' })

    expect(await screen.findByText('0 results')).toBeInTheDocument()
    expect(screen.getByText('No Data Available')).toBeInTheDocument()
    expect(
      screen.getByRole('button', { name: 'Clear filters' }),
    ).toBeInTheDocument()
  })
})

describe('Job details page', () => {
  it('renders the overview, outputs, attributes, query and criteria', async () => {
    renderWithProviders(<JobDetails id='job-0001' />)

    expect(
      await screen.findByRole('link', { name: 'spark-sql' }),
    ).toBeInTheDocument()
    expect(screen.getByTestId('job-details-status')).toHaveTextContent(
      'SUCCEEDED',
    )
    expect(screen.getByRole('link', { name: 'spark-sql' })).toHaveAttribute(
      'href',
      '/commands/cmd-spark',
    )
    expect(screen.getByRole('link', { name: 'emr-prod' })).toHaveAttribute(
      'href',
      '/clusters/cluster-emr',
    )
    expect(screen.getByText('Async')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /View stdout/ })).toHaveAttribute(
      'href',
      '/api/v1/job/job-0001/stdout',
    )
    expect(screen.getByRole('link', { name: /spark_ui/ })).toHaveAttribute(
      'href',
      'https://spark.example.com/app/1',
    )
    expect(screen.getByText('rows_written')).toBeInTheDocument()
    expect(screen.getByText('spark.executor.memory')).toBeInTheDocument()
    expect(screen.getByText('type:spark')).toBeInTheDocument()
    expect(screen.queryByText('Error')).not.toBeInTheDocument()
  })

  it('shows the error message of a failed job', async () => {
    renderWithProviders(<JobDetails id='job-0003' />)

    expect(
      await screen.findByText('java.lang.OutOfMemoryError: Java heap space'),
    ).toBeInTheDocument()
  })

  it('shows the error state when the job does not exist and links back to the list', async () => {
    const { user } = renderWithProviders(<JobDetails id='missing' />)

    const errorState = await screen.findByTestId('job-details-error')
    expect(
      within(errorState).getByText("We couldn't load this job"),
    ).toBeInTheDocument()
    await user.click(
      within(errorState).getByRole('button', { name: 'Back to jobs' }),
    )
    expect(mockRouter.push).toHaveBeenCalledWith('/jobs')
  })

  it('cancels a running job after confirmation', async () => {
    let canceled = ''
    server.use(
      http.post(api('/job/:id/cancel'), ({ params }) => {
        canceled = String(params.id)
        return HttpResponse.json({ status: 'CANCELING' })
      }),
    )
    const { user } = renderWithProviders(<JobDetails id='job-0002' />)

    await user.click(await screen.findByRole('button', { name: 'Cancel job' }))
    const dialog = await screen.findByTestId('job-details-cancel-confirm')
    expect(within(dialog).getByText('Cancel job?')).toBeInTheDocument()
    await user.click(within(dialog).getByRole('button', { name: 'Cancel job' }))

    await waitFor(() => expect(canceled).toBe('job-0002'))
  })

  it('does not allow canceling a finished job', async () => {
    renderWithProviders(<JobDetails id='job-0001' />)

    expect(
      await screen.findByRole('button', { name: 'Cancel job' }),
    ).toBeDisabled()
  })
})
