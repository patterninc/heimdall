'use client'
import dynamic from 'next/dynamic'
import { use } from 'react'

const ClusterDetails = dynamic(
  () => import('../../../modules/Clusters/ClustersDetails/ClustersDetails'),
  {
    ssr: false,
  },
)

const ClusterDetailsLayout = (props: { params: Promise<{ id: string }> }) => {
  const params = use(props.params)
  return <ClusterDetails id={params.id} />
}

export default ClusterDetailsLayout
