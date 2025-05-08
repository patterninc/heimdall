import { buildQueryString } from '@/common/Services'
import { ApiParams } from '@/modules/Clusters/Helper'

export const getClusters = async (params: ApiParams) => {
  const queryString = buildQueryString(params)
  const response = await fetch(`/api/v1/clusters?${queryString}`)
  const data = await response.json()
  return data?.data
}

export const getClusterStatus = async () => {
  const response = await fetch('/api/v1/cluster/statuses')
  const data = await response.json()
  return data
}

export const getCluster = async (id: string) => {
  const response = await fetch(`/api/v1/clusters?id=${id}`)
  const data = await response.json()
  return data.data
}
