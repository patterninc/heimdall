# Heimdall web UI

The Heimdall web UI is a Next.js 16 app. It shows jobs, commands and clusters from a Heimdall API.

In production, the Heimdall server proxies the UI at `http://localhost:9090/`. For UI work, run the Next.js dev server on port 4000 and point it at a Heimdall API. You can use a local API or a remote one.

All commands on this page run from the `web/` directory, unless the step says something different.

## Prerequisites

- **Node.js 26.** `web/.nvmrc` pins the version. With [nvm](https://github.com/nvm-sh/nvm), run `nvm use` in `web/`. Other version managers can also read `.nvmrc` (for example `fnm use` or `mise use node@26`).
- **pnpm 12.** The `packageManager` field in `web/package.json` pins the version (`pnpm@12.15.0`). Node 25 and later do not include corepack. Install corepack, then enable it:

  ```bash
  npm install --global corepack
  corepack enable
  ```

  Corepack then uses the pinned pnpm version in this directory. If you do not use corepack, install the same version directly: `npm install --global pnpm@12.15.0`.

- **Docker** with Docker Compose v2 (`docker compose`). You need Docker only for `pnpm dev:stack`.

## Install

```bash
cd web
pnpm install
```

## Run the dev server

```bash
pnpm dev
```

This runs `next dev` on [http://localhost:4000](http://localhost:4000).

### Point the UI at a Heimdall API

The browser calls relative paths, `/api/v1/...` (see `src/common/Services/index.ts`). The Next.js dev server forwards these calls to the Heimdall API. `next.config.ts` sets this rewrite and reads the API location from the `API_URL` environment variable:

```ts
const apiUrl = process.env.API_URL || 'http://localhost:9090'
```

`API_URL` is the base URL of the API, without `/api/v1`. When `API_URL` is not set, the UI uses `http://localhost:9090`. That is the port of the API from `docker compose`.

You can set `API_URL` in two ways:

- **In `web/.env.local`.** Next.js loads this file. Git ignores it (see `web/.gitignore`). Example:

  ```bash
  # web/.env.local
  # Local API (docker compose, or `go run` on the default port):
  API_URL=http://localhost:9090
  # Or a remote API:
  # API_URL=https://heimdall.example.com
  ```

- **In the shell, for one run.** A value in the shell has priority over `.env.local`:

  ```bash
  API_URL=https://heimdall.example.com pnpm dev
  ```

Next.js reads `next.config.ts` when the dev server starts. If you change `API_URL`, stop `pnpm dev` and start it again.

The rewrite runs on the dev server, so the browser does not make cross-origin calls. The rewrite sends the request headers from the browser to the API. With `configs/local.yaml`, the API reads the user name from the `X-Heimdall-User` header (the `auth` section).

## One-command local stack

```bash
pnpm dev:stack
```

This command runs `scripts/dev-stack.sh`. The script:

1. Starts the `postgres` and `heimdall` services from the repository's `docker-compose.yaml` with `docker compose up -d`. Postgres is on host port 5433. The Heimdall API is on host port 9090. The API container creates the database schema when it starts.
2. Waits until `http://localhost:9090/api/v1/commands` answers. The default wait is 180 seconds.
3. Runs `next dev` on [http://localhost:4000](http://localhost:4000) with `API_URL=http://localhost:9090`.

Press `Ctrl+C` to stop the dev server. The containers continue to run, so the next `pnpm dev:stack` starts quickly.

The first run builds the Heimdall Docker image from the repository's `Dockerfile`. This can take several minutes. Docker Compose does not rebuild the image automatically. After you change Go code or `configs/local.yaml`, rebuild the image:

```bash
pnpm dev:stack --build
```

Options:

| Setting             | Default                 | Description                                                                                                                                              |
| ------------------- | ----------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--build`           | not set                 | Rebuild the Heimdall image before the start.                                                                                                             |
| `API_URL`           | `http://localhost:9090` | The address that the script waits for and that the dev server uses. Set it in the shell. Change it only when Docker publishes ports on a different host. |
| `DEV_STACK_TIMEOUT` | `180`                   | The maximum number of seconds to wait for the API.                                                                                                       |

`docker-compose.yaml` gives the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN` and `AWS_REGION` variables from your shell to the API container. Only the plugins that use AWS need them. The `ping` command in `configs/local.yaml` does not need them.

To make test data, submit a `ping` job. See "Submit a Test Job" in the [repository README](../README.md).

### Stop the stack

Stop and remove the containers. The database data stays in the `postgres-db` volume:

```bash
pnpm dev:stack:down
```

Also delete the database volume, for a clean start:

```bash
pnpm dev:stack:down --volumes
```

These commands run `docker compose down` (with `--volumes` for the second command) on the repository's `docker-compose.yaml`. You can also run them from the repository root:

```bash
docker compose down            # keep the data
docker compose down --volumes  # delete the data
```

To see the API logs, run `docker compose logs -f heimdall` from the repository root.

## Lint, type-check and format

| Command               | What it does                                                                                          |
| --------------------- | ----------------------------------------------------------------------------------------------------- |
| `pnpm lint`           | Runs ESLint (`eslint.config.mjs`, Next.js core-web-vitals and TypeScript rules).                      |
| `pnpm lint-fix`       | Runs ESLint and fixes the problems that it can fix.                                                   |
| `pnpm type-check`     | Generates the Next.js route types (`next typegen`), then runs strict TypeScript (`tsc --noEmit`).     |
| `pnpm format`         | Checks the Prettier format of all files (`.prettierrc`, `.prettierignore`). It does not change files. |
| `pnpm prettier-check` | Checks the Prettier format of the files in `src/` only.                                               |
| `pnpm prettier`       | Formats all files with Prettier. This changes files.                                                  |
| `pnpm build`          | Makes a production build (`next build`).                                                              |

CI runs `./build.sh --go --ui --test` from the repository root. That command installs the dependencies and runs `pnpm run build`.
