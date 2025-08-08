'use client'

import React, { useContext } from 'react'
import { Breadcrumbs } from '@patterninc/react-ui'

import { useRouter } from 'next/navigation'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { AutoRefreshSelect } from '@/common/AutoRefreshSelect/AutoRefreshSelect'

const Header = () => {
  const { breadcrumbs, breadcrumbCallout } = useContext(BreadcrumbContext)
  const router = useRouter()

  return (
    <div className="header sticky top-0 bg-white z-[99]">
      <Breadcrumbs
        breadcrumbs={breadcrumbs}
        callout={(breadcrumb) => {
          breadcrumbCallout(breadcrumb)
          router.push(breadcrumb.link)
        }}
        backButtonProps={{ text: 'Back' }}
      />
      {breadcrumbs.length <= 1 && <AutoRefreshSelect />}
    </div>
  )
}

export default Header
