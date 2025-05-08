'use client'

import dynamic from 'next/dynamic'

const BreadcrumbProvider = dynamic(() =>
  import('./context').then((module) => ({
    default: module.BreadcrumbProvider,
  })),
)

export function BreadcrumbProvidercontainer({
  children,
}: {
  children: React.ReactNode
}) {
  return <BreadcrumbProvider>{children}</BreadcrumbProvider>
}
