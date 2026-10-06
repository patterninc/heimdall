'use client'

import { useQuery } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import { useQueryState } from 'nuqs'
import React, { useContext, useState } from 'react'
import { useDebounceValue } from 'usehooks-ts'

import { fetchJobs, getJobStatus } from '@/app/api/jobs/jobs'
import { AutoRefreshContext } from '@/common/AutoRefreshProvider/context'
import { SortBy, toggleSort } from '@/common/Services'
import { FilterMenu } from '@/components/FilterMenu/FilterMenu'
import type { FilterDimension } from '@/components/FilterMenu/filterTypes'
import ListEmptyState from '@/components/ListPage/ListEmptyState'
import ListPage from '@/components/ListPage/ListPage'
import {
  ApiParams,
  JOB_COLUMNS,
  JOBS_PAGE_SIZE,
  TagPair,
  parseTags,
  serializeTags,
  toJobRows,
} from './Helper'
import TagFilter from './TagFilter'

const FILTER_DEBOUNCE_MS = 300

const isSameJson = (a: unknown, b: unknown) =>
  JSON.stringify(a) === JSON.stringify(b)

type CursorPaging = {
  key: string
  index: number
  cursors: (string | null)[]
}

const Jobs = (): React.JSX.Element => {
  const router = useRouter()
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
  const [tags, setTags] = useQueryState<TagPair[]>('tags', {
    defaultValue: [],
    parse: (value) => (value ? parseTags(value) : []),
    serialize: (value) => serializeTags(value),
  })

  const filterParams: ApiParams = {}
  if (jobId) filterParams.id = jobId
  if (name) filterParams.name = name
  if (user) filterParams.username = user
  if (version) filterParams.version = version
  if (clusterId) filterParams.cluster = clusterId
  if (commandId) filterParams.command = commandId
  if (status.length > 0) filterParams.status = status
  const serializedTags = serializeTags(tags)
  if (serializedTags) filterParams.tags = serializedTags

  // Text filters apply as you type; debounce so each keystroke doesn't hit the API.
  const [debouncedParams] = useDebounceValue(filterParams, FILTER_DEBOUNCE_MS, {
    equalityFn: isSameJson,
  })

  // `flip` is the descending flag; the server owns ORDER BY.
  const [sortBy, setSortBy] = useState<SortBy>({
    prop: 'created_at',
    flip: true,
  })

  // Keyset pagination: cursors[i] fetches page i. Any filter or sort change restarts at page 1.
  const pagingKey = JSON.stringify([debouncedParams, sortBy])
  const [paging, setPaging] = useState<CursorPaging>({
    key: pagingKey,
    index: 0,
    cursors: [null],
  })
  const { index: pageIndex, cursors } =
    paging.key === pagingKey ? paging : { index: 0, cursors: [null] }
  const cursor = cursors[pageIndex] ?? null

  const { data, isPending, isPlaceholderData } = useQuery({
    queryKey: ['jobs', debouncedParams, sortBy, cursor],
    queryFn: () => fetchJobs(debouncedParams, cursor, sortBy),
    refetchInterval: refreshInterval.value,
    placeholderData: (prev) => prev,
  })

  const { data: jobStatuses } = useQuery<string[]>({
    queryKey: ['jobStatuses'],
    queryFn: getJobStatus,
  })

  const jobs = data?.data ?? []
  const hasMore = Boolean(data?.has_more && data?.next_cursor)
  const rowsSoFar = pageIndex * JOBS_PAGE_SIZE + jobs.length
  const resultCount = isPending
    ? undefined
    : hasMore
      ? `${rowsSoFar}+`
      : `${rowsSoFar}`

  const goToPage = (page: number) => {
    const target = page - 1
    if (target <= pageIndex) {
      setPaging({ key: pagingKey, index: Math.max(0, target), cursors })
      return
    }
    // Keyset paging can only step forward one page, and only from settled (non-placeholder) data.
    const next = data?.next_cursor
    if (!next || isPlaceholderData) return
    const nextCursors = [...cursors]
    nextCursors[pageIndex + 1] = next
    setPaging({ key: pagingKey, index: pageIndex + 1, cursors: nextCursors })
  }

  const clearAll = () => {
    setJobId(null)
    setName(null)
    setUser(null)
    setVersion(null)
    setClusterId(null)
    setCommandId(null)
    setStatus(null)
    setTags(null)
  }

  const dimensions: FilterDimension[] = [
    {
      key: 'id',
      label: 'Job ID',
      type: 'text',
      value: jobId,
      onChange: setJobId,
    },
    {
      key: 'name',
      label: 'Name',
      type: 'text',
      value: name,
      onChange: setName,
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
      key: 'clusterId',
      label: 'Cluster ID',
      type: 'text',
      value: clusterId,
      onChange: setClusterId,
    },
    {
      key: 'commandId',
      label: 'Command ID',
      type: 'text',
      value: commandId,
      onChange: setCommandId,
    },
    {
      key: 'status',
      label: 'Status',
      type: 'multi',
      value: status,
      onChange: setStatus,
      options: jobStatuses ?? [],
    },
    {
      key: 'tags',
      label: 'Tags',
      type: 'custom',
      isActive: tags.length > 0,
      content: <TagFilter tags={tags} onChange={setTags} />,
    },
  ]
  const hasFilters = Object.keys(filterParams).length > 0

  return (
    <ListPage
      title='Jobs'
      resultCount={resultCount}
      qaTestId='jobs-page'
      filters={<FilterMenu dimensions={dimensions} onClearAll={clearAll} />}
      tableProps={{
        columns: JOB_COLUMNS,
        rows: toJobRows(jobs),
        isLoading: isPending,
        sortedColumn: sortBy.prop,
        sortDirection: sortBy.flip ? 'desc' : 'asc',
        onSort: (columnKey) => setSortBy(toggleSort(sortBy, columnKey)),
        onRowClick: (rowIndex) => {
          const job = jobs[rowIndex]
          if (job) router.push(`/jobs/${job.id}`)
        },
        stickyColumns: { start: 1 },
        emptyState: (
          <ListEmptyState onClearFilters={hasFilters ? clearAll : undefined} />
        ),
        paginationProps: {
          currentPage: pageIndex + 1,
          // Keyset paging has no total; expose exactly one page ahead while more rows exist.
          totalPages: pageIndex + (hasMore ? 2 : 1),
          onPageChange: goToPage,
        },
        qaTestId: 'jobs-table',
      }}
    />
  )
}

export default Jobs
