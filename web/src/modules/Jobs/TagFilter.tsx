import React from 'react'
import { Button, FormLabel, TextInput } from '@patterninc/react-ui'
import { TagPair } from './Helper'

type TagFilterProps = {
  tags: TagPair[]
  onChange: (tags: TagPair[]) => void
}

const emptyPair: TagPair = { key: '', value: '' }

const TagFilter = ({ tags, onChange }: TagFilterProps): React.JSX.Element => {
  const rows = tags.length > 0 ? tags : [emptyPair]

  const updateRow = (index: number, field: keyof TagPair, value: string) => {
    onChange(
      rows.map((row, i) => (i === index ? { ...row, [field]: value } : row)),
    )
  }

  const addRow = () => onChange([...rows, { ...emptyPair }])

  const removeRow = (index: number) => {
    onChange(rows.filter((_, i) => i !== index))
  }

  return (
    <div className='flex flex-col gap-2'>
      <FormLabel label='Tags' />
      {rows.map((row, index) => (
        <div key={index} className='flex items-center gap-2'>
          <div className='flex-1'>
            <TextInput
              fullWidth
              value={row.key}
              placeholder='Key'
              stateName={`tag-key-${index}`}
              callout={(_, value) => updateRow(index, 'key', value as string)}
            />
          </div>
          <span>:</span>
          <div className='flex-1'>
            <TextInput
              fullWidth
              value={row.value}
              placeholder='Value'
              stateName={`tag-value-${index}`}
              callout={(_, value) => updateRow(index, 'value', value as string)}
            />
          </div>
          <Button as='unstyled' icon='trash' onClick={() => removeRow(index)} />
        </div>
      ))}
      <div>
        <Button as='button' styleType='tertiary' onClick={addRow}>
          Add tag
        </Button>
      </div>
    </div>
  )
}

export default TagFilter
