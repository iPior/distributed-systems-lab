# Distributed Systems Lab

Learn distributed systems by **building one small application in stages** and running focused experiments. Each experiment begins with a research question, a prediction, a controlled procedure, and real evidence. We examine correctness, consistency, resilience, latency, or throughput according to the concept being studied.

- [Roadmap](docs/ROADMAP.md)
- [Experiment template](docs/EXPERIMENT_TEMPLATE.md)
- [Module 00: Infrastructure](curriculum/00-infrastructure/README.md)

## Current status: Stage 01 complete; Stage 02 pending

The Go health API, [systemd supervision experiment](docs/experiments/01b-service-supervision.md), [Ansible deployment experiment](docs/experiments/01c-ansible-deployment.md), and [in-memory task experiment](docs/experiments/01d-in-memory-tasks.md) checkpoints are complete. Task tests and executable deployment passed; the old task returned 404 after a controlled restart with an unchanged executable. The [HTTP baseline checkpoint](docs/experiments/01e-http-baseline.md) is complete: corrected eight-user runs through SSH reached roughly 12,000 health requests/second with all checks passing. SSH timing obfuscation was isolated as a measurement confounder; sampled resource observations did not establish a bottleneck or capacity limit. The [higher-load observations](docs/experiments/01f-higher-load.md) showed a path plateau around 22k-23k requests/second and Windows SSH accept/file-limit errors at 256 active users. API maximum capacity remains unknown. Stage 01 is closed at the learner's request. Next session: design Stage 02 with three API instances behind HAProxy; no Stage 02 infrastructure has been implemented.

Post-closure access experiments: [host-local k6](docs/experiments/01g-host-local-k6.md) encountered shared-host CPU contention, so keep load generation on Windows. The learner then approved [moving the existing VM onto the LAN](docs/experiments/01h-lan-access.md) through bridge `br0`. Current DHCP address is `192.168.2.44`; host health requests and Ansible access pass, and Terraform reports no changes. The learner confirmed direct Windows health access. The [direct-LAN observations](docs/experiments/01i-direct-lan-baseline.md) reached 17,170.11 req/s at eight VUs and plateaued around 22.7k-23.1k req/s at 128-512 VUs, then fell to 21.3k at 1024 VUs; all reported checks passed with zero HTTP failures. Later 3000/5000-VU samples reached about 21.8k/20.1k req/s with p95 158.02/286.49 ms, also without failures. Under the learner's limits (p95 <=500 ms, p99 <=1 s, HTTP errors <=1%), two 7500-VU runs passed at about 19.6k-20.0k req/s; 7750 and 8000 failed on p95. Exact p99 is bounded by reported maxima rather than printed directly. This is an observed concurrency range, not a stable maximum. The limiting component and API capacity remain unknown; stop increasing load. Stage 02 remains unprovisioned.

**Stage 01, including the LAN and load follow-ups, is now closed at the learner's request.** Continue with a separate agent using the [Stage 02 handoff](docs/STAGE_02_HANDOFF.md). First agree the three-instance/HAProxy topology and investigate independent-store task correctness; performance comes afterward. No Stage 02 provisioning is approved by this handoff.

The guided infrastructure rebuild and lifecycle validation are complete. The active implementation lives in [`infrastructure/`](infrastructure/), with measured results and a repeatable checklist in [the learning run](docs/experiments/00-learning-run.md).

The original agent-built implementation was removed from the working tree after completing Module 00. Its source, scripts, and historical reports remain available in [Git history at commit `67c0f75`](https://github.com/iPior/distributed-systems-lab/tree/67c0f75/reference/experiment-00). Those results are separate from the new learning run.

The [seven-lesson curriculum](curriculum/00-infrastructure/README.md) remains available for another guided rebuild. Stage 01 currently uses the experiment notes rather than a separate curriculum. Continue one lesson or experiment at a time, with predictions, learner-written steps, and observed evidence.

The Optiplex runs unrelated Docker services: keep lab networks, disks, and VM resources isolated; require approval for changes that could affect the host.

## Future direction

After mastering VM provisioning, gradually introduce a Go API, HAProxy and replicas, Redis sessions, PostgreSQL primary and read replicas, caching and invalidation, RabbitMQ workers, and network fault injection. Every component must support a specific learning question.
