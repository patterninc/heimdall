import { Section } from '@patterninc/pattern-ui/section'
import React from 'react'

import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import DetailFields from '@/components/DetailPage/DetailFields'
import ExternalLinkList, {
  type ExternalLink,
} from '@/components/DetailPage/ExternalLinkList'
import KeyValueList from '@/components/DetailPage/KeyValueList'
import TagList from '@/components/DetailPage/TagList'
import { CommandType } from '../Helper'

type CommandInformationPaneProps = {
  command?: CommandType
  isLoading: boolean
}

const CommandInformationPane = ({
  command,
  isLoading,
}: CommandInformationPaneProps): React.JSX.Element => {
  const context = command?.context
  const links: ExternalLink[] = [
    { label: 'Logs', href: context?.logs_uri },
    { label: 'Queries', href: context?.queries_uri },
    { label: 'Results', href: context?.results_uri },
  ].filter((link): link is ExternalLink => Boolean(link.href))
  const properties = Object.entries(context?.properties ?? {}).map(
    ([key, value]): [string, string] => [key, String(value)],
  )

  return (
    <>
      <Section label='Overview' qaTestId='command-details-overview'>
        <DetailFields
          isLoading={isLoading}
          fields={[
            { label: 'Command ID', value: command?.id, mono: true },
            { label: 'User', value: command?.user },
            { label: 'Plugin', value: command?.plugin },
            { label: 'Version', value: command?.version },
            {
              label: 'Execution',
              value:
                command?.is_sync === undefined
                  ? undefined
                  : command.is_sync
                    ? 'Sync'
                    : 'Async',
            },
            {
              label: 'Created at',
              value: command?.created_at
                ? formatDateWithTimeZone(command.created_at, myTimezone)
                : undefined,
            },
            {
              label: 'Updated at',
              value: command?.updated_at
                ? formatDateWithTimeZone(command.updated_at, myTimezone)
                : undefined,
            },
            { label: 'Description', value: command?.description },
          ]}
        />
      </Section>

      {links.length > 0 ? (
        <Section label='Links' qaTestId='command-details-links'>
          <ExternalLinkList links={links} qaTestId='command-details-link' />
        </Section>
      ) : null}

      {properties.length > 0 ? (
        <Section
          label='Context properties'
          count={properties.length}
          qaTestId='command-details-properties'
        >
          <KeyValueList entries={properties} />
        </Section>
      ) : null}

      {command?.tags?.length ? (
        <Section
          label='Tags'
          count={command.tags.length}
          qaTestId='command-details-tags'
        >
          <TagList values={command.tags} />
        </Section>
      ) : null}

      {command?.cluster_tags?.length ? (
        <Section
          label='Cluster tags'
          count={command.cluster_tags.length}
          qaTestId='command-details-cluster-tags'
        >
          <TagList values={command.cluster_tags} />
        </Section>
      ) : null}
    </>
  )
}

export default CommandInformationPane
