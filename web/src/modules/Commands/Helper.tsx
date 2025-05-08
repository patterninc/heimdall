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

export type ApiParams = {
  id?: string
  username?: string
  name?: string
  plugin?: string
  version?: string
  status?: string[]
}

export type FilterType = {
  id: string
  user: string
  name: string
  version: string
  plugin: string
  status: string[]
}

export type CommandContext = {
  logs_uri?: string
  properties?: {
    'spark.driver.cores'?: string
    'spark.executor.cores'?: string
    'spark.executor.instances'?: string
    'spark.executor.memory'?: string
  }
  queries_uri?: string
  results_uri?: string
}

export type CommandType = {
  id: string
  name: string
  version: string
  user: string
  description: string
  tags: string[]
  created_at: number
  updated_at: number
  status: string
  plugin: string
  cluster_tags?: string[]
  is_sync?: boolean
  context?: CommandContext
}

type CommandConfigProps = {
  sortBy: SortByProps
}

export const useCommandConfig = ({
  sortBy,
}: CommandConfigProps): ConfigItemType<
  CommandType,
  Record<string, unknown>
>[] => {
  return useMemo(
    () => [
      {
        name: 'name',
        label: 'Name',
        cell: {
          children: (row: CommandType) => {
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
          children: (row: CommandType) => {
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
          children: (row: CommandType) => {
            return (
              <div className={sortBy.prop === 'user' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.user}</MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'description',
        label: 'Description',
        cell: {
          children: (row: CommandType) => {
            return (
              <div
                className={sortBy.prop === 'description' ? 'fw-semi-bold' : ''}
              >
                <MdashCheck check={!!row.id}>{row.description}</MdashCheck>
              </div>
            )
          },
        },
      },

      {
        name: 'created_at',
        label: 'Created At',
        cell: {
          children: (row: CommandType) => {
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
          children: (row: CommandType) => {
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
        name: 'plugin',
        label: 'Plugin',
        cell: {
          children: (row: CommandType) => {
            return (
              <div className={sortBy.prop === 'plugin' ? 'fw-semi-bold' : ''}>
                <MdashCheck check={!!row.id}>{row.plugin}</MdashCheck>
              </div>
            )
          },
        },
      },
      {
        name: 'status',
        label: 'Status',
        cell: {
          children: (row: CommandType) => {
            return (
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
            )
          },
        },
      },
      {
        isButton: true,
        name: '',
        label: '',
        noSort: true,
        cell: {
          children: (row: CommandType) => (
            <Button
              as='link'
              routerComponent={Link}
              href={`/commands/${row.id}`}
            >
              Details
            </Button>
          ),
        },
      },
    ],
    [sortBy],
  )
}
