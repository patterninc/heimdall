'use client'

import { useQuery } from '@tanstack/react-query'
import React from 'react'

import { getJobDetails } from '@/app/api/jobs/jobs'
import { jobStatusVariant } from '@/common/Services/status'
import DetailPage from '@/components/DetailPage/DetailPage'
import ExternalLinkButton from '@/components/ExternalLinkButton/ExternalLinkButton'
import { JobType } from '../Helper'
import CancelJobButton from './CancelJobButton'
import JobInformationPane from './JobInformationPane'

type JobDetailsProp = {
  id: string
}

const JobDetails = ({ id }: JobDetailsProp): React.JSX.Element => {
  const {
    data: jobData,
    isPending,
    isError,
  } = useQuery<JobType>({
    queryKey: ['job', id],
    queryFn: () => getJobDetails(id),
  })

  return (
    <DetailPage
      backHref='/jobs'
      backLabel='Back to jobs'
      resourceLabel='job'
      title={jobData?.name || id}
      subtitle={jobData?.name && jobData.name !== id ? id : undefined}
      status={jobData?.status}
      statusVariant={jobStatusVariant(jobData?.status)}
      isError={isError}
      qaTestId='job-details'
      right={
        <>
          <ExternalLinkButton
            href={`/api/v1/job/${id}`}
            qaTestId='job-details-api-response'
          >
            API response
          </ExternalLinkButton>
          {jobData ? <CancelJobButton job={jobData} /> : null}
        </>
      }
    >
      <JobInformationPane jobData={jobData} isLoading={isPending} />
    </DetailPage>
  )
}

export default JobDetails
