'use client'

import React, { useContext, useEffect, useMemo } from 'react'
import { APP_LOGOS, LeftNavLinkObj, LeftNav, Icon } from '@patterninc/react-ui'
import Link from 'next/link'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { getVersion } from '@/common/Services'


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

  const version = getVersion()

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
      footer={({ sidebarExpanded }) =>
        sidebarExpanded ? (
          <>
            <div className='flex items-center text-[10px] text-purple pl-4 h-5'>
              <>
                &copy; {new Date().getFullYear()}, v{version}-pattern. Built
                with &nbsp;
                <Icon icon='heart' iconSize='12px' color='dark-purple' /> &nbsp;
                at Pattern
              </>
            </div>
            <LeftNav.Divider />
          </>
        ) : (
          <></>
        )
      }
      hideCopyright={true}
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
