import LeftNavContainer from '@/components/NavigationBar/LeftNavContainer'

import dynamic from 'next/dynamic'
import { NuqsAdapter } from 'nuqs/adapters/next/app'

const ReactQueryProvider = dynamic(
  () => import('../common/ReactQueryProvider/ReactQueryProvider'),
  { ssr: false },
)

const HeaderContent = dynamic(() => import('../components/Header/Header'), {
  ssr: false,
})

const BreadcrumbProvidercontainer = dynamic(
  () =>
    import('../common/BreadCrumbsProvider/BreadCrumbsProvideContainer').then(
      (mod) => mod.BreadcrumbProvidercontainer,
    ),
  { ssr: false },
)

export const metadata = {
  title: 'Heimdall',
  description: 'Welcome to the Heimdall application',
  icons: {
    icon: '/favicon.png',
  },
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang='en' suppressHydrationWarning>
      <body>
        <NuqsAdapter>
          <ReactQueryProvider>
            <div className='App'>
              <BreadcrumbProvidercontainer>
                <LeftNavContainer />
                <main className='app-content-layout'>
                  <HeaderContent />
                  <div className='App-content'>{children}</div>
                </main>
              </BreadcrumbProvidercontainer>
            </div>
          </ReactQueryProvider>
        </NuqsAdapter>
      </body>
    </html>
  )
}
