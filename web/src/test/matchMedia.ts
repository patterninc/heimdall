import { vi } from 'vitest'

export const createMediaQueryList = (
  query: string,
  matches: boolean,
): MediaQueryList => ({
  matches,
  media: query,
  onchange: null,
  addListener: () => {},
  removeListener: () => {},
  addEventListener: () => {},
  removeEventListener: () => {},
  dispatchEvent: () => false,
})

/**
 * Replaces `window.matchMedia` for the current test; `matches` decides each query.
 * The setup file installs a never-matching default; `vi.restoreAllMocks()` puts it back.
 */
export const mockMatchMedia = (matches: (query: string) => boolean) =>
  vi
    .spyOn(window, 'matchMedia')
    .mockImplementation((query: string) =>
      createMediaQueryList(query, matches(query)),
    )
