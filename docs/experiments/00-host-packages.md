# Experiment 00 host installation ledger

## Status

Installation completed on 2026-10-14 UTC (2026-10-13 local time), run from
the user's host terminal. Across two transactions: **84 newly installed packages, zero upgraded packages,
and zero removed packages**. The first transaction added 82 packages and the
second added two XML-tool packages. Both exact lists and versions are below.
The lab VM, disks, pool, and NAT network were destroyed at wrap-up; host tools
and group memberships remain for recreation.

Run from the host terminal as `pior`:

```sh
cd /opt/distributed-systems-lab
scripts/install-host-packages
```

The script refreshes APT, shows a simulation, and uses APT's confirmation prompt.
It installs without recommended packages and rejects package removals. Libvirt's
required configuration package may create/start the default NAT network. The
experiment's dedicated NAT network and VM are not created by this installer.

Explicitly requested packages:

| Package | Purpose |
| --- | --- |
| `qemu-system-x86` | Run the x86 VM with KVM |
| `qemu-utils` | Inspect qcow2 disk images |
| `libvirt-daemon-system` | Manage system VMs and virtual networks |
| `libvirt-clients` | Supply `virsh` and host validation tooling |
| `ansible-core` | Configure the guest through SSH |
| `xorriso` | Build the cloud-init seed ISO |
| `xsltproc` | Apply the explicit qcow2 file-attachment XML transform |

`curl` and `unzip` were already installed and are not requested by this script.
The initial simulation proposed 82 new packages with no upgrades or removals.
This used stale APT metadata and is not an installation record.

Every attempt that reaches installation writes its own private record directory
under `.local/install-records/`, including:

- Before/after installed package names, versions, and installation status.
- Before/after APT manual and automatic package selections.
- Requested packages, APT simulation, update log, and installation log.
- `summary.md` with actual additions, version changes, and removals, even if APT fails.
- Installed tool versions and post-installation libvirt network definitions.

Do not run another package transaction concurrently, which would make attribution
ambiguous. Keep these records until any eventual removal is complete. The verified package
additions are copied below so they can be retained in Git.

## Terraform

Terraform 1.13.5 is installed separately as a standalone executable at
`/opt/distributed-systems-lab/.local/bin/terraform`. It does not register an APT
package, add an APT repository, change the shell profile, or start a service.
Use it by adding the directory to PATH for the current terminal:

```sh
export PATH="/opt/distributed-systems-lab/.local/bin:$PATH"
terraform version
```

The official Linux amd64 release archive was checked against HashiCorp's HTTPS
SHA256SUMS. Its file checksum and installation details are recorded privately in
`.local/install-records/terraform-1.13.5.md`. HTTPS checksums are not independent
signature verification.

## Eventual removal

Removal is a separate destructive operation requiring explicit approval. First
teardown any approved lab resources, then review the actual installation record
and current reverse dependencies. Start with a simulated removal of the requested
packages and assess dependencies individually. Do not run broad
`autoremove`, purge packages indiscriminately, or assume every dependency added
today will remain unused by other services.

Preparation also creates `.local/bin/mkisofs`, a symlink to the installed
`/usr/bin/xorrisofs`. Provider 0.8.3 invokes that compatibility name. This is a
project-local file, not another installed package or system-wide change.

The Terraform executable and compatibility symlink can be removed separately
after approval. Do not delete
the entire `.local/` directory: it also contains SSH keys, image cache, and these
installation records. Package removal alone does not reverse generated networks,
service changes, group memberships, or VM storage.


## Verified installation: 2026-10-14 UTC

Record directory: `.local/install-records/20261014T004514Z-2830965/`.

APT completed installation. The installer then exited with status 1 during
tool verification because its `LC_ALL=C` setting was rejected by Ansible.
Both the installer and lab runner now use `C.UTF-8`; package installation
does not need to be repeated. Preserve the original exit status in the raw record.

Follow-up checks passed: Terraform schema validation and formatting, the installed
Ansible 2.19.11 playbook syntax check, and shell syntax. `dpkg --audit` produced
no findings. QEMU reports 10.0.13; libvirt reports 11.3.0; xorriso reports 1.5.6.

The following describes the initial installation check before provisioning.
Final lifecycle and service verification is in the experiment document.

Libvirt is socket-activated; `libvirtd.socket` and `virtlogd.socket` are active.
Read-only libvirt inspection shows no VMs or storage pools. The package-created
`default` network is defined for `192.168.122.0/24`, inactive, with autostart off.
No lab network or VM has been created. Docker is active and visible IPv4 routes
match the earlier inspection. Container and application health checks still need
privileged access; Docker being active does not establish their health.

At the initial post-install check, user `pior` was outside `libvirt` and `kvm`;
normal management access failed with a polkit authentication error.
After the user's group change and session restart, both groups list `pior` in
the account database, and normal `virsh --connect qemu:///system` access succeeds.
This account change is separate from APT and must be considered during removal.

### Actual package additions

| Package | Installed version |
| --- | --- |
| ansible-core | 2.19.11-0+deb13u1 |
| dnsmasq-base | 2.91-1+deb13u2 |
| ieee-data | 20240722 |
| ipxe-qemu | 1.21.1+git20250501.dad20602+dfsg-1 |
| libaio1t64:amd64 | 0.3.113-8+b1 |
| libasound2-data | 1.2.14-1+deb13u1 |
| libasound2t64:amd64 | 1.2.14-1+deb13u1 |
| libatomic1:amd64 | 14.2.0-19 |
| libbrlapi0.8:amd64 | 6.7-3.1+deb13u3 |
| libburn4t64:amd64 | 1.5.6-1.1+b1 |
| libcacard0:amd64 | 1:2.8.0-3+b2 |
| libcapstone5:amd64 | 5.0.7-1~deb13u1 |
| libdaxctl1:amd64 | 81-1 |
| libduktape207:amd64 | 2.7.0-2+b2 |
| libexecs1:amd64 | 1.4-2+b2 |
| libfdt1:amd64 | 1.7.2-2+b1 |
| libglib2.0-0t64:amd64 | 2.84.4-3~deb13u5 |
| libibverbs1:amd64 | 56.1-1 |
| libisoburn1t64:amd64 | 1.5.6-1.2+b1 |
| libisofs6t64:amd64 | 1.5.6.pl01-1.1+b1 |
| libjpeg62-turbo:amd64 | 1:2.1.5-4 |
| libjte2:amd64 | 1.22-4+b2 |
| libndctl6:amd64 | 81-1 |
| libnl-3-200:amd64 | 3.7.0-2 |
| libnl-route-3-200:amd64 | 3.7.0-2 |
| libnspr4:amd64 | 2:4.36-1 |
| libnss3:amd64 | 2:3.110-1+deb13u4 |
| libnuma1:amd64 | 2.0.19-1 |
| libparted2t64:amd64 | 3.6-5 |
| libpcap0.8t64:amd64 | 1.10.5-2 |
| libpciaccess0:amd64 | 0.17-3+b3 |
| libpcsclite1:amd64 | 2.3.3-1 |
| libpixman-1-0:amd64 | 0.44.0-3 |
| libpmem1:amd64 | 1.13.1-1.1+b1 |
| libpolkit-agent-1-0:amd64 | 126-2 |
| libpolkit-gobject-1-0:amd64 | 126-2 |
| librdmacm1t64:amd64 | 56.1-1 |
| libslirp0:amd64 | 4.8.0-1+deb13u1 |
| libssh-4:amd64 | 0.11.5-0+deb13u1 |
| liburing2:amd64 | 2.9-1 |
| libusbredirparser1t64:amd64 | 0.15.0-1 |
| libvdeplug2t64:amd64 | 4.0.1-5.1+b1 |
| libvirt-clients | 11.3.0-3+deb13u3 |
| libvirt-common | 11.3.0-3+deb13u3 |
| libvirt-daemon | 11.3.0-3+deb13u3 |
| libvirt-daemon-common | 11.3.0-3+deb13u3 |
| libvirt-daemon-config-network | 11.3.0-3+deb13u3 |
| libvirt-daemon-config-nwfilter | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-network | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-nodedev | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-nwfilter | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-qemu | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-secret | 11.3.0-3+deb13u3 |
| libvirt-daemon-driver-storage | 11.3.0-3+deb13u3 |
| libvirt-daemon-log | 11.3.0-3+deb13u3 |
| libvirt-daemon-system | 11.3.0-3+deb13u3 |
| libvirt0:amd64 | 11.3.0-3+deb13u3 |
| libxml2-utils | 2.12.7+dfsg+really2.9.14-2.1+deb13u3 |
| libyaml-0-2:amd64 | 0.2.5-2 |
| mdevctl | 1.4.0-2+b1 |
| polkitd | 126-2 |
| python3-bcrypt | 4.2.0-2.1+b1 |
| python3-cffi-backend:amd64 | 1.17.1-3 |
| python3-cryptography | 43.0.0-3+deb13u1 |
| python3-dnspython | 2.7.0-1 |
| python3-httplib2 | 0.22.0-1+deb13u1 |
| python3-jinja2 | 3.1.6-1 |
| python3-markupsafe | 2.1.5-1+b3 |
| python3-netaddr | 1.3.0-1 |
| python3-packaging | 25.0-1 |
| python3-pyparsing | 3.1.2-1 |
| python3-resolvelib | 1.1.0-1 |
| python3-yaml | 6.0.2-1+b2 |
| qemu-system-common | 1:10.0.13+ds-0+deb13u1 |
| qemu-system-data | 1:10.0.13+ds-0+deb13u1 |
| qemu-system-x86 | 1:10.0.13+ds-0+deb13u1 |
| qemu-utils | 1:10.0.13+ds-0+deb13u1 |
| seabios | 1.16.3-2 |
| sgml-base | 1.31+nmu1 |
| systemd-container | 257.13-1~deb13u1 |
| xml-core | 0.19 |
| xorriso | 1.5.6-1.2+b1 |


## Verified additional installation: 2026-10-14 UTC

Record: `.local/install-records/20261014T011429Z-2855526/`. Installer exit status
was 0. Exactly two packages were added; none were upgraded or removed:

| Package | Installed version |
| --- | --- |
| `xsltproc` | `1.1.35-1.2+deb13u3` |
| `libxslt1.1:amd64` | `1.1.35-1.2+deb13u3` |

The pre-install simulation matched this transaction. Total additions
across both transactions: 84 APT packages. Terraform remains a separate local
executable. The original and prepared standalone image caches are local files,
not installed packages.
