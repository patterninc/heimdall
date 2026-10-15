'use client'

import { Button } from '@patterninc/pattern-ui/button'
import { ConfirmationDialog } from '@patterninc/pattern-ui/confirmation-dialog'
import StopIcon from '@patterninc/pattern-ui/icons/StopIcon'
import { pulse } from '@patterninc/pattern-ui/pulse'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import React, { useState } from 'react'

import { cancelJob } from '@/app/api/jobs/jobs'
import { JobType } from '../Helper'

const CANCELABLE_STATUSES = ['NEW', 'ACCEPTED', 'RUNNING']

type CancelJobButtonProps = {
  job: JobType
}

const CancelJobButton = ({ job }: CancelJobButtonProps): React.JSX.Element => {
  const queryClient = useQueryClient()
  const [confirmOpen, setConfirmOpen] = useState(false)

  const cancelMutation = useMutation({
    mutationFn: (id: string) => cancelJob(id),
    onSuccess: (response) => {
      setConfirmOpen(false)
      if (response.status === 'CANCELING') {
        pulse.info({ title: 'Job is being canceled…' })
        queryClient.invalidateQueries({ queryKey: ['job', job.id] })
      } else {
        pulse.error({
          title: 'Failed to cancel job',
          description: response.error,
        })
      }
    },
    onError: () => {
      setConfirmOpen(false)
      pulse.error({ title: 'Failed to cancel job' })
    },
  })

  // Only async jobs in active states can be canceled.
  const isCancelable =
    CANCELABLE_STATUSES.includes(job.status) && job.is_sync !== true
  const isCanceling = job.status === 'CANCELING'

  return (
    <>
      <Button
        variant='strong'
        size='medium'
        prefix={{ icon: StopIcon }}
        disabled={!isCancelable || isCanceling}
        isLoading={isCanceling}
        tooltipProps={
          !isCancelable && !isCanceling
            ? { content: 'Only running or queued async jobs can be canceled' }
            : undefined
        }
        onClick={() => setConfirmOpen(true)}
        qaTestId='job-details-cancel'
      >
        {isCanceling ? 'Canceling job…' : 'Cancel job'}
      </Button>
      <ConfirmationDialog
        open={confirmOpen}
        onOpenChange={setConfirmOpen}
        title='Cancel job?'
        description={`This stops "${job.name || job.id}" and can't be undone.`}
        confirmText='Cancel job'
        cancelText='Keep running'
        variant='destructive'
        isLoading={cancelMutation.isPending}
        onConfirm={() => cancelMutation.mutate(job.id)}
        qaTestId='job-details-cancel-confirm'
      />
    </>
  )
}

export default CancelJobButton
