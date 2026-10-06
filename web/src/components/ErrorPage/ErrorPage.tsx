'use client'

import { EmptyState } from '@patterninc/pattern-ui/empty-state'
import { useRouter } from 'next/navigation'
import React from 'react'

const NotFoundPage = (): React.JSX.Element => {
  const router = useRouter()

  return (
    <div className='flex flex-1 items-center justify-center p-4'>
      <EmptyState
        title="Oops! We couldn't find that page."
        subtitle="It looks like the page you're looking for doesn't exist or has been moved. Please check the URL."
        action={{
          children: 'Return to the dashboard',
          onClick: () => router.push('/jobs'),
        }}
        qaTestId='not-found'
      />
    </div>
  )
}

export default NotFoundPage
