# Heimdall

Heimdall is a data orchestration and job execution platform that sits on the boundary between the data infrastructure and the client (systems and users that access and process data). Heimdall provides a restful API to submit job requests and processes those jobs in both, sync and async (incoming feature) fashion. The project is inspired by [Netflix Genie](https://github.com/Netflix/genie), however changes the architecture to support plugins for types of the commands along with support of job queues and synchronous commands. Tools like Genie and Heimdall exist to abstract data infrastructure from the consumers and while they minimize the friction of interacting with Big Data infrastructures, eliminating the need of specific data application implementations on the client, they also eliminate the need of bringing shared credentials to the clients, specifically users' local environments.

At this point, a single plugin ("ping") is fully supported end to end using synchronous approach. To test it, follow these steps:

1. Pull heimdall to your local development environment: `git clone git@github.com:patterninc/heimdall.git`

2. Start heimdall: `podman compose up --build -d`: this will start heimdall server on local port `9090` and also will deploy the database objects to support jobs submission.

3. Submit a job: `curl -X POST -H "X-Pattern-User: test_user" -d '{"name":"ping-test", "version":"0.0.1", "context":{}, "command_criteria":["type:ping"], "cluster_criteria":["type:localhost"]}' http://127.0.0.1:9090/api/v1/job`. In this example we do a POST call to `/api/v1/job` endpoint to submit the job. As part of the request, we provide the following data:
   - `name`: this is the name of the job, which can be any string, however in production we'll usually request the service clients to have a specific names for their jobs so we could easily pull their jobs when we need to debug them.
   - `version`: this is a version of the implementation that client has for a specific job type. For the adhoc workloads it does not play a big role, however in case of a client service, we'll be able to distinguish different version of the client implementations.
   - `user`: this property will be set by the heimdall itself once we integrate with an identity provider, so we could have true authentication in the system for clients submitting the job. At this point, for testing, we can use any string.
   - `command_criteria`: is a set of tags that we want heimdall to use to select a command for our job to execute. We intentionally do not specify the exact command (even though it's possible when specifying `_id:...` tag) to allow heimdall to "orchestrate" the traffic: for example, we may have two commands with the same tags where one command implements the current behavior we have in production, and another command that implements the future behavior. Without user telling us, we can send 1% of the traffic to the future version of the command, while 99% of the traffic is still routed to the current production behavior. More details about command and cluster "pair" selection is outlined in the following section.
   - `cluster_criteria`: similar to the command, this is a set of tags we want a cluster to have to be selected for our job. This behavior is useful, so we could deploy multiple clusters for the same type of the command, so we could distribute the jobs among multiple clusters for availability (if one cluster is unhealthy, we can remove it from accepting new jobs, but also another cluster can pick up the failed jobs from the first cluster after a retry).
   - `context`: is a flexible data context that "understood" only by specific plugin that will execute our job. For the ping command, an empty context is sufficient as it's a simple test command. For other types of commands, we'll have different contexts (for example, for `snowflake` jobs, we'll be passing `query` as part of the context).

4. Since it's a sync command, the response from the server will contain the `status` of the command, which command and cluster it was executed on along with `id` of the command. In case of the sync execution of the query, response will contain the data as well. Here's how to get this data:

  - status of the job: `GET /api/v1/job/<job_id>`
  - `stdout`: `GET /api/v1/job/<job_id>/stdout`
  - `stderr`: `GET /api/v1/job/<job_id>/stderr`

## Heimdall Data Model

Heimdall operates with only three types of the objects:

- Commands: each command will execute a plugin call for a job. It represents a piece of the work that we want heimdall to do using specific `context` passed as part of a job.
- Cluster: cluster represents an abstraction of a compute component, that is represented by its context: it could be a localhost, EMR cluster, EKS cluster, connection string to a database, etc...
- Job: is a workload that we want heimdall to perform. For that we'll specify command and cluster criteria and context for the command we want to execute. Command and cluster criteria are used to identify a specific command and specific cluster on which that command will be executed.

## Heimdall configuration

At first, we'll be configuring Commands and Clusters using heimdall configuration file. In the future, those can be configured in a more sophisticated way using the APIs (clusters being self registering and deregistering based on their health metrics). The example of such configuration can be found in this [config file](https://github.com/patterninc/heimdall/blob/main/configs/local.yaml).

## Command and Cluster selection for a Job

Each job provides two criteria to identify a command heimdall will execute and on which cluster that command will be executed. The logic to identify this pair starts by looking for the commands that satisfy the criteria: a command must by `active` and must have all the `tags` specified in `command_criteria`. Once a command is identified, command's `cluster_criteria` is used to identify all the clusters that are configured / "compatible" with the command. All the identified clusters must be `active` as well. the identified pairs of the commands and clusters are filtered further with `cluster_criteria` of the job (in this case we can "route" a specific job to a cluster that processes critical production workloads, or a cluster to run adhoc requests from a user). At this point of the selection process, we can have one of the three scenarios:

- Given the job criteria, no pair was selected: in this case the job will fail with the error specifying that no pain was found.
- A single pair was identified: in this scenario, the job is associated with the command and cluster from the pair and sent for the execution.
- Multiple pairs of commands and clusters were identified: at this point, heimdall will randomly select one pair from the set and associate the job with it, following by sending the job for the execution on that command-cluster pair.

## Credits

Heimdall was created at Pattern, Inc by Stan Babourine, with contributions from Will Graham, Gaurav Warale and Josh Diaz.
