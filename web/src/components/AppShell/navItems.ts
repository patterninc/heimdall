import BriefcaseIcon from '@patterninc/pattern-ui/icons/BriefcaseIcon'
import GridIcon from '@patterninc/pattern-ui/icons/GridIcon'
import MagicWandIcon from '@patterninc/pattern-ui/icons/MagicWandIcon'
import type { NavItemProps } from '@patterninc/pattern-ui/nav'

export type AppRoute = {
  id: string
  label: string
  href: string
  icon: NavItemProps['icon']
}

export const APP_ROUTES: AppRoute[] = [
  { id: 'jobs', label: 'Jobs', href: '/jobs', icon: BriefcaseIcon },
  { id: 'commands', label: 'Commands', href: '/commands', icon: MagicWandIcon },
  { id: 'clusters', label: 'Clusters', href: '/clusters', icon: GridIcon },
]

export const isRouteActive = (pathname: string, href: string): boolean =>
  pathname === href || pathname.startsWith(`${href}/`)
