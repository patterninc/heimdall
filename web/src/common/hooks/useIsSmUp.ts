import { useMediaQuery } from 'usehooks-ts'

/** True at Tailwind's `sm` breakpoint and up. False until mounted, so SSR renders the mobile layout. */
export const useIsSmUp = (): boolean =>
  useMediaQuery('(min-width: 640px)', { initializeWithValue: false })
