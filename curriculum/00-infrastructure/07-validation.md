# Lesson 07 — Lifecycle, safety, and evidence

## Learn
A reproducible deployment must prove more than a successful first boot. Lifecycle testing includes stop/start, no-change plans, teardown, clean recreation, and verification that unrelated services still work.

## Predict
Which resources must remain after stopping the VM? Which should remain after destroy? How could a network change harm an existing container despite correctly scoped Terraform resources?

## Build and verify
Create a repeatable, narrowly scoped lifecycle checklist. After explicit approval, stop, start, destroy and recreate *only* the lab resources. Verify Terraform state and libvirt agree, and check media-stack services before and after. Record actual commands/timings/outputs, including failures.

## Document
Write `docs/experiments/00-learning-run.md` using the [experiment template](../..//docs/EXPERIMENT_TEMPLATE.md). Distinguish new observations from the original agent-run reference. Evaluate reproducibility, not VM throughput.

## Checkpoint
Can you explain and reproduce the lifecycle without copying the reference scripts? Review before progressing to the Go API stage.
