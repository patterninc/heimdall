'use client'

import { AutoRefreshProvider } from './context'

export const AutoRefreshProvidercontainer: React.FC<{
  children: React.ReactNode
}> = ({ children }) => {
  return <AutoRefreshProvider>{children}</AutoRefreshProvider>
}
