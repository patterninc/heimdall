import React, { useMemo } from 'react'
import {
  ConfigItemType,
  Button,
  MdashCheck,
  Tag,
  TagProps,
  SortByProps,
} from '@patterninc/react-ui'
import Link from 'next/link'
import { formatDateWithTimeZone, myTimezone } from '@/common/Services'

export type ApiParams = {
  id?: string
  username?: string
  name?: string
  version?: string
  cluster?: string
  command?: string
  status?: string[]
  tags?: string
  page?: string
}

// A single tag filter row, e.g. { key: 'owner', value: '#dev-core' }.
export type TagPair = {
  key: string
  value: string
}

// Join tag pairs into the backend's `tags` param format: `key:value,key:value`.
// Blank rows are dropped; the full `key:value` string is matched exactly server-side.
export const serializeTags = (pairs: TagPair[]): string =>
  pairs
    .filter((p) => p.key.trim() && p.value.trim())
    .map((p) => `${p.key.trim()}:${p.value.trim()}`)
    .join(',')

// Parse a `key:value,key:value` string back into tag pairs. Splits each tag on
// its first colon so values may themselves contain colons.
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

export type JobDataTypesProps = {
  jobData?: JobType
  isLoading: boolean
}

type JobConfigProps = {
  sortBy: SortByProps
}

export const getStatusColor = (status: string): TagProps['color'] => {
  const statusColors: Record<string, TagProps['color']> = {
    SUCCEEDED: 'green',
    FAILED: 'red',
    RUNNING: 'yellow',
    KILLED: 'dark-gray',
    NEW: 'blue',
    ACCEPTED: 'orange',
  }

  return statusColors[status] || 'gray' // Default to gray if status is unknown
}

export const useJobConfig = ({
  sortBy,
}: JobConfigProps): ConfigItemType<JobType, Record<string, unknown>>[] => {
  return useMemo(
    () => [
      {
        name: 'id',
        label: 'Job ID',

        cell: {
          children: (row: JobType) => {
            return (
              <div className={sortBy.prop === 'id' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.id}</MdashCheck>
              </div>
            )
          },
        },
        mainColumn: true,
      },
      {
        name: 'name',
        label: 'Name',
        cell: {
          children: (row: JobType) => {
            return (
              <div className={sortBy.prop === 'name' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.name}>{row.name}</MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'version',
        label: 'Version',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'version' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.version}>{row.version}</MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'user',
        label: 'User',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'user' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.user}>{row.user}</MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'cluster_id',
        label: 'Cluster ID',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'cluster_id' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.cluster_id}>{row.cluster_id}</MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'command_id',
        label: 'Command ID',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'command_id' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.command_id}>{row.command_id}</MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'created_at',
        label: 'Created At',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'created_at' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.created_at}>
                {formatDateWithTimeZone(row.created_at, myTimezone)}
              </MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'updated_at',
        label: 'Updated At',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'updated_at' ? 'fw-semi-bold' : ''}>
              <MdashCheck check={!!row.updated_at}>
                {formatDateWithTimeZone(row.updated_at, myTimezone)}
              </MdashCheck>
            </div>
          ),
        },
      },
      {
        name: 'status',
        label: 'Status',
        cell: {
          children: (row: JobType) => (
            <div className={sortBy.prop === 'status' ? 'fw-semi-bold' : ''}>
              <Tag color={getStatusColor(row?.status)}>{row.status}</Tag>
            </div>
          ),
        },
      },
      {
        isButton: true,
        name: '',
        label: '',
        noSort: true,
        cell: {
          children: (row: JobType) => (
            <Button as='link' routerComponent={Link} href={`/jobs/${row.id}`}>
              Details
            </Button>
          ),
        },
      },
    ],
    [sortBy],
  )
}
