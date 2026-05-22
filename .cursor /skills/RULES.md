# ConnorLLM Engineering Rules

## Project Vision

ConnorLLM is an AI Reliability Benchmark Runner focused on:

- runtime reliability
- benchmarking
- observability
- regression detection
- structured output validation
- stress testing for production LLM systems

The project is NOT:
- a chatbot
- a prompt playground
- an AI wrapper SaaS

The project MUST prioritize:
- reliability engineering
- runtime systems thinking
- observability
- deterministic benchmarking
- clean architecture
- maintainability

---

# Architecture Rules

## Runtime Layer

The Go runtime is responsible for:
- provider execution
- retries
- timeout propagation
- fallback policies
- runtime metrics
- tracing
- benchmark orchestration

The runtime layer MUST NOT:
- contain ML scoring logic
- contain embedding logic
- contain hallucination detection algorithms

---

## Evaluation Layer

The Python evaluation service is responsible for:
- semantic similarity
- hallucination scoring
- regression analysis
- groundedness evaluation
- structured output validation

The evaluation layer MUST remain stateless.

---

# Design Principles

## Clean Architecture

Use strict separation between:
- domain
- application
- infrastructure
- interfaces

Business rules MUST stay in domain/application.

Infrastructure code MUST remain isolated.

---

## Domain Rules

Domain entities must:
- be framework-agnostic
- contain core business semantics
- avoid infrastructure coupling

---

## Reliability First

Every runtime feature should consider:
- retries
- timeout propagation
- cancellation
- observability
- metrics
- graceful degradation

---

## Observability First

Every important runtime operation must expose:
- logs
- metrics
- traces

Important metrics:
- latency
- p95/p99
- TTFT
- retries
- throughput
- timeout rate
- fallback rate

---

# Benchmark Philosophy

Benchmarks must be:
- deterministic
- reproducible
- observable
- replayable

Benchmark scenarios should measure:
- runtime stability
- structured output reliability
- regression behavior
- provider consistency

---

# Coding Rules

## Go

- Prefer composition over inheritance
- Keep interfaces small
- Use context.Context everywhere
- Avoid global state
- Prefer explicit dependency injection
- Use zerolog for structured logging

---

## Python

- Use typed Pydantic models
- Keep evaluation services modular
- Avoid business logic inside routes
- Keep evaluation deterministic

---

# CLI Philosophy

The CLI should feel:
- operational
- infra-oriented
- benchmark-centric

Examples:

connor run benchmark.yaml
connor compare run_a run_b
connor stress benchmark.yaml
connor replay latest

---

# Performance Philosophy

ConnorLLM is a runtime reliability platform.

Prioritize:
- low overhead
- predictable latency
- stable execution
- clean failure handling

Do not optimize prematurely for:
- distributed orchestration
- Kubernetes complexity
- multi-region execution

Focus first on:
- correctness
- observability
- benchmark reproducibility

---

# Forbidden Anti-Patterns

Avoid:
- giant service files
- hidden shared state
- tightly coupled infrastructure
- business logic inside handlers
- unstructured logs
- magic retry behavior
- silent failures

---

# Long-Term Direction

ConnorLLM should evolve toward:
- AI runtime reliability engineering
- AI observability
- benchmark infrastructure
- AI SRE tooling
- production LLM validation systems