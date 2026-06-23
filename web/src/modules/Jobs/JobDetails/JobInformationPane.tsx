import {
  Button,
  Icon,
  InformationPane,
  ListLoading,
  TrimText,
} from '@patterninc/react-ui'
import { getStatusColor, JobDataTypesProps } from '../Helper'
import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import { useRouter } from 'next/navigation'
const JobInformationPane = ({
  jobData,
  isLoading,
}: JobDataTypesProps): React.JSX.Element => {
  const router = useRouter()
  const jobdetailsData = [
    { label: 'Version', data: jobData?.version, check: !!jobData?.version },
    {
      label: 'Command Name',
      data: (
        <Button
          as='unstyled'
          onClick={() => router.push(`/commands/${jobData?.command_name}`)}
        >
          {jobData?.command_name}
        </Button>
      ),
      check: !!jobData?.command_name,
    },

    {
      label: 'Cluster Name',
      data: (
        <Button
          as='unstyled'
          onClick={() => router.push(`/clusters/${jobData?.cluster_name}`)}
        >
          {jobData?.cluster_name}
        </Button>
      ),
      check: !!jobData?.cluster_name,
    },
    {
      label: 'Created At',
      data: jobData?.created_at
        ? formatDateWithTimeZone(jobData?.created_at, myTimezone)
        : undefined,
      check: !!jobData?.created_at,
    },
    {
      label: 'Updated At',
      data: jobData?.updated_at
        ? formatDateWithTimeZone(jobData?.updated_at, myTimezone)
        : undefined,
      check: !!jobData?.updated_at,
    },
  ]
  return (
    <div className='min-w-[300px]'>
      <InformationPane
        header={{
          labelAndData: {
            label: 'Job Name',
            data: (
              <TrimText
                text={jobData?.name ?? ''}
                limit={25}
                customClass='whitespace-nowrap'
              />
            ),
            check: !!jobData?.name,
          },
          tag: {
            color: getStatusColor(jobData?.status || ''),
            children: jobData?.status,
          },
        }}
      >
        {isLoading ? (
          <ListLoading />
        ) : (
          <div>
            <InformationPane.Section
              data={[
                { label: 'User', data: jobData?.user, check: !!jobData?.user },
                { label: 'Job ID', data: jobData?.id, check: !!jobData?.id },
              ]}
            />
            <InformationPane.Divider />
            <InformationPane.Section data={jobdetailsData} isTwoColumns />
            <InformationPane.Divider />
            <InformationPane.Section
              data={[
                {
                  label: '',
                  data: (
                    <Button
                      styleType='text-blue'
                      as='externalLink'
                      href={`/api/v1/job/${jobData?.id}/stdout`}
                      className='gap-1'
                    >
                      <span>View Stdout</span>
                      <Icon icon='launch' color='dark-blue' iconSize='12px' />
                    </Button>
                  ),
                  check: !!jobData?.id,
                },
                {
                  label: '',
                  data: (
                    <Button
                      styleType='text-red'
                      as='externalLink'
                      href={`/api/v1/job/${jobData?.id}/stderr`}
                      className='gap-1'
                    >
                      <span>View Stderr</span>
                      <Icon icon='launch' color='dark-red' iconSize='12px' />
                    </Button>
                  ),
                  check: !!jobData?.id,
                },
              ]}
              isTwoColumns
            />
            {!!jobData?.spark_application_id && (
              <InformationPane.Section
                data={[
                  {
                    label: '',
                    data: (
                      <Button
                        styleType='text-blue'
                        as='externalLink'
                        href={`https://spark-history.data-platform.aws.pattern.com/history/${jobData?.spark_application_id}/jobs/`}
                        className='gap-1'
                      >
                        <span>Spark History</span>
                        <Icon icon='launch' color='dark-blue' iconSize='12px' />
                      </Button>
                    ),
                    check: true,
                  },
                ]}
              />
            )}
          </div>
        )}
      </InformationPane>
    </div>
  )
}
export default JobInformationPane
