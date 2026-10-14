# Distributed Systems Lab

Learn distributed systems by **building one small application in stages** and running focused experiments. Each experiment begins with a research question, a prediction, a controlled procedure, and real evidence. We examine correctness, consistency, resilience, latency, or throughput according to the concept being studied.

- [Roadmap](docs/ROADMAP.md)
- [Experiment template](docs/EXPERIMENT_TEMPLATE.md)
- [Start Module 00 — Infrastructure](curriculum/00-infrastructure/README.md)

## Current work: Rebuild Experiment 00 by hand

The first infrastructure experiment was originally completed by an agent. This branch preserves that known-good implementation under [`reference/experiment-00/`](reference/experiment-00/) while reserving **`infrastructure/`** for the learner's new implementation. The old reference scripts retain their original paths internally and must **not be executed in-place**.

Work through the [seven lessons](curriculum/00-infrastructure/README.md) one at a time. Do not skip ahead or ask an agent to finish the module. Record new observations in [the learning run](docs/experiments/00-learning-run.md), separately from the historical results archived with the reference.

The Optiplex runs unrelated Docker services: keep lab networks, disks, and VM resources isolated; require approval for changes that could affect the host.

## Future direction

After mastering VM provisioning, gradually introduce a Go API, HAProxy and replicas, Redis sessions, PostgreSQL primary and read replicas, caching and invalidation, RabbitMQ workers, and network fault injection. Every component must support a specific learning question.
