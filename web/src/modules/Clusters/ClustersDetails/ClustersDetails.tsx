'use client'

import { useQuery } from '@tanstack/react-query'
import React from 'react'

import { getCluster } from '@/app/api/clusters/clusters'
import { resourceStatusVariant } from '@/common/Services/status'
import DetailPage from '@/components/DetailPage/DetailPage'
import ExternalLinkButton from '@/components/ExternalLinkButton/ExternalLinkButton'
import { ClusterType } from '../Helper'
import ClusterInformationPane from './ClusterInformationPane'

type ClusterDetailsProps = {
  /** Cluster ID or name — job details link here by name. */
  id: string
}

const ClustersDetails = ({ id }: ClusterDetailsProps): React.JSX.Element => {
  const { data, isPending, isError } = useQuery<ClusterType[]>({
    queryKey: ['cluster', id],
    queryFn: () => getCluster(id),
  })
  const cluster = data?.[0]

  return (
    <DetailPage
      backHref='/clusters'
      backLabel='Back to clusters'
      resourceLabel='cluster'
      title={cluster?.name || id}
      subtitle={
        cluster?.name && cluster.name !== cluster.id ? cluster.id : undefined
      }
      status={cluster?.status}
      statusVariant={resourceStatusVariant(cluster?.status)}
      isError={isError || (!isPending && !cluster)}
      qaTestId='cluster-details'
      right={
        <ExternalLinkButton
          href={`/api/v1/clusters?id=${cluster?.id ?? id}`}
          qaTestId='cluster-details-api-response'
        >
          API response
        </ExternalLinkButton>
      }
    >
      <ClusterInformationPane cluster={cluster} isLoading={isPending} />
    </DetailPage>
  )
}

export default ClustersDetails
