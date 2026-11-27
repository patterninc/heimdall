import type { Job, JobStatus, JobsPage } from '@/common/api/types'

/** 2026-01-15 10:30:00 UTC */
export const FIXTURE_CREATED_AT = 1768473000
/** 2026-01-15 10:45:00 UTC */
export const FIXTURE_UPDATED_AT = 1768473900

export const buildJob = (overrides: Partial<Job> = {}): Job => ({
  id: 'job-0001',
  name: 'daily-revenue-rollup',
  version: '1.2.0',
  user: 'alice',
  tags: ['team:data', 'env:prod'],
  created_at: FIXTURE_CREATED_AT,
  updated_at: FIXTURE_UPDATED_AT,
  status: 'SUCCEEDED',
  is_sync: false,
  command_criteria: ['type:spark'],
  cluster_criteria: ['type:emr'],
  command_id: 'cmd-spark',
  command_name: 'spark-sql',
  cluster_id: 'cluster-emr',
  cluster_name: 'emr-prod',
  job_attributes: {
    spark_ui: { kind: 'link', value: 'https://spark.example.com/app/1' },
    rows_written: { kind: 'text', value: '1024' },
  },
  context: {
    properties: {
      'spark.driver.cores': '2',
      'spark.driver.memory': '4g',
      'spark.executor.cores': '4',
      'spark.executor.instances': '8',
      'spark.executor.memory': '8g',
    },
    query: 'SELECT 1',
  },
  ...overrides,
})

export const jobs: Job[] = [
  buildJob(),
  buildJob({
    id: 'job-0002',
    name: 'hourly-ingest',
    user: 'bob',
    status: 'RUNNING',
    tags: [],
    job_attributes: undefined,
  }),
  buildJob({
    id: 'job-0003',
    name: 'backfill-orders',
    user: 'carol',
    status: 'FAILED',
    error: 'java.lang.OutOfMemoryError: Java heap space',
  }),
]

export const jobsPage = (
  data: Job[] = jobs,
  next_cursor?: string,
): JobsPage => ({
  data,
  has_more: Boolean(next_cursor),
  next_cursor,
})

export const jobStatuses: JobStatus[] = [
  'NEW',
  'ACCEPTED',
  'RUNNING',
  'SUCCEEDED',
  'FAILED',
  'KILLED',
]
