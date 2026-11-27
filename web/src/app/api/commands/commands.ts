import type {
  Command,
  CommandsPage,
  CommandsQuery,
  ResourceStatus,
} from '@/common/api/types'
import { API_URL, buildQueryString } from '@/common/Services'

export const getCommands = async (
  params: CommandsQuery,
): Promise<Command[]> => {
  const queryString = buildQueryString(params)
  const response = await fetch(`${API_URL}/commands?${queryString}`)
  if (!response.ok) return []
  const page: CommandsPage = await response.json()
  return page.data ?? []
}

export const getCommandStatus = async (): Promise<ResourceStatus[]> => {
  const response = await fetch(`${API_URL}/command/statuses`)
  if (!response.ok) return []
  return await response.json()
}

export const getCommandDetails = (id: string): Promise<Command[]> =>
  getCommands({ id })
