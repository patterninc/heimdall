import { screen } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import { describe, expect, it } from 'vitest'

import { api, server } from '@/test/fixtures'
import { renderWithProviders } from '@/test/render'
import ApiResponseButton from './ApiResponseButton'

const openModal = async () => {
  const view = renderWithProviders(
    <ApiResponseButton href={api('/job/job-0001')} qaTestId='api' />,
  )
  await view.user.click(screen.getByRole('button', { name: 'API response' }))
  return view
}

describe('ApiResponseButton', () => {
  it('does not fetch until opened, then shows the JSON response', async () => {
    let calls = 0
    server.use(
      http.get(api('/job/:id'), () => {
        calls += 1
        return HttpResponse.json({ id: 'job-0001', status: 'RUNNING' })
      }),
    )
    const { user } = renderWithProviders(
      <ApiResponseButton href={api('/job/job-0001')} qaTestId='api' />,
    )
    expect(calls).toBe(0)

    await user.click(screen.getByRole('button', { name: 'API response' }))

    expect(await screen.findByTestId('api-code')).toHaveTextContent(
      '"status": "RUNNING"',
    )
    expect(calls).toBe(1)
  })

  it('shows an empty state for an empty response', async () => {
    server.use(http.get(api('/job/:id'), () => HttpResponse.json([])))
    await openModal()

    expect(
      await screen.findByText(
        'The API returned an empty response for this request.',
      ),
    ).toBeInTheDocument()
  })

  it('shows an error state when the request fails', async () => {
    server.use(
      http.get(api('/job/:id'), () =>
        HttpResponse.json({ error: 'boom' }, { status: 500 }),
      ),
    )
    await openModal()

    expect(
      await screen.findByText("We couldn't load the API response"),
    ).toBeInTheDocument()
  })
})
