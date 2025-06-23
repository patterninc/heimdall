'use client'

import React from 'react'
import { ErrorPage, Button } from '@patterninc/react-ui'
import Link from 'next/link'

const NotFoundPage = () => {
  return (
    <div>
      <ErrorPage logo={<></>} title='Page not found!!' supportMessage='' />
      <div className='flex flex-direction-column align-items-center justify-content-center pat-pt-2'>
        <Link href='/'>
          <Button styleType='primary-green'>Reload Heimdall</Button>
        </Link>
      </div>
    </div>
  )
}

export default NotFoundPage
