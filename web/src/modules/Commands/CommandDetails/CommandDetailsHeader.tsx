'use client'

import {
  Button,
  ListLoading,
  PageHeader,
  SectionHeader,
} from '@patterninc/react-ui'
import { CommandDetailsProps } from './CommandInformationPane'

import React, { useMemo } from 'react'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const CommandDetailsHeader = ({
  commandData,
  isLoading,
}: CommandDetailsProps): React.JSX.Element => {
  const data = useMemo(() => commandData[0], [commandData])

  return (
    <div className='w-full'>
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
              className='rounded overflow-auto bgc-white pat-border-t bdrc-medium-purple'
            >
              <div className='flex flex-direction-column pat-gap-4 pat-p-4'>
                {data?.context?.properties ? (
                  <div>
                    <SectionHeader title='Context' />
                    <ul>
                      {Object.entries(data?.context?.properties || {}).map(
                        ([key, value]) => (
                          <li key={key}>
                            {key}:{String(value)}
                          </li>
                        ),
                      )}
                    </ul>
                  </div>
                ) : null}
                {data?.tags ? (
                  <div>
                    <SectionHeader title='Tags' />
                    <ul>
                      {data?.tags.map((value: string) => (
                        <li key={value}>{value}</li>
                      ))}
                    </ul>
                  </div>
                ) : null}
                {data?.cluster_tags ? (
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
