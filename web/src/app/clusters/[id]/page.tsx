'use client'
import dynamic from 'next/dynamic'
import { use } from 'react'

const ClusterDetails = dynamic(
  () => import('../../../modules/Clusters/ClustersDetails/ClustersDetails'),
  {
    ssr: false,
    loading: () => <div>Loading cluster details...</div>,
  },
)

const ClusterDetailsLayout = ({
  params,
}: {
  params: Promise<{ id: string }>
}) => {
  const resolvedParams = use(params)
  const id = resolvedParams.id

  if (!id) return <div>No cluster ID provided</div>

  return <ClusterDetails id={id} />
}

export default ClusterDetailsLayout
