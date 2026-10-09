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

// Applies the stored theme before hydration so dark mode doesn't flash light first.
// usehooks-ts stores the value JSON-encoded under the `theme` key.
const noFlashThemeScript = `(function(){try{var t=JSON.parse(localStorage.getItem('theme')||'"system"');var d=t==='dark'||(t==='system'&&window.matchMedia('(prefers-color-scheme: dark)').matches);document.documentElement.classList.toggle('dark',d)}catch(e){}})()`

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang='en' suppressHydrationWarning>
      <head>
        <script dangerouslySetInnerHTML={{ __html: noFlashThemeScript }} />
      </head>
      <body>
        <NuqsAdapter>
          <ClientLayout>{children}</ClientLayout>
        </NuqsAdapter>
      </body>
    </html>
  )
}
