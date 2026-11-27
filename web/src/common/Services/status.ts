import type { TagVariant } from '@patterninc/pattern-ui/tag'

import type { JobStatus } from '@/common/api/types'

const JOB_STATUS_VARIANTS: Partial<Record<JobStatus, TagVariant>> = {
  SUCCEEDED: 'positive',
  FAILED: 'negative',
  RUNNING: 'caution',
  KILLED: 'stroke',
  NEW: 'primary',
  ACCEPTED: 'notice',
}

export const jobStatusVariant = (status?: string): TagVariant =>
  (status && JOB_STATUS_VARIANTS[status as JobStatus]) || 'weak'

/** Commands and clusters share the ACTIVE / INACTIVE / other lifecycle. */
export const resourceStatusVariant = (status?: string): TagVariant => {
  if (status === 'ACTIVE') return 'positive'
  if (status === 'INACTIVE') return 'weak'
  return 'negative'
}
