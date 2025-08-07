import ClientLayout from '@/common/ClientLayout/ClientLayout'
import { NuqsAdapter } from 'nuqs/adapters/next/app'
import './globals.css'

export const metadata = {
  title: 'Heimdall',
  description: 'Welcome to the Heimdall application',
  icons: {
    icon: '/favicon.png',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang='en' suppressHydrationWarning>
      <body>
        <NuqsAdapter>
          <ClientLayout>{children}</ClientLayout>
        </NuqsAdapter>
      </body>
    </html>
  )
}
