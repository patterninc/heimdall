import { afterEach, describe, expect, it } from 'vitest'

import { FIXTURE_CREATED_AT, commands } from '@/test/fixtures'
import {
  buildQueryString,
  formatDateWithTimeZone,
  getVersion,
  orDash,
  sortData,
  toggleSort,
} from './index'

describe('buildQueryString', () => {
  it('encodes scalar values', () => {
    expect(buildQueryString({ name: 'spark sql', limit: 20 })).toBe(
      'name=spark+sql&limit=20',
    )
  })

  it('joins array values with commas', () => {
    expect(buildQueryString({ status: ['RUNNING', 'FAILED'] })).toBe(
      'status=RUNNING%2CFAILED',
    )
  })

  it('returns an empty string for no params', () => {
    expect(buildQueryString({})).toBe('')
  })
})

describe('orDash', () => {
  it.each([undefined, null, ''])('renders %j as an em dash', (value) => {
    expect(orDash(value)).toBe('—')
  })

  it('keeps zero and other values', () => {
    expect(orDash(0)).toBe('0')
    expect(orDash('alice')).toBe('alice')
  })
})

describe('toggleSort', () => {
  it('flips direction when the active column is clicked again', () => {
    expect(toggleSort({ prop: 'name', flip: false }, 'name')).toEqual({
      prop: 'name',
      flip: true,
    })
  })

  it('starts a new column ascending', () => {
    expect(toggleSort({ prop: 'name', flip: true }, 'user')).toEqual({
      prop: 'user',
      flip: false,
    })
  })
})

describe('sortData', () => {
  it('sorts ascending and descending by the given prop', () => {
    const names = (sorted: typeof commands) => sorted.map((c) => c.name)

    expect(names(sortData(commands, { prop: 'name', flip: false }))).toEqual([
      'ping',
      'shell',
      'spark-sql',
    ])
    expect(names(sortData(commands, { prop: 'name', flip: true }))).toEqual([
      'spark-sql',
      'shell',
      'ping',
    ])
  })

  it('does not mutate the input array', () => {
    const input = [...commands]
    sortData(input, { prop: 'name', flip: false })
    expect(input).toEqual(commands)
  })

  it('keeps the original order for equal values', () => {
    const sorted = sortData(commands, { prop: 'created_at', flip: false })
    expect(sorted.map((c) => c.id)).toEqual(commands.map((c) => c.id))
  })

  it('treats missing data as an empty list', () => {
    expect(
      sortData(undefined as unknown as [], { prop: 'id', flip: false }),
    ).toEqual([])
  })
})

describe('formatDateWithTimeZone', () => {
  it('formats a unix timestamp in the given zone', () => {
    expect(formatDateWithTimeZone(FIXTURE_CREATED_AT, 'UTC')).toBe(
      '2026-01-15 10:30:00',
    )
    expect(formatDateWithTimeZone(FIXTURE_CREATED_AT, 'Asia/Kolkata')).toBe(
      '2026-01-15 16:00:00',
    )
  })
})

describe('getVersion', () => {
  afterEach(() => {
    document.cookie = 'heimdall-version=; expires=Thu, 01 Jan 1970 00:00:00 GMT'
  })

  it('reads the heimdall-version cookie', () => {
    document.cookie = 'other=1'
    document.cookie = 'heimdall-version=1.4.2'
    expect(getVersion()).toBe('1.4.2')
  })

  it('returns null when the cookie is missing', () => {
    expect(getVersion()).toBeNull()
  })
})
