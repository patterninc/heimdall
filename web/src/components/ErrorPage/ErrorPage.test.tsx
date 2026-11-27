import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { renderWithProviders } from '@/test/render'
import { mockRouter } from '@/test/router'
import NotFoundPage from './ErrorPage'

describe('NotFoundPage', () => {
  it('sends the user back to the jobs list', async () => {
    const { user } = renderWithProviders(<NotFoundPage />)

    await user.click(
      screen.getByRole('button', { name: 'Return to the dashboard' }),
    )

    expect(mockRouter.push).toHaveBeenCalledWith('/jobs')
  })
})
