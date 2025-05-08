import { JobDataTypesProps } from '../Helper'
import { HeaderMetricGroup } from '@patterninc/react-ui'

const JobContext = ({ jobData, isLoading }: JobDataTypesProps) => {
  if (!jobData) return null
  const propertiesArray = jobData?.context?.properties
    ? Object?.entries(jobData.context.properties ?? []).map(([key, value]) => ({
        title: key,
        secondaryInfo: value,
      }))
    : []
  const emptyPropertiesArray = [
    { title: 'spark.driver.core', secondaryInfo: undefined },
    { title: 'spark.driver.memory', secondaryInfo: undefined },
    { title: 'spark.executor.cores', secondaryInfo: undefined },
    { title: 'spark.executor.instances', secondaryInfo: undefined },
    { title: 'spark.executor.memory', secondaryInfo: undefined },
  ]
  return (
    <HeaderMetricGroup
      data={jobData?.context ? propertiesArray : emptyPropertiesArray}
      loading={isLoading}
    />
  )
}
export default JobContext
