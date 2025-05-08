import { useEffect, useMemo, useState } from 'react'
import { ClusterDataTypeProps } from './ClusterInformationPane'
import yaml from 'js-yaml'
import styles from './_clustersDetails.module.scss'
import { ListLoading, PageHeader, SectionHeader } from '@patterninc/react-ui'
import SyntaxHighlighter from 'react-syntax-highlighter'
import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const ClustersDetailsHeader = ({
  clusterData,
  isLoading,
}: ClusterDataTypeProps): JSX.Element => {
  const [yamlContext, setYamlContext] = useState<string>()
  const [tags, setTags] = useState<string>()

  const data = useMemo(() => clusterData?.[0], [clusterData])

  useEffect(() => {
    if (data?.context) {
      const yamlString = yaml?.dump(data.context)
      setYamlContext(yamlString)
    }
    if (data?.tags) {
      const tags = yaml?.dump(data.tags)
      setTags(tags)
    }
  }, [clusterData, data?.context, data?.tags])
  return (
    <div className={styles.clustersDetailsHeader}>
      <PageHeader
        rightSectionChildren={
          <ApiResponseButton
            link={`/api/v1/clusters?id=${clusterData?.[0]?.id}`}
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
              </div>
            </div>
          )
        }
        header={{
          name: 'Cluster Details',
          value: <></>,
        }}
      />
    </div>
  )
}
export default ClustersDetailsHeader
