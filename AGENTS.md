# AGENTS.md — Distributed Systems Lab

## Purpose
Teach and build a small distributed application through **one concept-driven experiment at a time**. The learner writes and tests each step. Measurements must fit the question; never invent results.

## Read before acting
- `README.md` and `docs/ROADMAP.md`
- `curriculum/00-infrastructure/README.md` for the current module
- **Only the current lesson** in `curriculum/00-infrastructure/`
- `docs/EXPERIMENT_TEMPLATE.md` when documenting results

## Working method
1. Explain one concept, ask the learner to predict behavior, then propose a small practical task.
2. Let the learner attempt implementation; review or debug their work. Do not provide a full working solution unless requested.
3. Verify observable behavior and record real evidence. Stop at each lesson checkpoint. Never automatically advance through multiple lessons.
4. Keep active infrastructure code in `infrastructure/`. The original agent-built implementation was removed after Module 00; it remains in Git at `67c0f75:reference/experiment-00/`. Consult it only as a troubleshooting fallback, without copying it wholesale or running historical scripts against current resources.
5. Protect unrelated Optiplex services. Get explicit approval before host installs, network/firewall changes, creating VMs, destructive commands, or disruptive tests. No broad cleanup.
6. Keep resources isolated and disposable, and source/state/secrets separate. Don't claim an environment check happened unless it did.

## Current focus
**Module 00 is complete; Stage 01 is in progress.** Health API, systemd supervision, Ansible artifact deployment, and in-memory task experiment checkpoints are complete. See `docs/experiments/01d-in-memory-tasks.md`: learner formatted/tested/built in the guest, deployed the executable, and verified task loss after a controlled restart with the same artifact. Guest formatting has not been synchronized back to repository source. Last deployed SHA-256 is `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`; last service active/running, PID 27439, InvocationID `d48ba32e6aed4185aea429563bc3bc67`. Rediscover guest address/identity after a break. Stop at this checkpoint; agree the next experiment with the learner before implementing or advancing. Read `docs/experiments/00-learning-run.md` for infrastructure evidence/limitations and Stage 01 notes for application/deployment state. The host is a Debian 13 Optiplex (i5-6500, 32 GB RAM) running existing media services. k6 will be run later from a separate Windows desktop.
