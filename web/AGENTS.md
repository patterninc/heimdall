<!-- BEGIN:nextjs-agent-rules -->

## This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` (resolved from this file's directory; in monorepos the `next` package may not be visible from the repo root) before writing any code. Heed deprecation notices.

This block is written and re-added by `next dev` — verify at `node_modules/next/dist/server/lib/generate-agent-files.js`. Removing it from a diff only re-creates the uncommitted change; committing it with your work keeps the tree clean.

<!-- END:nextjs-agent-rules -->

## Tests

Vitest, React Testing Library and MSW. Run from `web/`:

| Command                                          | What it does                                                                                 |
| ------------------------------------------------ | -------------------------------------------------------------------------------------------- |
| `pnpm test`                                      | Runs every test once with coverage. Fails on a failing test or coverage below the threshold. |
| `pnpm test:watch`                                | Watch mode, no coverage.                                                                     |
| `pnpm vitest run src/test/integration -t "Jobs"` | Runs a subset by path and test name.                                                         |
| `pnpm vitest run -u`                             | Rewrites snapshots after an intended UI change. Review the `.snap` diff before you commit.   |

CI: the **Build and Test** workflow runs `./build.sh --go --ui --test`. With `--ui --test`, `build.sh` runs `pnpm test` before `pnpm run build`, so a failing test or a coverage drop fails the required `build` check.

Layout:

- **Unit tests** sit next to the code: `*.test.ts(x)` (for example `src/common/Services/index.test.ts`, `src/common/hooks/*.test.ts`).
- **Integration tests** for the list and detail pages are in `src/test/integration/`. They render the real page components and serve the API with MSW.
- **Snapshot tests** for the shared `DetailPage` and `ListPage` are next to those components, with `__snapshots__/`. Wrap the container in `withStableIds()` (`src/test/snapshot.ts`) so React ids do not change the snapshot.
- **Shared fixtures** are in `src/test/fixtures/`: typed builders (`buildJob`, `buildCommand`, `buildCluster`), default data, and the MSW `handlers` and `server`. Every test type uses them. To change a response in one test, call `server.use(http.get(api('/jobs'), ...))`. The handlers reset after each test.
- `src/test/render.tsx` (`renderWithProviders`) wraps a component in the same providers as `AppRoot` and returns a `user` from `userEvent`. `src/test/setup.ts` mocks `next/navigation` (assert on `mockRouter` and set `mockLocation` from `src/test/router.ts`), adds the jsdom stubs, and starts MSW. A request with no handler fails the test.
- `*ByTestId` queries read the `qa-test-id` attribute.

Coverage thresholds are in `vitest.config.mts` (`coverage.thresholds`). Raise them when coverage goes up. Do not lower them to make a change pass.
