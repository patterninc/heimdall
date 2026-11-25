import { act, render, screen } from '@testing-library/react'
import { useContext } from 'react'
import { describe, expect, it, vi } from 'vitest'

import { AutoRefreshContext, AutoRefreshProvider } from './context'

const Probe = () => {
  const { refreshInterval, updateRefreshInterval } =
    useContext(AutoRefreshContext)
  return (
    <button
      type='button'
      onClick={() => updateRefreshInterval({ label: '5s', value: 5000 })}
    >
      {refreshInterval.label}:{refreshInterval.value}
    </button>
  )
}

describe('AutoRefreshProvider', () => {
  it('defaults to off and persists changes to localStorage', () => {
    render(
      <AutoRefreshProvider>
        <Probe />
      </AutoRefreshProvider>,
    )
    expect(screen.getByRole('button')).toHaveTextContent('OFF:0')

    act(() => screen.getByRole('button').click())

    expect(screen.getByRole('button')).toHaveTextContent('5s:5000')
    expect(JSON.parse(localStorage.getItem('refreshInterval')!)).toEqual({
      label: '5s',
      value: 5000,
    })
  })

  it('restores the saved interval', () => {
    localStorage.setItem(
      'refreshInterval',
      JSON.stringify({ label: '30s', value: 30000 }),
    )
    render(
      <AutoRefreshProvider>
        <Probe />
      </AutoRefreshProvider>,
    )

    expect(screen.getByRole('button')).toHaveTextContent('30s:30000')
  })

  it('falls back to off when the saved value is not valid JSON', () => {
    const error = vi.spyOn(console, 'error').mockImplementation(() => {})
    localStorage.setItem('refreshInterval', '{not json')
    render(
      <AutoRefreshProvider>
        <Probe />
      </AutoRefreshProvider>,
    )

    expect(screen.getByRole('button')).toHaveTextContent('OFF:0')
    expect(error).toHaveBeenCalled()
  })
})
