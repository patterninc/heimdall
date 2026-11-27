import type { ApiContext } from './types'

/** `context` is free-form in the API contract (each plugin sets its own keys), so read keys with a runtime check. */
export const contextString = (
  context: ApiContext | undefined,
  key: string,
): string | undefined => {
  const value = context?.[key]
  return typeof value === 'string' && value !== '' ? value : undefined
}

/** Reads a flat key/value map, such as Spark `properties`, as display strings. */
export const contextEntries = (
  context: ApiContext | undefined,
  key: string,
): [string, string][] => {
  const value = context?.[key]
  if (!value || typeof value !== 'object' || Array.isArray(value)) return []
  return Object.entries(value).map(([k, v]): [string, string] => [
    k,
    typeof v === 'string' ? v : JSON.stringify(v),
  ])
}
