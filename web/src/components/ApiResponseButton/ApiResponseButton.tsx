'use client'

import { CodeBlock } from '@patterninc/pattern-ui/artifact/CodeBlock'
import { Button } from '@patterninc/pattern-ui/button'
import { EmptyState } from '@patterninc/pattern-ui/empty-state'
import CodeIcon from '@patterninc/pattern-ui/icons/CodeIcon'
import { Modal } from '@patterninc/pattern-ui/modal'
import { NoDataAvailable } from '@patterninc/pattern-ui/no-data-available'
import { Skeleton } from '@patterninc/pattern-ui/skeleton'
import { useQuery } from '@tanstack/react-query'
import React, { useState } from 'react'

import { useIsSmUp } from '@/common/hooks/useIsSmUp'
import { noDataAvailable } from '@/common/Services'

type ApiResponseButtonProps = {
  /** Same-origin API path, e.g. `/api/v1/job/<id>`. */
  href: string
  qaTestId: string
}

const fetchResponse = async (href: string): Promise<unknown> => {
  const response = await fetch(href)
  if (!response.ok) throw new Error(`Request failed (${response.status})`)
  return response.json()
}

const isEmptyResponse = (data: unknown): boolean =>
  data === null ||
  data === undefined ||
  (Array.isArray(data) && data.length === 0) ||
  (typeof data === 'object' &&
    !Array.isArray(data) &&
    Object.keys(data).length === 0)

/** Weak button that shows the raw API response for a resource in a modal CodeBlock. */
const ApiResponseButton = ({
  href,
  qaTestId,
}: ApiResponseButtonProps): React.JSX.Element => {
  const [open, setOpen] = useState(false)
  const isSmUp = useIsSmUp()
  const { data, isPending, isError } = useQuery({
    queryKey: ['api-response', href],
    queryFn: () => fetchResponse(href),
    enabled: open,
  })

  return (
    <>
      <Button
        variant='weak'
        size='medium'
        // Icon-only on phones, where the header has no room for labels.
        prefix={isSmUp ? undefined : { icon: CodeIcon }}
        aria-label='API response'
        tooltipProps={isSmUp ? undefined : { content: 'API response' }}
        onClick={() => setOpen(true)}
        qaTestId={qaTestId}
      >
        {isSmUp ? 'API response' : undefined}
      </Button>
      <Modal
        open={open}
        onOpenChange={setOpen}
        size='lg'
        headerProps={{ title: 'API response', description: href }}
        qaTestId={`${qaTestId}-modal`}
      >
        <div className='p-4'>
          {isPending ? (
            <Skeleton className='h-48 w-full' />
          ) : isError ? (
            <EmptyState
              title="We couldn't load the API response"
              subtitle='The server may be unavailable. Close this and try again.'
            />
          ) : isEmptyResponse(data) ? (
            <div className='flex flex-col items-center gap-1 py-10 text-center'>
              <NoDataAvailable text={noDataAvailable} />
              <p className='pui-text-body-medium-alt text-text-sub'>
                The API returned an empty response for this request.
              </p>
            </div>
          ) : (
            <CodeBlock
              code={JSON.stringify(data, null, 2)}
              language='json'
              filename='response.json'
              className='[&_code]:wrap-break-word! [&_code]:whitespace-pre-wrap! [&_pre]:whitespace-pre-wrap!'
              qaTestId={`${qaTestId}-code`}
            />
          )}
        </div>
      </Modal>
    </>
  )
}

export default ApiResponseButton
