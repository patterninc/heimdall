'use client'

import { getCluster } from '@/app/api/clusters/clusters'
import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { useQuery } from '@tanstack/react-query'
import { useContext, useEffect } from 'react'
import ClusterInformationPane from './ClusterInformationPane'
import ClustersDetailsHeader from './ClustersDetailsHeader'

type ClusterDetailsProps = {
  id: string
}

const ClustersDetails = ({ id }: ClusterDetailsProps): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)
  useEffect(() => {
    updateBreadcrumbs({
      name: id,
      link: `/clusters/${id}`,
    })
  }, [id, updateBreadcrumbs])
  const { data, isLoading } = useQuery({
    queryKey: ['cluster', id],
    queryFn: () => getCluster(id),
  })

  return (
    <div className='flex min-w-[300px] flex-col gap-4 pt-4 md:flex-row'>
      <ClusterInformationPane clusterData={data} isLoading={isLoading} />
      <ClustersDetailsHeader clusterData={data} isLoading={isLoading} />
    </div>
  )
}

export default ClustersDetails
