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
  page?: string
}

export type JobType = {
  id: string
  name: string
  version: string
  user: string
  tags: string[]
  created_at: number
  updated_at: number
  status: string
  command_criteria: string[]
  cluster_criteria: string[]
  command_id: string
  command_name: string
  cluster_id: string
  cluster_name: string
  cancelled_by: string
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
