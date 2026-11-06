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
**Module 00 is complete; Stage 01 is in progress and paused at the learner's request.** Completed experiments cover the health API, systemd supervision, and Ansible artifact deployment. Resume from `docs/experiments/01d-in-memory-tasks.md`: POST/GET task source and tests are authored and source-reviewed, but formatting, test execution, build, deployment, and restart/state verification are pending until guest access is available. Restart prediction and HTTP behavior are recorded. Do not implement or advance automatically. Read `docs/experiments/00-learning-run.md` for infrastructure evidence/limitations and the Stage 01 experiment notes for application/deployment state. The host is a Debian 13 Optiplex (i5-6500, 32 GB RAM) running existing media services. k6 will be run later from a separate Windows desktop.
