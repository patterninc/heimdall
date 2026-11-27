import { screen } from '@testing-library/react'
import { describe, expect, it, vi } from 'vitest'

import { renderWithProviders } from '@/test/render'
import TagFilter from './TagFilter'

describe('TagFilter', () => {
  it('adds a trimmed tag from the inputs and clears the draft', async () => {
    const onChange = vi.fn()
    const { user } = renderWithProviders(
      <TagFilter tags={[{ key: 'team', value: 'data' }]} onChange={onChange} />,
    )
    const add = screen.getByRole('button', { name: 'Add tag' })
    expect(add).toBeDisabled()

    await user.type(screen.getByPlaceholderText('Key'), ' env ')
    await user.type(screen.getByPlaceholderText('Value'), ' prod ')
    await user.click(add)

    expect(onChange).toHaveBeenCalledWith([
      { key: 'team', value: 'data' },
      { key: 'env', value: 'prod' },
    ])
    expect(screen.getByPlaceholderText('Key')).toHaveValue('')
  })

  it('adds a tag with Enter', async () => {
    const onChange = vi.fn()
    const { user } = renderWithProviders(
      <TagFilter tags={[]} onChange={onChange} />,
    )

    await user.type(screen.getByPlaceholderText('Key'), 'env')
    await user.type(screen.getByPlaceholderText('Value'), 'prod{Enter}')

    expect(onChange).toHaveBeenCalledWith([{ key: 'env', value: 'prod' }])
  })

  it('ignores Enter while the draft is incomplete', async () => {
    const onChange = vi.fn()
    const { user } = renderWithProviders(
      <TagFilter tags={[]} onChange={onChange} />,
    )

    await user.type(screen.getByPlaceholderText('Key'), 'env{Enter}')

    expect(onChange).not.toHaveBeenCalled()
  })

  it('removes a tag', async () => {
    const onChange = vi.fn()
    const { user } = renderWithProviders(
      <TagFilter
        tags={[
          { key: 'team', value: 'data' },
          { key: 'env', value: 'prod' },
        ]}
        onChange={onChange}
      />,
    )

    await user.click(
      screen.getByRole('button', { name: 'Remove tag team:data' }),
    )

    expect(onChange).toHaveBeenCalledWith([{ key: 'env', value: 'prod' }])
  })
})
