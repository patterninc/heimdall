'use client'

import React, {
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'
import {
  Button,
  SortByProps,
  SortColumnProps,
  StandardTable,
} from '@patterninc/react-ui'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { fetchJobs, getJobStatus } from '@/app/api/jobs/jobs'
import { useQuery } from '@tanstack/react-query'
import { ApiParams, JOBS_PAGE_SIZE, useJobConfig } from './Helper'
import { useQueryState } from 'nuqs'
import { FilterStatesType } from '@patterninc/react-ui'
import { noDataAvailable, noDataAvailableDescription } from '@/common/Services'
import { AutoRefreshContext } from '@/common/AutoRefreshProvider/context'

type FilterType = {
  id: string
  name: string
  user: string
  version: string
  clusterId: string
  commandId: string
  status: string[]
}

const Jobs = (): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)
  const { refreshInterval } = useContext(AutoRefreshContext)

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

  const [pageIndex, setPageIndex] = useState(0)
  const [cursors, setCursors] = useState<(string | null)[]>([null])

  // Sort state drives the server-side ORDER BY. `flip` is the descending flag.
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
    setPageIndex(0)
    setCursors([null])
  }

  // `placeholderData` keeps the previous page visible while the next one loads.
  const { data, isLoading, isFetching, isSuccess } = useQuery({
    queryKey: ['jobs', filterParams, sortBy, pageIndex, cursors[pageIndex]],
    queryFn: () => fetchJobs(filterParams, cursors[pageIndex] ?? null, sortBy),
    enabled: !!filterParams,
    refetchInterval: refreshInterval.value,
    placeholderData: (prev) => prev,
  })

  const { data: jobStatus } = useQuery({
    queryKey: ['jobs'],
    queryFn: getJobStatus,
  })

  const jobs = useMemo(() => data?.data ?? [], [data])

  const rowsSoFar = pageIndex * JOBS_PAGE_SIZE + jobs.length
  const totalResults = data?.has_more ? `${rowsSoFar}+` : `${rowsSoFar}`

  const goToNextPage = useCallback(() => {
    const next = data?.next_cursor
    if (!next) return
    setCursors((prev) => {
      const copy = [...prev]
      copy[pageIndex + 1] = next
      return copy
    })
    setPageIndex((i) => i + 1)
  }, [data?.next_cursor, pageIndex])

  const goToPrevPage = useCallback(() => {
    setPageIndex((i) => Math.max(0, i - 1))
  }, [])

  const updateFilter = useCallback(() => {
    const queryParams = new URLSearchParams()
    setPageIndex(0)
    setCursors([null])
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
    // eslint-disable-next-line react-hooks/preserve-manual-memoization
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
        options:
          jobStatus?.map((status: string) => ({
            status: status,
            key: `status-${status}`,
          })) || [],
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
    setPageIndex(0)
    setCursors([null])
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
    <div className='pt-4'>
      <StandardTable
        data={jobs}
        config={useJobConfig({ sortBy })}
        stickyTableConfig={{ right: 1 }}
        dataKey={'id'}
        hasData={jobs.length > 0}
        successStatus={isSuccess}
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
            value: totalResults,
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
      <div className='flex items-center justify-end gap-2 pt-3'>
        <Button
          styleType='secondary'
          onClick={goToPrevPage}
          disabled={pageIndex === 0 || isFetching}
        >
          Previous
        </Button>
        <span className='fw-semi-bold'>Page {pageIndex + 1}</span>
        <Button
          styleType='secondary'
          onClick={goToNextPage}
          disabled={!data?.next_cursor || isFetching}
        >
          Next
        </Button>
      </div>
    </div>
  )
}

export default Jobs
