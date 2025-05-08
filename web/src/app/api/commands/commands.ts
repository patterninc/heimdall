import { API_URL, buildQueryString } from '@/common/Services'
import { ApiParams } from '@/modules/Commands/Helper'

export const getCommands = async (params: ApiParams) => {
  const queryString = buildQueryString(params)
  const response = await fetch(`${API_URL}/commands?${queryString}`)
  const data = await response.json()
  return data.data
}

export const getCommandStatus = async () => {
  const response = await fetch(`${API_URL}/command/statuses`)
  const data = await response.json()
  return data
}

export const getCommandDetails = async (id: string) => {
  const queryString = buildQueryString({ id: id })
  const response = await fetch(`${API_URL}/commands?${queryString}`)
  const data = await response.json()
  return data.data
}
