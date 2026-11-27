import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import type { Cluster } from '@/common/api/types'
import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { resourceStatusVariant } from '@/common/Services/status'
import RowActionButton from '@/components/ListPage/RowActionButton'

/** Sorted client-side, so every column key must be a `Cluster` field. */
export const CLUSTER_COLUMNS: TableColumn[] = [
  { key: 'name', header: 'Name', sortable: true, minWidth: '200px' },
  { key: 'version', header: 'Version', sortable: true },
  { key: 'user', header: 'User', sortable: true },
  { key: 'created_at', header: 'Created At', sortable: true, wrap: 'nowrap' },
  { key: 'updated_at', header: 'Updated At', sortable: true, wrap: 'nowrap' },
  { key: 'status', header: 'Status', sortable: true },
  { key: 'actions', header: '', minWidth: '88px' },
]

export const toClusterRows = (clusters: Cluster[]): TableDataRow[] =>
  clusters.map((cluster) => ({
    name: cluster.name || cluster.id,
    version: orDash(cluster.version),
    user: orDash(cluster.user),
    created_at: cluster.created_at
      ? formatDateWithTimeZone(cluster.created_at, myTimezone)
      : orDash(),
    updated_at: cluster.updated_at
      ? formatDateWithTimeZone(cluster.updated_at, myTimezone)
      : orDash(),
    actions: (
      <RowActionButton
        href={`/clusters/${cluster.id}`}
        qaTestId={`cluster-view-${cluster.id}`}
      >
        View
      </RowActionButton>
    ),
    _cellProps: {
      status: {
        tag: (
          <Tag variant={resourceStatusVariant(cluster.status)}>
            {cluster.status}
          </Tag>
        ),
      },
    },
    _qaTestId: `cluster-row-${cluster.id}`,
  }))
