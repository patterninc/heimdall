import dynamic from 'next/dynamic'

const ClusterDetails = dynamic(
  () => import('../../../modules/Clusters/ClustersDetails/ClustersDetails'),
  {
    ssr: false,
  },
)

const ClusterDetailsLayout = ({ params }: { params: { id: string } }) => {
  return <ClusterDetails id={params.id} />
}

export default ClusterDetailsLayout
