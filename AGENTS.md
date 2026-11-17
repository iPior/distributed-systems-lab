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
**Module 00 is complete; Stage 01 is in progress.** Health API, systemd supervision, Ansible deployment, task-state/restart, concurrency tests, and HTTP baseline checkpoints are complete. Latest evidence: `docs/experiments/01e-http-baseline.md`. Corrected eight-user runs reached 11743.16 and 12222.68 health requests/second, p95 1.08/1.00 ms, all checks passed. Windows k6 2.4.0 uses an SSH tunnel via media-server (192.168.2.103), local port 28080, to guest 192.168.78.136:8080. Use per-invocation `ObscureKeystrokeTiming=no`: disable/re-enable testing strongly supported that setting as a cause of earlier low-concurrency delay. Alias optiplex points to another host (192.168.2.101). Learner declined changing the network/path. Guest resource sample had 33% CPU idle; Windows report: 12 logical processors, total CPU below 32%, no core sustained near 100%. Host had competing work and idle CPU in one sample. No bottleneck or saturation established; resource samples are not synchronized to specific k6 summaries. Stop at this saved checkpoint. Before further load, agree a separate bounded experiment with monitoring, fixed controls, stopping conditions, and required approval for potentially disruptive tests. Keep original/corrected tunnel results separate. Source formatting was done in the guest and not copied back; last recorded deployed checksum `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af` was not rechecked for performance runs. Latest guest top showed health-api PID 638; older PID/InvocationID observations are historical. Rediscover guest address/identity after a break. Host reboot recovery started the existing lab network/guest with learner approval and started the static API service; autostart policy unchanged. Read `docs/experiments/00-learning-run.md` and Stage 01 notes for evidence/limitations. Host: Debian 13 Optiplex, i5-6500, 32 GB RAM, existing media services. Keep k6 on the separate Windows desktop.
