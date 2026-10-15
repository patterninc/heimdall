'use client'

import { Button } from '@patterninc/pattern-ui/button'
import { InputField } from '@patterninc/pattern-ui/fields/input-field'
import { MenuSelectItem } from '@patterninc/pattern-ui/menu'
import { NoDataAvailable } from '@patterninc/pattern-ui/no-data-available'
import React from 'react'

import type {
  FilterDimension,
  MultiFilterDimension,
  TextFilterDimension,
} from './filterTypes'

// Menus run typeahead on keydown; stop it so typing in a field doesn't move menu focus.
const stopMenuKeys = (event: React.KeyboardEvent) => event.stopPropagation()

const TextFilterControl = ({
  dimension,
}: {
  dimension: TextFilterDimension
}): React.JSX.Element => (
  <InputField
    label={dimension.label}
    variant='weak'
    value={dimension.value}
    placeholder={
      dimension.placeholder ?? `Enter ${dimension.label.toLowerCase()}`
    }
    onChange={(event) => dimension.onChange(event.target.value)}
    onKeyDown={stopMenuKeys}
    autoFocus
    qaTestId={`filter-${dimension.key}-input`}
  />
)

const MultiFilterControl = ({
  dimension,
}: {
  dimension: MultiFilterDimension
}): React.JSX.Element => {
  const toggle = (option: string) =>
    dimension.onChange(
      dimension.value.includes(option)
        ? dimension.value.filter((value) => value !== option)
        : [...dimension.value, option],
    )

  return (
    <div className='flex flex-col gap-0.5'>
      <div className='flex items-center justify-between gap-2 px-2.5 py-1'>
        <span className='pui-text-label-small text-text-soft'>
          {dimension.label}
        </span>
        {dimension.value.length > 0 ? (
          <Button
            variant='inverse'
            size='small'
            onClick={() => dimension.onChange([])}
            qaTestId={`filter-${dimension.key}-clear`}
          >
            Clear
          </Button>
        ) : null}
      </div>
      {dimension.options.length === 0 ? (
        <div className='px-2.5 py-2'>
          <NoDataAvailable text='No options available' />
        </div>
      ) : (
        dimension.options.map((option) => (
          <MenuSelectItem
            key={option}
            selected={dimension.value.includes(option)}
            className='focus-visible:outline-primary-base focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-solid'
            tabIndex={0}
            onPointerDown={(event) => event.stopPropagation()}
            onKeyDown={(event) => {
              if (event.key === 'Enter' || event.key === ' ') {
                event.preventDefault()
                event.stopPropagation()
                toggle(option)
              }
            }}
            onClick={() => toggle(option)}
            qaTestId={`filter-${dimension.key}-${option}`}
          >
            {option}
          </MenuSelectItem>
        ))
      )}
    </div>
  )
}

export const FilterControl = ({
  dimension,
}: {
  dimension: FilterDimension
}): React.JSX.Element => {
  switch (dimension.type) {
    case 'text':
      return <TextFilterControl dimension={dimension} />
    case 'multi':
      return <MultiFilterControl dimension={dimension} />
    case 'custom':
      return <>{dimension.content}</>
  }
}

export { stopMenuKeys }
