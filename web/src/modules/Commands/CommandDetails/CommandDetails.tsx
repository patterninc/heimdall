'use client'

import { useQuery } from '@tanstack/react-query'
import React from 'react'

import { getCommandDetails } from '@/app/api/commands/commands'
import { resourceStatusVariant } from '@/common/Services/status'
import DetailPage from '@/components/DetailPage/DetailPage'
import ExternalLinkButton from '@/components/ExternalLinkButton/ExternalLinkButton'
import { CommandType } from '../Helper'
import CommandInformationPane from './CommandInformationPane'

type CommandDetailsProp = {
  /** Command ID — the details API filters on `command_id`. */
  id: string
}

export const CommandDetails = ({
  id,
}: CommandDetailsProp): React.JSX.Element => {
  const { data, isPending, isError } = useQuery<CommandType[]>({
    queryKey: ['commandDetails', id],
    queryFn: () => getCommandDetails(id),
  })
  const command = data?.[0]

  return (
    <DetailPage
      backHref='/commands'
      backLabel='Back to commands'
      resourceLabel='command'
      title={command?.name || id}
      subtitle={
        command?.name && command.name !== command.id ? command.id : undefined
      }
      status={command?.status}
      statusVariant={resourceStatusVariant(command?.status)}
      isError={isError || (!isPending && !command)}
      qaTestId='command-details'
      right={
        <ExternalLinkButton
          href={`/api/v1/commands?id=${command?.id ?? id}`}
          qaTestId='command-details-api-response'
        >
          API response
        </ExternalLinkButton>
      }
    >
      <CommandInformationPane command={command} isLoading={isPending} />
    </DetailPage>
  )
}
