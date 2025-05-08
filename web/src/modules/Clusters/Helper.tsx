import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import {
  Button,
  ConfigItemType,
  MdashCheck,
  SortByProps,
  Tag,
} from '@patterninc/react-ui'
import Link from 'next/link'
import { useMemo } from 'react'

export type ClusterContext = {
  emr_release_label: string
  execution_role_arn: string
  properties: Record<string, string>
  role_arn: string
}

export type ClusterType = {
  id: string
  name: string
  version: string
  user: string
  description: string
  tags: string[]
  context: ClusterContext
  created_at: number
  updated_at: number
  status: string
}

export type ApiParams = {
  id?: string
  username?: string
  name?: string
  version?: string
  status?: string[]
}

export type FilterType = {
  id: string
  name: string
  user: string
  version: string
  status: string[]
}

type CommandConfigProps = {
  sortBy: SortByProps
}

export const useClusterConfig = ({
  sortBy,
}: CommandConfigProps): ConfigItemType<
  ClusterType,
  Record<string, unknown>
>[] => {
  return useMemo(
    () => [
      {
        name: 'name',
        label: 'Name',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div className={sortBy.prop === 'name' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.name}</MdashCheck>
              </div>
            )
          },
        },
        mainColumn: true,
      },
      {
        name: 'version',
        label: 'Version',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div className={sortBy.prop === 'version' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.version}</MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'user',
        label: 'User',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div className={sortBy.prop === 'user' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.user}</MdashCheck>
              </div>
            )
          },
        },
      },

      {
        name: 'created_at',
        label: 'Created At',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div
                className={sortBy.prop === 'created_at' ? 'fw-semi-bold' : ''}
              >
                <MdashCheck check={!!row.id}>
                  {formatDateWithTimeZone(row.created_at, myTimezone)}
                </MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'updated_at',
        label: 'Updated At',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div
                className={sortBy.prop === 'updated_at' ? 'fw-semi-bold' : ''}
              >
                <MdashCheck check={!!row.id}>
                  {formatDateWithTimeZone(row.updated_at, myTimezone)}
                </MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'status',
        label: 'Status',
        cell: {
          children: (row: ClusterType) => {
            return (
              <div className={sortBy.prop === 'status' ? 'fw-semi-bold' : ''}>
                <Tag
                  color={
                    row.status === 'ACTIVE'
                      ? 'green'
                      : row.status === 'INACTIVE'
                        ? 'gray'
                        : 'red'
                  }
                >
                  {row.status}
                </Tag>
              </div>
            )
          },
        },
      },
      {
        name: '',
        label: '',
        isButton: true,
        noSort: true,
        cell: {
          children: (row: ClusterType) => {
            return (
              <Button
                as='link'
                routerComponent={Link}
                href={`/clusters/${row.id}`}
              >
                Details
              </Button>
            )
          },
        },
      },
    ],
    [sortBy],
  )
}
