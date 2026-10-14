# AGENTS.md — Distributed Systems Lab

## Purpose

This is a learning project: progressively build a small distributed application and use **focused experiments** to understand its behavior. Each experiment asks one question, tests a hypothesis, and documents actual observations. Don't optimize every stage for throughput; choose measurements that fit the question (latency, correctness, consistency, recovery, etc.). Never invent results.

## Environment

- Repo: `ipior/distributed-systems-lab`, hosted at `/opt/distributed-systems-lab` on `optiplex-media-server`.
- Host: Debian 13, Intel i5-6500 (4 cores), 32 GB RAM, 1 TB SSD; KVM support confirmed.
- The host already runs important Docker media/homelab services. **Do not modify, disrupt, or clean up unrelated resources.**
- Intended stack: Terraform + libvirt for disposable VMs, Ansible for configuration, Docker Compose and Go for the application. Later experiments may add HAProxy, Redis, PostgreSQL replicas, RabbitMQ, and fault injection.
- k6 load tests will run on a separate Windows PC, not on the Optiplex.

## Working rules

1. Work on **one experiment at a time**. Do not implement future stages prematurely.
2. Keep code and infrastructure simple, readable, reproducible, and version-controlled.
3. Document each experiment in its own Markdown file: **question, hypothesis, setup, procedure, observations, interpretation, and trade-offs**. Leave results blank until measured.
4. Change one significant variable at a time where practical; record relevant hardware, resource limits, and test conditions.
5. Keep lab VMs, networks, disks, and containers isolated and clearly named. Lab infrastructure must be safe to stop, destroy, and rebuild without affecting existing services or deleting the repository.
6. **Ask before destructive commands, host networking/firewall changes, or disruptive tests.** Never run broad Docker prune or host cleanup commands.
7. Inspect the repository and current environment before acting. Explain meaningful decisions and report what was actually verified.

## Current focus: Experiment 00 — Reproducible Infrastructure

**Question:** Can we provision, access, stop, destroy, and recreate an isolated Debian VM entirely from code?

Start with a single VM named `dsl-node-01` (1 vCPU, 2 GB RAM, 20 GB thin-provisioned disk). Verify tooling first; do not assume libvirt or Terraform has been installed. Provision with Terraform/libvirt and configure access with Ansible or cloud-init. Document the setup, VM lifecycle, validation, and safe teardown in `docs/experiments/00-reproducible-infrastructure.md`.

**Done when:** the VM can be created, accessed, stopped, destroyed, and recreated without affecting other Optiplex services. Do not deploy the Go API or any distributed components yet.
