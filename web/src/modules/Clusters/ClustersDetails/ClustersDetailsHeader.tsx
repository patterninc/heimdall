import { useMemo } from 'react'
import { ClusterDataTypeProps } from './ClusterInformationPane'

import styles from './_clustersDetails.module.scss'
import { ListLoading, PageHeader, SectionHeader } from '@patterninc/react-ui'
import ApiResponseButton from '@/components/ApiResponseButton/ApiResponseButton'

const ClustersDetailsHeader = ({
  clusterData,
  isLoading,
}: ClusterDataTypeProps): React.JSX.Element => {
  const data = useMemo(() => clusterData?.[0], [clusterData])

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
