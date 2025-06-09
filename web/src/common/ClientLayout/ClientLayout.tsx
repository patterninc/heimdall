'use client'

import dynamic from 'next/dynamic'
import { ReactNode } from 'react'

const ReactQueryProvider = dynamic(
  () => import('../ReactQueryProvider/ReactQueryProvider'),
  { ssr: false },
)

const HeaderContent = dynamic(() => import('../../components/Header/Header'), {
  ssr: false,
})

const BreadcrumbProvidercontainer = dynamic(
  () =>
    import('../BreadCrumbsProvider/BreadCrumbsProvideContainer').then(
      (mod) => mod.BreadcrumbProvidercontainer,
    ),
  { ssr: false },
)

const AutoRefreshProvidercontainer = dynamic(
  () =>
    import('../AutoRefreshProvider/AutoRefreshProviderContainer').then(
      (mod) => mod.AutoRefreshProvidercontainer,
    ),
  { ssr: false },
)

const LeftNavContainer = dynamic(
  () => import('../../components/NavigationBar/LeftNavContainer'),
  { ssr: false },
)

const ClientLayout = ({
  children,
  user,
}: {
  children: ReactNode
  user: string | null
}) => {
  return (
    <ReactQueryProvider>
      <div className='App'>
        <BreadcrumbProvidercontainer>
          <AutoRefreshProvidercontainer>
            <LeftNavContainer user={user} />
            <main className='app-content-layout'>
              <HeaderContent />
              <div className='App-content'>{children}</div>
            </main>
          </AutoRefreshProvidercontainer>
        </BreadcrumbProvidercontainer>
      </div>
    </ReactQueryProvider>
  )
}

export default ClientLayout
