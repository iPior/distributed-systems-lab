# Lesson 01 — Virtualization and isolation

## Learn
A **process** runs in a host OS. A **container** isolates processes but shares the host kernel. A **VM** has its own guest kernel and virtual hardware. **KVM** accelerates CPU virtualization, **QEMU** runs/emulates the virtual machine, and **libvirt** manages VM definitions, storage, and networks. Terraform will later declare resources through libvirt; Ansible will configure the guest over SSH.

## Predict
Before running anything, write down: which layer would fail if the Optiplex lost power? Which layer allows a Debian guest to have its own kernel? What would `/dev/kvm` indicate?

## Explore (read-only)
```sh
lscpu | grep -E 'Model name|Virtualization|CPU\(s\)'
ls -l /dev/kvm
virsh --connect qemu:///system list --all
```
If `virsh` is missing, record that fact rather than installing packages automatically.

## Checkpoint
Sketch host → KVM/QEMU → libvirt-managed VM → guest. Explain in your own words why three VMs on one physical host are not three independent physical failure domains. **Stop here and review before Lesson 02.**
