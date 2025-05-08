'use client'

import { useState } from 'react'

import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

import { useIsClient } from '@patterninc/react-ui'

export const ReactQueryProvider = ({
  children,
}: {
  children: React.ReactNode
}) => {
  const isClient = useIsClient()
  const [queryClient] = useState(
      () =>
        new QueryClient({
          defaultOptions: {
            queries: {
              refetchOnWindowFocus: false,
            },
          },
        }),
    ),
    enableReactQueryDevTools =
      isClient && localStorage.getItem('react-query-devtools') === 'true'

  return (
    <QueryClientProvider client={queryClient}>
      {children}
      {enableReactQueryDevTools ? <ReactQueryDevtools /> : null}
    </QueryClientProvider>
  )
}

export default ReactQueryProvider
