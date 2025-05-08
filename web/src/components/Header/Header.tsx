'use client'

import React, { useContext } from 'react'
import { Breadcrumbs } from '@patterninc/react-ui'

import { useRouter } from 'next/navigation'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import styles from './_header.module.scss'

const Header = () => {
  const { breadcrumbs, breadcrumbCallout } = useContext(BreadcrumbContext)
  const router = useRouter()

  return (
    <div className={`header ${styles.headerContainer}`}>
      <Breadcrumbs
        breadcrumbs={breadcrumbs}
        callout={(breadcrumb) => {
          breadcrumbCallout(breadcrumb)
          router.push(breadcrumb.link)
        }}
        backButtonProps={{ text: 'Back' }}
      />
    </div>
  )
}

export default Header
