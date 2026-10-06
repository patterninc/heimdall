'use client'

import dynamic from 'next/dynamic'
import { ReactNode } from 'react'

const AppRoot = dynamic(() => import('./AppRoot'), { ssr: false })

const ClientLayout = ({ children }: { children: ReactNode }) => {
  return <AppRoot>{children}</AppRoot>
}

export default ClientLayout
