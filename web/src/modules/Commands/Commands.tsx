'use client'

import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { useQuery } from '@tanstack/react-query'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { getCommands, getCommandStatus } from '@/app/api/commands/commands'
import { ApiParams, useCommandConfig, CommandType, FilterType } from './Helper' // Adjust the import path as needed
import {
  StandardTable,
  ConfigItemType,
  SortColumnProps,
  SortByProps,
} from '@patterninc/react-ui'
import { useQueryState } from 'nuqs'
import { FilterStatesType } from '@patterninc/react-ui/dist/components/Filter/FilterMenu'
import {
  noDataAvailable,
  noDataAvailableDescription,
  sortData,
} from '@/common/Services'
import { AutoRefreshContext } from '@/common/AutoRefreshProvider/context'

const Commands = (): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)
  const { refreshInterval } = useContext(AutoRefreshContext)
  const [commandId, setCommandId] = useQueryState('id', { defaultValue: '' })
  const [commandName, setCommandName] = useQueryState('name', {
    defaultValue: '',
  })
  const [user, setUser] = useQueryState('user', { defaultValue: '' })
  const [plugin, setPlugin] = useQueryState('plugin', { defaultValue: '' })
  const [version, setVersion] = useQueryState('version', { defaultValue: '' })
  const [status, setStatus] = useQueryState<string[]>('status', {
    defaultValue: [],
    parse: (value) => (value ? value.split(',') : []),
    serialize: (value) => value?.join(',') ?? '',
  })
  const [filter, setFilter] = useState<FilterType>({
    id: commandId,
    name: commandName,
    user: user,
    plugin: plugin,
    version: version,
    status: status,
  })

  useEffect(() => {
    updateBreadcrumbs({
      name: 'Commands',
      link: '/commands',
      changeType: 'rootLevel',
    })
  }, [updateBreadcrumbs])

  const filterParams = useMemo(() => {
    const params: ApiParams = {}
    if (commandId) params.id = commandId
    if (user) params.username = user
    if (commandName) params.name = commandName
    if (plugin) params.plugin = plugin
    if (version) params.version = version
    if (status.length > 0) params.status = status
    return params
  }, [commandId, commandName, user, status, plugin, version])

  const { data, isSuccess, isLoading } = useQuery({
    queryKey: ['commands', filterParams],
    queryFn: () => getCommands(filterParams),
    refetchInterval: refreshInterval.value,
  })

  const { data: statusData } = useQuery({
    queryKey: ['commandStatuses'],
    queryFn: getCommandStatus,
  })

  const updateFilter = useCallback(() => {
    const queryParams = new URLSearchParams()
    setCommandId(filter?.id)
    setCommandName(filter?.name)
    setUser(filter?.user)
    setPlugin(filter?.plugin)
    setVersion(filter?.version)
    setStatus(filter?.status)

    if (filter?.id) queryParams.append('id', filter.id)
    if (filter?.name) queryParams.append('name', filter.name)
    if (filter?.version) queryParams.append('version', filter.version)
    if (filter?.user) queryParams.append('user', filter.user)
    if (filter?.plugin) queryParams.append('plugin', filter.plugin)
    if (filter?.status && filter.status.length > 0) {
      queryParams.append('status', filter.status.join(','))
    }
    updateBreadcrumbs({
      name: 'Commands',
      link: `/commands?${queryParams.toString()}`,
      changeType: 'rootLevel',
    })
  }, [
    filter.id,
    filter.name,
    filter.plugin,
    filter.status,
    filter.user,
    filter.version,
    setCommandId,
    setCommandName,
    setPlugin,
    setStatus,
    setUser,
    setVersion,
    updateBreadcrumbs,
  ])

  const filters: FilterStatesType<unknown> = useMemo(
    () => ({
      id: {
        type: 'text',
        defaultValue: filter?.id,
        placeholder: 'Enter ID',
        stateName: 'id',
        labelText: 'Command ID',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },

      name: {
        type: 'text',
        defaultValue: filter?.name,
        placeholder: 'Enter Name',
        stateName: 'name',
        labelText: 'Name',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      user: {
        type: 'text',
        defaultValue: filter?.user,
        placeholder: 'Enter User',
        stateName: 'user',
        labelText: 'User',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      version: {
        type: 'text',
        defaultValue: filter?.version,
        placeholder: 'Enter Version',
        stateName: 'version',
        labelText: 'Version',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      plugin: {
        type: 'text',
        defaultValue: filter?.plugin,
        placeholder: 'Enter Plugin',
        stateName: 'plugin',
        labelText: 'Plugin',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      status: {
        type: 'multi-select',
        formLabelProps: { label: 'Status' },
        options: statusData
          ? statusData?.map((status: string) => ({
              status: status,
              key: `status-${status}`,
            }))
          : [],
        selectPlaceholder: '--Select Status--',
        labelKey: 'status',
        selectedOptions: filter?.status.map((s) => ({
          status: s,
          key: `selected-${s}`,
        })),
        stateName: 'status',
      },
    }),
    [
      filter?.id,
      filter?.name,
      filter?.plugin,
      filter?.status,
      filter?.user,
      filter?.version,
      statusData,
      updateFilter,
    ],
  )

  const updateFormField = useCallback(
    (...params: unknown[]) => {
      const stateAttr = params[0] as keyof FilterType
      const value = params[1] as string | string[]
      setFilter((prevFilter) => ({
        ...prevFilter,
        [stateAttr]: value,
      }))
    },
    [setFilter],
  )

  const filterCount = useMemo(() => {
    return [
      commandId,
      commandName,
      user,
      version,
      plugin,
      status?.length > 0,
    ].filter(Boolean).length
  }, [commandId, commandName, user, version, plugin, status])

  const resetFilters = useCallback(() => {
    setCommandId(null)
    setCommandName(null)
    setUser(null)
    setPlugin(null)
    setVersion(null)
    setStatus(null)
    setFilter({
      id: '',
      name: '',
      user: '',
      plugin: '',
      version: '',
      status: [],
    })
  }, [setCommandId, setCommandName, setUser, setPlugin, setVersion, setStatus])

  const [sortBy, setSort] = useState<SortByProps>({
    prop: 'name',
    flip: false,
  })

  const setSortBy: SortColumnProps['sorter'] = (obj: {
    activeColumn: string
    direction: boolean
    lowerCaseParams?: boolean
  }) => {
    setSort({
      prop: obj.activeColumn,
      flip: obj.direction,
    })
  }

  const sortedData: CommandType[] = useMemo(() => {
    return sortData(data, sortBy)
  }, [data, sortBy])

  return (
    <div className='pat-pt-4'>
      <StandardTable
        data={sortedData || []}
        dataKey={'id'}
        loading={isLoading}
        successStatus={isSuccess}
        tableId='commands'
        config={
          useCommandConfig({ sortBy }) as ConfigItemType<
            unknown,
            Record<string, unknown>
          >[]
        }
        hasData={(data && data.length > 0) || false}
        stickyTableConfig={{ right: 1 }}
        sort={setSortBy}
        sortBy={sortBy}
        noDataFields={{
          primaryText: noDataAvailable,
          secondaryText: noDataAvailableDescription,
        }}
        tableHeaderProps={{
          header: {
            name: 'Results',
            value: data?.length,
          },
          pageFilterProps: {
            filterStates: filters,
            filterCallout: updateFilter,
            appliedFilters: filterCount,
            resetCallout: resetFilters,
            cancelCallout: () => {},
            onChangeCallout: updateFormField,
          },
        }}
      />
    </div>
  )
}

export default Commands
