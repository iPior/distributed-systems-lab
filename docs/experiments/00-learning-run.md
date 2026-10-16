# Experiment 00 — Independent learning run

**Status:** Lesson 01 checkpoint complete. Lesson 02 paused at the learner's request. These are not results from the original agent-built run.

Use [the experiment template](../EXPERIMENT_TEMPLATE.md) and the [Module 00 lessons](../../curriculum/00-infrastructure/README.md).

## Resume here

- Current lesson: Lesson 02, manual VM lifecycle. Its checkpoint is not complete.
- Next exercise: inspect `qemu-img --help` and draft a `convert` command, using `-f` and `-O`, to copy `.local/debian13.qcow2` to `/var/lib/libvirt/images/dsl-lesson02.qcow2` as qcow2. Review the draft before execution.
- Before writing the working disk, verify that the destination is unused and review permissions. Disk creation and resizing have not been performed. The intended working disk capacity is 20 GiB.
- Then finish guest access and boot configuration, review the full `virt-install` command, and obtain explicit approval before VM creation or any host networking changes. No virtual network should be activated automatically.
- Verify the domain definition, CPU/RAM, disk association, and shutdown behavior. Finish the lesson by explaining removal of both the domain and its disk; get explicit approval before destructive commands.
- No VM has been provisioned in this learning run. The learner installed `virt-install` themselves; the instructor has only inspected the environment and updated this log.

Latest learner-written command draft (not executed and not ready to run):

```sh
virt-install -n dsl-lesson02 --vcpus 1 --memory 2048 --connect qemu:///system --disk path=/var/lib/libvirt/images/dsl-lesson02.qcow2,size=20,format=qcow2
```

This draft originally described a blank disk. The verified cloud-image route requires preparing a working copy, selecting `--import`, and completing access and boot configuration before execution.

## Research question
Can I independently explain, provision, configure, stop, destroy, and recreate a disposable Debian VM without disrupting existing Optiplex services?

## Hypothesis
_To write before executing the experiment._

## Environment and controls
_To record._

## Procedure
### Lesson 01 - Virtualization and isolation

- Read the repository instructions, module overview, and current lesson.
- Learner predictions recorded; learner completed read-only host inspection and supplied output.
- Checkpoint complete: learner revised the virtualization sketch to include host Linux and the guest kernel. Instructor clarified that libvirt manages QEMU/KVM from the side rather than forming an execution layer.
- Learner explained that a guest crash can leave the other two replicas available, and that an Optiplex failure takes down all three VMs.

#### Learner predictions before host inspection

1. All three VMs would go down if the Optiplex lost power because they depend on its hardware.
2. A VM can run a different kernel; the initial explanation contrasted this with containers sharing host hardware. Instructor clarification: both share physical hardware, but containers share the host kernel while VMs run their own guest kernels.
3. `/dev/kvm` would indicate KVM is installed. Instructor clarification: its presence normally indicates the kernel exposes the KVM interface; it does not prove a VM is running or that the current user can access it.
4. An empty VM list would mean no VMs are running. Instructor clarification: with `--all`, an empty list means no domains are currently known to that libvirt connection, including stopped domains.

### Lesson 02 - Manual VM lifecycle

- Read the current lesson; introduced domains, storage volumes, and virtual networks as separate objects.
- Learner correctly predicted that stopped domains remain visible with `--all` and that removing a domain definition does not automatically delete its disk. Instructor clarification: definitions and storage have separate lifecycles; explicit storage-removal options can alter this behavior.
- Learner supplied read-only storage pool and network inventory and correctly identified the basic roles of `shutdown`, `destroy`, and `undefine`.
- Learner predicted that `destroy` deletes the disk. Instructor correction: plain `virsh destroy` abruptly stops execution, retaining disks and, for a persistent domain, its saved definition.
- Learner subsequently confirmed that a destroyed persistent VM remains listed with `--all` and can be restarted. Instructor clarified that its retained definition explains the listing, while its retained disk preserves the installed OS.
- Learner initially attempted `virt-install --version` and `virt-install --help`; both returned `bash: virt-install: command not found`. Package inspection preceded the learner's subsequent installation.
- Learner predicted that installing the tool would not create a VM. Instructor clarification: `virt-install` and `virsh` both use libvirt; `virt-install` is not simply a wrapper around `virsh`.
- Learner supplied package policy and an installation simulation. Instructor proposed `sudo apt-get install virtinst` and requested approval under the repository's explicit-approval rule. The learner then reported completing installation themselves.
- Learner subsequently reported installing the tool themselves and reported `virt-install` version 5.0.0. The exact installation command and transaction output were not supplied.
- Learner identified `-n`, `--vcpus`, `--memory`, and `--disk` in the help output; requested help locating installation source and network options. Instructor inspected local `virt-install --help` and explained `--cdrom`, `--location`, `--import`, and `--network`.
- Learner selected `--cdrom` for a blank disk plus installer ISO and identified that `--import` requires an existing OS image. Instructor clarified that the image must contain a bootable installed OS.
- Learner drafted (without executing) `virt-install -n dsl-lesson02 --vcpus 1 --memory 2048 --connect qemu:///system`, correcting an initial spelling error in `--memory`.
- Disk argument exercise: draft a 20 GiB qcow2 disk using the suboptions exposed by `virt-install --disk help`. Candidate path: `/var/lib/libvirt/images/dsl-lesson02.qcow2`; existence and creation permissions have not yet been verified.
- Learner corrected the disk draft to `--disk path=/var/lib/libvirt/images/dsl-lesson02.qcow2,size=20,format=qcow2`. Full draft has not been executed. Instructor clarified that qcow2 allocation grows with use but includes metadata overhead.
- At the learner's request, instructor inspected the repository for downloaded media. No `resources/` directory exists; the search outside `reference/` found two qcow2 files in `.local/` and no ISO. Read-only `qemu-img info` confirmed their formats and capacities, not OS contents, provenance, or bootability. No images were modified or booted.
- Learner clarified that they meant `reference/`. Instructor inspected its file inventory and image-related documentation. The historical report identifies Debian 13 genericcloud amd64 build `20261001-2618`; the archived Terraform points to a `.local/debian13-root.qcow2` path relative to its original location. No ISO or qcow2 binary was found inside `reference/`. The `.local/` files' provenance remains unverified; no archived scripts were executed or implementation copied. Installation method selection remains open: installer ISO versus a verified clean cloud image.
- At the learner's request, instructor verified `.local/debian13.qcow2` against Debian's published SHA-512 checksum fetched over HTTPS. The complete file hash matches the Debian 13 genericcloud amd64 build `20261001-2618` entry. This establishes that the local file matches the published image rather than a modified guest disk. No guest boot was tested. The existing `.local/debian13-root.qcow2` remains unverified and untouched. Next: learner prediction about preserving the base image before preparing a separate 20 GiB working disk.
- Learner learned that a SHA-512 checksum is a fingerprint of file contents and correctly predicted that guest writes change the working disk's checksum. Instructor explained preserving the verified base image by booting a separate writable copy.
- Remaining guest design, full command review, creation, lifecycle verification, and checkpoint are pending. No VM creation or host changes have been performed by the instructor. Session paused before the learner attempted the `qemu-img convert` draft.

## Observations
### Lesson 01 - Learner-supplied host inspection

The learner ran these commands on the Optiplex and supplied the output; the instructor did not independently execute the checks.

- `lscpu | grep -E 'Model name|Virtualization|CPU\(s\)'`: Intel Core i5-6500 CPU @ 3.20GHz, 4 CPUs online (0-3), virtualization reported as VT-x.
- `ls -l /dev/kvm`: `crw-rw---- 1 root kvm 10, 232 Oct 13 21:23 /dev/kvm`.
- `virsh --connect qemu:///system list --all`: successful response with an empty domain table, repeated with the same result.

These checks show CPU virtualization support, the KVM device interface, and a working system libvirt connection for listing domains. They do not demonstrate VM creation, guest startup, or the learner's direct access to `/dev/kvm`. No VM lifecycle tests have been performed in this learning run.

### Lesson 02 - Learner-supplied libvirt inventory

- `virsh --connect qemu:///system pool-list --all`: empty table, showing no pools defined on this connection. This does not establish that no disk images exist on the host.
- `virsh --connect qemu:///system net-list --all`: `default` is inactive, has autostart disabled, and is persistent.
- No pool creation or network activation has been performed in this lesson. The inactive network's topology has not been inspected.
- `virt-install --version` and `virt-install --help`: the learner's shell reported `command not found` for both. This establishes command unavailability in that shell, not a complete package inventory.
- `apt-cache policy virtinst`: not installed; candidate `1:5.0.0-5+deb13u1` from Debian trixie/main.
- `apt-get --simulate install virtinst`: 33 newly installed packages (including both `virtinst` and `virt-install`), 0 upgrades, 0 removals, and 62 not upgraded. This is a simulation using the host's current package metadata, not evidence that installation occurred.
- Learner supplied `virt-install --disk help`, which lists `path`, `size`, `format`, and `sparse` among its suboptions.
- Learner supplied `ls -ld /var/lib/libvirt /var/lib/libvirt/images`: both directories exist and are owned by `root:root`; modes are `drwxr-xr-x` and `drwx--x--x`, respectively. This does not establish that the proposed disk path is unused.
- Instructor inspected `.local/debian13.qcow2`: qcow2, 3 GiB virtual capacity, approximately 326 MiB allocated. `.local/debian13-root.qcow2`: qcow2, 20 GiB virtual capacity, approximately 955 MiB allocated. Neither reported a backing file. These metadata checks do not establish a clean reusable guest image.
- Image checksum verification: `sha512sum .local/debian13.qcow2` returned `f46f0671a6e5bdec5291ab8972bae2f10e5408c2f64a74078f11efc2f06a436a9d0313ed50e0472542eeabf780e9f7c792ac0a314c6c20507fcd9fd81b468c3d`. This matches the `debian-13-genericcloud-amd64-20261001-2618.qcow2` entry in [Debian's published SHA512SUMS](https://cloud.debian.org/images/cloud/trixie/20261001-2618/SHA512SUMS), retrieved directly with Python's HTTPS client after the browser tool could not access the URL. It also matches the archived script's pinned checksum; no archived script was executed.

## Analysis and limitations

### Lesson 01 reflection

The learner identified the shared physical failure domain: "because if the optiplex fails, all three VMs are going down as well". Separate VMs provide guest OS isolation but share the physical host and its power. Continued request handling after one guest fails also depends on traffic reaching the surviving replicas; this has not been tested.

VM lifecycle validation and the remaining experiment analysis are pending.

## Conclusion
_To complete after testing._
