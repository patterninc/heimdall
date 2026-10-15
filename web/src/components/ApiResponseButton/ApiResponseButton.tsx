'use client'

import { CodeBlock } from '@patterninc/pattern-ui/artifact/CodeBlock'
import { Button } from '@patterninc/pattern-ui/button'
import { EmptyState } from '@patterninc/pattern-ui/empty-state'
import CodeIcon from '@patterninc/pattern-ui/icons/CodeIcon'
import { Modal } from '@patterninc/pattern-ui/modal'
import { Skeleton } from '@patterninc/pattern-ui/skeleton'
import { useQuery } from '@tanstack/react-query'
import React, { useState } from 'react'

import { useIsSmUp } from '@/common/hooks/useIsSmUp'

type ApiResponseButtonProps = {
  /** Same-origin API path, e.g. `/api/v1/job/<id>`. */
  href: string
  qaTestId: string
}

const fetchRaw = async (href: string): Promise<string> => {
  const response = await fetch(href)
  if (!response.ok) throw new Error(`Request failed (${response.status})`)
  return JSON.stringify(await response.json(), null, 2)
}

/** Weak button that shows the raw API response for a resource in a modal CodeBlock. */
const ApiResponseButton = ({
  href,
  qaTestId,
}: ApiResponseButtonProps): React.JSX.Element => {
  const [open, setOpen] = useState(false)
  const isSmUp = useIsSmUp()
  const { data, isPending, isError } = useQuery({
    queryKey: ['api-response', href],
    queryFn: () => fetchRaw(href),
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
          {isError ? (
            <EmptyState
              title="We couldn't load the API response"
              subtitle='The server may be unavailable. Close this and try again.'
            />
          ) : isPending || data === undefined ? (
            <Skeleton className='h-48 w-full' />
          ) : (
            <CodeBlock
              code={data}
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
