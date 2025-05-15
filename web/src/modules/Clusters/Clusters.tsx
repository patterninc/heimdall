'use client'

import { getClusters, getClusterStatus } from '@/app/api/clusters/clusters'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { useQuery } from '@tanstack/react-query'
import { useCallback, useContext, useEffect, useMemo, useState } from 'react'
import {
  ApiParams,
  ClusterType,
  FilterType,
  useClusterConfig,
} from '@/modules/Clusters/Helper'
import {
  ConfigItemType,
  SortByProps,
  SortColumnProps,
  StandardTable,
} from '@patterninc/react-ui'
import {
  noDataAvailable,
  noDataAvailableDescription,
  sortData,
} from '@/common/Services'
import { useQueryState } from 'nuqs'
import { FilterStatesType } from '@patterninc/react-ui/dist/components/Filter/FilterMenu'

const Cluster = (): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)

  const [clusterId, setClusterId] = useQueryState('id', { defaultValue: '' })
  const [clusterName, setClusterName] = useQueryState('name', {
    defaultValue: '',
  })
  const [user, setUser] = useQueryState('user', { defaultValue: '' })
  const [version, setVersion] = useQueryState('version', { defaultValue: '' })
  const [status, setStatus] = useQueryState<string[]>('status', {
    defaultValue: [],
    parse: (value) => (value ? value.split(',') : []),
    serialize: (value) => value?.join(',') ?? '',
  })

  const [filter, setFilter] = useState<FilterType>({
    id: clusterId,
    name: clusterName,
    user: user,
    version: version,
    status: status,
  })

  useEffect(() => {
    updateBreadcrumbs({
      name: 'Clusters',
      link: '/clusters',
      changeType: 'rootLevel',
    })
  }, [updateBreadcrumbs])

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

  const filterParams = useMemo(() => {
    const params: ApiParams = {}
    if (clusterId) params.id = clusterId
    if (user) params.username = user
    if (clusterName) params.name = clusterName
    if (version) params.version = version
    if (status.length > 0) params.status = status
    return params
  }, [clusterId, user, clusterName, version, status])

  const { data, isLoading, isSuccess } = useQuery({
    queryKey: ['clusters', filterParams],
    queryFn: () => getClusters(filterParams),
  })

  const { data: statusData } = useQuery({
    queryKey: ['clusterStatuses'],
    queryFn: getClusterStatus,
  })

  const sortedData: ClusterType[] = useMemo(() => {
    return sortData(data, sortBy)
  }, [data, sortBy])

  const updateFilter = useCallback(() => {
    const queryParams = new URLSearchParams()
    setClusterId(filter?.id)
    setClusterName(filter?.name)
    setUser(filter?.user)
    setVersion(filter?.version)
    setStatus(filter?.status)

    if (filter?.id) queryParams.append('id', filter.id)
    if (filter?.name) queryParams.append('name', filter.name)
    if (filter?.version) queryParams.append('version', filter.version)
    if (filter?.user) queryParams.append('user', filter.user)
    if (filter?.status && filter.status.length > 0) {
      queryParams.append('status', filter.status.join(','))
    }
    updateBreadcrumbs({
      name: 'Clusters',
      link: `/clusters?${queryParams.toString()}`,
      changeType: 'rootLevel',
    })
  }, [
    setClusterId,
    filter.id,
    filter.name,
    filter.user,
    filter.version,
    filter.status,
    setClusterName,
    setUser,
    setVersion,
    setStatus,
    updateBreadcrumbs,
  ])

  const filters: FilterStatesType<unknown> = useMemo(
    () => ({
      id: {
        type: 'text',
        defaultValue: filter?.id,
        placeholder: 'Enter ID',
        stateName: 'id',
        labelText: 'Cluster ID',
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
      status: {
        type: 'multi-select',
        formLabelProps: { label: 'Status' },
        options: statusData?.map((s: string) => ({
          status: s,
        })),
        selectPlaceholder: '--Select Status--',
        labelKey: 'status',
        selectedOptions: filter?.status?.map((s, index) => ({
          status: s,
          id: index,
        })),
        stateName: 'status',
      },
    }),
    [
      filter?.id,
      filter?.name,
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
    return [clusterId, clusterName, user, version, status?.length > 0].filter(
      Boolean,
    ).length
  }, [clusterId, clusterName, user, version, status?.length])

  const resetFilters = useCallback(() => {
    setClusterId(null)
    setClusterName(null)
    setUser(null)
    setVersion(null)
    setStatus(null)
    setFilter({
      id: '',
      name: '',
      user: '',
      version: '',
      status: [],
    })
  }, [setClusterId, setClusterName, setUser, setVersion, setStatus])

  return (
    <div className='pat-pt-4'>
      <StandardTable
        data={sortedData || []}
        dataKey={'id'}
        loading={isLoading}
        successStatus={isSuccess}
        tableId='clusters'
        config={
          useClusterConfig({ sortBy }) as ConfigItemType<
            unknown,
            Record<string, unknown>
          >[]
        }
        hasData={(data && data.length > 0) || false}
        sort={setSortBy}
        sortBy={sortBy}
        noDataFields={{
          primaryText: noDataAvailable,
          secondaryText: noDataAvailableDescription,
        }}
        tableHeaderProps={{
          header: {
            name: 'Results',
            value: data?.length > 100 ? '100+' : data?.length,
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
export default Cluster
