import { API_URL, buildQueryString } from '@/common/Services'
import { ApiParams, JOBS_PAGE_SIZE, JobsResponse } from '@/modules/Jobs/Helper'

export const fetchJobs = async (
  params: ApiParams,
  cursor: string | null,
  sort: { prop: string; flip: boolean },
): Promise<JobsResponse> => {
  const queryParams: ApiParams = {
    ...params,
    limit: JOBS_PAGE_SIZE.toString(),
    order_by: sort.prop,
    direction: sort.flip ? 'desc' : 'asc',
  }
  // keyset pagination: send the opaque cursor for the previous page's last row.
  // Omitted on the first page.
  if (cursor) queryParams.cursor = cursor
  const queryString = buildQueryString(queryParams)

  const response = await fetch(`${API_URL}/jobs?${queryString}`)

  return await response.json()
}

export const getJobDetails = async (id: string) => {
  const response = await fetch(`${API_URL}/job/${id}`)
  return await response.json()
}

export const getJobStatus = async () => {
  const response = await fetch(`${API_URL}/job/statuses`)
  return await response.json()
}

export const cancelJob = async (id: string) => {
  const response = await fetch(`${API_URL}/job/${id}/cancel`, {
    method: 'POST',
  })
  return await response.json()
}
