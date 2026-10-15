'use client'

import { Button } from '@patterninc/pattern-ui/button'
import {
  Drawer,
  DrawerFooter,
  DrawerHeader,
  DrawerPanel,
  DrawerPopup,
} from '@patterninc/pattern-ui/drawer'
import ChevronRightIcon from '@patterninc/pattern-ui/icons/ChevronRightIcon'
import FiltersIcon from '@patterninc/pattern-ui/icons/FiltersIcon'
import { ListItem } from '@patterninc/pattern-ui/list-item'
import { Menu, MenuSubItem } from '@patterninc/pattern-ui/menu'
import { NumberBadge } from '@patterninc/pattern-ui/number-badge'
import React, { useState } from 'react'
import { useMediaQuery } from 'usehooks-ts'

import { FilterControl } from './FilterControl'
import { type FilterDimension, isDimensionActive } from './filterTypes'

type FilterMenuProps = {
  dimensions: FilterDimension[]
  onClearAll: () => void
}

const FilterActiveDot = (): React.JSX.Element => (
  <span className='bg-primary-base h-2 w-2 rounded-full' aria-hidden />
)

export const FilterMenu = ({
  dimensions,
  onClearAll,
}: FilterMenuProps): React.JSX.Element => {
  const isDesktop = useMediaQuery('(min-width: 768px)', {
    defaultValue: true,
    initializeWithValue: false,
  })
  const [menuOpen, setMenuOpen] = useState(false)
  const [mobileDimensionKey, setMobileDimensionKey] = useState<string | null>(
    null,
  )

  const activeCount = dimensions.filter(isDimensionActive).length
  const activeBadge =
    activeCount > 0 ? (
      <NumberBadge value={activeCount} size={20} variant='stroke' />
    ) : null

  const trigger = (
    <Button
      variant='strong'
      size='medium'
      prefix={{ icon: FiltersIcon }}
      onClick={!isDesktop ? () => setMenuOpen(true) : undefined}
      aria-label='Filters'
      qaTestId='filter-menu-button'
    >
      {isDesktop ? (
        <span className='flex items-center gap-1.5'>
          <span>Filters</span>
          {activeBadge}
        </span>
      ) : (
        activeBadge
      )}
    </Button>
  )

  const clearAllButton =
    activeCount > 0 ? (
      <Button
        variant='weak'
        size='medium'
        onClick={onClearAll}
        qaTestId='filter-clear-all'
      >
        Clear all
      </Button>
    ) : null

  if (isDesktop) {
    return (
      <div className='flex items-center gap-2'>
        {clearAllButton}
        {/* Anchored in the header's right slot, so the menu and its submenus open inward. */}
        <Menu
          open={menuOpen}
          onOpenChange={setMenuOpen}
          position='bottom-end'
          className='w-52 p-2'
          preventScrollClose
          trigger={trigger}
        >
          {dimensions.map((dimension) => (
            <MenuSubItem
              key={dimension.key}
              listItemType='medium'
              subMenuPosition='left-start'
              subMenuClassName='max-h-none'
              accessory={
                isDimensionActive(dimension) ? <FilterActiveDot /> : undefined
              }
              subMenuContent={
                <div className='flex w-64 flex-col gap-2 p-1'>
                  <FilterControl dimension={dimension} />
                </div>
              }
              qaTestId={`filter-${dimension.key}`}
            >
              {dimension.label}
            </MenuSubItem>
          ))}
        </Menu>
      </div>
    )
  }

  const activeDimension = dimensions.find((d) => d.key === mobileDimensionKey)

  return (
    <div className='flex items-center gap-2'>
      {/* Clearing lives in the drawer footer; the header row has no room for a second button. */}
      {trigger}

      {/* The category drawer is nested inside the first popup so the two stack as one flow. */}
      <Drawer
        open={menuOpen}
        onOpenChange={(open) => {
          if (!open) {
            setMenuOpen(false)
            setMobileDimensionKey(null)
          }
        }}
      >
        <DrawerPopup showBar showCloseButton surface='weak'>
          <DrawerHeader>
            <h2 className='pui-text-label-large text-text-strong'>Filters</h2>
          </DrawerHeader>
          <DrawerPanel>
            <div className='flex flex-col py-1'>
              {dimensions.map((dimension) => (
                <ListItem
                  key={dimension.key}
                  listItemType='medium'
                  listItemSuffix={
                    <div className='flex items-center gap-1.5'>
                      {isDimensionActive(dimension) ? (
                        <FilterActiveDot />
                      ) : null}
                      <ChevronRightIcon size={16} color='icon-sub' />
                    </div>
                  }
                  onClick={() => setMobileDimensionKey(dimension.key)}
                >
                  {dimension.label}
                </ListItem>
              ))}
            </div>
          </DrawerPanel>
          {activeCount > 0 ? (
            <DrawerFooter>
              <Button
                variant='stroke'
                size='medium'
                onClick={() => {
                  onClearAll()
                  setMenuOpen(false)
                }}
              >
                Clear all filters
              </Button>
            </DrawerFooter>
          ) : null}

          <Drawer
            open={activeDimension !== undefined}
            onOpenChange={(open) => {
              if (!open) setMobileDimensionKey(null)
            }}
          >
            <DrawerPopup showBar showCloseButton surface='weak'>
              <DrawerHeader>
                <h2 className='pui-text-label-large text-text-strong'>
                  {activeDimension?.label ?? ''}
                </h2>
              </DrawerHeader>
              <DrawerPanel>
                <div className='flex flex-col gap-2 py-1'>
                  {activeDimension ? (
                    <FilterControl dimension={activeDimension} />
                  ) : null}
                </div>
              </DrawerPanel>
            </DrawerPopup>
          </Drawer>
        </DrawerPopup>
      </Drawer>
    </div>
  )
}
