import type { Cluster, ResourceStatus } from '@/common/api/types'
import { FIXTURE_CREATED_AT, FIXTURE_UPDATED_AT } from './jobs'

export const buildCluster = (overrides: Partial<Cluster> = {}): Cluster => ({
  id: 'cluster-emr',
  name: 'emr-prod',
  version: '7.1.0',
  user: 'alice',
  description: 'Production EMR on EKS',
  tags: ['type:emr', 'env:prod'],
  context: {
    emr_release_label: 'emr-7.1.0',
    execution_role_arn: 'arn:aws:iam::123456789012:role/emr-exec',
    role_arn: 'arn:aws:iam::123456789012:role/emr',
    properties: { 'spark.executor.memory': '8g' },
  },
  created_at: FIXTURE_CREATED_AT,
  updated_at: FIXTURE_UPDATED_AT,
  status: 'ACTIVE',
  ...overrides,
})

export const clusters: Cluster[] = [
  buildCluster(),
  buildCluster({
    id: 'cluster-local',
    name: 'local',
    version: '1.0.0',
    user: 'bob',
    description: '',
    tags: [],
    status: 'INACTIVE',
    context: {
      emr_release_label: '',
      execution_role_arn: '',
      role_arn: '',
      properties: {},
    },
  }),
]

export const clusterStatuses: ResourceStatus[] = ['ACTIVE', 'INACTIVE']
