import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import type { Command } from '@/common/api/types'
import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { resourceStatusVariant } from '@/common/Services/status'
import RowActionButton from '@/components/ListPage/RowActionButton'
import TruncatedText from '@/components/TruncatedText/TruncatedText'

/** Sorted client-side, so every column key must be a `Command` field. */
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
  { key: 'actions', header: '', minWidth: '88px' },
]

export const toCommandRows = (commands: Command[]): TableDataRow[] =>
  commands.map((command) => ({
    name: command.name || command.id,
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
    actions: (
      <RowActionButton
        href={`/commands/${command.id}`}
        qaTestId={`command-view-${command.id}`}
      >
        View
      </RowActionButton>
    ),
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
