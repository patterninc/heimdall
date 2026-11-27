import { screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'

import { renderWithProviders } from '@/test/render'
import TruncatedText from './TruncatedText'

describe('TruncatedText', () => {
  it('renders short text as is', () => {
    renderWithProviders(<TruncatedText text='short' />)

    expect(screen.getByText('short')).toBeInTheDocument()
  })

  it('cuts long text at the limit and keeps the full text for assistive tech and the tooltip', async () => {
    const text = 'a very long description that goes past the limit'
    const { user } = renderWithProviders(
      <TruncatedText text={text} limit={10} />,
    )

    const trigger = screen.getByLabelText(text)
    expect(trigger).toHaveTextContent('a very lon…')

    await user.hover(trigger)
    expect(await screen.findByRole('tooltip')).toHaveTextContent(text)
  })
})
