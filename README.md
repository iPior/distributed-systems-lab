# Distributed Systems Lab

Learn distributed systems by **building one small application in stages** and running focused experiments. Each experiment begins with a research question, a prediction, a controlled procedure, and real evidence. We examine correctness, consistency, resilience, latency, or throughput according to the concept being studied.

- [Roadmap](docs/ROADMAP.md)
- [Experiment template](docs/EXPERIMENT_TEMPLATE.md)
- [Module 00: Infrastructure](curriculum/00-infrastructure/README.md)

## Current status: Module 00 complete

The guided infrastructure rebuild and lifecycle validation are complete. The active implementation lives in [`infrastructure/`](infrastructure/), with measured results and a repeatable checklist in [the learning run](docs/experiments/00-learning-run.md).

The original agent-built implementation was removed from the working tree after completing Module 00. Its source, scripts, and historical reports remain available in [Git history at commit `67c0f75`](https://github.com/iPior/distributed-systems-lab/tree/67c0f75/reference/experiment-00). Those results are separate from the new learning run.

The [seven-lesson curriculum](curriculum/00-infrastructure/README.md) remains available for another guided rebuild. Continue one lesson or experiment at a time, with predictions, learner-written steps, and observed evidence. Begin the Go API stage only when requested.

The Optiplex runs unrelated Docker services: keep lab networks, disks, and VM resources isolated; require approval for changes that could affect the host.

## Future direction

After mastering VM provisioning, gradually introduce a Go API, HAProxy and replicas, Redis sessions, PostgreSQL primary and read replicas, caching and invalidation, RabbitMQ workers, and network fault injection. Every component must support a specific learning question.
