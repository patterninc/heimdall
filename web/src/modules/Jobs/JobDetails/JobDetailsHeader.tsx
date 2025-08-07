'use client'

import { Alert, PageHeader, SectionHeader } from '@patterninc/react-ui'
import { JobDataTypesProps } from '../Helper'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const JobDetailsHeader = ({
  jobData,
}: JobDataTypesProps): React.JSX.Element => {
  return (
    <div className='w-full'>
      <PageHeader
        rightSectionChildren={
          <ApiResponseButton link={`/api/v1/job/${jobData?.id}`} />
        }
        bottomSectionChildren={
          <div
            className='rounded overflow-auto bg-white border-t border-medium-purple'
          >
            <div className='flex flex-col gap-4 p-4'>
              {jobData?.status === 'FAILED' ? (
                <Alert type='error' text={jobData?.error} />
              ) : null}
              {jobData?.context?.properties ? (
                <div>
                  <SectionHeader title='Context' />
                  <ul>
                    {Object.entries(jobData?.context?.properties || {}).map(
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
                    {jobData?.tags.map((value) => <li key={value}>{value}</li>)}
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
    </div>
  )
}

export default JobDetailsHeader
