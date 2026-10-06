'use client'

import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { useQueryState } from 'nuqs'
import React, { useContext, useState } from 'react'
import { useDebounceValue } from 'usehooks-ts'

import { getClusters, getClusterStatus } from '@/app/api/clusters/clusters'
import { AutoRefreshContext } from '@/common/AutoRefreshProvider/context'
import { useClientPagination } from '@/common/hooks/useClientPagination'
import { SortBy, sortData, toggleSort } from '@/common/Services'
import { FilterMenu } from '@/components/FilterMenu/FilterMenu'
import type { FilterDimension } from '@/components/FilterMenu/filterTypes'
import ListEmptyState from '@/components/ListPage/ListEmptyState'
import ListPage from '@/components/ListPage/ListPage'
import {
  ApiParams,
  CLUSTER_COLUMNS,
  ClusterType,
  toClusterRows,
} from './Helper'

const FILTER_DEBOUNCE_MS = 300

/** The API caps cluster results, so a full page means "at least this many". */
const CLUSTER_RESULT_CAP = 100

const isSameJson = (a: unknown, b: unknown) =>
  JSON.stringify(a) === JSON.stringify(b)

const Clusters = (): React.JSX.Element => {
  const router = useRouter()
  const { refreshInterval } = useContext(AutoRefreshContext)

  const [clusterId, setClusterId] = useQueryState('id', { defaultValue: '' })
  const [clusterName, setClusterName] = useQueryState('name', {
    defaultValue: '',
  })
  const [user, setUser] = useQueryState('user', { defaultValue: '' })
  const [version, setVersion] = useQueryState('version', {
    defaultValue: '',
  })
  const [status, setStatus] = useQueryState<string[]>('status', {
    defaultValue: [],
    parse: (value) => (value ? value.split(',') : []),
    serialize: (value) => value?.join(',') ?? '',
  })

  const filterParams: ApiParams = {}
  if (clusterId) filterParams.id = clusterId
  if (user) filterParams.username = user
  if (clusterName) filterParams.name = clusterName
  if (version) filterParams.version = version
  if (status.length > 0) filterParams.status = status

  const [debouncedParams] = useDebounceValue(filterParams, FILTER_DEBOUNCE_MS, {
    equalityFn: isSameJson,
  })

  const { data, isPending } = useQuery<ClusterType[]>({
    queryKey: ['clusters', debouncedParams],
    queryFn: () => getClusters(debouncedParams),
    refetchInterval: refreshInterval.value,
    placeholderData: (prev) => prev,
  })

  const { data: statusData } = useQuery<string[]>({
    queryKey: ['clusterStatuses'],
    queryFn: getClusterStatus,
  })

  const [sortBy, setSortBy] = useState<SortBy>({ prop: 'name', flip: false })
  const sortedClusters = sortData(data ?? [], sortBy)
  const { pageRows, paginationProps } = useClientPagination(
    sortedClusters,
    JSON.stringify([debouncedParams, sortBy]),
  )

  const total = sortedClusters.length
  const resultCount = isPending
    ? undefined
    : total > CLUSTER_RESULT_CAP
      ? `${CLUSTER_RESULT_CAP}+`
      : String(total)

  const clearAll = () => {
    setClusterId(null)
    setClusterName(null)
    setUser(null)
    setVersion(null)
    setStatus(null)
  }

  const dimensions: FilterDimension[] = [
    {
      key: 'id',
      label: 'Cluster ID',
      type: 'text',
      value: clusterId,
      onChange: setClusterId,
    },
    {
      key: 'name',
      label: 'Name',
      type: 'text',
      value: clusterName,
      onChange: setClusterName,
    },
    {
      key: 'user',
      label: 'User',
      type: 'text',
      value: user,
      onChange: setUser,
    },
    {
      key: 'version',
      label: 'Version',
      type: 'text',
      value: version,
      onChange: setVersion,
    },
    {
      key: 'status',
      label: 'Status',
      type: 'multi',
      value: status,
      onChange: setStatus,
      options: statusData ?? [],
    },
  ]
  const hasFilters = Object.keys(filterParams).length > 0

  return (
    <ListPage
      title='Clusters'
      resultCount={resultCount}
      qaTestId='clusters-page'
      filters={<FilterMenu dimensions={dimensions} onClearAll={clearAll} />}
      tableProps={{
        columns: CLUSTER_COLUMNS,
        rows: toClusterRows(pageRows),
        isLoading: isPending,
        sortedColumn: sortBy.prop,
        sortDirection: sortBy.flip ? 'desc' : 'asc',
        onSort: (columnKey) => setSortBy(toggleSort(sortBy, columnKey)),
        onRowClick: (rowIndex) => {
          const cluster = pageRows[rowIndex]
          if (cluster) router.push(`/clusters/${cluster.id}`)
        },
        stickyColumns: { start: 1 },
        emptyState: (
          <ListEmptyState onClearFilters={hasFilters ? clearAll : undefined} />
        ),
        paginationProps,
        qaTestId: 'clusters-table',
      }}
    />
  )
}

export default Clusters
