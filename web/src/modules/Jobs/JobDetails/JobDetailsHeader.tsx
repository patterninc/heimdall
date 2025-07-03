'use client'

import { Alert, PageHeader, SectionHeader } from '@patterninc/react-ui'
import styles from './_job-Details.module.scss'
import { JobDataTypesProps } from '../Helper'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { useEffect, useState } from 'react'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const JobDetailsHeader = ({
  jobData,
}: JobDataTypesProps): React.JSX.Element => {
  const [sql, setSql] = useState<string>()

  useEffect(() => {
    if (jobData?.context) {
      const rawSql = jobData.context?.query
      setSql(rawSql)
    }
  }, [jobData])

  return (
    <div className={styles.pageHeaderContainer}>
      <PageHeader
        rightSectionChildren={
          <ApiResponseButton link={`/api/v1/job/${jobData?.id}`} />
        }
        bottomSectionChildren={
          <div
            className={`${styles.bottomSectionContainer} bgc-white pat-border-t bdrc-medium-purple`}
          >
            <div className='flex flex-direction-column pat-gap-4 pat-p-4'>
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
                {sql ? (
                  <>
                    <SectionHeader title='SQL Query' />
                    <SyntaxHighlighter language='sql' style={github}>
                      {sql}
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
