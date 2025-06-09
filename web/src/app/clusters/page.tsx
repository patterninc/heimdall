'use client'
import dynamic from 'next/dynamic'

const Cluster = dynamic(() => import('../../modules/Clusters/Clusters'), {
  ssr: false,
})

const ClusterLayout = () => {
  return <Cluster />
}
export default ClusterLayout
