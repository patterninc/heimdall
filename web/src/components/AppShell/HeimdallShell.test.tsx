import { ThemeProvider } from '@patterninc/pattern-ui/theme'
import { screen, within } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { renderWithProviders } from '@/test/render'
import { mockLocation, mockRouter } from '@/test/router'
import HeimdallShell from './HeimdallShell'

const renderShell = () =>
  renderWithProviders(
    <ThemeProvider>
      <HeimdallShell>
        <p>Page content</p>
      </HeimdallShell>
    </ThemeProvider>,
  )

describe('HeimdallShell', () => {
  it('renders the page inside the shell with links to every resource page', () => {
    renderShell()

    expect(screen.getByText('Page content')).toBeInTheDocument()
    expect(screen.getByRole('link', { name: 'Heimdall home' })).toHaveAttribute(
      'href',
      '/jobs',
    )
    const hrefs = screen
      .getAllByRole('link')
      .map((link) => link.getAttribute('href'))
    expect(hrefs).toEqual(
      expect.arrayContaining(['/jobs', '/commands', '/clusters']),
    )
  })

  it('redirects the root path to the jobs list', () => {
    mockLocation.pathname = '/'
    renderShell()

    expect(mockRouter.replace).toHaveBeenCalledWith('/jobs')
  })

  it('shows the signed-in user, version and theme switch in the user menu', async () => {
    document.cookie = `heimdall-user=${encodeURIComponent('jane.doe@example.com')}`
    document.cookie = 'heimdall-version=1.4.2'
    const { user } = renderShell()

    await user.click(
      screen.getByRole('button', {
        name: 'jane.doe@example.com — user menu',
      }),
    )

    expect(await screen.findByText('v1.4.2-pattern')).toBeInTheDocument()
    const themes = screen.getByTestId('theme-segment')
    await user.click(within(themes).getByText('Dark'))
    expect(localStorage.getItem('theme')).toContain('dark')
  })

  it('opens the mobile navigation drawer with the readable user name', async () => {
    document.cookie = 'heimdall-user=jane.doe'
    const { user } = renderShell()
    expect(screen.queryByText('jane doe')).not.toBeInTheDocument()

    await user.click(screen.getByRole('button', { name: 'Open navigation' }))

    expect(await screen.findByText('jane doe')).toBeInTheDocument()
  })
})
