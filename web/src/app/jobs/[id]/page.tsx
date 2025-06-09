'use client'
import dynamic from 'next/dynamic'
import { useEffect, useState } from 'react'

const JobDetails = dynamic(
  () => import('../../../modules/Jobs/JobDetails/JobDetails'),
  {
    ssr: false,
  },
)

type PageProps = {
  params: Promise<{ id: string }>
}

const JobDetailLayout = ({ params }: PageProps) => {
  const [id, setId] = useState<string>('')

  useEffect(() => {
    params.then((resolvedParams) => {
      setId(resolvedParams.id)
    })
  }, [params])

  if (!id) return null

  return <JobDetails id={id} />
}

export default JobDetailLayout
