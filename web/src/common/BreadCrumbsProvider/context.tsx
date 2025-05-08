'use client'

import React, { createContext, useCallback, useEffect, useState } from 'react'

import {
  addNewBreadcrumb,
  breadcrumbNavigation,
  BreadcrumbType,
} from '@patterninc/react-ui'

type ContextStateType = {
  breadcrumbs: Array<BreadcrumbType>
  updateBreadcrumbs: (breadcrumb: BreadcrumbType) => void
  breadcrumbCallout: (breadcrumb: BreadcrumbType) => void
}

export const BreadcrumbContext = createContext<ContextStateType>({
  breadcrumbs: [
    {
      name: '',
      link: '',
    },
  ],
  updateBreadcrumbs: () => null,
  breadcrumbCallout: () => null,
})

const { Provider } = BreadcrumbContext

export const BreadcrumbProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbType[]>(() => {
    if (typeof window !== 'undefined') {
      const savedBreadcrumbs = localStorage.getItem('breadcrumbs')
      try {
        return savedBreadcrumbs ? JSON.parse(savedBreadcrumbs) : []
      } catch (error) {
        console.error('Error parsing breadcrumbs from localStorage', error)
        return []
      }
    }
    return []
  })

  const updateBreadcrumbs = useCallback((breadcrumb: BreadcrumbType) => {
    setBreadcrumbs((prevState) =>
      addNewBreadcrumb({
        breadcrumb,
        breadcrumbs: prevState,
      }),
    )
  }, [])

  useEffect(() => {
    localStorage.setItem('breadcrumbs', JSON.stringify(breadcrumbs))
  }, [breadcrumbs])

  const breadcrumbCallout = (breadcrumb: BreadcrumbType) => {
    setBreadcrumbs((prevState) =>
      breadcrumbNavigation({
        breadcrumb,
        breadcrumbs: prevState,
      }),
    )
  }

  return (
    <Provider
      value={{
        breadcrumbs,
        updateBreadcrumbs,
        breadcrumbCallout,
      }}
    >
      {children}
    </Provider>
  )
}
