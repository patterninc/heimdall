# Skill: Heimdall UI changes

Agent-agnostic guide for making and self-verifying **UI-only** changes to the Heimdall web app
(Next.js + `@patterninc/react-ui`). Any coding agent or engineer can follow it.

The loop is: **edit → gate on lint/type-check/build → verify the rendered result in a real browser → done.**
Do not call a UI task done on a passing type-check alone — the bugs that matter here
(missing element, broken layout, stray placeholder, bad responsiveness, false-success handlers) only
show up when rendered.

> **This skill applies to a UI change the user requests — it is not a standing instruction to go
> find work.** If no specific task was given, **ask what to change and stop** — do **not** scan the
> codebase for issues to fix on your own. The standards below (accessibility, failure handling,
> etc.) are things to **apply to the requested change**, not a mandate to hunt for and fix
> violations elsewhere.

## Setup (one-time, per machine)

To arm this skill for smooth use, **ask your agent: "set up the heimdall-ui skill."** The agent should:
- **Claude Code:** scaffold `.claude/skills/heimdall-ui/SKILL.md` — a thin wrapper that points back to
  this file (so it auto-triggers on UI tasks and stays in sync with this single source of truth).
- **Browser MCP:** register the chrome-devtools MCP for you — **Step 0** has the order: try the
  client's official install command first, fall back to writing the config directly, then surface
  as a last resort. Then **restart** your agent so the new skill + MCP connect, and **approve** it.

You don't have to pre-run setup — Step 0 does the same registration on demand when a task first
needs the browser.

## Scope — read this first

- **In scope:** button/link additions, labels/text, layout & spacing, responsiveness,
  truncation/overflow, tooltips, table columns, and rendering data the API **already returns**.
- **Out of scope — STOP and flag for a human:** any backend/Go change, a new API field, DB
  schema/migrations, plugin logic, or anything where the data the UI needs isn't already
  available from the API. Do not improvise server-side changes.

## Repo orientation

- Web app lives in `web/` (Next.js, React, TypeScript). UI built from `@patterninc/react-ui`
  components — e.g. `Button`, `InformationPane`, `TrimText`, `MdashCheck`, `Icon`, `Tag`.
- Feature modules: `web/src/modules/<Feature>/` (e.g. `Jobs`, `Clusters`, `Commands`).
  Detail panes under `web/src/modules/<Feature>/<Feature>Details/`.
- API access: the browser calls **relative** `/api/v1/...` paths; `API_URL` is `/api/v1`
  (`web/src/common/Services`). Next proxies `/api/v1/*` to heimdall.
- Conventions to respect:
  - `MdashCheck` renders an em-dash (`—`) when a value is missing — to show *nothing*, omit the
    item, don't rely on `check: false`.
  - `TrimText` truncates with a built-in tooltip (`limit`, `customClass`).
  - `InformationPane.Section` with `isTwoColumns` is a 2-column grid — **mixed cell heights or odd
    item counts wrap/space unevenly**; keep rows uniform (e.g. all one-line, or all label+value).
  - Icon-only controls need an `aria-label`; handlers that can fail (e.g. `navigator.clipboard`)
    must handle rejection and only show a success toast on success.

## Step 0 — browser MCP preflight (do BEFORE any browser verification)

Browser verification (required below) needs a browser-automation tool. This project uses the
**Chrome DevTools MCP** (`chrome-devtools-mcp`).

> **STOP if it isn't connected.** Check whether your browser MCP tools are available. If not, **do
> not proceed to verification** — get it registered first. **Do not hardcode or guess the config**
> (it drifts and is per-client); take the exact `chrome-devtools` server config from the **official
> docs** (fetch the page if you can):
> **https://github.com/ChromeDevTools/chrome-devtools-mcp** (per-client — Claude Code, Cursor, VS
> Code/Copilot, Gemini CLI, Codex, Windsurf, JetBrains, …).
>
> **Register it — try these in order:**
> 1. **Run the client's official install command** (from the docs above) — for Claude Code,
>    `claude mcp add …`. This is the sanctioned path and writes the config correctly.
> 2. **If the command is missing or fails** (e.g. the CLI isn't on PATH), **fall back to writing the
>    client's MCP config directly** — for Claude, MERGE the `chrome-devtools` entry into
>    `~/.claude.json` → `mcpServers` (read-modify-write; **never overwrite the file** — it holds the
>    user's entire config).
> 3. **If neither is possible** (unknown client / user declines), **surface** the command for the
>    user to run.
>
> All three are **permission-gated by the user's own permission mode** — whether they see a prompt is
> their setting, not this guide's to force. Prereqs: **Node.js** (`npx`) and **Google Chrome**. MCP
> servers connect only at **startup**, so the user must **restart** and **approve** the server
> afterward. No browser tool at all → do the checks manually.

## Make the change

1. Find the component under `web/src/modules/...`. Match the patterns and the react-ui component
   already used nearby.
2. Edit **only** frontend files. Reuse existing components; don't hand-roll styling a react-ui
   component already provides.

## Local dev loop (from repo root)

- Full stack: `docker compose up --build -d`. First build is slow (~8 min). UI + API at
  `http://127.0.0.1:9090`; Postgres at `localhost:5433` (user/db/pass all `heimdall`).
- **Faster UI iteration:** `cd web && pnpm dev` — Next dev server on **:4000**, proxying `/api/v1/*`
  to the running `:9090` API. Hot-reload, **no image rebuild**. Rebuild the Docker image only for
  Go / plugin / DB-schema changes; config-only changes just need a container restart.
- API calls need `-H "X-Heimdall-User: test_user"` (note the space after the colon).
- Submit a test job (ping plugin, completes locally):
  ```sh
  curl -X POST -H "X-Heimdall-User: test_user" -H "Content-Type: application/json" \
    -d '{"name":"ui-test","version":"0.0.1","context":{},"command_criteria":["type:ping"],"cluster_criteria":["type:localhost"]}' \
    http://127.0.0.1:9090/api/v1/job     # → returns the job id
  ```
- Seed a specific rendering case:
  ```sh
  docker exec heimdall_postgres psql -U heimdall -d heimdall -c "update jobs set ... where job_id='<id>';"
  ```

## Gate: lint / type-check / build (must pass) — run in `web/`

- `npx prettier --check <changed files>`
- `npx eslint <changed files>`
- `npx tsc --noEmit 2>&1 | grep -E "<your files>"`
  (Type errors inside `node_modules/...` are pre-existing noise — filter to the files you touched.)

## Verify in a real browser (required)

Do **not** skip this — a passing type-check does not prove the UI renders. (Step 0 must be satisfied.)

1. Open **every view your change affects** — whichever feature area it lives in (Jobs, Clusters,
   Commands, …) and whichever view types apply (list, detail, …). If the change is in a **shared
   component** (react-ui wrappers, panes, table cells), verify **each** view it renders on, not only
   the one you edited. Serve at `http://127.0.0.1:9090/...` (or `:4000` with the dev server).
2. Screenshot each; confirm the change renders as intended — element present/absent, text correct,
   layout not broken, alignment/spacing sane.
3. Check the browser **console** for new errors/warnings (a change can pass `tsc` and still throw).
4. If the change is responsive, check key widths (mobile ~375px, tablet ~768px, desktop).
5. Verify the **negative/empty case** (data absent → element correctly hidden, no stray `—`).

## Definition of done

- Frontend-only edits; `prettier` + `eslint` + `tsc` clean for the touched files; build succeeds.
- Browser-verified: target renders correctly (positive **and** negative case), no new console
  errors, responsive breakpoints checked when relevant — backed by screenshots.
- Accessibility + failure paths handled (aria-labels on icon buttons; fallible handlers don't
  claim false success).
- If the task turned out to need a backend/API/schema/plugin change, you stopped and flagged it
  rather than improvising it.
