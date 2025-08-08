import { useMemo } from 'react'
import { ClusterType } from '../Helper'
import { InformationPane, ListLoading } from '@patterninc/react-ui'
import { formatDateWithTimeZone, myTimezone } from '@/common/Services'

export type ClusterDataTypeProps = {
  clusterData: ClusterType[]
  isLoading?: boolean
}

const ClusterInformationPane = ({
  clusterData,
  isLoading,
}: ClusterDataTypeProps): React.JSX.Element => {
  const data = useMemo(() => clusterData?.[0], [clusterData])
  const clusterDetailsData = [
    {
      label: 'User',
      data: data?.user,
      check: !!data?.user,
    },
    {
      label: 'Version',
      data: data?.version,
      check: !!data?.version,
    },
    {
      label: 'Created At',
      data: formatDateWithTimeZone(data?.created_at, myTimezone),
      check: !!data?.created_at,
    },
    {
      label: 'Updated At',
      data: formatDateWithTimeZone(data?.updated_at, myTimezone),
      check: !!data?.updated_at,
    },
    {
      label: 'Description',
      data: data?.description,
      check: !!data?.description,
    },
  ]
  return (
    <div className='min-w-[300px]'>
      <InformationPane
        header={{
          labelAndData: {
            label: 'Cluster Name',
            data: <span>{data?.name}</span>,
            check: true,
          },
          tag: {
            color:
              data?.status === 'ACTIVE'
                ? 'green'
                : data?.status === 'INACTIVE'
                  ? 'gray'
                  : 'red',
            children: data?.status,
          },
        }}
      >
        {isLoading ? (
          <ListLoading />
        ) : (
          <>
            <InformationPane.Section
              data={[
                {
                  label: 'Cluster ID',
                  data: data?.id,
                  check: !!data?.id,
                },
              ]}
            />
            <InformationPane.Divider />
            <InformationPane.Section data={clusterDetailsData} isTwoColumns />
          </>
        )}
      </InformationPane>
    </div>
  )
}

export default ClusterInformationPane
