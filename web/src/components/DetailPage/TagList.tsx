import { Tag } from '@patterninc/pattern-ui/tag'
import React from 'react'

type TagListProps = {
  values: string[]
}

const TagList = ({ values }: TagListProps): React.JSX.Element => (
  <div className='flex flex-wrap gap-2 p-4'>
    {values.map((value) => (
      <Tag key={value} variant='neutral'>
        {value}
      </Tag>
    ))}
  </div>
)

export default TagList
