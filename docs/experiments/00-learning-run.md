# Experiment 00 — Independent learning run

**Status:** Lessons 01 and 02 checkpoints complete. Lesson 03 not started. These are not results from the original agent-built run.

Use [the experiment template](../EXPERIMENT_TEMPLATE.md) and the [Module 00 lessons](../../curriculum/00-infrastructure/README.md).

## Resume here

- Lesson 02 checkpoint complete. Stop here; begin Lesson 03 only when the learner requests it.
- Guest access verified: learner logged in as `lab` to `dsl-lesson02`, confirmed Debian 13.7, one CPU, a 20 GiB disk, an expanded root filesystem, and completed cloud-init with no reported errors.
- Learner correctly explained that shutdown retains the VM definition and disk, while plain undefine removes the definition but retains the disk. Instructor clarified that destroy is an abrupt power-off and disk deletion is a separate explicit action. No teardown was performed for this checkpoint.
- The VM was running at the last verified boot and has no network interface. Its definition includes the VGA device that allowed boot to progress. Guest networking is deferred to Lesson 03.
- The working disk is `/var/lib/libvirt/images/dsl-lesson02.qcow2`; the login-repair CD is `/var/lib/libvirt/images/dsl-lesson02-seed.iso`, attached at `sda`. Keep both until a deliberate lifecycle decision. Graceful shutdown was already verified before the login-repair boot.
- Source metadata and the credential-free user-data template are in `infrastructure/lesson02/cloud-init/`. Rendered user-data, seed ISO, and the `lab` account password are in Git-ignored `.local/lesson02-cloud-init/`. Do not commit credentials.
- The learner installed the tool and performed VM creation, configuration edits, media attachment, and boots. At explicit request, the instructor prepared cloud-init files and the seed ISO. The instructor has not changed host networking or started/stopped the VM.

Original reviewed creation command (historical; do not rerun against the existing VM):

```sh
sudo virt-install --name dsl-lesson02 --vcpus 1 --memory 2048 --connect qemu:///system --virt-type kvm --osinfo debian13 --disk path=/var/lib/libvirt/images/dsl-lesson02.qcow2,format=qcow2 --network none --graphics none --console pty,target_type=serial --autoconsole text --import --cloud-init root-password-generate=on,disable=on
```

The learner executed this command. The current saved definition additionally contains VGA hardware and a replacement login-repair CD. The original generated root login did not work; the repair created the `lab` account used for successful console access.

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
- Session resumed on 2026-10-16. Learner drafted `qemu-img convert -f qcow2 -O qcow2 .local/debian13.qcow2 /var/lib/libvirt/images/dsl-lesson02.qcow`; instructor corrected the destination suffix to `.qcow2`.
- Instructor requested the read-only check `sudo ls -l /var/lib/libvirt/images/dsl-lesson02.qcow2`. Learner reported that no such file exists, without supplying raw output. Working disk creation remains pending; this check was not independently repeated by the instructor.
- Learner predicted that conversion would preserve the image size but give it a 20 GB maximum capacity. Instructor clarified that conversion alone retains the source's 3 GiB virtual capacity because the command specifies no new size.
- Learner executed `sudo qemu-img convert -f qcow2 -O qcow2 .local/debian13.qcow2 /var/lib/libvirt/images/dsl-lesson02.qcow2` and supplied successful inspection output. No VM is using the newly created working disk in this learning run. Resize command drafting is next.
- Learner correctly drafted `sudo qemu-img resize /var/lib/libvirt/images/dsl-lesson02.qcow2 20G`. Learner predicted that disk capacity and filesystem size should be the same; instructor explained the separate virtual disk, partition, and filesystem layers. Guest partition/filesystem expansion remains unverified.
- After review, learner executed the resize command, received `Image resized.`, and supplied metadata confirming 20 GiB virtual capacity. Instructor proposed a console-only guest with `--network none` for Lesson 02, leaving host networking unchanged and guest networking for Lesson 03. No VM has been created yet.
- Learner revised the VM draft to import the existing disk, omit `size=20`, and use `--network none --graphics none`. Instructor clarified that `--connect` selects the libvirt connection while `--import` takes no value; learner corrected both this and a trailing disk suboption comma. Latest draft is saved above and remains unexecuted.
- Instructor inspected local `virt-install --cloud-init help` and its installed manual. For this disposable offline guest, the documented `root-password-generate=on` option generates a guest root password and displays it locally; `disable=on` disables cloud-init on subsequent boots. Guest login configuration is the next exercise; neither option has been executed.
- At the learner's explicit request, instructor appended `--cloud-init root-password-generate=on,disable=on` to the saved VM command draft. This is a documentation edit only; final review and explicit approval for VM creation are still pending.
- At the learner's request for final review, instructor performed read-only preflight checks. Added explicit `--virt-type kvm`, `--osinfo debian13`, `--console pty,target_type=serial`, and `--autoconsole text` to the reviewed draft, with `sudo` for system-managed resource creation. No command creating a guest or cloud-init ISO was executed. Approval for VM creation/start remains pending.
- Learner said they would execute the reviewed command themselves, then reported a console remaining at `Booting 'Debian GNU/Linux'`. Instructor performed read-only troubleshooting: libvirt and QEMU report the persistent KVM domain running with 1 vCPU, 2 GiB RAM, the intended virtio working disk, a serial console at `/dev/pts/2`, and no network interface. These host-side states do not establish a healthy guest boot.
- Disk statistics during troubleshooting report reads but zero writes to `vda`. The running definition has no medium attached to its CD-ROM. The reason for the apparent boot stall and the cloud-init media's lifecycle are not established. QEMU logs are unreadable to the instructor's user; `sudo -n` also requires a password. No restart, disk mutation, or forced console takeover was attempted.
- Learner supplied the QEMU log tail: launch arguments match the expected disk, KVM, CPU/RAM, serial console, and no graphics; the excerpt contains no explicit startup error. Its `custom-monitor` taint was caused by the instructor's diagnostic monitor queries, not established as a guest failure. Repeated read-only CPU snapshots show 32-bit execution with paging disabled and PCI configuration reads; the guest is not observed running normal 64-bit userspace. This suggests early boot activity but does not establish its cause. Exact executed creation command and any intervening restarts need clarification before a controlled boot test.
- Learner recovered shell history entry 787, confirming execution of the reviewed `sudo virt-install` command, including cloud-init. Learner reported that the boot message appeared once and remained unchanged.
- Instructor monitored QEMU RESET events read-only with `virsh --connect qemu:///system qemu-monitor-event dsl-lesson02 --event RESET --loop --timeout 12 --timestamp`. The listener received 32 events in 12 seconds, in closely spaced pairs recurring approximately every 0.75 seconds; each reported `guest: true` and `reason: guest-reset`. This confirms repeated resets, not a merely slow boot. Root cause remains unknown. No stop, restart, configuration change, or disk write was performed by the instructor.
- Learner reported running `virsh --connect qemu:///system destroy dsl-lesson02`. Instructor independently verified `shut off` and inspected the retained inactive definition, including the working disk path. No domain or disk deletion occurred.
- Troubleshooting hypothesis: the missing virtual video device may interact with the image's GRUB graphical/serial terminal configuration. [Debian bug 1111240](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=1111240;msg=22) reports a similar Trixie cloud-image boot failure and a workaround changing GRUB terminal output to serial-only. Historical working local XML also contains a video device, but differs in other settings, so this is not proof of cause. Proposed controlled test adds only a virtual VGA device, retaining serial access and no graphics listener; it has not been applied or tested.
- Learner edited the stopped domain with `virsh edit`, adding `<video><model type='vga'/></video>` inside `<devices>`, and reported successful validation. After the proposed `virsh --connect qemu:///system start dsl-lesson02 --console` test, learner reported reaching the login prompt. Instructor confirmed the running definition includes VGA, no graphics listener, and no network interface. This supports missing video hardware as a contributor to the previous boot failure; it does not establish the underlying GRUB defect definitively.
- Instructor observed guest disk writes after the VGA change (452 writes totaling 140,521,984 bytes at inspection), unlike the previous zero-write reset loop. A five-second read-only RESET listener received zero events. This establishes no observed resets during that short window, not long-term stability or successful authentication.
- Learner reported `login incorrect` when attempting guest login. Missing cloud-init media makes unapplied initial credentials a plausible explanation, but guest logs have not been inspected to confirm this. Proposed next step is a graceful shutdown before preparing replacement first-boot login media; no password or password hash has been collected by the instructor.
- Learner ran `virsh --connect qemu:///system shutdown dsl-lesson02` and received the shutdown-request acknowledgement, then ran `domstate` and supplied `shut off`. This verifies successful graceful shutdown for this boot without guest authentication; no shutdown duration was measured. Replacement cloud-init metadata drafting is next.
- At the learner's explicit request, instructor created `infrastructure/lesson02/cloud-init/meta-data` with `instance-id: dsl-lesson02-login-repair-01` and `local-hostname: dsl-lesson02`. Only the identification file was created; user-data, seed ISO creation, media attachment, and guest restart remain pending. No credentials were added, no VM changes were made, and no commit was requested for this step.
- At the learner's explicit request to also prepare user-data, instructor created `infrastructure/lesson02/cloud-init/user-data.template` defining guest account `lab`, Bash shell, `sudo` group membership, and unlocked password authentication via a hash placeholder. This changes the planned console login from root to lab; it has not yet configured any account inside the guest.
- Instructor generated a random password and SHA-512 crypt password hash locally, rendered `.local/lesson02-cloud-init/user-data`, and saved the password in `.local/lesson02-cloud-init/console-password`. The directory has mode 0700; both files have mode 0600. Neither password nor hash was printed. Confirmed `.local/` is ignored by Git and parsed the rendered YAML to verify its structure and substitution. Full guest cloud-init validation and authentication remain untested. No ISO, VM change, boot, or commit was performed.
- Learner requested skipping the extra ISO-building exercise. Instructor prepared `.local/lesson02-cloud-init/seed.iso` with `xorrisofs`, using only the metadata and rendered user-data files. Read-only inspection verified volume ID `cidata` and exactly the two root files `meta-data` and `user-data`; the plaintext password file is not included. The lowercase-volume warning is expected for this NoCloud label. ISO is Git-ignored. Instructor confirmed VM `shut off` and empty CD-ROM target `sda`. No VM configuration change or boot was performed by the instructor; learner attachment and login verification remain pending.
- After the learner's next boot, console output stopped visibly at `systemd-networkd-wait-online`. Instructor verified read-only that the domain is running, the seed ISO is attached at `sda`, and the guest has read from it. Guest login/account creation remain unverified. Network readiness timeout is a plausible explanation for the delay in this deliberately NIC-less guest; the host's installed wait-online manual documents a default of 120 seconds, but the guest's unit configuration has not been inspected. No boot duration was measured and no service settings were changed.

## Observations
### Lesson 01 - Learner-supplied host inspection

The learner ran these commands on the Optiplex and supplied the output; the instructor did not independently execute the checks.

- `lscpu | grep -E 'Model name|Virtualization|CPU\(s\)'`: Intel Core i5-6500 CPU @ 3.20GHz, 4 CPUs online (0-3), virtualization reported as VT-x.
- `ls -l /dev/kvm`: `crw-rw---- 1 root kvm 10, 232 Oct 13 21:23 /dev/kvm`.
- `virsh --connect qemu:///system list --all`: successful response with an empty domain table, repeated with the same result.

These Lesson 01 checks show CPU virtualization support, the KVM device interface, and a working system libvirt connection for listing domains. VM creation and lifecycle checks were performed later in Lesson 02, as recorded below.

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
- Learner-supplied working disk inspection after conversion: `/var/lib/libvirt/images/dsl-lesson02.qcow2`, format qcow2, virtual size 3 GiB (3,221,225,472 bytes), disk size 954 MiB, file length 1,000,800,256 bytes, cluster size 65,536 bytes, `corrupt: false`, and no backing file reported. The conversion command returned without an error. This verifies image creation and metadata, not guest bootability.
- Learner-supplied working disk inspection after resize: virtual size 20 GiB (21,474,836,480 bytes), disk size 955 MiB, file length 1,000,800,768 bytes, format qcow2, `corrupt: false`, and no backing file reported. Guest partitions and filesystems have not been inspected or resized explicitly.
- Instructor read-only preflight on 2026-10-16: system libvirt domain and pool lists empty; `default` network still inactive, persistent, and not autostarted. Working disk exists, owned by `root:root` with mode `0644`, virtual capacity 20 GiB, disk size 955 MiB, and `corrupt: false`. Host reports about 413 GiB available disk space and 26 GiB available RAM. `virt-install --osinfo list` includes `debian13`; QEMU and `xorrisofs` executables are available. These checks do not prove guest startup, console login, or filesystem expansion.

### Lesson 02 - Successful guest access and resource checks

Learner supplied authenticated console output from `lab@dsl-lesson02`:

| Check | Observed output |
| --- | --- |
| `hostname` | `dsl-lesson02` |
| `/etc/os-release` | Debian GNU/Linux 13 (trixie), `DEBIAN_VERSION_FULL=13.7` |
| Login banner kernel | `6.12.111+deb13-cloud-amd64`, x86_64 |
| `nproc` | `1` |
| `lsblk` | `vda` 20G, root partition `vda1` 19.9G mounted at `/`, EFI partition `vda15` 124M, `vda14` 3M, CD-ROM `sr0` 370K |
| `df -h /` | `/dev/vda1` size 20G, used 595M, available 19G, 4% used (human-readable rounded values) |
| `cloud-init status --long` | `status: done`, `extended_status: done`, `DataSourceNoCloud [seed=/dev/sr0]`, `errors: []`, `recoverable_errors: {}` |

The reported cloud-init `last_update` was `Thu, 01 Jan 1970 00:02:08 +0000`; it is not used as a measured boot duration or session timestamp.

## Analysis and limitations

### Lesson 01 reflection

The learner identified the shared physical failure domain: "because if the optiplex fails, all three VMs are going down as well". Separate VMs provide guest OS isolation but share the physical host and its power. Continued request handling after one guest fails also depends on traffic reaching the surviving replicas; this has not been tested.

### Lesson 02 findings so far

- Adding VGA while retaining the serial console allowed the previously resetting guest to reach login. This supports the video/GRUB compatibility hypothesis, without proving the precise underlying defect.
- Attaching the replacement NoCloud seed enabled successful `lab` console login. The original failed password setup was not diagnosed through guest logs, so its exact failure mechanism remains uncertain.
- The working disk's partition and root filesystem expanded to use its larger virtual capacity during guest initialization. Merely resizing the image was not itself evidence of filesystem growth; the later guest outputs supply that evidence.
- A graceful guest shutdown was observed and a later boot reached authenticated userspace. Networking, full teardown/recreation, and the remaining experiment conclusions are still pending. No service-health comparison was performed during this lesson, so unrelated host-service health is not claimed as measured.

### Lesson 02 checkpoint reflection

The learner initially confused shutdown and deletion. After reviewing the command roles, the learner correctly stated: "the vm definition and the disk, after undefine the disk remains but the definition doesnt". This completes the conceptual checkpoint: stopping execution, removing the persistent definition, and deleting storage are separate operations. Actual disk deletion and full teardown/recreation remain untested in this learning run. No VM changes or commits accompanied checkpoint completion.

## Conclusion
_To complete after testing._
