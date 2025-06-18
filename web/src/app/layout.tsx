import ClientLayout from '@/common/ClientLayout/ClientLayout'
import { headers } from 'next/headers'
import { NuqsAdapter } from 'nuqs/adapters/next/app'

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
  const headerList = await headers()
  const user = headerList.get('X-Heimdall-User')
  return (
    <html lang='en' suppressHydrationWarning>
      <body>
        <NuqsAdapter>
          <ClientLayout user={user}>{children}</ClientLayout>
        </NuqsAdapter>
      </body>
    </html>
  )
}
