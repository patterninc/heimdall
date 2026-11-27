import { describe, expect, it } from 'vitest'

import { recordRequests } from '@/test/requests'
import { jobs } from '@/test/fixtures'
import { cancelJob, fetchJobs, getJobDetails, getJobStatus } from './jobs'

describe('jobs API client', () => {
  it('sends filters, page size, sort and cursor to /jobs', async () => {
    const requests = recordRequests('/api/v1/jobs')

    const response = await fetchJobs(
      { name: 'rollup', status: ['RUNNING', 'FAILED'] },
      'cursor-2',
      { prop: 'created_at', flip: true },
    )

    expect(response.data).toHaveLength(jobs.length)
    expect(Object.fromEntries(requests.last()!)).toEqual({
      name: 'rollup',
      status: 'RUNNING,FAILED',
      limit: '20',
      order_by: 'created_at',
      direction: 'desc',
      cursor: 'cursor-2',
    })
    requests.stop()
  })

  it('omits the cursor on the first page and sorts ascending when not flipped', async () => {
    const requests = recordRequests('/api/v1/jobs')

    await fetchJobs({}, null, { prop: 'id', flip: false })

    expect(requests.last()!.has('cursor')).toBe(false)
    expect(requests.last()!.get('direction')).toBe('asc')
    requests.stop()
  })

  it('reads a job, its statuses, and cancels it', async () => {
    await expect(getJobDetails('job-0001')).resolves.toMatchObject({
      id: 'job-0001',
    })
    await expect(getJobStatus()).resolves.toContain('RUNNING')
    await expect(cancelJob('job-0002')).resolves.toEqual({
      status: 'CANCELING',
    })
  })
})
