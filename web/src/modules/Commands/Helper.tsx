import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { resourceStatusVariant } from '@/common/Services/status'
import RowLink from '@/components/ListPage/RowLink'
import TruncatedText from '@/components/TruncatedText/TruncatedText'

export type ApiParams = {
  id?: string
  username?: string
  name?: string
  plugin?: string
  version?: string
  status?: string[]
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

/** Sorted client-side, so every column key must be a `CommandType` field. */
export const COMMAND_COLUMNS: TableColumn[] = [
  { key: 'name', header: 'Name', sortable: true, minWidth: '200px' },
  { key: 'version', header: 'Version', sortable: true },
  { key: 'user', header: 'User', sortable: true },
  {
    key: 'description',
    header: 'Description',
    sortable: true,
    minWidth: '240px',
  },
  { key: 'created_at', header: 'Created At', sortable: true, wrap: 'nowrap' },
  { key: 'updated_at', header: 'Updated At', sortable: true, wrap: 'nowrap' },
  { key: 'plugin', header: 'Plugin', sortable: true },
  { key: 'status', header: 'Status', sortable: true },
]

export const toCommandRows = (commands: CommandType[]): TableDataRow[] =>
  commands.map((command) => ({
    name: (
      <RowLink href={`/commands/${command.id}`}>
        {command.name || command.id}
      </RowLink>
    ),
    version: orDash(command.version),
    user: orDash(command.user),
    description: command.description ? (
      <TruncatedText text={command.description} />
    ) : (
      orDash()
    ),
    created_at: command.created_at
      ? formatDateWithTimeZone(command.created_at, myTimezone)
      : orDash(),
    updated_at: command.updated_at
      ? formatDateWithTimeZone(command.updated_at, myTimezone)
      : orDash(),
    plugin: orDash(command.plugin),
    _cellProps: {
      status: {
        tag: (
          <Tag variant={resourceStatusVariant(command.status)}>
            {command.status}
          </Tag>
        ),
      },
    },
    _qaTestId: `command-row-${command.id}`,
  }))
