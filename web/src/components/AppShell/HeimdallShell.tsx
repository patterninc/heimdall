'use client'

import { AppShell } from '@patterninc/pattern-ui/app-shell'
import { Avatar } from '@patterninc/pattern-ui/avatar'
import { Button } from '@patterninc/pattern-ui/button'
import { Divider } from '@patterninc/pattern-ui/divider'
import { getInitials } from '@patterninc/pattern-ui/get-initials'
import MenuIcon from '@patterninc/pattern-ui/icons/MenuIcon'
import Moon2Icon from '@patterninc/pattern-ui/icons/Moon2Icon'
import PatternLogoIcon from '@patterninc/pattern-ui/icons/PatternLogoIcon'
import SunIcon from '@patterninc/pattern-ui/icons/SunIcon'
import { Nav, type NavItem } from '@patterninc/pattern-ui/nav'
import { NavMobileDrawer } from '@patterninc/pattern-ui/nav-mobile-drawer'
import { Segment, type SegmentItem } from '@patterninc/pattern-ui/segment'
import { type Theme, useTheme } from '@patterninc/pattern-ui/theme'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import React, {
  type ReactElement,
  type ReactNode,
  useEffect,
  useState,
} from 'react'

import { useUser } from '@/common/hooks/useUser'
import { getVersion } from '@/common/Services'
import { APP_ROUTES, isRouteActive } from './navItems'

const THEME_ITEMS: SegmentItem[] = [
  { value: 'light', label: 'Light', icon: SunIcon },
  { value: 'dark', label: 'Dark', icon: Moon2Icon },
  { value: 'system', label: 'System' },
]

const isTheme = (value: string | number): value is Theme =>
  value === 'light' || value === 'dark' || value === 'system'

const HeimdallShell = ({
  children,
}: {
  children: ReactNode
}): React.JSX.Element => {
  const pathname = usePathname()
  const router = useRouter()
  const user = useUser()
  const { theme, setTheme } = useTheme()
  const [mobileNavOpen, setMobileNavOpen] = useState(false)
  const [userMenuOpen, setUserMenuOpen] = useState(false)
  const [version] = useState(getVersion)

  useEffect(() => {
    if (pathname === '/') router.replace('/jobs')
  }, [pathname, router])

  const userName = user || 'X-Pattern-User'
  // getInitials splits on spaces; usernames here are hyphen/underscore/dot separated. First + last
  // word matches the initials NavMobileDrawer derives from `navProfile.displayName`.
  const nameWords = userName.split(/[-_.\s]+/).filter(Boolean)
  const initials =
    nameWords.length > 1
      ? getInitials(`${nameWords[0]} ${nameWords[nameWords.length - 1]}`)
      : getInitials(userName)

  const navItems: NavItem[] = APP_ROUTES.map((route) => ({
    id: route.id,
    type: 'link',
    to: route.href,
    icon: route.icon,
    tooltip: route.label,
    active: isRouteActive(pathname, route.href),
  }))

  const renderItem = (
    item: NavItem,
    defaultElement: ReactElement,
  ): ReactElement => {
    if (item.type === 'link' && item.to) {
      return <Link href={item.to}>{defaultElement}</Link>
    }
    return defaultElement
  }

  const userMenuContent = (
    <div className='flex flex-col gap-0.5'>
      <div className='flex items-center gap-3 p-2'>
        <Avatar initials={initials} size='medium' />
        <div className='flex min-w-0 flex-1 flex-col'>
          <span className='pui-text-label-medium text-text-strong truncate'>
            {userName}
          </span>
          {version ? (
            <span className='pui-text-mono-tiny-alt text-text-soft truncate'>
              v{version}-pattern
            </span>
          ) : null}
        </div>
      </div>

      <Divider type='dashed' spacing />

      {/* Keeps Segment clicks from reaching the nav's own click handlers. */}
      <div
        onClick={(e) => e.stopPropagation()}
        onPointerDown={(e) => e.stopPropagation()}
      >
        <Segment
          ariaLabel='Theme'
          items={THEME_ITEMS}
          value={theme}
          onChange={(value) => {
            if (isTheme(value)) setTheme(value)
            setUserMenuOpen(false)
          }}
          fullWidth
          qaTestId='theme-segment'
        />
      </div>
    </div>
  )

  const avatarPopover = {
    trigger: (
      <button
        type='button'
        className='cursor-pointer'
        aria-label={`${userName} — user menu`}
      >
        <Avatar initials={initials} size='medium' />
      </button>
    ),
    children: userMenuContent,
    position: 'top-start' as const,
    variant: 'inverse' as const,
    className: 'min-w-64 max-h-none overflow-visible',
    open: userMenuOpen,
    onOpenChange: setUserMenuOpen,
  }

  const sharedNavProps = {
    renderItem,
    avatarPopover,
    navProfile: { displayName: userName },
  }

  const navRail = (
    <aside className='absolute top-0 left-0 z-40 h-dvh w-0 md:w-16'>
      <div className='hidden h-full md:flex'>
        <Nav
          {...sharedNavProps}
          items={navItems}
          appIcon={
            <Link href='/jobs' aria-label='Heimdall home'>
              <PatternLogoIcon size={20} color='icon-strong' />
            </Link>
          }
        />
      </div>
      <div className='md:hidden'>
        <NavMobileDrawer
          {...sharedNavProps}
          items={navItems}
          open={mobileNavOpen}
          onOpenChange={setMobileNavOpen}
          pathname={pathname}
        />
      </div>
    </aside>
  )

  return (
    <AppShell nav={navRail}>
      <div className='bg-base-inverse sticky top-0 z-50 md:hidden'>
        <div className='flex h-14 items-center gap-3 px-4'>
          <Button
            variant='weak'
            className='flex h-9 w-9 items-center justify-center'
            prefix={{ icon: MenuIcon }}
            aria-label='Open navigation'
            onClick={() => setMobileNavOpen(true)}
            qaTestId='mobile-nav-button'
          />
          <span className='pui-text-label-medium text-text-strong'>
            Heimdall
          </span>
        </div>
        <Divider />
      </div>

      {/* No overflow-hidden: it would create a containing block that stops ContentHeader sticking. */}
      <div className='flex min-h-0 min-w-0 flex-1 flex-col'>{children}</div>
    </AppShell>
  )
}

export default HeimdallShell
