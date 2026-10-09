import { Skeleton } from '@patterninc/pattern-ui/skeleton'
import React, { type ReactNode } from 'react'

export type DetailField = {
  label: string
  /** Empty values (`undefined`, `null`, `''`) render as an em dash. */
  value?: ReactNode
  mono?: boolean
}

type DetailFieldsProps = {
  fields: DetailField[]
  isLoading?: boolean
}

const isEmpty = (value: ReactNode) =>
  value === undefined || value === null || value === ''

const DetailFields = ({
  fields,
  isLoading,
}: DetailFieldsProps): React.JSX.Element => (
  <dl className='grid grid-cols-1 gap-x-6 gap-y-4 p-4 sm:grid-cols-2 xl:grid-cols-3'>
    {fields.map(({ label, value, mono }) => (
      <div key={label} className='flex min-w-0 flex-col gap-1'>
        <dt className='pui-text-label-small text-text-sub'>{label}</dt>
        <dd
          className={`${mono ? 'pui-text-mono-small' : 'pui-text-body-medium'} text-text-strong min-w-0 break-words`}
        >
          {isLoading ? (
            <Skeleton className='h-4 w-32' />
          ) : isEmpty(value) ? (
            '—'
          ) : (
            value
          )}
        </dd>
      </div>
    ))}
  </dl>
)

export default DetailFields
