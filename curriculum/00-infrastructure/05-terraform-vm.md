# Lesson 05 — Provision a VM using Terraform

## Learn
A libvirt guest depends on a storage pool, guest disk, cloud-init seed, network, and domain definition. **Cloud-init** bootstraps guest identity, SSH access, and initial network setup; it does not replace ongoing configuration management.

## Predict
List resource dependencies before writing HCL. What happens if a disk pool is destroyed while unrelated data is placed inside it? Why should the virtual disk and Terraform state be kept separate from the media stack?

## Build
Implement one resource at a time in `infrastructure/terraform/`: dedicated pool, isolated network, Debian cloud image/disk, cloud-init, then VM. Inspect each plan and explain resource counts. Never apply network changes without approval. If libvirt requires special disk XML handling, understand the reference workaround before reusing it.

## Verify
After review and approval, boot the guest, confirm disk capacity/CPU/RAM/SSH, then test a no-change plan.

## Checkpoint
Explain how to safely destroy only lab-owned resources. Stop before Ansible.
