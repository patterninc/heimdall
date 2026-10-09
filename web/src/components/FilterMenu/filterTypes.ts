import type { ReactNode } from 'react'

type FilterDimensionBase = {
  key: string
  label: string
}

export type TextFilterDimension = FilterDimensionBase & {
  type: 'text'
  value: string
  onChange: (value: string) => void
  placeholder?: string
}

export type MultiFilterDimension = FilterDimensionBase & {
  type: 'multi'
  value: string[]
  onChange: (value: string[]) => void
  options: string[]
}

/** Escape hatch for dimensions with their own editor (e.g. key:value tag pairs). */
export type CustomFilterDimension = FilterDimensionBase & {
  type: 'custom'
  isActive: boolean
  content: ReactNode
}

export type FilterDimension =
  | TextFilterDimension
  | MultiFilterDimension
  | CustomFilterDimension

export const isDimensionActive = (dimension: FilterDimension): boolean => {
  switch (dimension.type) {
    case 'text':
      return dimension.value.trim() !== ''
    case 'multi':
      return dimension.value.length > 0
    case 'custom':
      return dimension.isActive
  }
}
