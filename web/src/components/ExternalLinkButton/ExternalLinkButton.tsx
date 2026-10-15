import { Button, type ButtonProps } from '@patterninc/pattern-ui/button'
import ArrowTopRightSquareIcon from '@patterninc/pattern-ui/icons/ArrowTopRightSquareIcon'
import React from 'react'

type ExternalLinkButtonProps = {
  href: string
  children: string
  variant?: ButtonProps['variant']
  qaTestId?: string
}

const ExternalLinkButton = ({
  href,
  children,
  variant = 'stroke',
  qaTestId,
}: ExternalLinkButtonProps): React.JSX.Element => {
  return (
    <Button
      asChild
      variant={variant}
      size='medium'
      suffix={{ icon: ArrowTopRightSquareIcon, size: 16 }}
      qaTestId={qaTestId}
    >
      <a href={href} target='_blank' rel='noreferrer'>
        {children}
      </a>
    </Button>
  )
}

export default ExternalLinkButton
