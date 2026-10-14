# Lesson 02 — Manual VM lifecycle

## Learn
A virtual machine requires a guest image, virtual CPU/RAM, virtual disk, and access to a network or console. `virsh` manages defined libvirt domains; `virt-install` can create one. A stopped domain can still be defined and retain a disk. A domain definition, storage volume, and virtual network are separate objects.

## Predict
How will you tell the difference between a VM that is stopped and one that has been deleted? Why shouldn't you reuse production VM names or paths?

## Build
Design a **temporary, lab-named guest**: 1 vCPU, 2 GiB RAM, 20 GiB disk. Before creating it, inspect `virsh help`, available storage pools, and networks. Decide whether you can safely reuse the previously installed libvirt tools. Propose the exact `virt-install` command for review rather than running one generated automatically.

## Verify
After approval, inspect the guest's domain definition, CPU/RAM, storage association, and shutdown behavior. Record the actual output.

## Checkpoint
Explain what must be removed to delete the VM *and its disk*. Stop before modifying network state in Lesson 03.
