import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { jobStatusVariant } from '@/common/Services/status'
import RowLink from '@/components/ListPage/RowLink'
import TruncatedText from '@/components/TruncatedText/TruncatedText'

export type ApiParams = {
  id?: string
  username?: string
  name?: string
  version?: string
  cluster?: string
  command?: string
  status?: string[]
  tags?: string
  limit?: string
  cursor?: string
  order_by?: string
  direction?: string
}

export const JOBS_PAGE_SIZE = 20

export type JobsResponse = {
  data: JobType[]
  has_more?: boolean
  next_cursor?: string
}

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

export type JobType = {
  id: string
  name: string
  version: string
  user: string
  tags: string[]
  created_at: number
  updated_at: number
  status: string
  is_sync: boolean
  command_criteria: string[]
  cluster_criteria: string[]
  command_id: string
  command_name: string
  cluster_id: string
  cluster_name: string
  canceled_by?: string
  job_attributes?: Record<string, { kind?: string; value?: string }>
  error?: string
  context?: {
    properties: {
      'spark.driver.cores': string
      'spark.driver.memory': string
      'spark.executor.cores': string
      'spark.executor.instances': string
      'spark.executor.memory': string
    }
    query: string
    return_result?: boolean
  }
}

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
]

export const toJobRows = (jobs: JobType[]): TableDataRow[] =>
  jobs.map((job) => ({
    id: <RowLink href={`/jobs/${job.id}`}>{job.id}</RowLink>,
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
    _cellProps: {
      status: {
        tag: <Tag variant={jobStatusVariant(job.status)}>{job.status}</Tag>,
      },
    },
    _qaTestId: `job-row-${job.id}`,
  }))
