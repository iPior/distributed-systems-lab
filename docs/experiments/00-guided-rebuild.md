# Experiment 00: Guided rebuild

**Status:** learning in progress — not a claim that this branch's rebuild is complete.

## Purpose

Rebuild understanding of the VM lifecycle one file at a time, instead of asking Codex to build an entire working lab in one pass. The completed implementation and measurements are preserved on `main` in `docs/experiments/00-reproducible-infrastructure.md`. Treat them as a reference, **not** as measurements from this rebuild.

**Question:** Can we independently explain, create, configure, stop, destroy, and recreate one disposable Debian 13 guest while the existing Optiplex services remain unaffected?

## Guided sequence (stop at each checkpoint)

1. **Architecture first** — Sketch the host, hypervisor, libvirt, NAT network, VM disk, guest and SSH paths. Explain process vs VM isolation. Compare NAT and bridge access. Do not change the host.
2. **Read the existing `infra/experiment-00/versions.tf`** — Explain Terraform's role, provider pinning, state, and the `qemu:///system` connection. Do not change the provider yet.
3. **Study the Terraform resources in `main.tf` one by one** — pool → virtual network → disk → cloud-init → domain. Understand their dependencies and what each `plan` would create. Keep the known-good XSL template until its reason is understood.
4. **Study `scripts/lab`** — Identify what `prepare`, `validate`, `plan`, `provision`, `stop`, `start`, and `destroy` do and which commands modify the host. Prefer manually stepping through a command before accepting convenience automation.
5. **Study `ansible/configure.yml`** — Explain what cloud-init does only on first boot and what Ansible can converge repeatedly. Verify two Ansible runs for idempotence.
6. **Approved live lifecycle** — After checking current host/network/storage state and getting permission for NAT changes: plan → provision → SSH → configure → stop/start → recheck → destroy → recreate. Use safe scopes; no host-wide pruning or firewall resets.
7. **Document your own run** — Record the actual commands, outputs, differences from the original run, failures, and lessons in a *new run record*. Do not overwrite historical evidence.

## Intended resources

- Host: Optiplex Debian 13, Intel i5-6500, 32 GB RAM, also hosting unrelated services.
- One VM: `dsl-node-01`, 1 vCPU, 2 GiB RAM, 20 GiB sparse disk.
- Guest: Debian 13, host-only access through the isolated libvirt NAT network and SSH.
- All lab resources must be named and destroyable without affecting non-lab resources.

## Checkpoint rule

Ask the learner to explain the responsibility of the current file and predict a command's effect. Offer a focused explanation or small edit, then stop for review. **Do not execute the entire plan or move to the next file automatically.**

## Completion criteria

- Explain the roles of KVM, QEMU, libvirt, Terraform, cloud-init, and Ansible.
- Predict what `terraform plan` and `destroy` will change.
- Show working guest access and idempotent configuration.
- Show a verified stop, destroy, and fresh recreation.
- Validate the existing media stack and host networking remain intact.
- Distinguish observations from prior run records.
