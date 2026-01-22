import { API_URL, buildQueryString } from '@/common/Services'
import { ApiParams } from '@/modules/Jobs/Helper'

export const fetchJobs = async (params: ApiParams, pageParam?: number) => {
  if (pageParam) {
    params.page = pageParam.toString()
  }
  const queryString = buildQueryString(params)

  const response = await fetch(`${API_URL}/jobs?${queryString}`)

  const data = await response.json()

  return data?.data
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
