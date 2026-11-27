import type {
  Cluster,
  ClustersPage,
  ClustersQuery,
  ResourceStatus,
} from '@/common/api/types'
import { API_URL, buildQueryString } from '@/common/Services'

export const getClusters = async (
  params: ClustersQuery,
): Promise<Cluster[]> => {
  const queryString = buildQueryString(params)
  const response = await fetch(`${API_URL}/clusters?${queryString}`)
  if (!response.ok) return []
  const page: ClustersPage = await response.json()
  return page.data ?? []
}

export const getClusterStatus = async (): Promise<ResourceStatus[]> => {
  const response = await fetch(`${API_URL}/cluster/statuses`)
  if (!response.ok) return []
  return await response.json()
}

export const getCluster = (id: string): Promise<Cluster[]> =>
  getClusters({ id })
