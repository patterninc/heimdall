'use client'

import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { useQueryState } from 'nuqs'
import React, { useContext, useState } from 'react'
import { useDebounceValue } from 'usehooks-ts'

import { getCommands, getCommandStatus } from '@/app/api/commands/commands'
import { AutoRefreshContext } from '@/common/AutoRefreshProvider/context'
import { useClientPagination } from '@/common/hooks/useClientPagination'
import { SortBy, sortData, toggleSort } from '@/common/Services'
import { FilterMenu } from '@/components/FilterMenu/FilterMenu'
import type { FilterDimension } from '@/components/FilterMenu/filterTypes'
import ListEmptyState from '@/components/ListPage/ListEmptyState'
import ListPage from '@/components/ListPage/ListPage'
import {
  ApiParams,
  COMMAND_COLUMNS,
  CommandType,
  toCommandRows,
} from './Helper'

const FILTER_DEBOUNCE_MS = 300

/** The API caps command results, so a full page means "at least this many". */
const COMMAND_RESULT_CAP = 100

const isSameJson = (a: unknown, b: unknown) =>
  JSON.stringify(a) === JSON.stringify(b)

const Commands = (): React.JSX.Element => {
  const router = useRouter()
  const { refreshInterval } = useContext(AutoRefreshContext)

  const [commandId, setCommandId] = useQueryState('id', { defaultValue: '' })
  const [commandName, setCommandName] = useQueryState('name', {
    defaultValue: '',
  })
  const [user, setUser] = useQueryState('user', { defaultValue: '' })
  const [plugin, setPlugin] = useQueryState('plugin', { defaultValue: '' })
  const [version, setVersion] = useQueryState('version', {
    defaultValue: '',
  })
  const [status, setStatus] = useQueryState<string[]>('status', {
    defaultValue: [],
    parse: (value) => (value ? value.split(',') : []),
    serialize: (value) => value?.join(',') ?? '',
  })

  const filterParams: ApiParams = {}
  if (commandId) filterParams.id = commandId
  if (user) filterParams.username = user
  if (commandName) filterParams.name = commandName
  if (plugin) filterParams.plugin = plugin
  if (version) filterParams.version = version
  if (status.length > 0) filterParams.status = status

  const [debouncedParams] = useDebounceValue(filterParams, FILTER_DEBOUNCE_MS, {
    equalityFn: isSameJson,
  })

  const { data, isPending, isPlaceholderData } = useQuery<CommandType[]>({
    queryKey: ['commands', debouncedParams],
    queryFn: () => getCommands(debouncedParams),
    refetchInterval: refreshInterval.value,
    placeholderData: (prev) => prev,
  })

  const { data: statusData } = useQuery<string[]>({
    queryKey: ['commandStatuses'],
    queryFn: getCommandStatus,
  })

  const [sortBy, setSortBy] = useState<SortBy>({ prop: 'name', flip: false })
  const sortedCommands = sortData(data ?? [], sortBy)
  const { pageRows, paginationProps } = useClientPagination(
    sortedCommands,
    JSON.stringify([debouncedParams, sortBy]),
  )

  const resultCount =
    isPending || isPlaceholderData
      ? undefined
      : sortedCommands.length > COMMAND_RESULT_CAP
        ? `${COMMAND_RESULT_CAP}+`
        : String(sortedCommands.length)

  const clearAll = () => {
    setCommandId(null)
    setCommandName(null)
    setUser(null)
    setPlugin(null)
    setVersion(null)
    setStatus(null)
  }

  const dimensions: FilterDimension[] = [
    {
      key: 'id',
      label: 'Command ID',
      type: 'text',
      value: commandId,
      onChange: setCommandId,
    },
    {
      key: 'name',
      label: 'Name',
      type: 'text',
      value: commandName,
      onChange: setCommandName,
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
      key: 'plugin',
      label: 'Plugin',
      type: 'text',
      value: plugin,
      onChange: setPlugin,
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
      title='Commands'
      resultCount={resultCount}
      qaTestId='commands-page'
      filters={<FilterMenu dimensions={dimensions} onClearAll={clearAll} />}
      tableProps={{
        columns: COMMAND_COLUMNS,
        rows: toCommandRows(pageRows),
        isLoading: isPending || isPlaceholderData,
        sortedColumn: sortBy.prop,
        sortDirection: sortBy.flip ? 'desc' : 'asc',
        onSort: (columnKey) => setSortBy(toggleSort(sortBy, columnKey)),
        onRowClick: (rowIndex) => {
          const command = pageRows[rowIndex]
          if (command) router.push(`/commands/${command.id}`)
        },
        stickyColumns: { start: 1 },
        emptyState: (
          <ListEmptyState onClearFilters={hasFilters ? clearAll : undefined} />
        ),
        paginationProps,
        qaTestId: 'commands-table',
      }}
    />
  )
}

export default Commands
