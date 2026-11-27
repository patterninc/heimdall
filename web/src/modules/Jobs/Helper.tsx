import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import type { Job } from '@/common/api/types'
import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { jobStatusVariant } from '@/common/Services/status'
import RowActionButton from '@/components/ListPage/RowActionButton'
import TruncatedText from '@/components/TruncatedText/TruncatedText'

export const JOBS_PAGE_SIZE = 20

export type TagPair = {
  key: string
  value: string
}

export const serializeTags = (pairs: TagPair[]): string =>
  pairs
    .filter((p) => p.key.trim() && p.value.trim())
    .map((p) => `${p.key.trim()}:${p.value.trim()}`)
    .join(',')

export const parseTags = (value: string): TagPair[] =>
  value
    .split(',')
    .map((raw) => raw.trim())
    .reduce<TagPair[]>((acc, tag) => {
      const idx = tag.indexOf(':')
      if (idx > 0) {
        acc.push({
          key: tag.slice(0, idx).trim(),
          value: tag.slice(idx + 1).trim(),
        })
      }
      return acc
    }, [])

/** Column keys double as the server's `order_by` values for sortable columns. */
export const JOB_COLUMNS: TableColumn[] = [
  { key: 'id', header: 'Job ID', sortable: true, wrap: 'nowrap' },
  { key: 'name', header: 'Name', minWidth: '200px' },
  { key: 'version', header: 'Version' },
  { key: 'user', header: 'User' },
  { key: 'cluster_id', header: 'Cluster ID', wrap: 'nowrap' },
  { key: 'command_id', header: 'Command ID', wrap: 'nowrap' },
  { key: 'created_at', header: 'Created At', sortable: true, wrap: 'nowrap' },
  { key: 'updated_at', header: 'Updated At', sortable: true, wrap: 'nowrap' },
  { key: 'status', header: 'Status' },
  { key: 'actions', header: '', minWidth: '88px' },
]

export const toJobRows = (jobs: Job[]): TableDataRow[] =>
  jobs.map((job) => ({
    id: job.id,
    name: job.name ? (
      <TruncatedText text={job.name} className='whitespace-nowrap' />
    ) : (
      orDash(job.name)
    ),
    version: orDash(job.version),
    user: orDash(job.user),
    cluster_id: orDash(job.cluster_id),
    command_id: orDash(job.command_id),
    created_at: job.created_at
      ? formatDateWithTimeZone(job.created_at, myTimezone)
      : orDash(),
    updated_at: job.updated_at
      ? formatDateWithTimeZone(job.updated_at, myTimezone)
      : orDash(),
    actions: (
      <RowActionButton href={`/jobs/${job.id}`} qaTestId={`job-view-${job.id}`}>
        View
      </RowActionButton>
    ),
    _cellProps: {
      status: {
        tag: <Tag variant={jobStatusVariant(job.status)}>{job.status}</Tag>,
      },
    },
    _qaTestId: `job-row-${job.id}`,
  }))
