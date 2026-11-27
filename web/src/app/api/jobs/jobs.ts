import type {
  ApiError,
  Job,
  JobCancelResponse,
  JobStatus,
  JobsPage,
  JobsQuery,
} from '@/common/api/types'
import { API_URL, buildQueryString } from '@/common/Services'
import { JOBS_PAGE_SIZE } from '@/modules/Jobs/Helper'

export type JobFilters = Omit<
  JobsQuery,
  'limit' | 'cursor' | 'order_by' | 'direction'
>

/** `order_by` values that the API sorts on. Other columns sort on the internal row ID. */
const SORTABLE_COLUMNS: readonly string[] = [
  'id',
  'created_at',
  'updated_at',
] satisfies NonNullable<JobsQuery['order_by']>[]

const isSortColumn = (
  prop: string,
): prop is NonNullable<JobsQuery['order_by']> => SORTABLE_COLUMNS.includes(prop)

export const fetchJobs = async (
  params: JobFilters,
  cursor: string | null,
  sort: { prop: string; flip: boolean },
): Promise<JobsPage> => {
  const queryParams: JobsQuery = {
    ...params,
    limit: JOBS_PAGE_SIZE,
    direction: sort.flip ? 'desc' : 'asc',
  }
  if (isSortColumn(sort.prop)) queryParams.order_by = sort.prop
  // keyset: send the opaque cursor for the previous page (omitted on first page)
  if (cursor) queryParams.cursor = cursor
  const queryString = buildQueryString(queryParams)

  const response = await fetch(`${API_URL}/jobs?${queryString}`)
  if (!response.ok) return { data: [] }
  return await response.json()
}

/** Resolves to `null` when the API returns an error (for example an unknown job ID). */
export const getJobDetails = async (id: string): Promise<Job | null> => {
  const response = await fetch(`${API_URL}/job/${id}`)
  if (!response.ok) return null
  return await response.json()
}

export const getJobStatus = async (): Promise<JobStatus[]> => {
  const response = await fetch(`${API_URL}/job/statuses`)
  if (!response.ok) return []
  return await response.json()
}

export const cancelJob = async (
  id: string,
): Promise<JobCancelResponse | ApiError> => {
  const response = await fetch(`${API_URL}/job/${id}/cancel`, {
    method: 'POST',
  })
  return await response.json()
}
