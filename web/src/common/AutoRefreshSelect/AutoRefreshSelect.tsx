'use client'

import { Select } from '@patterninc/react-ui'
import { useContext, useState } from 'react'
import { AutoRefreshContext } from '../AutoRefreshProvider/context'


export type RefreshInterval = {
  label: string
  value: number
}

const refreshOptions = [
  { label: 'OFF', value: 0 },
  { label: '5 seconds', value: 5000 },
  { label: '15 seconds', value: 15000 },
  { label: '30 seconds', value: 30000 },
  { label: '1 minute', value: 60000 },
  { label: '5 minutes', value: 300000 },
]

export const AutoRefreshSelect = () => {
  const { refreshInterval, updateRefreshInterval } =
    useContext(AutoRefreshContext)
  const [intervalValue, setIntervalValue] = useState<RefreshInterval>(
    refreshInterval || { label: 'OFF', value: 0 },
  )

  return (
    <div className='flex flex-row gap-2 items-center'>
      <span>Auto Refresh:</span>
      <div className='w-[104px]'>{/* Fixed width ensures the popover menu inherits the same width, preventing options from being cut off or extending beyond the window when selecting longer options like "OFF" */}
        <Select
          options={refreshOptions}
          optionKeyName={'label'}
          labelKeyName={'label'}
          selectedItem={intervalValue}
          onChange={(option) => {
            setIntervalValue(option)
            updateRefreshInterval(option)
          }}
        />
      </div>
    </div>
  )
}
