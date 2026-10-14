# Distributed Systems Lab

Controlled experiments for learning distributed systems. Current scope:
[Experiment 00: Reproducible Infrastructure](docs/experiments/00-reproducible-infrastructure.md).

One disposable Debian 13 VM, managed by Terraform/libvirt and configured by
Ansible. No application components yet.

Experiment 00 is complete. Creation, SSH, configuration, stop/start, teardown,
and clean recreation were verified, with existing service checks passing.
The lab VM, disks, storage pool, and NAT network have been destroyed for now.
Installed tools, cached images, and the SSH key are retained for recreation.

The [host package ledger](docs/experiments/00-host-packages.md) records all 84
added APT packages and their versions, plus the separate Terraform installation.

Read the experiment document for host prerequisites, approval boundaries,
and verification. To recreate after reviewing and approving NAT creation:

```sh
scripts/lab prepare
scripts/lab validate
scripts/lab plan
scripts/lab provision
scripts/lab configure
scripts/lab ssh
```

To stop and restart the existing VM:

```sh
scripts/lab stop
scripts/lab start
```

`scripts/lab destroy` removes the VM, its managed disks, pool, and NAT network.
It requires explicit approval and Terraform's interactive confirmation.
