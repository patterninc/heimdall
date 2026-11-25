import { describe, expect, it } from 'vitest'

import { parseTags, serializeTags } from './Helper'

describe('serializeTags', () => {
  it('joins trimmed key:value pairs and drops incomplete ones', () => {
    expect(
      serializeTags([
        { key: ' team ', value: ' data ' },
        { key: 'env', value: '' },
        { key: '', value: 'x' },
        { key: 'env', value: 'prod' },
      ]),
    ).toBe('team:data,env:prod')
  })
})

describe('parseTags', () => {
  it('splits on the first colon so values may contain colons', () => {
    expect(parseTags('team:data, url:http://x:1')).toEqual([
      { key: 'team', value: 'data' },
      { key: 'url', value: 'http://x:1' },
    ])
  })

  it('ignores entries without a key', () => {
    expect(parseTags(':nokey,novalue,ok:1')).toEqual([
      { key: 'ok', value: '1' },
    ])
  })

  it('round-trips with serializeTags', () => {
    const value = 'team:data,env:prod'
    expect(serializeTags(parseTags(value))).toBe(value)
  })
})
