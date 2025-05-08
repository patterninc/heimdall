'use client'

import {
  Button,
  ListLoading,
  PageHeader,
  SectionHeader,
} from '@patterninc/react-ui'
import { CommandDetailsProps } from './CommandInformationPane'
import styles from './_commandDetails.module.scss'
import { useEffect, useMemo, useState } from 'react'
import yaml from 'js-yaml'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const CommandDetailsHeader = ({
  commandData,
  isLoading,
}: CommandDetailsProps): JSX.Element => {
  const [yamlContext, setYamlContext] = useState<string>()
  const [tags, setTags] = useState<string>()
  const [clusterTag, setClusterTag] = useState<string>()

  const data = useMemo(() => commandData[0], [commandData])

  useEffect(() => {
    if (data?.context) {
      const yamlString = yaml?.dump(data.context?.properties)
      setYamlContext(yamlString)
    }
    if (data?.tags) {
      const tags = yaml?.dump(data.tags)
      setTags(tags)
    }

    if (data?.cluster_tags) {
      setClusterTag(yaml?.dump(data.cluster_tags))
    }
  }, [commandData, data?.cluster_tags, data?.context, data?.tags])
  return (
    <div className={`${styles.commandDetailsHeader}`}>
      <PageHeader
        rightSectionChildren={
          <ApiResponseButton
            link={`/api/v1/commands?id=${commandData?.[0]?.id}`}
          />
        }
        bottomSectionChildren={
          isLoading ? (
            <ListLoading />
          ) : (
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
                {clusterTag && data.cluster_tags ? (
                  <div>
                    <SectionHeader title='Cluster Tags' />
                    <ul>
                      {data?.cluster_tags.map((tag: string) => (
                        <li key={tag}>{tag}</li>
                      ))}
                    </ul>
                  </div>
                ) : null}

                {data?.context?.queries_uri ? (
                  <div>
                    <SectionHeader title='Logs' />
                    <ul>
                      <li>
                        <Button
                          styleType='text-blue'
                          as='externalLink'
                          href={data?.context?.logs_uri}
                        >
                          Logs
                        </Button>
                      </li>
                    </ul>
                  </div>
                ) : null}
                {data?.context?.results_uri ? (
                  <div>
                    <SectionHeader title='Queries' />
                    <ul>
                      <li>
                        <Button
                          styleType='text-blue'
                          as='externalLink'
                          href={data?.context?.queries_uri}
                        >
                          Queries
                        </Button>
                      </li>
                    </ul>
                  </div>
                ) : null}
                {data?.context?.results_uri ? (
                  <div>
                    <SectionHeader title='Results' />
                    <ul>
                      <li>
                        <Button
                          styleType='text-blue'
                          as='externalLink'
                          href={data?.context?.results_uri}
                        >
                          Results
                        </Button>
                      </li>
                    </ul>
                  </div>
                ) : null}
              </div>
            </div>
          )
        }
        header={{
          name: 'Command Details',
          value: <></>,
        }}
      />
    </div>
  )
}
export default CommandDetailsHeader
