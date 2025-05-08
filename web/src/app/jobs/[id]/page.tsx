import dynamic from 'next/dynamic'

const JobDetails = dynamic(
  () => import('../../../modules/Jobs/JobDetails/JobDetails'),
  {
    ssr: false,
  },
)

const JobDetailLayout = ({ params }: { params: { id: string } }) => {
  const { id } = params

  return <JobDetails id={id} />
}
export default JobDetailLayout
