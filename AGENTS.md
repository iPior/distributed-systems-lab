# Agent instructions — Distributed Systems Lab

This is a **learning project**, not a race to complete infrastructure. Build one small distributed application through focused experiments. Each experiment begins with a question and hypothesis, then tests the result using relevant evidence (correctness, latency, consistency, throughput, or failure behavior).

## Read first

- `README.md` for the project overview.
- `docs/ROADMAP.md` for the planned architecture stages.
- `docs/EXPERIMENT_TEMPLATE.md` for documenting experiments.
- `docs/experiments/00-guided-rebuild.md` for the **current task**.

## Safety and workflow

- Host: `optiplex-media-server`, Debian 13, i5-6500, 32 GB RAM; repository in `/opt/distributed-systems-lab`. Existing Docker media services must remain untouched. k6 will run on a separate Windows desktop.
- All lab VMs, networking, and disks must be isolated, named, and disposable. Keep the repository, Git credentials, and records outside disposable storage.
- **Teach before implementing.** Explain the next file's role and relevant trade-offs; ask the learner to predict what it does. Work on **one file or checkpoint at a time**, then stop for feedback.
- **Do not automatically implement an entire experiment, execute the full lifecycle, or advance to another stage.**
- Inspect current state before making changes; never assume the live host matches past reports.
- Obtain explicit approval before destructive operations, changing host networking/firewall rules, installing host packages, or running disruptive tests. Never run broad Docker prune/cleanup commands.
- Use IaC, small readable changes, and actual measurements. Never invent successful checks or experimental outcomes.

## Current focus

**Experiment 00: Guided rebuild.** Start with the architecture and `infra/experiment-00/versions.tf`, **without provisioning anything**. Explain KVM, QEMU, libvirt, Terraform, cloud-init, and Ansible in their respective roles. Then proceed through the checkpoints in `docs/experiments/00-guided-rebuild.md`, waiting for learner review at each step. The known-good initial implementation and historical measurements are preserved on `main`; this learning run must collect its own evidence.
