'use client'

import { Button } from '@patterninc/pattern-ui/button'
import { ButtonInset } from '@patterninc/pattern-ui/button-inset'
import { InputField } from '@patterninc/pattern-ui/fields/input-field'
import PlusIcon from '@patterninc/pattern-ui/icons/PlusIcon'
import TrashIcon from '@patterninc/pattern-ui/icons/TrashIcon'
import { ListItem } from '@patterninc/pattern-ui/list-item'
import React, { useState } from 'react'

import { stopMenuKeys } from '@/components/FilterMenu/FilterControl'
import { TagPair } from './Helper'

type TagFilterProps = {
  tags: TagPair[]
  onChange: (tags: TagPair[]) => void
}

const TagFilter = ({ tags, onChange }: TagFilterProps): React.JSX.Element => {
  const [draft, setDraft] = useState<TagPair>({ key: '', value: '' })
  const canAdd = draft.key.trim() !== '' && draft.value.trim() !== ''

  const addTag = () => {
    if (!canAdd) return
    onChange([...tags, { key: draft.key.trim(), value: draft.value.trim() }])
    setDraft({ key: '', value: '' })
  }

  const onFieldKeyDown = (event: React.KeyboardEvent) => {
    stopMenuKeys(event)
    if (event.key === 'Enter') addTag()
  }

  return (
    <div className='flex flex-col gap-2'>
      <span className='pui-text-label-small text-text-soft px-2.5 pt-1'>
        Tags
      </span>
      {tags.length > 0 ? (
        <div className='flex flex-col gap-0.5'>
          {tags.map((tag, index) => (
            <ListItem
              key={`${tag.key}:${tag.value}:${index}`}
              listItemType='small'
              listItemSuffix={
                <ButtonInset
                  icon={TrashIcon}
                  variant='inverse'
                  aria-label={`Remove tag ${tag.key}:${tag.value}`}
                  onClick={() => onChange(tags.filter((_, i) => i !== index))}
                  qaTestId='tag-filter-remove'
                />
              }
            >
              <span className='pui-text-mono-small'>
                {tag.key}:{tag.value}
              </span>
            </ListItem>
          ))}
        </div>
      ) : null}
      <div className='flex items-end gap-2'>
        <InputField
          className='min-w-0 flex-1'
          label='Key'
          variant='weak'
          value={draft.key}
          placeholder='Key'
          onChange={(event) => setDraft({ ...draft, key: event.target.value })}
          onKeyDown={onFieldKeyDown}
          qaTestId='tag-filter-key'
        />
        <InputField
          className='min-w-0 flex-1'
          label='Value'
          variant='weak'
          value={draft.value}
          placeholder='Value'
          onChange={(event) =>
            setDraft({ ...draft, value: event.target.value })
          }
          onKeyDown={onFieldKeyDown}
          qaTestId='tag-filter-value'
        />
      </div>
      <Button
        variant='weak'
        size='medium'
        prefix={{ icon: PlusIcon }}
        disabled={!canAdd}
        onClick={addTag}
        qaTestId='tag-filter-add'
      >
        Add tag
      </Button>
    </div>
  )
}

export default TagFilter
