'use client'
import dynamic from 'next/dynamic'
import { use } from 'react'

const JobDetails = dynamic(
  () => import('../../../modules/Jobs/JobDetails/JobDetails'),
  {
    ssr: false,
  },
)

const JobDetailLayout = (props: { params: Promise<{ id: string }> }) => {
  const params = use(props.params)
  return <JobDetails id={params.id} />
}
export default JobDetailLayout
