import DetailSection from '@/components/DetailPage/DetailSection'
import React from 'react'

import { contextEntries, contextString } from '@/common/api/context'
import type { Cluster } from '@/common/api/types'
import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import DetailFields from '@/components/DetailPage/DetailFields'
import KeyValueList from '@/components/DetailPage/KeyValueList'
import TagList from '@/components/DetailPage/TagList'

type ClusterInformationPaneProps = {
  cluster?: Cluster
  isLoading: boolean
}

const ClusterInformationPane = ({
  cluster,
  isLoading,
}: ClusterInformationPaneProps): React.JSX.Element => {
  const context = cluster?.context
  const properties = contextEntries(context, 'properties')
  const emrReleaseLabel = contextString(context, 'emr_release_label')
  const roleArn = contextString(context, 'role_arn')
  const executionRoleArn = contextString(context, 'execution_role_arn')
  const hasRuntime = Boolean(emrReleaseLabel || roleArn || executionRoleArn)

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
                value: emrReleaseLabel,
                mono: true,
              },
              { label: 'Role ARN', value: roleArn, mono: true },
              {
                label: 'Execution role ARN',
                value: executionRoleArn,
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
