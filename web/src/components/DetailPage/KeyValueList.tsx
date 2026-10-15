import React from 'react'

type KeyValueListProps = {
  entries: [string, string][]
}

/** Config-style key/value pairs (e.g. Spark properties), monospaced so values align and copy cleanly. */
const KeyValueList = ({ entries }: KeyValueListProps): React.JSX.Element => (
  <dl className='flex flex-col'>
    {entries.map(([key, value], index) => (
      <div
        key={key}
        className={`flex flex-col gap-1 px-4 py-3 sm:flex-row sm:items-baseline sm:gap-4 ${index > 0 ? 'border-border-soft border-t' : ''}`}
      >
        <dt className='pui-text-mono-small text-text-sub min-w-0 sm:w-72 sm:shrink-0'>
          {key}
        </dt>
        <dd className='pui-text-mono-small text-text-strong min-w-0 break-all'>
          {value}
        </dd>
      </div>
    ))}
  </dl>
)

export default KeyValueList
