# ConnorLLM — Engineering Rules

## Project Identity

ConnorLLM is an:

```txt
AI Runtime Reliability & Benchmark Infrastructure
```

for production LLM systems.

The platform focuses on:

- runtime reliability
- deterministic benchmarking
- observability
- regression detection
- structured output validatio
- stress testing
- QoS enforcement

ConnorLLM is designed as an operational reliability system for AI workloads — not as a chatbot framework or prompt playground.

---

# Core Philosophy

ConnorLLM treats LLM systems as:

```txt
probabilistic distributed runtime systems
```

rather than simple API wrappers.

The project prioritizes:

- reliability engineering
- runtime systems thinking
- observability
- reproducibility
- deterministic validation
- clean architecture
- operational correctness

---

# Non-Goals

ConnorLLM is NOT:

- a chatbot application
- a prompt playground
- an agent framework
- a no-code AI platform
- a generic OpenAI wrapper
- a RAG framework

Avoid product decisions that drift toward these directions.

---

# System Architecture

ConnorLLM is split into two major layers:

```txt
Runtime Reliability Layer
vs
Evaluation Intelligence Layer
```

---

# Runtime Layer (Go)

The Go runtime acts as the operational control plane.

Responsibilities:

- provider execution
- benchmark orchestration
- retries
- timeout propagation
- cancellation handling
- fallback policies
- runtime metrics
- tracing
- observability
- replay execution
- stress execution
- QoS validation

The runtime layer MUST NOT contain:

- embedding logic
- semantic similarity algorithms
- hallucination detection
- ML scoring logic
- evaluation heuristics

The runtime layer must remain:

- deterministic
- operational
- observable
- low-overhead

---

# Evaluation Layer (Python)

The Python evaluation service handles semantic evaluation.

Responsibilities:

- semantic similarity
- hallucination scoring
- groundedness evaluation
- regression analysis
- structured output scoring

The evaluation layer MUST:

- remain stateless
- remain deterministic
- expose typed contracts
- avoid runtime orchestration logic

The evaluation service is a pure evaluation engine.

---

# Architecture Rules

ConnorLLM follows strict Clean Architecture principles.

Use strict separation between:

- domain
- application
- infrastructure
- interfaces

---

## Domain Layer

The domain layer contains:

- business semantics
- reliability rules
- QoS policies
- benchmark entities
- validation logic

Domain code MUST:

- remain framework-agnostic
- avoid infrastructure coupling
- avoid HTTP logic
- avoid provider-specific logic

---

## Application Layer

The application layer orchestrates:

- benchmark execution
- retry flows
- replay execution
- QoS evaluation
- runtime coordination

Application code may depend only on:

- domain
- ports/interfaces

---

## Infrastructure Layer

Infrastructure contains:

- provider clients
- storage
- tracing
- metrics
- HTTP clients
- OpenTelemetry integration

Infrastructure MUST remain isolated from business rules.

---

## Interface Layer

The interface layer contains:

- CLI commands
- HTTP handlers
- external adapters

No business logic inside handlers.

---

# Reliability Engineering Rules

Every runtime feature MUST consider:

- retries
- timeout propagation
- cancellation
- observability
- metrics
- graceful degradation
- deterministic failure handling

Avoid:

- silent failures
- hidden retries
- implicit fallback behavior
- retry storms
- unbounded concurrency

---

# Observability Rules

Observability is mandatory.

Every important runtime operation must expose:

- structured logs
- metrics
- traces

---

## Mandatory Metrics

ConnorLLM must support:

- latency
- p50
- p95
- p99
- TTFT
- retries
- throughput
- timeout rate
- fallback rate
- error rate
- cost/request

---

## Logging Rules

Use:

- structured logging only
- zerolog for Go services
- correlation IDs
- request IDs
- benchmark run IDs

Avoid:

- unstructured logs
- println debugging
- hidden runtime behavior

---

# Benchmark Philosophy

Benchmarks must be:

- deterministic
- reproducible
- observable
- replayable
- versioned

Benchmark suites should evaluate:

- runtime stability
- provider consistency
- structured output reliability
- latency regressions
- timeout behavior
- fallback correctness

---

# QoS Philosophy

ConnorLLM treats LLM systems as services with operational guarantees.

QoS policies may define:

- maximum p95 latency
- maximum timeout rate
- minimum JSON validity
- maximum fallback rate
- minimum success rate

QoS evaluation must produce deterministic PASS / FAIL decisions.

---

# Coding Rules — Go

## General Principles

- prefer composition over inheritance
- keep interfaces small
- avoid premature abstractions
- avoid hidden shared state
- prefer explicit dependency injection
- use context.Context everywhere

---

## Concurrency Rules

Concurrency must remain:

- bounded
- observable
- cancellable

Avoid:

- uncontrolled goroutines
- goroutine leaks
- unbounded worker pools
- hidden async execution

---

## Error Handling

Errors must:

- be explicit
- propagate context
- remain observable

Avoid:

- swallowed errors
- panic-driven flows
- hidden retry loops

---

# Coding Rules — Python

## General Principles

- use typed Pydantic models
- keep services modular
- isolate evaluation logic
- avoid business logic inside routes
- keep evaluation deterministic

---

# CLI Philosophy

The CLI should feel:

```txt
operational
infra-oriented
benchmark-centric
```

Example commands:

```bash
connor run benchmark.yaml
connor compare run_a run_b
connor stress benchmark.yaml
connor replay latest
connor gate latest
```

The CLI is an operational engineering tool, not a consumer product interface.

---

# Performance Philosophy

ConnorLLM is a runtime reliability platform.

Prioritize:

- correctness
- reproducibility
- observability
- predictable latency
- clean failure handling
- low runtime overhead

Avoid premature optimization for:

- Kubernetes orchestration
- distributed scheduling
- multi-region execution
- large-scale clustering

Focus first on:

- local correctness
- benchmark reproducibility
- deterministic execution

---

# Forbidden Anti-Patterns

Avoid:

- giant service files
- hidden shared state
- tightly coupled infrastructure
- business logic inside handlers
- magic retry behavior
- silent failures
- implicit side effects
- global mutable state
- provider-specific domain logic

---

# Long-Term Direction

ConnorLLM should evolve toward:

- AI Runtime Reliability Engineering
- AI Observability
- Benchmark Infrastructure
- AI SRE Tooling
- Production LLM Validation Systems
- QoS Validation Infrastructure
- Runtime Governance Systems
- Deterministic Replay Infrastructure

---

# Engineering Mindset

ConnorLLM should be engineered like:

- a distributed systems tool
- an observability platform
- a runtime validation engine
- an SRE-oriented infrastructure system

not like:

- a consumer AI application
- a prompt engineering demo
- a generic AI wrapper