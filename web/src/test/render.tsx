import { i18n } from '@patterninc/pattern-ui/i18n'
import { SimpleProvider } from '@patterninc/pattern-ui/simple-provider'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, type RenderOptions } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { NuqsTestingAdapter } from 'nuqs/adapters/testing'
import React, { type ReactElement, type ReactNode } from 'react'

// Same as AppRoot: pattern-ui formats dates through the global Lingui locale.
i18n.load('en', {})
i18n.activate('en')

type ProviderOptions = {
  /** Initial URL query string for nuqs, e.g. `'?status=FAILED'`. */
  searchParams?: string
  queryClient?: QueryClient
}

export const createTestQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: { retry: false, gcTime: Infinity, refetchOnWindowFocus: false },
      mutations: { retry: false },
    },
  })

/** Renders with the same providers AppRoot wires up (pattern-ui, React Query, nuqs), isolated per test. */
export const renderWithProviders = (
  ui: ReactElement,
  {
    searchParams = '',
    queryClient = createTestQueryClient(),
    ...options
  }: ProviderOptions & Omit<RenderOptions, 'wrapper'> = {},
) => {
  const Wrapper = ({ children }: { children: ReactNode }) => (
    <SimpleProvider>
      <QueryClientProvider client={queryClient}>
        <NuqsTestingAdapter searchParams={searchParams}>
          {children}
        </NuqsTestingAdapter>
      </QueryClientProvider>
    </SimpleProvider>
  )

  return {
    user: userEvent.setup(),
    queryClient,
    ...render(ui, { wrapper: Wrapper, ...options }),
  }
}
