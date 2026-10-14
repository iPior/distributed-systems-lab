# Experiment 00: Reproducible Infrastructure

## Objective and question

Can we provision, access, stop, destroy, and recreate an isolated Debian 13 VM
entirely from code without affecting the Optiplex's existing services?

**Status: complete**, measured 2026-10-14 UTC (2026-10-13 local time).
The VM and all five lab resources were destroyed at the user's request after
verification. Host tooling, image caches, SSH key, and local evidence remain
available for recreation. No application, database, or container was deployed.

## Hypothesis

A pinned Debian cloud image, Terraform/libvirt definitions, and an idempotent
Ansible playbook will recreate the same guest configuration. Dedicated resource
names, storage, and a nonoverlapping NAT subnet should allow the lab to coexist
with existing services. Success requires measured lifecycle and service checks,
not just a successful Terraform plan.

## Setup

| Item | Final configuration |
| --- | --- |
| Host | Optiplex, Debian 13.6 at inspection, i5-6500, 4 cores, 32 GB RAM |
| Available resources at inspection | About 26 GiB RAM and 431 GiB disk space |
| VM | `dsl-node-01`, KVM, 1 vCPU, 2048 MiB RAM |
| Guest image | Debian 13 genericcloud amd64, build `20261001-2618` |
| Observed guest OS | Debian 13.7 |
| Root disk | Standalone sparse qcow2, 20 GiB virtual capacity |
| Storage pool | `dsl-exp00`, `/var/lib/libvirt/images/dsl-exp00` |
| NAT network | `dsl-exp00`, bridge `virbr-dsl00`, `192.168.100.0/24` |
| Addresses | Host gateway `192.168.100.1`; DHCP reservation `192.168.100.10` |
| Guest access | SSH as `lab`, dedicated Ed25519 key in `.local/` |
| Terraform / provider | Terraform 1.13.5; `dmacvicar/libvirt` 0.8.3 |
| Other verified tools | libvirt 11.3.0, QEMU 10.0.13, Ansible 2.19.11 |
| Image/XML tools | xorriso 1.5.6, xsltproc 1.1.35 |

### Responsibilities and ownership

KVM provides hardware virtualization; QEMU provides the virtual hardware.
Libvirt manages QEMU, storage, and networking through `qemu:///system`.
Terraform owns exactly five resources: storage pool, NAT network, root volume,
cloud-init seed disk, and VM. It deletes the VM before its disks and pool.
The repository is outside the pool. Never place unrelated files in the pool.

`scripts/lab prepare` verifies the original image against a pinned SHA-512,
converts a standalone qcow2 copy, and expands its virtual capacity to 20 GiB.
Terraform uploads the prepared image into the dedicated pool. There are no shell
provisioners in Terraform and no external backing-file dependency.

Cloud-init supplies the hostname, `lab` user, public SSH key, DHCP configuration,
and filesystem expansion. Password SSH and root login are disabled. The
administrative `lab` account has passwordless sudo inside the disposable guest.
Ansible waits for cloud-init, verifies Debian 13 and one CPU, sets the hostname,
and ensures the basic SSH/Python packages and SSH service are present.

The pinned provider requires two compatibility choices:

- Preparation creates `.local/bin/mkisofs` as an alias to xorriso's `xorrisofs`.
- `disk-file.xsl.tftpl` emits an explicit qcow2 file attachment while Terraform
  retains the managed-volume ID. This handles the observed pool-volume access
  failure and the provider's file-format/state behavior without changing host
  security profiles. VNC's listener is disabled; access is SSH and serial console.

### Host changes and isolation

KVM modules and `/dev/kvm` already existed. Terraform, Ansible, QEMU, and libvirt
were initially absent. Two user-run installation transactions added 84 APT
packages, with no upgrades or removals. Exact names and versions are in the
[host package ledger](00-host-packages.md). Terraform is a separate executable
at `.local/bin/terraform`; no HashiCorp APT repository or shell profile was added.
The user added `pior` to `libvirt` and `kvm`; group membership is a separate host
change to consider during any future removal.

The package-created libvirt `default` network uses `192.168.122.0/24` and remained
inactive with autostart off. It was not adopted, modified, or destroyed.
Existing Docker networks use `172.17.0.0/16` through `172.24.0.0/16`; the LAN is
`192.168.2.0/24`. Inspected Docker IPAM and LAN/Tailscale routes did not overlap
the lab subnet. No physical interface was bridged and no SSH host port published.

NAT adds host routing/firewall rules and permits guest outbound LAN access; it
is resource separation, not a security sandbox. The host uses iptables-nft for
Docker/Tailscale; libvirt chose compatible iptables rules. The existing IPv4
FORWARD policy is DROP, so guest egress required actual verification.

AppArmor remained enforcing. No Docker configuration, host security profile,
global security setting, or unrelated resource was changed. The VM and lab network
do not autostart on host boot. The storage pool autostarts by provider default;
after teardown it no longer exists.

### State and local files

Commit the provider lockfile. Keep `.local/`, agent metadata, Terraform state,
plans, downloaded providers, keys, and raw host snapshots out of Git.
Retain local state through teardown and back it up privately while resources
exist. Never delete state as a substitute for destroying resources. Concurrent
Terraform operations are not supported for this single-operator lab.

## Methodology and procedure

Keep image, guest size, network, and playbook fixed when comparing successful
creations. Record significant setup changes and failures. Before and after
lifecycle operations, inspect container status, Docker IPAM, routes, firewall
rules, and representative service health. Container uptime alone does not prove
application health.

### 1. Inspect and install prerequisites

The retained installer records package inventories, APT selections, transaction
logs, and version differences before/after each attempt:

```sh
cd /opt/distributed-systems-lab
scripts/install-host-packages
```

Run as the normal user with interactive sudo. Review the displayed transaction;
libvirt installation may start services or create its default network. Avoid
recommended packages and unexpected upgrades/removals. Install Terraform 1.13.5
from its official release archive and verify its SHA256SUMS as described in the
package ledger. `.local/bin/terraform` already exists on this host.

Verify KVM and libvirt access, then check names, storage, network conflicts, and
existing services before provisioning. Read-only host checks include:

```sh
sudo virt-host-validate qemu
virsh --connect qemu:///system list --all
virsh --connect qemu:///system net-list --all
virsh --connect qemu:///system pool-list --all
sudo docker ps
sudo docker network inspect $(sudo docker network ls -q)
sudo nft list ruleset
ip -4 route show table all
ip -6 route show table all
ss -lntup
free -h
```

Never delete an existing resource to resolve a name collision. Obtain approval
before any required group, networking, firewall, or destructive change.

### 2. Prepare, plan, and create

From the project directory:

```sh
scripts/lab prepare
scripts/lab validate
scripts/lab plan
scripts/lab provision
scripts/lab configure
```

Preparation reuses the checked original image and SSH key, regenerates the sparse
20 GiB copy, and initializes the provider. SHA-512 checks detect corruption and
unexpected image content; the checksum was fetched over Debian HTTPS, without
independent GPG signature verification.

On empty state, expect **five additions, zero changes, zero deletions**. Review
names, paths, network, resources, and provider URI. Provision shows a fresh plan
and requires Terraform's interactive confirmation. Create NAT only after explicit
approval. Avoid running Terraform as sudo.

After an approved teardown/recreation, the new guest has a new SSH host key.
The destroy helper removes only that guest's saved host-key entry after successful
teardown. If teardown was performed directly through Terraform, remove the old
lab entry before configuring or reconnecting:

```sh
ssh-keygen -f .local/known_hosts -R 192.168.100.10
```

The SSH helper accepts a key on first use and rejects a changed key afterward.

### 3. Verify the guest and repeatability

```sh
scripts/lab status
scripts/lab ssh
```

Inside the guest, inspect `hostname`, `/etc/os-release`, `nproc`, `free -m`,
`lsblk -b`, `cloud-init status --long`, and `sudo systemctl is-active ssh`.
Confirm the 2048 MiB memory allocation in libvirt domain info; guest usable memory
is smaller due to kernel overhead. The disk should be 21,474,836,480 bytes, with
a root partition/filesystem slightly smaller. Compare libvirt volume allocation
with capacity to confirm sparse storage.

Run `scripts/lab configure` a second time and record changed/failed counts.
Verify guest DNS and package access, then run `scripts/lab plan` again. The final
working configuration should report no changes. If egress fails, inspect rules
and routes before proposing a narrowly scoped fix. Never flush rules or restart
Docker as a workaround.

### 4. Stop and restart

```sh
scripts/lab stop
scripts/lab start
scripts/lab ssh
```

Stop requests graceful shutdown and checks for `shut off`, without forced
power-off on timeout. Disks and network persist. Start uses the existing lab
network; restarting an inactive network also recreates its host rules and needs
appropriate networking approval. Terraform normally starts a stopped VM on apply.

### 5. Destroy and recreate

`AGENTS.md` rule 6 says: "Ask before destructive commands, host networking/firewall
changes, or disruptive tests." The user also requires explicit approval for
all destructive commands. Obtain that approval before teardown:

```sh
scripts/lab stop
scripts/lab destroy
```

Review exactly five lab deletions before confirming. Verify absence of the domain,
network, pool, disk files, pool directory, and lab bridge. Retain the repository,
SSH key, caches, package records, and Terraform state. Never use broad host
cleanup, Docker prune, manual volume deletion, or state deletion for teardown.

After approval for renewed NAT creation, repeat preparation, plan, provision,
known-host reset, configuration, and guest/service verification. A successful
initial creation alone does not demonstrate recreation.

## Observations

### Significant failures and corrections

| Observation | Correction / lesson |
| --- | --- |
| Initial APT metadata referenced an unavailable package version | Refresh APT and record actual installed versions |
| Installer set `LC_ALL=C`, which Ansible rejected | Use `C.UTF-8`; the installation itself had succeeded |
| xorriso was installed but provider invoked missing `mkisofs` | Create the project-local alias to `xorrisofs` |
| Provider defaulted to SPICE, absent from minimal QEMU | Select built-in VNC with its listener disabled |
| Pool-volume attachment failed with QEMU permission denied | File-path attachment booted without security changes; kernel logs showed helper capability denials but no image-file denial |
| File-attached qcow2 was generated as raw during clean recreation | Explicit qcow2 XML transform; xsltproc was added and recorded |
| Removing the backing file alone did not fix pool-volume access | Keep standalone storage for simplicity, and use the file-attachment transform |
| Provider refresh normalized file paths to volume IDs | Declare the volume ID and transform generated XML; fresh plans then settled |
| Failed domain creation left a defined VM outside Terraform state | Import that domain to track partial resources; final clean creation needed no imports |
| Temporary diagnostics assumed host `rg` existed | Use installed `grep`; remove the temporary helpers at wrap-up |
| Start helper used `grep -q` under pipefail | Consume full output to avoid an upstream SIGPIPE |

The initial six-resource overlay design was replaced with the final five-resource
standalone design. These are setup corrections, not performance comparisons.
An unsuccessful recreation was gracefully interrupted after 190.424 seconds
while waiting for a DHCP lease; this is a failure duration, not a provisioning
result. The nonbooted attempt was removed within the approved lab teardown scope.

### Measured lifecycle results

| Check | Initial repaired VM | Final clean recreation |
| --- | --- | --- |
| Guest OS/resources | Debian 13.7, 1 CPU, 2048 MiB, 20 GiB | Same; about 1974 MiB guest usable RAM |
| SSH / cloud-init | Passed | Passed, no cloud-init errors |
| DNS / Ansible package access | Passed | Passed |
| Ansible repeated runs | Both: ok=7, changed=0, failed=0 | Both: ok=7, changed=0, failed=0 |
| Clean creation from empty state | Initial layout needed repair | Five additions; 17.811 seconds |
| Graceful shutdown | 4.120 seconds | 2.107 seconds |
| Start / SSH readiness | 1.357 seconds; SSH port 10.0 seconds later | 1.322 seconds; SSH port 10.0 seconds later |
| Boot identity | Changed after stop/start | Changed after stop/start |
| Terraform after creation and stop/start | Partial import initially required replacement | No changes |
| Thin provisioning | Overlay allocation measured | About 958 MiB allocated for 20 GiB capacity at inspection |
| Teardown before final recreation | Six-resource teardown: 1.049 seconds | Five-resource teardown: 0.992 seconds; no domains/pools remained |
| Final wrap-up teardown | Not applicable | Graceful stop: 4.144 seconds; five resources destroyed: 0.974 seconds |

Timings measure individual commands on this host; creation excludes image
preparation and subsequent Ansible runs. SSH readiness was sampled every two
seconds. These are observations, not throughput benchmarks or guarantees.

### Existing service and firewall verification

Before/after privileged snapshots showed all 19 container names, reported
uptime/health status, and ports identical. All seven containers reporting health
remained healthy. Docker network names, subnets, and gateways were identical.

After ignoring traffic counters and added libvirt hook jumps, original firewall
chain rules, order, and policies were unchanged. No original rules were removed.
Added filter/NAT/checksum chains belonged to libvirt; guest-specific rules
referenced the lab subnet/bridge. Docker and Tailscale rule bodies were preserved.
All original IPv4 routes remained; additions were only the lab subnet and its
host-local/broadcast entries.

Host-origin Homepage root, Jellyfin `/health`, and Immich `/api/server/ping` HTTP
checks returned 200 before and after lifecycle operations. They also returned 200
after the final teardown. Docker remained active. Required QEMU host checks
passed; unknown Secure Guest support is outside this experiment's requirements.
These checks do not cover every application function or remote-client path.

### Final teardown and retained files

At wrap-up, libvirt showed no domains or storage pools and only its inactive,
non-autostart default network. Terraform state contained no resources. The
`dsl-exp00` storage directory, disk files, and `virbr-dsl00` bridge were absent.
No lab infrastructure remains running.

Host packages and user groups remain installed/configured. `.local/` retains
Terraform, the ISO-tool alias, original/prepared images, SSH keys/known-hosts,
installation inventories, plans, and raw evidence. These are intentionally
excluded from Git. Package removal was not requested or performed.
The destroyed guest's current known-host entry was removed for easy recreation;
the dedicated SSH key pair remains.

Temporary host snapshot and VM diagnostic scripts were removed after use. The
retained scripts are `scripts/lab` and `scripts/install-host-packages`.

Local evidence includes:

- `.local/install-records/20261014T004514Z-2830965/` and `20261014T011429Z-2855526/`.
- `.local/preflight/20261014T005539Z-2841027/` and `20261014T012022Z-2861344/`.
- `.local/experiment-00/host-comparison.json`, lifecycle/Ansible logs, guest/domain
  observations, HTTP snapshots, and `wrapup-verification.json`.

## Interpretation

The measured lifecycle supports the hypothesis for this host and final
configuration: we created, accessed, configured, stopped, destroyed, and recreated
the VM from code. Clean creation needed no manual domain changes or imports;
Terraform settled with no changes and repeated Ansible runs reported zero changes.
Observed container/network/firewall and HTTP checks support coexistence without
disruption under the tested conditions.

Static validation did not catch the live provider/tool/storage requirements.
Controlled changes and actual boot attempts identified a working configuration.
AppArmor remained enforcing throughout; disabling host security was unnecessary.
The experiment is complete and intentionally ends with the VM destroyed.

## Trade-offs

- Standalone qcow2 uses more physical space than an overlay but remains sparse
  and removes a backing-file dependency for a single disposable VM.
- NAT keeps lab addressing separate but shares host firewall machinery and
  permits outbound LAN access; coexistence must be measured.
- Pinned image/provider reduce drift. Provider upgrades, including the incompatible
  0.9 rewrite, require deliberate review and new verification.
- The small XML transform adds a tool dependency but preserves host security and
  stable Terraform state with this provider version.
- Local state is simple for one operator; concurrent runs and state recovery need
  operator care. Raw evidence is local; measured summaries and package versions
  are version-controlled.
- CPU/RAM allocation does not reserve a host CPU or isolate host I/O. No
  performance claim is made in this infrastructure experiment.
- Accepting the first SSH host key is convenient on this private host-to-guest
  network; it is a trust-on-first-use choice.

## References

- [Pinned libvirt provider source and documentation](https://github.com/dmacvicar/terraform-provider-libvirt/tree/v0.8.3)
- [Provider 0.9 breaking change notice](https://github.com/dmacvicar/terraform-provider-libvirt/discussions/1194)
- [Debian image build and checksums](https://cloud.debian.org/images/cloud/trixie/20261001-2618/)
- [Terraform 1.13.5 official archive](https://releases.hashicorp.com/terraform/1.13.5/)
- [Libvirt networking/firewall behavior](https://libvirt.org/firewall.html)
- [Libvirt report of volume/file attachment differences](https://lists.libvirt.org/archives/list/users@lists.libvirt.org/message/NRHWRC7TZS2VQDPQCLJRF5WM7YFUOV44/)
