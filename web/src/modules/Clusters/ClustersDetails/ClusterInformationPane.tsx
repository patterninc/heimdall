import { Section } from '@patterninc/pattern-ui/section'
import React from 'react'

import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import DetailFields from '@/components/DetailPage/DetailFields'
import KeyValueList from '@/components/DetailPage/KeyValueList'
import TagList from '@/components/DetailPage/TagList'
import { ClusterType } from '../Helper'

type ClusterInformationPaneProps = {
  cluster?: ClusterType
  isLoading: boolean
}

const ClusterInformationPane = ({
  cluster,
  isLoading,
}: ClusterInformationPaneProps): React.JSX.Element => {
  const context = cluster?.context
  const properties = Object.entries(context?.properties ?? {}).map(
    ([key, value]): [string, string] => [key, String(value)],
  )
  const hasRuntime = Boolean(
    context?.emr_release_label ||
    context?.role_arn ||
    context?.execution_role_arn,
  )

  return (
    <>
      <Section label='Overview' qaTestId='cluster-details-overview'>
        <DetailFields
          isLoading={isLoading}
          fields={[
            { label: 'Cluster ID', value: cluster?.id, mono: true },
            { label: 'User', value: cluster?.user },
            { label: 'Version', value: cluster?.version },
            {
              label: 'Created at',
              value: cluster?.created_at
                ? formatDateWithTimeZone(cluster.created_at, myTimezone)
                : undefined,
            },
            {
              label: 'Updated at',
              value: cluster?.updated_at
                ? formatDateWithTimeZone(cluster.updated_at, myTimezone)
                : undefined,
            },
            { label: 'Description', value: cluster?.description },
          ]}
        />
      </Section>

      {hasRuntime ? (
        <Section label='Runtime' qaTestId='cluster-details-runtime'>
          <DetailFields
            fields={[
              {
                label: 'EMR release',
                value: context?.emr_release_label,
                mono: true,
              },
              { label: 'Role ARN', value: context?.role_arn, mono: true },
              {
                label: 'Execution role ARN',
                value: context?.execution_role_arn,
                mono: true,
              },
            ]}
          />
        </Section>
      ) : null}

      {properties.length > 0 ? (
        <Section
          label='Context properties'
          count={properties.length}
          qaTestId='cluster-details-properties'
        >
          <KeyValueList entries={properties} />
        </Section>
      ) : null}

      {cluster?.tags?.length ? (
        <Section
          label='Tags'
          count={cluster.tags.length}
          qaTestId='cluster-details-tags'
        >
          <TagList values={cluster.tags} />
        </Section>
      ) : null}
    </>
  )
}

export default ClusterInformationPane
