import { screen } from '@testing-library/react'
import { useState } from 'react'
import { describe, expect, it, vi } from 'vitest'

import { mockMatchMedia } from '@/test/matchMedia'
import { renderWithProviders } from '@/test/render'
import { FilterMenu } from './FilterMenu'
import { isDimensionActive, type FilterDimension } from './filterTypes'

const Harness = ({ onClearAll }: { onClearAll: () => void }) => {
  const [name, setName] = useState('')
  const [status, setStatus] = useState<string[]>([])
  const dimensions: FilterDimension[] = [
    {
      key: 'name',
      label: 'Name',
      type: 'text',
      value: name,
      onChange: setName,
    },
    {
      key: 'status',
      label: 'Status',
      type: 'multi',
      value: status,
      onChange: setStatus,
      options: ['ACTIVE', 'INACTIVE'],
    },
    {
      key: 'empty',
      label: 'Plugin',
      type: 'multi',
      value: [],
      onChange: () => {},
      options: [],
    },
  ]

  return (
    <>
      <FilterMenu
        dimensions={dimensions}
        onClearAll={() => {
          setName('')
          setStatus([])
          onClearAll()
        }}
      />
      <output qa-test-id='state'>{JSON.stringify({ name, status })}</output>
    </>
  )
}

const state = () => JSON.parse(screen.getByTestId('state').textContent!)

// The default jsdom matchMedia never matches, so FilterMenu renders its phone (drawer) layout.
describe('FilterMenu (drawer layout)', () => {
  it('edits text and multi-select dimensions and clears them all', async () => {
    const onClearAll = vi.fn()
    const { user } = renderWithProviders(<Harness onClearAll={onClearAll} />)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(await screen.findByText('Name'))
    await user.type(await screen.findByPlaceholderText('Enter name'), 'spark')
    expect(state().name).toBe('spark')
    await user.keyboard('{Escape}')

    await user.click(screen.getByText('Status'))
    await user.click(await screen.findByTestId('filter-status-ACTIVE'))
    await user.click(screen.getByTestId('filter-status-INACTIVE'))
    expect(state().status).toEqual(['ACTIVE', 'INACTIVE'])

    await user.click(screen.getByTestId('filter-status-ACTIVE'))
    expect(state().status).toEqual(['INACTIVE'])

    await user.click(screen.getByTestId('filter-status-clear'))
    expect(state().status).toEqual([])
    await user.click(screen.getByTestId('filter-status-INACTIVE'))
    await user.keyboard('{Escape}')

    await user.click(
      await screen.findByRole('button', { name: 'Clear all filters' }),
    )
    expect(onClearAll).toHaveBeenCalledOnce()
    expect(state()).toEqual({ name: '', status: [] })
  })

  it('toggles a multi-select option from the keyboard', async () => {
    const { user } = renderWithProviders(<Harness onClearAll={() => {}} />)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(await screen.findByText('Status'))
    const option = await screen.findByTestId('filter-status-ACTIVE')
    option.focus()
    await user.keyboard('{Enter}')

    expect(state().status).toEqual(['ACTIVE'])
  })

  it('shows a placeholder when a multi-select has no options', async () => {
    const { user } = renderWithProviders(<Harness onClearAll={() => {}} />)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(await screen.findByText('Plugin'))

    expect(await screen.findByText('No options available')).toBeInTheDocument()
  })
})

describe('isDimensionActive', () => {
  it('treats whitespace-only text as inactive', () => {
    const text: FilterDimension = {
      key: 'name',
      label: 'Name',
      type: 'text',
      value: '  ',
      onChange: () => {},
    }
    expect(isDimensionActive(text)).toBe(false)
    expect(isDimensionActive({ ...text, value: 'x' })).toBe(true)
  })

  it('uses the custom flag for custom dimensions', () => {
    expect(
      isDimensionActive({
        key: 'tags',
        label: 'Tags',
        type: 'custom',
        isActive: true,
        content: null,
      }),
    ).toBe(true)
  })
})

describe('FilterMenu (menu layout)', () => {
  it('shows the active count and a clear-all button on wide screens', async () => {
    mockMatchMedia(() => true)
    const onClearAll = vi.fn()
    const { user } = renderWithProviders(<Harness onClearAll={onClearAll} />)

    await user.click(screen.getByRole('button', { name: 'Filters' }))
    await user.click(await screen.findByTestId('filter-status'))
    await user.click(await screen.findByTestId('filter-status-ACTIVE'))
    expect(state().status).toEqual(['ACTIVE'])
    await user.keyboard('{Escape}')
    await user.keyboard('{Escape}')

    await user.click(await screen.findByTestId('filter-clear-all'))
    expect(onClearAll).toHaveBeenCalledOnce()
    expect(state().status).toEqual([])
  })
})
