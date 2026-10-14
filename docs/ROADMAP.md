# Learning roadmap

This is a single application that evolves through **architecture stages**. An experiment is a focused question evaluated against a stage; one stage can support multiple experiments.

| Stage | Change | Example question and evidence |
| --- | --- | --- |
| 00 | Disposable Debian VM with Terraform, cloud-init, Ansible | Can a VM be provisioned, configured twice without drift, stopped, destroyed, and recreated safely? Lifecycle checks, plans, host health |
| 01 | Single Go HTTP API | What is baseline behavior? Correctness, latency, resource usage |
| 02 | Three API instances + HAProxy | How does adding replicas affect throughput and error rates with controlled resources? k6 results |
| 03 | Redis-backed sessions | Does any replica preserve the same session? Routed requests, correctness |
| 04 | PostgreSQL primary | What are transaction semantics and write latency? Transaction tests |
| 05 | Two PostgreSQL read replicas | How does asynchronous replication affect read distribution and lag? Read results, WAL lag |
| 06 | Consistency policies | Can readers observe stale data after successful writes? Staleness windows |
| 07 | Redis application cache | How do hits/misses affect latency and DB load? p95, cache hit rate |
| 08 | Invalidation | Can an invalidated cache be repopulated with stale replica data? Deterministic ordering, version comparisons |
| 09 | RabbitMQ and Go workers | What guarantees do retries, acknowledgements, and idempotency provide? Job correctness |
| 10 | Fault injection | How do outages and partitions affect availability and consistency? Error/response timelines |

## Rules

- Start with the research question and expected system behavior. Pick measurements to suit it; not every experiment is a throughput test.
- Change one major variable at a time when practical. Record environment, resource limits, application revision, methodology, and limitations.
- **Performance:** latency/throughput; **correctness:** invariant checks and state history; **resilience:** error rates and recovery.
- Keep k6 on the separate Windows desktop when measuring service capacity.
- Separate **proposed procedures** from **observed results**. No invented measurements.
- Review a proposed architecture before Codex implements it. Only proceed one experiment at a time.
- Existing media-stack services and host networking are outside the project: request explicit approval for any operation that could affect them.

## Current learning objective

Revisit **Experiment 00** interactively, using [the guided rebuild plan](experiments/00-guided-rebuild.md). The original complete implementation and measured report remain available on `main` as a reference. Do not jump to Stage 01 until the learner can explain and reproduce the VM lifecycle.
