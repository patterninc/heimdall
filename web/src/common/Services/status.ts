import type { TagVariant } from '@patterninc/pattern-ui/tag'

const JOB_STATUS_VARIANTS: Record<string, TagVariant> = {
  SUCCEEDED: 'positive',
  FAILED: 'negative',
  RUNNING: 'caution',
  KILLED: 'stroke',
  NEW: 'primary',
  ACCEPTED: 'notice',
}

export const jobStatusVariant = (status?: string): TagVariant =>
  (status && JOB_STATUS_VARIANTS[status]) || 'weak'

/** Commands and clusters share the ACTIVE / INACTIVE / other lifecycle. */
export const resourceStatusVariant = (status?: string): TagVariant => {
  if (status === 'ACTIVE') return 'positive'
  if (status === 'INACTIVE') return 'weak'
  return 'negative'
}
