# Architecture Decision Records

Short records of why Heimdall is built the way it is. These were **written after the fact** (2026-11-19) from the code and
git history; where the original rationale is not recorded the ADR says `inferred`. Add new ADRs as `NNNN-title.md`; do not
edit an accepted ADR, supersede it with a new one.

| ADR | Decision |
| --- | --- |
| [0001](0001-plugin-architecture.md) | Backends are Go plugins loaded at runtime |
| [0002](0002-postgres-state-and-deploydb.md) | Postgres is the only state store; schema is applied by `deploydb` from a list file |
| [0003](0003-async-pool-janitor-and-awol.md) | Async jobs run in a worker pool; a janitor releases stale jobs and fails them as AWOL |
