'use client'

import {
  Alert,
  Ellipsis,
  PageFooter,
  PageHeader,
  SectionHeader,
  toast,
} from '@patterninc/react-ui'
import { JobDataTypesProps } from '../Helper'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'
import { cancelJob } from '@/app/api/jobs/jobs'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useMemo } from 'react'

const CANCELABLE_STATUSES = ['NEW', 'ACCEPTED', 'RUNNING']

const JobDetailsHeader = ({
  jobData,
}: JobDataTypesProps): React.JSX.Element => {
  const queryClient = useQueryClient()
  const cancelMutation = useMutation({
    mutationFn: (id: string) => cancelJob(id),
    onSuccess: (response) => {
      if (response.status === 'CANCELING') {
        toast({
          type: 'info',
          message: `Job is being canceled...`,
        })
        queryClient.invalidateQueries({ queryKey: ['job', jobData?.id] })
      } else {
        toast({
          type: 'error',
          message: response.error || 'Failed to cancel job',
        })
      }
    },
    onError: () => {
      toast({
        type: 'error',
        message: 'Failed to cancel job',
      })
    },
  })

  // Only async jobs in active states can be canceled
  const isCancelable = useMemo(
    () =>
      CANCELABLE_STATUSES.includes(jobData?.status ?? '') &&
      jobData?.is_sync !== true,
    [jobData?.status, jobData?.is_sync],
  )

  const isCanceling = jobData?.status === 'CANCELING'

  return (
    <div className='w-full'>
      <PageHeader
        rightSectionChildren={
          <ApiResponseButton link={`/api/v1/job/${jobData?.id}`} />
        }
        bottomSectionChildren={
          <div className='border-medium-purple overflow-auto rounded border-t bg-white'>
            <div className='flex flex-col gap-4 p-4'>
              {jobData?.status === 'FAILED' ? (
                <Alert type='error' text={jobData?.error} />
              ) : null}
              {jobData?.context?.parameters?.properties ? (
                <div>
                  <SectionHeader title='Context' />
                  <ul>
                    {Object.entries(jobData?.context?.parameters?.properties || {}).map(
                      ([key, value]) => (
                        <li key={key}>
                          {key}:{value}
                        </li>
                      ),
                    )}
                  </ul>
                </div>
              ) : null}

              {jobData?.tags ? (
                <div>
                  <SectionHeader title='Tags' />
                  <ul>
                    {jobData?.tags.map((value) => (
                      <li key={value}>{value}</li>
                    ))}
                  </ul>
                </div>
              ) : null}
              <div>
                {jobData?.context?.query ? (
                  <>
                    <SectionHeader title='SQL Query' />
                    <SyntaxHighlighter language='sql' style={github}>
                      {jobData.context.query}
                    </SyntaxHighlighter>
                  </>
                ) : null}
              </div>
              <div>
                {jobData?.command_criteria ? (
                  <>
                    <SectionHeader title='Command Criteria' />
                    <ul>
                      {jobData?.command_criteria.map((value) => (
                        <li key={value}>{value}</li>
                      ))}
                    </ul>
                  </>
                ) : null}
              </div>
              <div>
                {jobData?.cluster_criteria ? (
                  <>
                    <SectionHeader title='Cluster Criteria' />
                    <ul>
                      {jobData?.cluster_criteria.map((value) => (
                        <li key={value}>{value}</li>
                      ))}
                    </ul>
                  </>
                ) : null}
              </div>
            </div>
          </div>
        }
        header={{
          name: 'Job Details',
          value: <></>,
        }}
      />
      <PageFooter
        rightSection={[
          {
            as: 'confirmation',
            confirmation: {
              header: 'Cancel Job',
              body: 'Are you sure you want to cancel this job?',
              confirmCallout: () => {
                if (jobData?.id) cancelMutation.mutate(jobData.id)
              },
              type: 'red',
            },
            children: isCanceling ? (
              <span>
                Canceling job
                <Ellipsis />
              </span>
            ) : (
              'Cancel Job'
            ),
            styleType: 'primary-red',
            type: 'button',
            disabled: !isCancelable || isCanceling,
          },
        ]}
      />
    </div>
  )
}

export default JobDetailsHeader
