# Stage 02 handoff: three API instances behind HAProxy

## User decision and first checkpoint

On 2026-11-17, learner explicitly closed Stage 01, including its LAN/performance follow-ups, and requested a separate agent to continue Stage 02. Stop Stage 01 load exploration. Begin with a **design discussion**, not provisioning.

The goal is to compare a single API server with horizontal scaling. First investigate task correctness across independent in-memory stores; only then compare throughput. Agree how three API instances map to VMs, where HAProxy lives, their CPU/RAM allocations, and how Windows reaches the frontend. Those choices are not approved yet.

Explain the independent-store concept, ask the learner to predict create/read behavior when requests reach different replicas, and propose one small learner-written task. Stop at each checkpoint. Do not provide the full implementation unless asked. Later correctness evidence should distinguish missing data from per-replica numeric-ID collisions that can return a different task under the same ID. Make backend routing observable rather than assuming a request sequence.

## Read first

Check `git status --short` before edits. Read `AGENTS.md`, `README.md`, `docs/ROADMAP.md`, and `curriculum/00-infrastructure/README.md`. Stage 02 has no dedicated curriculum lesson yet; do not read all old lessons or automatically repeat Module 00.

Relevant evidence:

- [01d: in-memory task correctness and restart](experiments/01d-in-memory-tasks.md).
- [01i: final direct-LAN baseline and learner-defined failure criteria](experiments/01i-direct-lan-baseline.md).
- [01h: approved LAN migration, identity and rollback](experiments/01h-lan-access.md).
- [01e: original/corrected SSH measurements](experiments/01e-http-baseline.md), [01f: tunnel-side accept error](experiments/01f-higher-load.md), and [01g: host-local generator contention](experiments/01g-host-local-k6.md) explain historical limitations. Keep these paths separate.
- Use `docs/EXPERIMENT_TEMPLATE.md` for new experiment records.

## Current infrastructure, last verified

- Shared media-server: `192.168.2.103`, hostname `optiplex-media-server`, Debian 13, i5-6500, 32 GB RAM. Alias `optiplex` points to a different host, `192.168.2.101`.
- Existing single API guest `dsl-node-01`: one vCPU, 2048 MiB; UUID `bc9dc19c-d81b-4df5-85bf-0a94ac7b2be3`, MAC `52:54:00:56:8a:8a`.
- Current **DHCP LAN address `192.168.2.44`**, attached through existing host Ethernet bridge `br0`. Guest interface `ens3`; default gateway `192.168.2.1`. Rediscover identity/address after a break. The old `192.168.78.136` is historical.
- Direct API URL: `http://192.168.2.44:8080`. Windows desktop is `192.168.2.18`. Windows health access was learner-confirmed; host HTTP and Ansible checks passed.
- Keep k6 on Windows, using direct LAN HTTP. No SSH tunnel is needed for load traffic. The old tunnel workaround `ObscureKeystrokeTiming=no` matters only when revisiting historical tunnel access.
- Learner disconnected ProtonVPN during connectivity setup. Its state during later individual runs was not freshly confirmed. Record VPN/route state for future comparisons.
- Current Ansible inventory and Terraform source describe LAN access. Ordinary Terraform plan was clean after migration; recheck before making changes.
- Pinned libvirt provider 0.8.3 does not implement existing NIC-source changes through its update handler, despite an in-place plan. Migration used a reviewed, graceful lab-only shutdown, XML redefinition preserving UUID/disk/MAC, and restart, followed by a reviewed refresh-only state update. See 01h; do not blindly apply the obsolete migration plan.
- Private `dsl-tf00-net` still exists, active/persistent, but has no attached guest and its bridge route shows `linkdown`. No cleanup is requested.
- API systemd unit is **static**, not enabled at boot. Start it explicitly after guest restarts; autostart policies were not changed. Last recorded API PID 635, InvocationID `970322a9bd1f4832bbeb6de11f8c27c4`; these are historical until rechecked.
- Last verified executable SHA-256: `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`. Guest formatting was never copied back to repository source; do not claim reproducible source/binary equivalence.
- API listens on `0.0.0.0:8080`. Existing code/deployment live in `api/` and `infrastructure/`. No HAProxy or extra replicas exist.

## Closed Stage 01 observations

All results are health GETs, not representative task traffic or ordinary users. Direct-LAN summaries are learner-supplied, mostly without exact run timestamps or synchronized resource monitoring.

| Finding | Recorded observation |
| --- | --- |
| Highest observed passing throughput | 23,131.68 req/s at 512 VUs; p95 30.10 ms |
| Highest tested passing concurrency | Two 7500-VU runs: 20,011.59 / 19,568.18 req/s; p95 420.60 / 430.63 ms |
| Lowest tested failing concurrency | 7750 VUs; p95 549.62 ms; 8000 also failed at 550.04 ms |
| HTTP/check failures | Zero in all supplied direct-LAN summaries |
| Duration | Commands requested 30 seconds each; some final elapsed times included a short drain |

Learner defines failure as **p95 >500 ms OR p99 >1000 ms OR HTTP error rate >1%**. Equality passes. Use HTTP request duration and `http_req_failed`; keep response-content checks visible separately. Exact p99 was not printed, but every HTTP maximum was below 1000 ms, proving that criterion passed for those samples. Iteration-duration maxima are a different metric.

The observed 7500-7750-VU range is not a stable capacity limit. The limiting component and sustainable API maximum remain unidentified. Highest passing concurrency and highest observed throughput are different comparison points. No further Stage 01 diagnosis is required before the design discussion; any later capacity claim needs matched client, guest/API and host observations and repeatability.

`load-tests/health-baseline.js` still lacks thresholds and p99 reporting. Its instance check is hard-coded to `dsl-node-01`; a later replica benchmark must deliberately adapt instance validation while preserving correctness checks. The learner implements unless asking for code. Reused Windows export name `direct-lan-8vu.json` overwrote earlier exports at that path; use unique filenames in future. Local `.local/` monitoring, XML, state backups and plans are ignored and do not travel with Git.

## Scope, permissions and working-tree state

- The Stage 02 request authorizes teaching, read-only discovery and preparing reviewable designs. Obtain explicit approval before installs, VM creation, network/firewall changes, destructive actions or disruptive tests. Earlier approval covered only the existing guest's LAN migration, not new replicas.
- Protect unrelated media services, HAOS and devbox. Other VM persistent XML was unchanged during migration; media application health and privileged host firewall/container inspection remain unverified. Do not claim preservation beyond observed evidence.
- Three one-vCPU replicas increase total allocated API CPU relative to the one-vCPU baseline and still share the same physical host. Record that distinction, proxy overhead and competing work when interpreting later throughput. Do not promise a threefold gain.
- Learner requested **committing and pushing the Stage 01 closure**, with Stage 02 starting next session. The closure includes the following changes: `AGENTS.md`, `README.md`, `docs/ROADMAP.md`, `docs/experiments/01e-http-baseline.md`, `infrastructure/ansible/inventory.yaml`, and `infrastructure/terraform/vm.tf` are modified. Experiments 01f through 01i and this handoff are new files. Check current Git status and preserve any later changes; do not assume the closure files are still uncommitted.
- Use plain hyphens, never em dashes. Never add agent co-author attribution to commits. Keep agent-local directories ignored. Do not give development effort disproportionate weight in technical decisions.

Next action: review the three-instance/HAProxy topology with the learner and obtain the independent-store behavior prediction. No Stage 02 provisioning has happened.

## Fresh-agent opening proposal

Separate UI-task creation failed because its local MCP service could not be reached. A fresh-context subagent, `/root/stage_02`, reviewed the handoff instead. Its initial proposal is Windows -> dedicated HAProxy VM -> three API VMs: each API VM 1 vCPU / 2 GiB, HAProxy 1 vCPU / 1 GiB, total lab allocation 4 vCPUs / 7 GiB. This is a proposal only, pending shared-host headroom review and learner agreement; no infrastructure or service change was made. Separate VMs still share the Optiplex's physical CPU and failure domain.

Initial learner task: assume empty independent stores; create a task titled `alpha` on A, then predict status/title for its ID read on A, B and C. Create a different task titled `beta` on B, then predict another read of the original ID on B. Review the prediction before implementation. Backend identity must become observable: API `/healthz` currently hard-codes `dsl-node-01` and task responses include no backend identity, so neither existing response format alone distinguishes replicas.

## Final session decision

Learner requested committing and pushing the Stage 01 closure and deferring Stage 02 until next time. The topology and prediction above are a proposed opening only; no learner prediction or topology approval has been received. Resume with the design discussion next session, not provisioning.
