# ConnorLLM Architecture

> **ConnorLLM** is the CI/CD reliability toolkit for AI systems.

Connor runs **before merge** — not in production observability. It executes versioned YAML suites against any OpenAI-compatible endpoint and returns `exit 0` or `exit 1` for CI.

Your AI stack (Python agents, vLLM, gateways) stays unchanged. Connor only needs HTTP.

---

## Six engines

| # | Engine | Role | Today | Target |
|---|--------|------|-------|--------|
| 1 | **Execution** | Providers, retry, timeout, suite runs | `services/runtime/` (Go) | Agent runner, tools (v0.2) |
| 2 | **Evaluation** | JSON, schema, contains; later semantic eval | Go: syntax + schema + text | `services/evaluation/` (Python, v1) |
| 3 | **Benchmark** | Multi-case suites, model comparison | YAML suites | `connor compare` (v0.1) |
| 4 | **Quality Gates** | CI pass/fail, budgets | `exit 0/1`, fail reasons | Latency / pass-rate thresholds (v0.1) |
| 5 | **Observability** | Run history, replay (CI scope) | — | `run.json` + store (v0.1 → v1) |
| 6 | **Developer Experience** | CLI, parser, docs, GitHub Actions | `connor run`, human output | SDK, feature docs |

**Shipped today:** Execution (partial) + Evaluation (deterministic gates) + DX (CLI/YAML).

---

## Data flow

```text
YAML suite  →  benchmark.Parse  →  application.ExecuteSuite
                                        │
                                        ▼
                              application.ExecuteCase
                              (HTTP via openai_compatible)
                                        │
                                        ▼
                              domain.validation.Evaluate
                              (contains → JSON → schema)
                                        │
                                        ▼
                              entities.CaseResult
                                        │
                                        ▼
                              cli/output.PrintRun  →  exit 0 | 1
```

### Layering (DDD)

| Layer | Path | Responsibility |
|-------|------|----------------|
| Benchmark (infra) | `internal/benchmark/` | YAML → `Spec`, parse-time validation |
| Application | `internal/runtime/application/` | Orchestrate cases, wire expectations |
| Domain | `internal/runtime/domain/` | `Request`, `Response`, `Expectations`, gates |
| Infrastructure | `internal/runtime/infrastructure/` | OpenAI-compatible HTTP client |
| CLI | `internal/cli/` | `connor run`, terminal output |

Domain code does not import YAML or HTTP client types.

---

## Active gates (beta.2)

| YAML field | Fail reason | Question |
|------------|-------------|----------|
| *(implicit)* | `call_failed` | HTTP 2xx? Timeout? Retries exhausted? |
| `expect_json` | `invalid_json` | Valid JSON syntax? |
| `expect_json_schema` | `schema_mismatch` | Matches JSON Schema? (syntax implied) |
| `expect_contains` | `content_mismatch` | Body contains substring? |
| `expect_contains_ignore_case` | `content_mismatch` | Case-insensitive contains |

Evaluation order: **contains → JSON syntax → JSON schema**.

---

## Repository layout

```text
ConnorLLM/
├── services/runtime/           # Execution + Evaluation (Go)
│   ├── cmd/connor/             # CLI entrypoint
│   └── internal/
│       ├── benchmark/          # YAML parser
│       ├── cli/                # Commands + output
│       └── runtime/
│           ├── application/    # ExecuteSuite, EvaluateCase
│           ├── domain/         # Entities, validation, reliability
│           └── infrastructure/ # openai_compatible provider
├── benchmarks/examples/        # Runnable demo suites
├── docs/                       # Architecture, features (growing)
└── ROADMAP.md
```

**Planned:** `services/evaluation/` (Python) for semantic / groundedness checks (v1).

---

## Integration levels

| Level | Description | Status |
|-------|-------------|--------|
| L1 Serving | `POST /chat/completions` | ✅ beta.1 |
| L2 Gateway | Staging/prod OpenAI-compatible URL | ✅ config only |
| L3 Agent API | Custom agent endpoint + tool gates | 🔜 v0.2 |
| L4 Workflow | Multi-step, replay, semantic eval | 🔜 v1 |

---

## Non-goals

- Production APM / tracing (Langfuse territory)
- Academic model benchmarks (MMLU)
- Replacing your agent runtime language (Python, etc.)

---

## Further reading

- [ROADMAP.md](../ROADMAP.md) — releases, use-case matrix, exit criteria
- [README.md](../README.md) — quick start and demo
- Example suites: `benchmarks/examples/`
