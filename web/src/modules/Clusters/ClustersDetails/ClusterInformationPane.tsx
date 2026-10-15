import DetailSection from '@/components/DetailPage/DetailSection'
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
      <DetailSection label='Overview' qaTestId='cluster-details-overview'>
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
      </DetailSection>

      {hasRuntime ? (
        <DetailSection label='Runtime' qaTestId='cluster-details-runtime'>
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
        </DetailSection>
      ) : null}

      {properties.length > 0 ? (
        <DetailSection
          label='Context properties'
          count={properties.length}
          qaTestId='cluster-details-properties'
        >
          <KeyValueList entries={properties} />
        </DetailSection>
      ) : null}

      {cluster?.tags?.length ? (
        <DetailSection
          label='Tags'
          count={cluster.tags.length}
          qaTestId='cluster-details-tags'
        >
          <TagList values={cluster.tags} />
        </DetailSection>
      ) : null}
    </>
  )
}

export default ClusterInformationPane
