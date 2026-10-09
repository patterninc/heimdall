import Link from 'next/link'
import React, { type ReactNode } from 'react'

type RowLinkProps = {
  href: string
  children: ReactNode
}

/**
 * Primary-cell link for clickable table rows. Keeps cmd/ctrl-click and keyboard navigation working
 * while the row itself also routes on click; propagation is stopped so it doesn't navigate twice.
 */
const RowLink = ({ href, children }: RowLinkProps): React.JSX.Element => (
  <Link
    href={href}
    onClick={(event) => event.stopPropagation()}
    className='pui-text-body-medium text-text-strong whitespace-nowrap hover:underline focus-visible:underline'
  >
    {children}
  </Link>
)

export default RowLink
