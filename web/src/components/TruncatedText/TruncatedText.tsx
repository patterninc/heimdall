import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@patterninc/pattern-ui/tooltip'
import React from 'react'

type TruncatedTextProps = {
  text: string
  /** Character count past which the full text is offered in a tooltip. */
  limit?: number
  className?: string
}

const TruncatedText = ({
  text,
  limit = 60,
  className,
}: TruncatedTextProps): React.JSX.Element => {
  const isLong = text.length > limit
  const label = isLong ? `${text.slice(0, limit).trimEnd()}…` : text

  if (!isLong) return <span className={className}>{text}</span>

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <span className={`cursor-default ${className ?? ''}`}>{label}</span>
      </TooltipTrigger>
      <TooltipContent side='top' className='max-w-md break-words'>
        {text}
      </TooltipContent>
    </Tooltip>
  )
}

export default TruncatedText
