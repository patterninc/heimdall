'use client'

import React, { createContext, useCallback, useEffect, useState } from 'react'

type RefreshInterval = {
  label: string
  value: number
}

type ContextStateType = {
  refreshInterval: RefreshInterval
  updateRefreshInterval: (interval: RefreshInterval) => void
}

export const AutoRefreshContext = createContext<ContextStateType>({
  refreshInterval: { label: 'OFF', value: 0 }, // Default 0 seconds
  updateRefreshInterval: () => {},
})

const { Provider } = AutoRefreshContext

export const AutoRefreshProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [refreshInterval, setRefreshInterval] = useState<RefreshInterval>(
    () => {
      if (typeof window !== 'undefined') {
        const savedInterval = localStorage.getItem('refreshInterval')
        try {
          return savedInterval
            ? JSON.parse(savedInterval)
            : { label: 'OFF', value: 0 }
        } catch (error) {
          console.error(
            'Error parsing refresh interval from localStorage',
            error,
          )
          return 0
        }
      }
      return { label: 'OFF', value: 0 }
    },
  )

  const updateRefreshInterval = useCallback((interval: RefreshInterval) => {
    setRefreshInterval(interval)
  }, [])

  useEffect(() => {
    localStorage.setItem('refreshInterval', JSON.stringify(refreshInterval))
  }, [refreshInterval])

  return (
    <Provider
      value={{
        refreshInterval,
        updateRefreshInterval,
      }}
    >
      {children}
    </Provider>
  )
}
