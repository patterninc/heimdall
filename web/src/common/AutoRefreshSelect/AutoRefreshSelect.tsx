'use client'

import RefreshIcon from '@patterninc/pattern-ui/icons/RefreshIcon'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
} from '@patterninc/pattern-ui/select'
import { useContext } from 'react'

import { AutoRefreshContext } from '../AutoRefreshProvider/context'

const refreshOptions = [
  { label: 'Auto-refresh off', value: 0 },
  { label: 'Every 5 seconds', value: 5000 },
  { label: 'Every 15 seconds', value: 15000 },
  { label: 'Every 30 seconds', value: 30000 },
  { label: 'Every minute', value: 60000 },
  { label: 'Every 5 minutes', value: 300000 },
]

export const AutoRefreshSelect = () => {
  const { refreshInterval, updateRefreshInterval } =
    useContext(AutoRefreshContext)

  return (
    <Select
      value={String(refreshInterval.value)}
      onValueChange={(value) => {
        const option = refreshOptions.find((o) => String(o.value) === value)
        if (option) updateRefreshInterval(option)
      }}
      qaTestId='auto-refresh-select'
    >
      <SelectTrigger
        width={212}
        prefix={{ icon: RefreshIcon, color: 'icon-sub' }}
        aria-label='Auto-refresh interval'
      />
      <SelectContent position='bottom-end'>
        {refreshOptions.map((option) => (
          <SelectItem key={option.value} value={String(option.value)}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
