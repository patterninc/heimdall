import moment from 'moment-timezone'

export const API_URL = '/api/v1'

/** `flip` is the descending flag. */
export type SortBy = {
  prop: string
  flip: boolean
}

export const noDataAvailable = 'No Data Available'
export const noDataAvailableDescription =
  'Please update the filters to find what you need.'

export const buildQueryString = (
  params: Record<string, string | string[] | number>,
): string => {
  const queryParams = new URLSearchParams()

  Object.entries(params).forEach(([key, value]) => {
    if (Array.isArray(value)) {
      queryParams.append(key, value.join(','))
    } else {
      queryParams.append(key, value.toString())
    }
  })

  return queryParams.toString()
}

/** Em dash for empty values, matching the design system's empty-cell convention. */
export const orDash = (value?: string | number | null): string =>
  value === undefined || value === null || value === '' ? '—' : String(value)

/** Clicking the active column flips direction; a new column starts ascending. */
export const toggleSort = (current: SortBy, columnKey: string): SortBy =>
  current.prop === columnKey
    ? { prop: columnKey, flip: !current.flip }
    : { prop: columnKey, flip: false }

// Utility function to format dates with time zone
export const formatDateWithTimeZone = (date: number, timeZone: string) => {
  return moment.unix(date).tz(timeZone).format('YYYY-MM-DD HH:mm:ss')
}

export const myTimezone = Intl.DateTimeFormat().resolvedOptions().timeZone

// Utility function to sort data
export const sortData = <T>(data: T[], sortBy: SortBy): T[] => {
  // Copy first so the react-query cache isn't sorted in place.
  return [...(data ?? [])].sort((a, b) => {
    const aValue = a[sortBy.prop as keyof T]
    const bValue = b[sortBy.prop as keyof T]
    if (aValue < bValue) return sortBy.flip ? 1 : -1
    if (aValue > bValue) return sortBy.flip ? -1 : 1
    return 0
  })
}

// Utility function to get version
export const getVersion = () => {
  if (typeof window === 'undefined') {
    return null
  }
  const value = `; ${document.cookie}`
  const parts = value.split(`; heimdall-version=`)
  const version = parts.pop()?.split(';').shift() || null
  return version ? `${version}` : null
}
