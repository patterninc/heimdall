'use client'

import { InformationPane, ListLoading } from '@patterninc/react-ui'
import { CommandType } from '../Helper'
import styles from './_commandDetails.module.scss'

import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import { useMemo } from 'react'

export type CommandDetailsProps = {
  commandData: CommandType[]
  isLoading?: boolean
}
const CommandInformationPane = ({
  commandData,
  isLoading,
}: CommandDetailsProps): JSX.Element => {
  const data = useMemo(() => commandData[0], [commandData])
  const commandDetailsData = [
    {
      label: 'user',
      data: data?.user,
      check: !!data?.user,
    },
    {
      label: 'plugin',
      data: data?.plugin,
      check: !!data?.plugin,
    },
    {
      label: 'Description',
      data: data?.description,
      check: !!data?.description,
    },
    {
      label: 'version',
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
  ]
  return (
    <div className={styles.informationPaneContainer}>
      <InformationPane
        header={{
          labelAndData: {
            label: 'Command Name',
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
                  label: 'Command ID',
                  data: data?.id,
                  check: !!data?.id,
                },
              ]}
            />
            <InformationPane.Divider />
            <InformationPane.Section data={commandDetailsData} isTwoColumns />
          </>
        )}
      </InformationPane>
    </div>
  )
}

export default CommandInformationPane
