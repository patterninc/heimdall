'use client'

import React from 'react'
import { ErrorPage, APP_LOGOS } from '@patterninc/react-ui'
import Link from 'next/link'
import Image from 'next/image'
import styles from './_error-page.module.scss'

const NotFoundPage = () => {
  return (
    <div className={styles.errorPageContainer}>
      <ErrorPage
        logo={
          <div className='flex flex-direction-column align-items-center justify-content-center'>
            <Image
              src={APP_LOGOS.HEIMDALL.isolated}
              alt='Heimdall'
              width={136}
              height={43}
            />
          </div>
        }
        title={`Oops! We couldn't find that page.`}
        message={`It looks like the page you're looking for doesn't exist or has
              been moved.`}
        supportMessage={
          <span>
            Please check the URL or{' '}
            <Link href='/' className='fc-blue'>
              return to the dashboard.
            </Link>
          </span>
        }
        thankYouMessage='Thank you for using Heimdall!'
      />
    </div>
  )
}

export default NotFoundPage
