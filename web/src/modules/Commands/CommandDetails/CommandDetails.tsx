'use client'

import { BreadcrumbContext } from '@/common/BreadCrumbsProvider/context'
import { useContext, useEffect } from 'react'
import CommandInformationPane from './CommandInformationPane'
import { useQuery } from '@tanstack/react-query'
import { getCommandDetails } from '@/app/api/commands/commands'
import { CommandType } from '../Helper'
import CommandDetailsHeader from './CommandDetailsHeader'

type CommandDetailsProp = {
  id: string
}
export const CommandDetails = ({
  id,
}: CommandDetailsProp): React.JSX.Element => {
  const { updateBreadcrumbs } = useContext(BreadcrumbContext)
  useEffect(() => {
    updateBreadcrumbs({
      name: id,
      link: `/commands/${id}`,
    })
  }, [updateBreadcrumbs, id])

  const { data, isLoading } = useQuery<CommandType[]>({
    queryKey: ['commandDetails', id],
    queryFn: () => getCommandDetails(id),
  })

  return (
    <div className='flex flex-col md:flex-row pat-gap-4 pat-pt-4 min-w-[300px]'>
      <CommandInformationPane
        commandData={data ? data : []}
        isLoading={isLoading}
      />
      <CommandDetailsHeader
        commandData={data ? data : []}
        isLoading={isLoading}
      />
    </div>
  )
}
