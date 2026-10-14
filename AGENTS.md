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
4. Keep the new implementation in `infrastructure/`. The completed original lives in `reference/experiment-00/` and on `main`. **Don't copy it wholesale or run archived scripts in-place.**
5. Protect unrelated Optiplex services. Get explicit approval before host installs, network/firewall changes, creating VMs, destructive commands, or disruptive tests. No broad cleanup.
6. Keep resources isolated and disposable, and source/state/secrets separate. Don't claim an environment check happened unless it did.

## Current focus
**Module 00, Lesson 01: Virtualization and isolation.** Start with its explanation and read-only checks; don't provision anything. The host is a Debian 13 Optiplex (i5-6500, 32 GB RAM) running existing media services. k6 will be run later from a separate Windows desktop.
