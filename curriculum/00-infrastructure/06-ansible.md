# Lesson 06 — Guest configuration with Ansible

## Learn
Cloud-init initializes a new guest; Ansible describes and converges ongoing configuration over SSH. Ansible **idempotency** means reapplying desired configuration should produce no unnecessary changes.

## Predict
Which guest changes belong in cloud-init versus Ansible? What should the second unchanged playbook run report?

## Build
Create an inventory and a small playbook under `infrastructure/ansible/` to verify Debian and install basic guest management packages. Read one module's documentation before using it. Do not automate host package changes.

## Verify
Run syntax checks, then (when the VM is approved and reachable) run the playbook twice. Record changed and failed counts, and inspect any drift.

## Checkpoint
Explain what Ansible did that Terraform did not. Stop before full lifecycle validation.
