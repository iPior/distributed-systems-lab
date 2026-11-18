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

**Latest access checkpoint:** the learner approved a LAN attachment for the existing single API VM after [host-local k6](experiments/01g-host-local-k6.md) showed shared-host CPU contention. [Experiment 01h](experiments/01h-lan-access.md) records the completed move to existing bridge `br0`, current DHCP address `192.168.2.44`, passing host HTTP/Ansible checks and a clean Terraform plan. Learner confirmed direct Windows health access. The [direct-LAN observations](experiments/01i-direct-lan-baseline.md) reached 17,170.11 req/s at eight VUs, 22,746.77 at 128, 22,889.96 at 256, 23,131.68 at 512, and 21,279.19 at 1024; all checks passed with no reported HTTP failures. Throughput plateaued at 128-512 and degraded at higher concurrency. Subsequent 3000/5000-VU samples reached 21,835.06/20,064.11 req/s with p95 158.02/286.49 ms, no failures. The learner's criteria are p95 <=500 ms, p99 <=1 s (bounded by maxima, not printed), and HTTP errors <=1%. Later tests found two passes at 7500 VUs (p95 420.60/430.63 ms) and failures at 7750/8000 (p95 549.62/550.04 ms), with zero HTTP/check failures. This is an observed concurrency range, not a stable maximum; highest observed passing throughput remains about 23.1k req/s at 512 VUs. Limiting component and API capacity remain unknown. Stop increasing load; keep k6 on Windows and agree the next focused experiment before proceeding. This later approved topology supersedes the earlier preference to retain the private guest attachment. Stage 02 remains unstarted. Learner has now explicitly closed Stage 01 including these follow-ups and requested a separate Stage 02 agent. Use the [handoff](STAGE_02_HANDOFF.md); do not automatically repeat or increase Stage 01 load. Start with topology review and task correctness across independent stores, then plan performance comparisons. Matched resource observations and repeatability remain requirements for later capacity claims.

**Module 00 is complete.** The [learning run](experiments/00-learning-run.md) records the guided rebuild, lifecycle verification, and limitations. Active infrastructure lives in `infrastructure/`; the original implementation is retained in [Git history at commit `67c0f75`](https://github.com/iPior/distributed-systems-lab/tree/67c0f75/reference/experiment-00).

**Stage 01 is complete at the learner's request.** Completed checkpoints cover [HTTP correctness and listener binding](experiments/01-single-api.md), [systemd supervision](experiments/01b-service-supervision.md), [repeatable Ansible artifact deployment](experiments/01c-ansible-deployment.md), and [in-memory task state](experiments/01d-in-memory-tasks.md). Task tests and executable deployment passed; the old task returned 404 after restart with an unchanged executable. Baseline performance and overload observations are recorded below; the API's maximum capacity is not established.

Latest checkpoint: [HTTP baseline through SSH](experiments/01e-http-baseline.md), complete with corrected 1/2/4/8-user runs and sampled resource observations. Roughly 12,000 health requests/second at eight users is an observed path baseline, not an API capacity limit. The [higher-load checkpoint](experiments/01f-higher-load.md) observed a path plateau around 22k-23k req/s and a Windows SSH accept/file-resource limit at 256 users. Stage 01 closes with these limitations. **Stage 02 is pending:** next session, agree the topology for three API instances behind HAProxy, predict per-replica in-memory state behavior, and proceed one learner checkpoint at a time. No replicas or HAProxy have been provisioned. The [guided rebuild plan](experiments/00-guided-rebuild.md) remains available for revisiting infrastructure.
