'use client'

import { i18n } from '@patterninc/pattern-ui/i18n'
import { I18nEnglishProvider } from '@patterninc/pattern-ui/i18n-english-provider'
import { PulseProvider } from '@patterninc/pattern-ui/pulse'
import { ThemeProvider } from '@patterninc/pattern-ui/theme'
import { TooltipProvider } from '@patterninc/pattern-ui/tooltip'
import { ReactNode } from 'react'

import HeimdallShell from '@/components/AppShell/HeimdallShell'
import { AutoRefreshProvider } from '../AutoRefreshProvider/context'
import ReactQueryProvider from '../ReactQueryProvider/ReactQueryProvider'

// pattern-ui formats dates through `Intl` with the global Lingui locale, which is empty until a
// locale is activated. Heimdall has no catalogs of its own, so English is activated once here.
i18n.load('en', {})
i18n.activate('en')

const AppRoot = ({ children }: { children: ReactNode }) => {
  return (
    <I18nEnglishProvider>
      <ThemeProvider>
        <TooltipProvider>
          <ReactQueryProvider>
            <AutoRefreshProvider>
              <PulseProvider />
              <HeimdallShell>{children}</HeimdallShell>
            </AutoRefreshProvider>
          </ReactQueryProvider>
        </TooltipProvider>
      </ThemeProvider>
    </I18nEnglishProvider>
  )
}

export default AppRoot
