# Connor

> **The CI/CD reliability toolkit for AI systems.**

Versioned YAML suites, deterministic gates, `exit 0` / `exit 1` — before you merge.

Your agent stays Python (or any stack). Connor only needs an OpenAI-compatible HTTP endpoint.

**Not** a production observability platform. Connor runs in CI — like unit tests for your LLM serving and agents.

## What it does (30 seconds)

- Run YAML suites against any OpenAI-compatible endpoint (OpenAI, OpenRouter, vLLM, LiteLLM…)
- Gate on HTTP health, JSON syntax, JSON Schema, text contains
- Block merges with `exit 1` in GitHub Actions, Argo, or any CI

## Demo in 60 seconds

```bash
git clone https://github.com/noodl-labs/ConnorLLM.git
cd ConnorLLM
make install

export CONNOR_BASE_URL=https://openrouter.ai/api/v1   # must include /v1
export CONNOR_API_KEY=sk-your-key

# Serving smoke (multi-model)
connor run benchmarks/examples/serving-smoke.yaml

# Agent JSON Schema gate (CI-friendly — expects GATE PASSED)
connor run benchmarks/examples/agent-json-smoke.yaml
echo $?   # 0 = pass, 1 = fail
```

Use `-v` to show response bodies on passed cases.

## Quick start

**Requirements:** Go 1.22+

Copy [`.env.example`](.env.example) or export `CONNOR_BASE_URL` and `CONNOR_API_KEY`.

**Single case:**

```bash
connor run \
  --model openai/gpt-4o-mini \
  --prompt "Reply with exactly one word: pong"
```

**YAML suite:**

```bash
connor run benchmarks/examples/serving-smoke.yaml
```

## Example output

```text
Connor  v0.1.0-beta.3
Target  https://openrouter.ai/api/v1
Suite   agent-json-smoke (2 cases)

✓  flight-gpt4o-mini  openai/gpt-4o-mini                 901ms  HTTP 200   schema ✓
✓  flight-gemini      google/gemini-2.5-flash-lite       481ms  HTTP 200   schema ✓
────────────────────────────────────
2/2 passed · slowest 901ms (flight-gpt4o-mini) · total 1.4s
GATE PASSED — safe to merge
exit 0
```

## GitHub Actions (sketch)

```yaml
jobs:
  llm-gate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: services/runtime/go.mod
      - run: go install ./services/runtime/cmd/connor
      - name: Connor gate
        env:
          CONNOR_BASE_URL: ${{ secrets.CONNOR_BASE_URL }}
          CONNOR_API_KEY: ${{ secrets.CONNOR_API_KEY }}
        run: connor run benchmarks/examples/agent-json-smoke.yaml
```

## Example suites

| File | Purpose | CI use |
|------|---------|--------|
| [`serving-smoke.yaml`](benchmarks/examples/serving-smoke.yaml) | Multi-model serving smoke | Post-deploy health |
| [`agent-json-smoke.yaml`](benchmarks/examples/agent-json-smoke.yaml) | JSON Schema agent gate (2 models) | **Recommended for CI** |
| [`agent-json.yaml`](benchmarks/examples/agent-json.yaml) | Demo: 1 pass + 2 intentional fails | Documentation / local demo |
| [`agent-json-compare.yaml`](benchmarks/examples/agent-json-compare.yaml) | Cross-model JSON syntax compare | Benchmark / exploration |

## Architecture

Connor is built as six engines:

| Engine | Role |
|--------|------|
| **Execution** | Run providers, retries, timeouts |
| **Evaluation** | JSON, schema, contains (semantic eval planned in Python) |
| **Benchmark** | Multi-case suites, model comparison |
| **Quality Gates** | CI rules: pass/fail, future latency/cost thresholds |
| **Observability** | Run artifacts, replay (CI scope — not prod APM) |
| **Developer Experience** | CLI, YAML, GitHub Actions |

**Today:** Execution + deterministic Evaluation + DX.

**Details:** [docs/architecture.md](docs/architecture.md) · [Roadmap](ROADMAP.md)

## Documentation

| Doc | Description |
|-----|-------------|
| [Architecture](docs/architecture.md) | Six engines, DDD layout, data flow |
| [Roadmap](ROADMAP.md) | Releases, use-case matrix, exit criteria |
| [Changelog](CHANGELOG.md) | Release notes |

## Status

**Current:** `v0.1.0-beta.3` — smoke, JSON schema, `run.json` export, `connor compare` (p95)  
**Next (v0.1.0):** `--min-pass-rate`, CI handbook

## License

License TBD.
