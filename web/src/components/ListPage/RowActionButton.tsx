import { Button } from '@patterninc/pattern-ui/button'
import Link from 'next/link'
import React from 'react'

type RowActionButtonProps = {
  href: string
  children: string
  qaTestId?: string
}

/** Table row action: a Button rendered as a Next link, so navigation stays a real anchor. */
const RowActionButton = ({
  href,
  children,
  qaTestId,
}: RowActionButtonProps): React.JSX.Element => (
  <Button asChild variant='weak' size='small' qaTestId={qaTestId}>
    <Link href={href}>{children}</Link>
  </Button>
)

export default RowActionButton
