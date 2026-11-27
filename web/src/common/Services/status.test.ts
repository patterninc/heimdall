import { describe, expect, it } from 'vitest'

import { jobStatusVariant, resourceStatusVariant } from './status'

describe('jobStatusVariant', () => {
  it.each([
    ['SUCCEEDED', 'positive'],
    ['FAILED', 'negative'],
    ['RUNNING', 'caution'],
    ['KILLED', 'stroke'],
    ['NEW', 'primary'],
    ['ACCEPTED', 'notice'],
    ['CANCELING', 'weak'],
    [undefined, 'weak'],
  ])('maps %s to %s', (status, variant) => {
    expect(jobStatusVariant(status)).toBe(variant)
  })
})

describe('resourceStatusVariant', () => {
  it.each([
    ['ACTIVE', 'positive'],
    ['INACTIVE', 'weak'],
    ['DELETED', 'negative'],
    [undefined, 'negative'],
  ])('maps %s to %s', (status, variant) => {
    expect(resourceStatusVariant(status)).toBe(variant)
  })
})
