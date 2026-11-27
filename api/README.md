# Heimdall API contract

`openapi.yaml` is the OpenAPI 3.1 spec of the `/api/v1` routes of the Go service.
It is the contract between the Go API and the web UI in `web/`.

> **Status: draft.** The spec was written by hand from the Go handlers. The Go API
> owners must confirm it.

## Files

| File | What it is |
| --- | --- |
| `api/openapi.yaml` | The spec. Edit this file by hand. |
| `api/redocly.yaml` | Lint rules for the spec (Redocly `recommended`). |
| `web/src/common/api/schema.d.ts` | TypeScript types made from the spec. Do not edit. |
| `web/src/common/api/types.ts` | Short type names that the web code imports. |
| `internal/pkg/heimdall/openapi_test.go` | Go test that compares the spec with the Go routes and structs. |

## How to change the API

1. Change the Go handler, route or struct.
2. Change `api/openapi.yaml` to match.
3. Make the web types again. Run this from `web/`:

   ```bash
   pnpm gen:api
   ```

4. Commit the spec and `web/src/common/api/schema.d.ts` together.

## Checks

`./build.sh --test` runs two checks. CI runs them in the required `build` check.

- `go test ./internal/pkg/...` runs `openapi_test.go`. It fails when:
  - a route in `apiRoutes()` is not in the spec, or a spec path has no route
  - a spec path does not reach its own route (the router matches path prefixes in order)
  - the properties of a schema are not the JSON fields of its Go struct
  - the `JobStatus` or `ResourceStatus` enum is not the same as the Go status values
- `pnpm check:api` (from `web/`) lints the spec and fails when `schema.d.ts` is not
  the output of `pnpm gen:api`.

The Go test cannot check field types, which fields are required, or status codes.
Review those by hand when you change a handler.
