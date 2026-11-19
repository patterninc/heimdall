# Skill: Heimdall release and rollback

Follow [`docs/runbooks/deploy-and-rollback.md`](../../docs/runbooks/deploy-and-rollback.md). Do not tag, push tags, or deploy without an explicit instruction from a person.

Before proposing a release: `CHANGELOG.md` has an entry for the version, schema changes are listed in `assets/databases/heimdall/build/heimdall.lst` and called out as a Migration note, and `main` is green on `Build and Test`.
