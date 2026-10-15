import { Button } from '@patterninc/pattern-ui/button'
import { NoDataAvailable } from '@patterninc/pattern-ui/no-data-available'
import React from 'react'

import { noDataAvailable, noDataAvailableDescription } from '@/common/Services'

type ListEmptyStateProps = {
  /** Shown only when filters are narrowing the result set. */
  onClearFilters?: () => void
}

const ListEmptyState = ({
  onClearFilters,
}: ListEmptyStateProps): React.JSX.Element => (
  <div className='flex flex-col items-center gap-3 py-10 text-center'>
    <div className='flex flex-col items-center gap-1'>
      <NoDataAvailable text={noDataAvailable} />
      <p className='pui-text-body-medium-alt text-text-sub'>
        {noDataAvailableDescription}
      </p>
    </div>
    {onClearFilters ? (
      <Button
        variant='weak'
        size='small'
        onClick={onClearFilters}
        qaTestId='empty-state-clear-filters'
      >
        Clear filters
      </Button>
    ) : null}
  </div>
)

export default ListEmptyState
