# ConnorLLM Architecture

ConnorLLM is a **CI-native runtime validator** for LLM endpoints and agents.  
It runs versioned YAML test suites, applies deterministic checks, and returns `exit 0` or `exit 1` for pipelines.

**Promise:** block merges when serving or agent runtime regresses — before production, not after.

See also: [ROADMAP.md](../ROADMAP.md) · [yaml-overview.md](yaml-overview.md) (planned)

---

## 1. System in one picture

```text
  Developer / CI
        │
        ▼
  ┌─────────────┐     YAML suite      ┌──────────────────┐
  │  connor CLI │ ◄────────────────── │ benchmarks/*.yaml │
  └──────┬──────┘                     └──────────────────┘
         │
         ▼
  ┌──────────────────────────────────────────────────────┐
  │              Connor Runtime (Go)                      │
  │                                                       │
  │   Parse → Execute → Evaluate → Gate → Report         │
  └──────────────────────────┬───────────────────────────┘
                             │ HTTP
                             ▼
              ┌──────────────────────────────┐
              │  CONNOR_BASE_URL (external)   │
              │  OpenAI-compatible API        │
              │  (vLLM, OpenRouter, gateway)  │
              └──────────────────────────────┘