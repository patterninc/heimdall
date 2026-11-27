import type { CommandType } from '@/modules/Commands/Helper'
import { FIXTURE_CREATED_AT, FIXTURE_UPDATED_AT } from './jobs'

export const buildCommand = (
  overrides: Partial<CommandType> = {},
): CommandType => ({
  id: 'cmd-spark',
  name: 'spark-sql',
  version: '3.5.0',
  user: 'alice',
  description: 'Runs Spark SQL on EMR',
  tags: ['type:spark'],
  created_at: FIXTURE_CREATED_AT,
  updated_at: FIXTURE_UPDATED_AT,
  status: 'ACTIVE',
  plugin: 'spark',
  cluster_tags: ['type:emr'],
  is_sync: false,
  context: {
    logs_uri: 's3://logs/spark',
    queries_uri: 's3://queries/spark',
    results_uri: 's3://results/spark',
    properties: { 'spark.executor.memory': '8g' },
  },
  ...overrides,
})

export const commands: CommandType[] = [
  buildCommand(),
  buildCommand({
    id: 'cmd-shell',
    name: 'shell',
    version: '1.0.0',
    user: 'bob',
    description: 'Runs a shell script',
    plugin: 'shell',
    status: 'INACTIVE',
    is_sync: true,
    context: undefined,
    cluster_tags: [],
  }),
  buildCommand({
    id: 'cmd-ping',
    name: 'ping',
    version: '0.1.0',
    user: 'carol',
    description: '',
    plugin: 'ping',
    status: 'DEPRECATED',
  }),
]

export const commandStatuses = ['ACTIVE', 'INACTIVE', 'DEPRECATED']
