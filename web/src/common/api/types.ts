import type { components, operations } from './schema'

type Schemas = components['schemas']

type QueryOf<Op extends keyof operations> = NonNullable<
  operations[Op]['parameters']['query']
>

export type ApiError = Schemas['Error']
export type ApiContext = Schemas['Context']

export type JobStatus = Schemas['JobStatus']
export type ResourceStatus = Schemas['ResourceStatus']

export type Job = Schemas['Job']
export type JobAttribute = Schemas['JobAttribute']
export type JobsPage = Schemas['JobsPage']
export type JobCancelResponse = Schemas['JobCancelResponse']
export type JobsQuery = QueryOf<'listJobs'>

export type Command = Schemas['Command']
export type CommandsPage = Schemas['CommandsPage']
export type CommandsQuery = QueryOf<'listCommands'>

export type Cluster = Schemas['Cluster']
export type ClustersPage = Schemas['ClustersPage']
export type ClustersQuery = QueryOf<'listClusters'>
