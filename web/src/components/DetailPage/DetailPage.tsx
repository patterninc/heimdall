'use client'

import { ContentHeader } from '@patterninc/pattern-ui/content-header'
import { EmptyState } from '@patterninc/pattern-ui/empty-state'
import ArrowLeftIcon from '@patterninc/pattern-ui/icons/ArrowLeftIcon'
import { Tag, type TagVariant } from '@patterninc/pattern-ui/tag'
import { useRouter } from 'next/navigation'
import React, { type ReactNode } from 'react'

type DetailPageProps = {
  /** List route the back button returns to. */
  backHref: string
  backLabel: string
  /** Singular resource noun for the error state, e.g. `"job"`. */
  resourceLabel: string
  title?: string
  subtitle?: string
  status?: string
  statusVariant?: TagVariant
  right?: ReactNode
  isError?: boolean
  children: ReactNode
  qaTestId: string
}

const DetailPage = ({
  backHref,
  backLabel,
  resourceLabel,
  title,
  subtitle,
  status,
  statusVariant,
  right,
  isError,
  children,
  qaTestId,
}: DetailPageProps): React.JSX.Element => {
  const router = useRouter()

  return (
    <div className='flex min-h-0 flex-1 flex-col' qa-test-id={qaTestId}>
      <ContentHeader
        variant='fixed'
        className='top-14 md:top-0'
        left={{
          buttonProps: {
            variant: 'weak',
            size: 'medium',
            prefix: { icon: ArrowLeftIcon },
            'aria-label': backLabel,
            tooltipProps: { content: backLabel },
            onClick: () => router.push(backHref),
            qaTestId: `${qaTestId}-back`,
          },
        }}
        title={title}
        subtitle={subtitle}
        titleSuffix={
          status ? (
            <Tag variant={statusVariant} qaTestId={`${qaTestId}-status`}>
              {status}
            </Tag>
          ) : undefined
        }
        right={
          isError ? undefined : (
            <div className='flex items-center gap-2'>{right}</div>
          )
        }
        qaTestId={`${qaTestId}-header`}
      />
      <div className='flex flex-col gap-4 px-4 pb-6'>
        {isError ? (
          <EmptyState
            title={`We couldn't load this ${resourceLabel}`}
            subtitle='It may have been removed, or the server is unavailable.'
            action={{
              children: backLabel,
              onClick: () => router.push(backHref),
            }}
            qaTestId={`${qaTestId}-error`}
          />
        ) : (
          children
        )}
      </div>
    </div>
  )
}

export default DetailPage
