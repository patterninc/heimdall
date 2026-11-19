# Deploy and roll back

**Images.** `.github/workflows/docker-image.yml` builds `patternoss/heimdall` for amd64 and arm64 and publishes one manifest:
a push to `main` publishes tag `dev`; a git tag `vX.Y.Z` publishes `X.Y.Z`; a tag named `latest` publishes `latest`.

## Release
1. Make sure `main` is green (`Build and Test`).
2. Add the release to `CHANGELOG.md` (include migration notes if the schema or config changed).
3. Tag: `git tag vX.Y.Z && git push origin vX.Y.Z`. Wait for the `Build and Push Docker image` run to finish.
4. If the release changes the schema, apply it first: see [database-schema-deploy.md](database-schema-deploy.md).
5. Roll the deployment to the new tag. `TODO(owner)`: where Heimdall is deployed and the command or pipeline used.

## Roll back
1. Redeploy the previous tag (`docker buildx imagetools inspect patternoss/heimdall:<previous>` confirms it exists).
2. Schema is additive and re-runnable (ADR 0002), so the previous server version normally runs on the newer schema. If a release dropped or renamed a column, restore it before rolling back. Check the CHANGELOG migration note for that release.
3. Resumable jobs (sparkeks, ecs) are released to the remaining agents when a pod stops (ADR 0003). Watch `GET /api/v1/jobs` for jobs that stay `RUNNING`.
