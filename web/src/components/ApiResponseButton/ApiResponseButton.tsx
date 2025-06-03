import { Button, Icon } from '@patterninc/react-ui'
import Link from 'next/link'
import React from 'react'

type ApiResponseButtonProps = {
  link: string
}

const ApiResponseButton = ({ link }: ApiResponseButtonProps): React.JSX.Element => {
  return (
    <Button
      as='link'
      target='_blank'
      href={link}
      routerComponent={Link}
      styleType='text-blue'
      className='pat-gap-2'
    >
      <span>API Response</span>
      <Icon icon='launch' iconSize='12px' color='dark-blue' />
    </Button>
  )
}

export default ApiResponseButton
