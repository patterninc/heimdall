'use client'
import dynamic from 'next/dynamic'
import { useEffect, useState } from 'react'

const ClusterDetails = dynamic(
  () => import('../../../modules/Clusters/ClustersDetails/ClustersDetails'),
  {
    ssr: false,
  },
)

const ClusterDetailsLayout = ({
  params,
}: {
  params: Promise<{ id: string }>
}) => {
  const [id, setId] = useState<string>('')

  useEffect(() => {
    params.then((resolvedParams) => {
      setId(resolvedParams.id)
    })
  }, [params])

  if (!id) return null

  return <ClusterDetails id={id} />
}

export default ClusterDetailsLayout
