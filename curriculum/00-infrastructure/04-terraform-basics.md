# Lesson 04 — Terraform fundamentals

## Learn
Terraform is a **desired-state** tool. A provider translates resources to an external API. Configuration describes desired resources, **state** maps configuration to managed objects, `plan` previews changes, `apply` makes them, and `destroy` removes managed resources. A lock file pins provider selections; never commit state files containing possibly sensitive values.

## Predict
If you add one managed resource, what should `plan` display? If you change only a comment, what should a subsequent plan show? Why is deleting state not the same as destroying infrastructure?

## Build
Create only a minimal provider and version configuration in the new `infrastructure/terraform/` directory. Use the reference provider version as a starting point only after understanding why it is pinned. Run `terraform fmt` and `terraform validate` and inspect the results.

## Checkpoint
Explain what `qemu:///system` refers to and what will be created by your current configuration (potentially nothing yet). Stop before adding VM resources.
