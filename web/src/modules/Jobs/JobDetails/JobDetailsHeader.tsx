'use client'

import { PageHeader, SectionHeader } from '@patterninc/react-ui'
import styles from './_job-Details.module.scss'
import { JobDataTypesProps } from '../Helper'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { useEffect, useState } from 'react'
import yaml from 'js-yaml'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const JobDetailsHeader = ({ jobData }: JobDataTypesProps): React.JSX.Element => {
  const [sql, setSql] = useState<string>()
  const [yamlContext, setYamlContext] = useState<string>()
  const [tags, setTags] = useState<string>()
  const [commandCriteria, setCommandCriteria] = useState<string>()
  const [clusterCriteria, setClusterCriteria] = useState<string>()

  useEffect(() => {
    if (jobData?.context) {
      const rawSql = jobData.context?.query
      setSql(rawSql)
      const yamlString = yaml.dump(jobData.context?.properties)
      setYamlContext(yamlString)
    }
    if (jobData?.tags) {
      const tags = yaml.dump(jobData.tags)
      setTags(tags)
    }
    if (jobData?.command_criteria) {
      const commandCriteria = yaml.dump(jobData.command_criteria)
      setCommandCriteria(commandCriteria)
    }
    if (jobData?.cluster_criteria) {
      const clusterCriteria = yaml.dump(jobData.cluster_criteria)
      setClusterCriteria(clusterCriteria)
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
              {yamlContext ? (
                <div>
                  <SectionHeader title='Context' />
                  <SyntaxHighlighter style={github}>
                    {yamlContext}
                  </SyntaxHighlighter>
                </div>
              ) : null}

              {tags ? (
                <div>
                  <SectionHeader title='Tags' />
                  <SyntaxHighlighter style={github}>{tags}</SyntaxHighlighter>
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
                {commandCriteria ? (
                  <>
                    <SectionHeader title='Command Criteria' />
                    <SyntaxHighlighter style={github}>
                      {commandCriteria}
                    </SyntaxHighlighter>
                  </>
                ) : null}
              </div>
              <div>
                {clusterCriteria ? (
                  <>
                    <SectionHeader title='Cluster Criteria' />
                    <SyntaxHighlighter style={github}>
                      {clusterCriteria}
                    </SyntaxHighlighter>
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
