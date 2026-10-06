import ArrowTopRightSquareIcon from '@patterninc/pattern-ui/icons/ArrowTopRightSquareIcon'
import { ListItem } from '@patterninc/pattern-ui/list-item'
import React from 'react'

export type ExternalLink = {
  label: string
  href: string
}

type ExternalLinkListProps = {
  links: ExternalLink[]
  qaTestId: string
}

const ExternalLinkList = ({
  links,
  qaTestId,
}: ExternalLinkListProps): React.JSX.Element => (
  <div className='flex flex-col p-1'>
    {links.map((link) => (
      <ListItem
        key={link.label}
        asChild
        listItemType='medium'
        listItemSuffix={<ArrowTopRightSquareIcon size={16} color='icon-sub' />}
        qaTestId={`${qaTestId}-${link.label.toLowerCase().replace(/\s+/g, '-')}`}
      >
        <a href={link.href} target='_blank' rel='noreferrer'>
          {link.label}
        </a>
      </ListItem>
    ))}
  </div>
)

export default ExternalLinkList
