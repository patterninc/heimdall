'use client'

import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'
import {
  SortByProps,
  SortColumnProps,
  StandardTable,
} from '@patterninc/react-ui'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { fetchJobs, getJobStatus } from '@/app/api/jobs/jobs'
import { useInfiniteQuery, useQuery } from '@tanstack/react-query'
import { ApiParams, JobType, useJobConfig } from './Helper'
import { useQueryState } from 'nuqs'
import { FilterStatesType } from '@patterninc/react-ui/dist/components/Filter/FilterMenu'
import {
  noDataAvailable,
  noDataAvailableDescription,
  sortData,
} from '@/common/Services'

type FilterType = {
  id: string
  name: string
  user: string
  version: string
  clusterId: string
  commandId: string
  status: string[]
}

const Jobs = (): JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)

  const [jobId, setJobId] = useQueryState('id', { defaultValue: '' })
  const [name, setName] = useQueryState('name', { defaultValue: '' })
  const [user, setUser] = useQueryState('user', { defaultValue: '' })
  const [version, setVersion] = useQueryState('version', {
    defaultValue: '',
  })
  const [clusterId, setClusterId] = useQueryState('clusterId', {
    defaultValue: '',
  })
  const [commandId, setCommandId] = useQueryState('commandId', {
    defaultValue: '',
  })
  const [status, setStatus] = useQueryState<string[]>('status', {
    defaultValue: [],
    parse: (value) => (value ? value.split(',') : []),
    serialize: (value) => value?.join(',') ?? '',
  })

  const [filter, setFilter] = useState<FilterType>({
    id: jobId,
    name: name,
    user: user,
    version: version,
    clusterId: clusterId,
    commandId: commandId,
    status: status,
  })

  useEffect(() => {
    updateBreadcrumbs({ name: 'Jobs', link: '/jobs', changeType: 'rootLevel' })
  }, [updateBreadcrumbs])

  // Build API params from query state
  const filterParams = useMemo(() => {
    const params: ApiParams = {}
    if (jobId) params.id = jobId
    if (name) params.name = name
    if (user) params.username = user
    if (version) params.version = version
    if (clusterId) params.cluster = clusterId
    if (commandId) params.command = commandId
    if (status.length > 0) params.status = status
    return params
  }, [jobId, name, user, version, clusterId, commandId, status])

  // Fetch Jobs with filters applied
  const { data, isLoading, fetchNextPage, hasNextPage, isSuccess } =
    useInfiniteQuery({
      queryKey: ['jobs', filterParams],
      queryFn: ({ pageParam }) => fetchJobs(filterParams, pageParam),
      getNextPageParam: (lastPage) => lastPage?.nextPage,
      enabled: !!filterParams,
      initialPageParam: 1,
    })

  const { data: jobStatus } = useQuery({
    queryKey: ['jobs'],
    queryFn: getJobStatus,
  })

  const [sortBy, setSort] = useState<SortByProps>({
    prop: 'created_at',
    flip: true,
  })

  const setSortBy: SortColumnProps['sorter'] = (obj: {
    activeColumn: string
    direction: boolean
  }) => {
    setSort({
      prop: obj.activeColumn,
      flip: obj.direction,
    })
  }

  const jobs = useMemo(() => data?.pages?.flatMap((page) => page) || [], [data])

  const sortedData: JobType[] = useMemo(() => {
    return sortData(jobs, sortBy)
  }, [jobs, sortBy])

  const updateFilter = useCallback(() => {
    const queryParams = new URLSearchParams()
    setJobId(filter.id)
    setName(filter.name)
    setUser(filter.user)
    setVersion(filter.version)
    setClusterId(filter.clusterId)
    setCommandId(filter.commandId)
    setStatus(filter.status)

    if (filter?.user) queryParams.append('user', filter.user)
    if (filter?.name) queryParams.append('name', filter.name)
    if (filter?.id) queryParams.append('id', filter.id)
    if (filter?.version) queryParams.append('version', filter.version)
    if (filter?.clusterId) queryParams.append('cluster_id', filter.clusterId)
    if (filter?.commandId) queryParams.append('command_id', filter.commandId)
    if (filter?.status && filter.status.length > 0) {
      queryParams.append('status', filter.status.join(','))
    }
    updateBreadcrumbs({
      name: 'Jobs',
      link: `/jobs?${queryParams.toString()}`,
      changeType: 'rootLevel',
    })
  }, [
    filter.clusterId,
    filter.commandId,
    filter.id,
    filter.name,
    filter.status,
    filter.user,
    filter.version,
    setClusterId,
    setCommandId,
    setJobId,
    setName,
    setStatus,
    setUser,
    setVersion,
    updateBreadcrumbs,
  ])

  const filters: FilterStatesType<unknown> = useMemo(
    () => ({
      id: {
        type: 'text',
        defaultValue: filter.id,
        placeholder: 'Enter Job ID',
        stateName: 'id',
        labelText: 'Job ID',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      name: {
        type: 'text',
        defaultValue: filter.name,
        placeholder: 'Enter Name',
        stateName: 'name',
        labelText: 'Name',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      user: {
        type: 'text',
        defaultValue: filter.user,
        placeholder: 'Enter User',
        stateName: 'user',
        labelText: 'User',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      version: {
        type: 'text',
        defaultValue: filter.version,
        placeholder: 'Enter Version',
        stateName: 'version',
        labelText: 'Version',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      clusterId: {
        type: 'text',
        defaultValue: filter.clusterId,
        placeholder: 'Enter Cluster ID',
        stateName: 'clusterId',
        labelText: 'Cluster ID',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      commandId: {
        type: 'text',
        defaultValue: filter.commandId,
        placeholder: 'Enter Command ID',
        stateName: 'commandId',
        labelText: 'Command ID',
        inputType: 'text',
        onReturnCallout: updateFilter,
        shouldClose: true,
      },
      status: {
        type: 'multi-select',
        formLabelProps: { label: 'Status' },
        options: jobStatus?.map((status: string) => ({ status: status })) || [],
        selectPlaceholder: '--Select Status--',
        labelKey: 'status',
        selectedOptions: filter?.status.map((s) => ({ status: s })),
        stateName: 'status',
      },
    }),
    [
      filter.id,
      filter.name,
      filter.user,
      filter.version,
      filter.clusterId,
      filter.commandId,
      filter?.status,
      updateFilter,
      jobStatus,
    ],
  )

  const resetFilters = useCallback(() => {
    setJobId(null)
    setName(null)
    setUser(null)
    setVersion(null)
    setClusterId(null)
    setCommandId(null)
    setStatus(null)
    setFilter({
      id: '',
      name: '',
      user: '',
      version: '',
      clusterId: '',
      commandId: '',
      status: [],
    })
  }, [
    setJobId,
    setName,
    setUser,
    setVersion,
    setClusterId,
    setCommandId,
    setStatus,
  ])

  // Compute applied filter count
  const filterCount = useMemo(() => {
    return [
      jobId,
      name,
      user,
      version,
      commandId,
      clusterId,
      status.length > 0,
    ].filter(Boolean).length
  }, [jobId, name, user, version, commandId, clusterId, status.length])

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

  return (
    <div className='pat-pt-4'>
      <StandardTable
        data={sortedData ?? []}
        config={useJobConfig({ sortBy })}
        dataKey={'id'}
        hasData={jobs.length > 0}
        successStatus={isSuccess}
        getData={fetchNextPage}
        hasMore={!!hasNextPage}
        loading={isLoading}
        tableId='tableId'
        sort={setSortBy}
        sortBy={sortBy}
        noDataFields={{
          primaryText: noDataAvailable,
          secondaryText: noDataAvailableDescription,
        }}
        tableHeaderProps={{
          header: {
            name: 'Results',
            value: jobs?.length > 100 ? '100+' : jobs?.length,
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

export default Jobs
