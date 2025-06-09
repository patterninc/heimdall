'use client'

import React, { useContext, useEffect, useMemo } from 'react'
import { APP_LOGOS, LeftNavLinkObj, LeftNav } from '@patterninc/react-ui'
import Link from 'next/link'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'

type LeftNavBarProps = {
  user: string | null
}

const LeftNavBar = ({ user }: LeftNavBarProps) => {
  const { breadcrumbs } = useContext(BreadcrumbContext)

  useEffect(() => {
    if (window.location.pathname === '/') {
      window.location.href = '/jobs'
    }
  }, [])

  const sidebarContent: LeftNavLinkObj[] = useMemo(
    () => [
      {
        name: 'Jobs',
        link: '/jobs',
        icon: 'briefcase',
      },
      {
        name: 'Commands',
        link: '/commands',
        icon: 'magicWand',
      },
      {
        name: 'Clusters',
        link: '/clusters',
        icon: 'gridIcon',
      },
    ],
    [],
  )

  const translatedSideContent = useMemo(() => {
    return sidebarContent.map((menu: LeftNavLinkObj) => {
      return {
        ...menu,
        name: menu.name,
        key: menu.link, // Add unique key based on link
      }
    })
  }, [sidebarContent])

  return (
    <LeftNav
      logo={{
        url: APP_LOGOS.HEIMDALL.isolated,
        abbreviatedUrl: APP_LOGOS.HEIMDALL.abbr,
        isolatedUrl: APP_LOGOS.HEIMDALL.isolated,
        maxHeight: 32,
      }}
      routerComponent={Link}
      routerProp='href'
      leftNavLinks={translatedSideContent}
      breadcrumbs={breadcrumbs}
      accountPopoverProps={{
        name: user || 'X-Pattern-User',
        options: [],
        disabled: true,
      }}
      userPermissions={[]}
      navigate={() => window.open('/jobs')?.focus()}
      hideLogo={true}
    />
  )
}
export default LeftNavBar
