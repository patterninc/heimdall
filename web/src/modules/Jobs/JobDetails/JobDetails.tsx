'use client'

import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { useContext, useEffect } from 'react'
import JobInformationPane from './JobInformationPane'
import styles from './_job-Details.module.scss'
import { JobType } from '../Helper'

import JobDetailsHeader from './JobDetailsHeader'
import { useQuery } from '@tanstack/react-query'
import { getJobDetails } from '@/app/api/jobs/jobs'
import { usePathname } from 'next/navigation'

type JobDetailsProp = {
  id: string
}

const JobDetails = ({ id }: JobDetailsProp): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)
  const pathname = usePathname()

  const { data: jobData, isLoading } = useQuery<JobType>({
    queryKey: ['job', id],
    queryFn: () => getJobDetails(id),
  })

  useEffect(() => {
    updateBreadcrumbs({
      name: id,
      link: pathname,
    })
  }, [updateBreadcrumbs, id, pathname])

  return (
    <div className={styles.jobDetailsContainer}>
      <JobInformationPane jobData={jobData} isLoading={isLoading} />
      <JobDetailsHeader jobData={jobData} isLoading={isLoading} />
    </div>
  )
}
export default JobDetails
