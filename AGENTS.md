# AGENTS.md

Agent-agnostic guidance for working in this repo. Human docs live in `README.md`.

This file is a **router** — deeper agent context lives in the **`.agents/`** folder, loaded on
demand. Read the matching entry before starting a task it covers.

## Task guides — `.agents/skills/`

| When the task is…                                                                                                                                                   | Read                                                              |
| -------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------- |
| **UI-only** change to the web app (`web/`) — buttons/links, labels, layout, responsiveness, tooltips, table columns, rendering data the API already returns | [`.agents/skills/heimdall-ui.md`](.agents/skills/heimdall-ui.md) |

## Personas — `.agents/personas/`

Adopt a persona when it fits the task:

| Persona           | Use for                 | Profile                                                                           |
| ----------------- | ----------------------- | --------------------------------------------------------------------------------- |
| Frontend reviewer | UI-only `web/` changes | [`.agents/personas/frontend-reviewer.md`](.agents/personas/frontend-reviewer.md) |

## Setup (one-time, per machine)

Some guides need a tool your agent provides — e.g. UI browser verification uses the
**chrome-devtools MCP**. Two ways it gets set up:

- **On demand:** the guide's **Step 0** fires when a task first needs the tool.
- **Proactively:** ask your agent to **"set up the &lt;name&gt; skill"** (e.g. *"set up the
  heimdall-ui skill"*).
