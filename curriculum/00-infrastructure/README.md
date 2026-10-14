# Module 00 — Build a disposable Debian VM

**Goal:** Understand, implement, and reproduce one isolated Debian VM on the Optiplex. Read [the experiment template](../../docs/EXPERIMENT_TEMPLATE.md), and record *your own* findings in `docs/experiments/00-learning-run.md`.

## How to use this module

Work through one lesson at a time. At every lesson:
1. **Learn:** Read the short explanation and inspect the relevant tool's own help.
2. **Predict:** Write what you expect before running commands.
3. **Build:** Implement a small step yourself. Codex may review, but must not supply a complete solution unprompted.
4. **Verify:** Gather command output and check assumptions.
5. **Reflect:** Note what you learned. Stop before proceeding.

**Safety:** The Optiplex runs unrelated media services. Request explicit approval before installs, host networking changes, VM creation, or destruction. Never clean up outside the project's resources. The original known-good implementation is available in `reference/experiment-00/` if you get stuck; avoid copying it wholesale.

## Lessons

1. [Virtualization and isolation](01-virtualization.md) — understand the layers, inspect the host.
2. [Manual VM lifecycle](02-manual-vm.md) — learn libvirt concepts using a temporary guest.
3. [Guest networking](03-networking.md) — NAT, DHCP, DNS, and SSH access.
4. [Terraform fundamentals](04-terraform-basics.md) — providers, plan, state, apply, destroy.
5. [Terraform VM](05-terraform-vm.md) — reproduce guest, network, disks, and cloud-init declaratively.
6. [Ansible](06-ansible.md) — configure and prove idempotency.
7. [Validation and evidence](07-validation.md) — teardown, recreate, verify service isolation and document.

The module is a guide, not a command dump. If a tool behaves differently from the lesson's assumptions, investigate and document why.
