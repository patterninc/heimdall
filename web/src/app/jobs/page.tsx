import dynamic from 'next/dynamic'

const Jobs = dynamic(() => import('@/modules/Jobs/Jobs'), {
  ssr: false,
})

const JobLayout = () => {
  return <Jobs />
}

export default JobLayout
