'use client'

import { ContentHeader } from '@patterninc/pattern-ui/content-header'
import { Table, type TableProps } from '@patterninc/pattern-ui/table'
import { useViewportTableHeight } from '@patterninc/pattern-ui/use-viewport-table-height'
import React, { type ReactNode } from 'react'

import { AutoRefreshSelect } from '@/common/AutoRefreshSelect/AutoRefreshSelect'

type ListPageProps = {
  title: string
  /** Pre-formatted count (e.g. `"20+"`); omitted while the first page loads. */
  resultCount?: string
  /** Rendered in a sticky, right-aligned row between the header and the table. */
  filters: ReactNode
  tableProps: TableProps
  qaTestId: string
}

const ListPage = ({
  title,
  resultCount,
  filters,
  tableProps,
  qaTestId,
}: ListPageProps): React.JSX.Element => {
  // Shell inset (8px) + page bottom padding (16px), so the table body scrolls instead of the page.
  const { ref: tableWrapperRef, height: tableHeight } = useViewportTableHeight({
    bottomOffset: 24,
  })

  return (
    <div className='flex min-h-0 flex-1 flex-col' qa-test-id={qaTestId}>
      <ContentHeader
        variant='fixed'
        className='top-14 md:top-0'
        title={title}
        subtitle={
          resultCount !== undefined ? `${resultCount} results` : undefined
        }
        right={<AutoRefreshSelect />}
        qaTestId={`${qaTestId}-header`}
      />
      <div className='bg-base-inverse sticky top-30 z-20 flex justify-end px-4 pb-3 md:top-16'>
        {filters}
      </div>
      <div className='flex min-h-0 flex-1 flex-col px-4 pb-4'>
        {/* A flex column, so the table shrinks to the measured cap and scrolls internally (keeping
            its header sticky), while short lists still collapse to their own height. */}
        <div
          ref={tableWrapperRef}
          className='flex min-h-0 flex-col'
          style={{ maxHeight: tableHeight }}
        >
          <Table
            rounded
            rowSize='medium'
            headerBackground='filled'
            className='min-h-0'
            {...tableProps}
          />
        </div>
      </div>
    </div>
  )
}

export default ListPage
