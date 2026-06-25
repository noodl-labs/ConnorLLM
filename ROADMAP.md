# ConnorLLM Roadmap

> Continuous integration for production AI — quality gates before merge.

**Current release:** [v0.1.0-beta.1](CHANGELOG.md#010-beta1)  
**Status:** Beta — suitable for LLM serving smoke tests in CI.

---

## Vision

Block merges when LLM/agent runtime regresses — latency, availability, structured output — using versioned YAML suites and `exit 0` / `exit 1`.

**Not in scope:** production observability (Langfuse), model intelligence benchmarks (MMLU).

---

## Release timeline

| Release | Target | Theme | CI value |
|---------|--------|-------|----------|
| **v0.1.0-beta.1** | Now | Serving smoke | "Does my endpoint respond?" |
| **v0.1.0** | +2–4 weeks | Regression compare | "Did we regress vs baseline?" |
| **v0.2.0** | +4–8 weeks | Agent gates | "Did the agent call the right tool?" |
| **v1.0.0** | +3–6 months | Full agent CI | Workflows, replay, semantic eval |

Dates are indicative — ship when **exit criteria** below are met.

---

## Shipped  — v0.1.0-beta.1

### Execution Engine
- [x] `connor run --model --prompt` (single case)
- [x] `connor run suite.yaml` (multi-case)
- [x] OpenAI-compatible provider (`CONNOR_BASE_URL`, `CONNOR_API_KEY`)
- [x] Per-attempt timeout
- [x] Retry on 429 / 5xx / transient network
- [x] Sequential suite execution

### Evaluation Engine
- [x] HTTP 2xx success check
- [x] JSON syntax gate (`expect_json`)

### Quality Gates
- [x] `exit 0` if all cases pass
- [x] `exit 1` if any case fails
- [x] Fail reasons: `call_failed`, `invalid_json`

### Developer Experience
- [x] YAML parser + validation
- [x] `benchmarks/examples/serving-smoke.yaml`
- [x] CLI output (case_id, latency, body preview, summary)

### Integration level
- [x] **L1** — direct `/chat/completions` (serving)
- [x] **L2** — same API, user staging gateway (config only)

---

## In progress  — toward v0.1.0

**Theme:** Regression testing & budget gates (latency / pass rate)

### Benchmark Engine
- [ ] Export `run.json` after suite run
- [ ] `connor compare baseline.json candidate.json`
- [ ] Suite summary: p50 / p95 latency

### Quality Gates
- [ ] `max_p95_regression` threshold (e.g. 15%)
- [ ] `min_pass_rate` threshold (e.g. 0.95)
- [ ] `connor compare` exits 1 on gate failure

### Developer Experience
- [ ] README complete (quick start + CI snippet)
- [ ] `docs/getting-started.md`, `docs/ci-github-actions.md`
- [ ] `.env.example`
- [ ] LICENSE

### Exit criteria for v0.1.0
- [ ] `connor run suite.yaml --out run.json` works
- [ ] `connor compare` blocks on latency regression in demo
- [ ] Docs + tag `v0.1.0` published

---

## Planned 📋 — v0.2.0

**Theme:** Agent-ready gates (L3 integration begins)

### Evaluation Engine
- [ ] `expect_contains` (text assertions)
- [ ] JSON Schema validation (L2)
- [ ] Parse `tool_calls` from API response
- [ ] `expect_tool` / `expect_tool_calls` (name, order)

### Execution Engine
- [ ] Agent HTTP provider (custom URL, not only `/chat/completions`)
- [ ] `context` in YAML (session_id, user_id)
- [ ] Streaming + TTFT measurement

### Quality Gates
- [ ] Token usage from `usage` field
- [ ] `max_cost_regression` vs baseline

### Developer Experience
- [ ] `benchmarks/examples/agent-support.yaml` (reference)
- [ ] `docs/yaml-overview.md` (L1–L4)

### Exit criteria for v0.2.0
- [ ] One demo: wrong tool → `exit 1`
- [ ] JSON schema gate on agent output

---

## Planned 📋 — v1.0.0

**Theme:** Full agent CI (L4)

### Execution Engine
- [ ] Multi-step scenarios (chained cases / workflow)
- [ ] `connor replay` from stored run

### Evaluation Engine (`services/evaluation/` Python)
- [ ] Semantic similarity
- [ ] Groundedness / hallucination checks
- [ ] Tool correctness (args matching)

### Benchmark Engine
- [ ] Prompt diff (version A vs B)
- [ ] Model leaderboard report

### Quality Gates
- [ ] Reliability score (objective sub-scores + explicit N/A)
- [ ] Composite gate on score

### Observability Engine
- [ ] Run history (local / SaaS-ready)
- [ ] Replay production run in CI

### Exit criteria for v1.0.0
- [ ] End-to-end PR demo: regression + tool + semantic gate
- [ ] Stable YAML schema v1

---

## Non-goals

- Replacing Langfuse / LangSmith in production
- Kubernetes operators / distributed runners (for now)
- MMLU and academic model leaderboards
- Bubble Tea TUI (optional `--tui` later)
- Explicit `retries: 0` in YAML (beta limitation; fix in v0.2)

---

## How releases map to integration levels

| Level | Description | Release |
|-------|-------------|---------|
| L1 Serving | `POST /chat/completions` | beta.1 ✅ |
| L2 Gateway | User staging OpenAI-compatible URL | beta.1 ✅ (config) |
| L3 Agent API | Custom agent endpoint + tool checks | v0.2 |
| L4 Workflow | Multi-step, replay, semantic eval | v1 |

---

## How to read feature docs

Detailed feature specs: [`docs/features/`](docs/features/)  
Architecture: [`docs/architecture.md`](docs/architecture.md)  
YAML levels: [`docs/yaml-overview.md`](docs/yaml-overview.md)

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md).