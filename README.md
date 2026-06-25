# ConnorLLM

> Quality gates & regression suites for LLM endpoints — CI-native, CLI-first, PASS/FAIL before merge.

**Not** an observability platform. Connor runs **before** you merge — like unit tests for your LLM serving and agents.

## What it does (30 seconds)

- Run versioned YAML suites against any OpenAI-compatible endpoint
- Gate on HTTP success, latency, JSON validity
- Exit `0` or `1` for GitHub Actions / Argo / any CI

## Quick start

**Requirements:** Go 1.22+

```bash
git clone https://github.com/noodl-labs/ConnorLLM.git
cd ConnorLLM
make install   # puts connor on your PATH ($HOME/go/bin)
```

Set your endpoint (copy [`.env.example`](.env.example) or export manually):

```bash
export CONNOR_BASE_URL=https://openrouter.ai/api/v1
export CONNOR_API_KEY=sk-your-key
```

`CONNOR_BASE_URL` must include `/v1` (OpenAI, OpenRouter, vLLM, LiteLLM, etc.).

**Single case:**

```bash
connor run \
  --model openai/gpt-4o-mini \
  --prompt "Reply with exactly one word: pong"
```

**YAML suite:**

```bash
connor run benchmarks/examples/serving-smoke.yaml
echo $?   # 0 = pass, 1 = fail (CI gate)
```

## Example output

```text
Connor  v0.1.0-beta.1
Target  https://openrouter.ai/api/v1
Suite   serving-smoke (3 cases)

✓  ping-gemini   google/gemini-2.5-flash-lite      479ms  HTTP 200
✓  ping-llama    meta-llama/llama-3.1-8b-instruct  328ms  HTTP 200
✓  json-ok       openai/gpt-4o-mini                  630ms  HTTP 200   json ✓

────────────────────────────────────
3/3 passed · slowest 630ms (json-ok) · total 1.4s
GATE PASSED — safe to merge
exit 0
```

Use `connor run -v` to show body and attempts on passed cases.

## Use cases

- Post-deploy smoke on vLLM / gateway
- JSON structured-output gate for agents
- Multi-model latency checks in one suite
- *(Planned v0.1)* Regression compare vs baseline

## Documentation

| Doc | Description |
|-----|-------------|
| [Architecture](docs/architecture.md) | Six engines, code layout, data flow |
| [Roadmap](ROADMAP.md) | Shipped vs planned releases |
| [serving-smoke.yaml](benchmarks/examples/serving-smoke.yaml) | Example regression suite |

## Status

**Current:** `v0.1.0-beta.1` — smoke suites, JSON gate, `connor run`  
**Next:** `connor compare`, `run.json`, budget gates (latency, pass rate)

## License

License TBD.
