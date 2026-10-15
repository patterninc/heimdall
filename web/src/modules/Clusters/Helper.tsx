import type { TableColumn, TableDataRow } from '@patterninc/pattern-ui/table'
import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

import { formatDateWithTimeZone, myTimezone, orDash } from '@/common/Services'
import { resourceStatusVariant } from '@/common/Services/status'
import RowLink from '@/components/ListPage/RowLink'

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

/** Sorted client-side, so every column key must be a `ClusterType` field. */
export const CLUSTER_COLUMNS: TableColumn[] = [
  { key: 'name', header: 'Name', sortable: true, minWidth: '200px' },
  { key: 'version', header: 'Version', sortable: true },
  { key: 'user', header: 'User', sortable: true },
  { key: 'created_at', header: 'Created At', sortable: true, wrap: 'nowrap' },
  { key: 'updated_at', header: 'Updated At', sortable: true, wrap: 'nowrap' },
  { key: 'status', header: 'Status', sortable: true },
]

export const toClusterRows = (clusters: ClusterType[]): TableDataRow[] =>
  clusters.map((cluster) => ({
    name: (
      <RowLink href={`/clusters/${cluster.id}`}>
        {cluster.name || cluster.id}
      </RowLink>
    ),
    version: orDash(cluster.version),
    user: orDash(cluster.user),
    created_at: cluster.created_at
      ? formatDateWithTimeZone(cluster.created_at, myTimezone)
      : orDash(),
    updated_at: cluster.updated_at
      ? formatDateWithTimeZone(cluster.updated_at, myTimezone)
      : orDash(),
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
