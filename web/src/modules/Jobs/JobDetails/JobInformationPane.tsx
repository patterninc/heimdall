import FileErrorIcon from '@patterninc/pattern-ui/icons/FileErrorIcon'
import { CodeBlock } from '@patterninc/pattern-ui/artifact/CodeBlock'
import { Section } from '@patterninc/pattern-ui/section'
import Link from 'next/link'
import React from 'react'

import { formatDateWithTimeZone, myTimezone } from '@/common/Services'
import DetailFields from '@/components/DetailPage/DetailFields'
import ExternalLinkList, {
  type ExternalLink,
} from '@/components/DetailPage/ExternalLinkList'
import KeyValueList from '@/components/DetailPage/KeyValueList'
import TagList from '@/components/DetailPage/TagList'
import { JobType } from '../Helper'

type JobInformationPaneProps = {
  jobData?: JobType
  isLoading: boolean
}

const inlineLinkClass =
  'text-text-strong underline decoration-border-sub underline-offset-2 hover:decoration-current'

const JobInformationPane = ({
  jobData,
  isLoading,
}: JobInformationPaneProps): React.JSX.Element => {
  const attributes = Object.entries(jobData?.job_attributes ?? {})
  const outputLinks: ExternalLink[] = jobData
    ? [
        { label: 'View stdout', href: `/api/v1/job/${jobData.id}/stdout` },
        { label: 'View stderr', href: `/api/v1/job/${jobData.id}/stderr` },
        ...attributes
          .filter(([, attr]) => attr?.kind === 'link' && attr.value)
          .map(([label, attr]) => ({ label, href: attr.value as string })),
      ]
    : []
  const textAttributes: [string, string][] = attributes
    .filter(([, attr]) => attr?.kind !== 'link')
    .map(([label, attr]) => [label, attr?.value ?? '—'])
  const properties = Object.entries(jobData?.context?.properties ?? {})

  return (
    <>
      {jobData?.status === 'FAILED' && jobData.error ? (
        <Section
          label='Error'
          prefix={{ icon: FileErrorIcon, color: 'negative-base' }}
          qaTestId='job-details-error'
        >
          <p className='pui-text-mono-small text-text-strong p-4 break-words whitespace-pre-wrap'>
            {jobData.error}
          </p>
        </Section>
      ) : null}

      <Section label='Overview' qaTestId='job-details-overview'>
        <DetailFields
          isLoading={isLoading}
          fields={[
            { label: 'Job ID', value: jobData?.id, mono: true },
            { label: 'User', value: jobData?.user },
            { label: 'Version', value: jobData?.version },
            {
              label: 'Command',
              value: jobData?.command_name ? (
                <Link
                  href={`/commands/${jobData.command_name}`}
                  className={inlineLinkClass}
                >
                  {jobData.command_name}
                </Link>
              ) : undefined,
            },
            {
              label: 'Cluster',
              value: jobData?.cluster_name ? (
                <Link
                  href={`/clusters/${jobData.cluster_name}`}
                  className={inlineLinkClass}
                >
                  {jobData.cluster_name}
                </Link>
              ) : undefined,
            },
            {
              label: 'Execution',
              value: jobData ? (jobData.is_sync ? 'Sync' : 'Async') : undefined,
            },
            {
              label: 'Created at',
              value: jobData?.created_at
                ? formatDateWithTimeZone(jobData.created_at, myTimezone)
                : undefined,
            },
            {
              label: 'Updated at',
              value: jobData?.updated_at
                ? formatDateWithTimeZone(jobData.updated_at, myTimezone)
                : undefined,
            },
            ...(jobData?.canceled_by
              ? [{ label: 'Canceled by', value: jobData.canceled_by }]
              : []),
          ]}
        />
      </Section>

      {outputLinks.length > 0 ? (
        <Section label='Outputs' qaTestId='job-details-outputs'>
          <ExternalLinkList links={outputLinks} qaTestId='job-details-link' />
        </Section>
      ) : null}

      {textAttributes.length > 0 ? (
        <Section
          label='Attributes'
          count={textAttributes.length}
          qaTestId='job-details-attributes'
        >
          <KeyValueList entries={textAttributes} />
        </Section>
      ) : null}

      {jobData?.context?.query ? (
        <Section label='SQL query' qaTestId='job-details-query'>
          {/* The theme sets inline `white-space: pre` on pre/code, hence the important modifiers to wrap long lines. */}
          <div className='px-3 pt-1 pb-3'>
            <CodeBlock
              code={jobData.context.query}
              language='sql'
              className='[&_.max-h-\[600px\]]:max-h-72 [&_code]:break-words! [&_code]:whitespace-pre-wrap! [&_pre]:whitespace-pre-wrap!'
            />
          </div>
        </Section>
      ) : null}

      {properties.length > 0 ? (
        <Section
          label='Context properties'
          count={properties.length}
          qaTestId='job-details-properties'
        >
          <KeyValueList entries={properties} />
        </Section>
      ) : null}

      {jobData?.tags?.length ? (
        <Section
          label='Tags'
          count={jobData.tags.length}
          qaTestId='job-details-tags'
        >
          <TagList values={jobData.tags} />
        </Section>
      ) : null}

      {jobData?.command_criteria?.length ? (
        <Section
          label='Command criteria'
          count={jobData.command_criteria.length}
          qaTestId='job-details-command-criteria'
        >
          <TagList values={jobData.command_criteria} />
        </Section>
      ) : null}

      {jobData?.cluster_criteria?.length ? (
        <Section
          label='Cluster criteria'
          count={jobData.cluster_criteria.length}
          qaTestId='job-details-cluster-criteria'
        >
          <TagList values={jobData.cluster_criteria} />
        </Section>
      ) : null}
    </>
  )
}

export default JobInformationPane
