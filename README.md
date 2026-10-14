# Distributed Systems Lab

An experiment-driven learning project: evolve one small application from a single VM toward a distributed system. Introduce components only to answer a focused question, measure what matters to that question, and document actual observations and trade-offs.

The planned progression includes a Go API, multiple replicas and HAProxy, Redis sessions and caching, PostgreSQL replication, RabbitMQ workers, and network fault experiments. See [the roadmap](docs/ROADMAP.md) and [experiment template](docs/EXPERIMENT_TEMPLATE.md).

## Current stage: Experiment 00 — guided rebuild

We are **revisiting** the initial infrastructure experiment to understand each file and decision before running it. Follow [the guided, checkpoint-by-checkpoint rebuild](docs/experiments/00-guided-rebuild.md). Do not let an agent perform every step at once.

The original Experiment 00 was already completed and documented: [reproducible infrastructure report](docs/experiments/00-reproducible-infrastructure.md). That record describes a Terraform/libvirt Debian VM, Ansible configuration, stop/start, destruction, and recreation, with host-service checks. **Those are historical results, not proof of a new guided run.** The known-good implementation remains in this branch as study material and on `main` as a reference.

### Existing implementation

- `infra/experiment-00/`: Terraform for one VM, private NAT, managed disks, and cloud-init.
- `ansible/`: guest configuration and verification.
- `scripts/lab`: explicit prepare, validate, plan, provision, configure, SSH, stop, start, destroy commands.
- `docs/experiments/00-host-packages.md`: existing host dependency ledger.

The Optiplex also hosts unrelated media services. Lab changes must not disrupt them. Require explicit approval for host-level, network-changing, and destructive operations.

### Working agreement

Study the next file, explain its purpose and expected behavior, make only the agreed change, validate it, and pause. When a live lifecycle is approved, record the new run separately rather than overwriting historical findings.
