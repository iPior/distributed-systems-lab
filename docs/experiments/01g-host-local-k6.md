# Experiment 01g: Host-local k6 path comparison

**Status:** one approved run complete; stop at comparison checkpoint

**Architecture stage:** 01 configuration, supplementary comparison after Stage 01 closure

**Source revision / configuration:** `fe72b64` plus existing uncommitted Stage 01 closure documentation; unchanged `load-tests/health-baseline.js`

## Concept and research question

Does moving k6 from Windows through SSH forwarding to media-server directly improve the observed eight-VU health workload?

## Hypothesis

Learner predicted similar performance because the existing path already traverses media-server. Removing SSH forwarding could reduce overhead, but host-local k6 also competes with the VM and unrelated services for physical resources.

## Architecture and controls

```text
Historical: Windows k6 -> LAN -> SSH forwarding on media-server -> guest API
This run:   media-server k6 -> existing private bridge -> guest API
```

- Host identity verified as `optiplex-media-server`; existing Debian 13 Optiplex, i5-6500, 32 GB RAM. No network configuration changed.
- Downloaded official Linux amd64 k6 v2.4.0 into ignored `.local/tools/k6`, with explicit learner approval. Version check: commit `dd9395906d`, Go 1.26.11. Matches recorded Windows k6 version/commit. No system package installation.
- Release archive SHA-256 verified against the published checksum: `5a5ecf23529ccf93135bc7e5b0c1cd5126ecf5f2672e22eed1556a0ab74abb87`.
- Fresh discovery: running `dsl-node-01`, UUID `bc9dc19c-d81b-4df5-85bf-0a94ac7b2be3`, MAC `52:54:00:56:8a:8a`, DHCP address `192.168.78.136` on `dsl-tf00-net`. Allocation of one vCPU / 2048 MiB is carried forward, not freshly inspected.
- Strict SSH confirmed hostname `dsl-node-01`, API active/running, PID 638, InvocationID `63183976a1ee4367ab2a5328a73eaaae`.
- Fresh executable checksum matched the prior artifact: `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`.
- Unchanged script, eight continuous VUs, 30-second active duration, five-second request timeout, three checks per request. Host and guest CPU sampled concurrently with `vmstat -t 1 35`.
- Historical comparator is the corrected Windows eight-VU baseline with `ObscureKeystrokeTiming=no`, not the earlier delayed tunnel measurements or higher-VU runs.

## Procedure

Learner explicitly approved one eight-VU, 30-second run after readiness verification. Instructor executed:

```sh
.local/tools/k6 run --vus 8 --duration 30s \
  --summary-export .local/k6-host-comparison-20261117-2038/summary.json \
  -e BASE_URL=http://192.168.78.136:8080 \
  load-tests/health-baseline.js
```

Run timestamps: 2026-11-17 20:37:59 to 20:38:29 America/Toronto. Direct health requests before and after the workload returned HTTP 200 with expected instance/status.

## Observations and evidence

| Metric | Host-local result |
| --- | --- |
| Requests | 449,893 |
| Requests/second | 14,996.219084 |
| Average HTTP latency | 0.433875 ms |
| p95 HTTP latency | 1.049147 ms |
| Maximum HTTP latency | 20.457789 ms |
| HTTP failures | 0/449,893 |
| Checks passed | 1,349,679/1,349,679 |
| Interrupted iterations | 0 |

CPU samples, excluding vmstat's first since-boot row: over 34 one-second intervals spanning the workload and adjacent time, host idle ranged 3-92% (mean 18.5%); guest idle ranged 13-100% (mean 30.1%). These are aggregate CPU observations, not per-process attribution or exact whole-run averages. Two earlier host intervals before the workload showed 67-70% idle.

Attempted unrelated-container status inspection failed: Docker socket access was denied, and `sudo -n docker ps` required a password. No container state or media application health was verified. No service was changed or restarted.

Raw output, JSON summary, timestamps, CPU samples, and post-run health response are saved under ignored `.local/k6-host-comparison-20261117-2038/`. They are local evidence and will not travel with a Git checkout.

## Analysis

Throughput was 22.7-27.7% above the two historical Windows eight-VU results (12,222.68 and 11,743.16 req/s). p95 was within their 1.00-1.08 ms range. This run observed a throughput improvement with broadly similar p95, rather than unchanged performance.

This supports investigating measurement-path overhead, but does not isolate SSH as the sole cause: client OS, client hardware placement, and background activity also differ. No same-session Windows repeat was performed. Aggregate host CPU had little idle capacity in some intervals, so collocating the generator can introduce resource competition. API capacity remains unknown.

## Trade-offs and limitations

One short host-local run compared with historical measurements, not an alternating or repeated controlled comparison. No contemporaneous media-service health evidence, per-process CPU attribution, or independent generator-capacity check. Guest checksum was verified before this run but was not freshly verified for the historical Windows measurements. No task traffic or replicas were tested.

## Conclusion and next question

Observed about 15,000 health req/s at eight VUs with no failures and similar p95 to the historical tunnel baseline. Stop here; no higher load is authorized by this checkpoint. Stage 02 remains pending. Before future performance work, explicitly choose the generator placement and record its shared-resource limitations. Existing Windows tunnel settings remain applicable when using that path.

## Follow-up: learner-approved 128-VU run

Learner subsequently requested 128 VUs and conditionally 256 if the first run was fine. Instructor explained that proceeding depended on host headroom as well as HTTP results. Ran the unchanged baseline directly from media-server for 30 seconds at 128 VUs, 2026-11-17 20:40:21 to 20:40:51 America/Toronto. Only `--vus` and evidence output paths changed from the command above. Concurrent host and guest `vmstat -t 1 35` monitoring was repeated.

| Metric | Host-local 128-VU result |
| --- | --- |
| Requests | 493,768 |
| Requests/second | 16,456.315908 |
| Average HTTP latency | 7.341781 ms |
| p95 HTTP latency | 17.479767 ms |
| Maximum HTTP latency | 89.909860 ms |
| HTTP failures | 0/493,768 |
| Checks passed | 1,481,304/1,481,304 |
| Interrupted iterations | 0 |

Excluding the since-boot row, 34 one-second samples spanning the run and adjacent time showed host idle 1-84% (mean 12.8%) and guest steal 0-21% (mean 9.9%). Several consecutive workload samples showed host idle 2-3% and guest steal 14-15%. Guest steal records time the VM waited for physical CPU scheduling. These samples support shared-host CPU contention; they do not identify which host processes consumed the CPU or prove media-service degradation. Media application health remains unverified. Direct health requests before and after returned HTTP 200 with expected instance/status.

At 128 VUs, throughput was about 28% below the historical Windows 128-VU observation (22,899.97 req/s), and p95 was higher (17.48 versus 9.36 ms). Sixteen times the host-local eight-VU concurrency produced only about 9.7% more throughput. Client placement and background activity differ from historical runs, so this does not isolate a single cause or establish API capacity.

The conditional 256-VU run was **not performed** because the shared host showed little CPU headroom and guest scheduling contention. No load increase followed. Raw summary, output, timestamps, CPU samples, and post-run health response are saved locally in ignored `.local/k6-host-128-20261117/`. Stage 02 remains pending.

## Follow-up: renewed 256-VU request, stopped early

Learner subsequently asked to try 256 VUs or check whether headroom remained insufficient. Five fresh unloaded host intervals showed 76-95% idle; guest intervals showed 100% idle and zero steal. Instructor announced a bounded 256-VU attempt with an early-stop rule: interrupt k6 if host idle stayed below 5% for three consecutive one-second samples. No other services were changed.

Started 2026-11-17 20:43:13 America/Toronto with the unchanged baseline script, `--vus 256 --duration 30s`, direct guest URL, and summary export under `.local/k6-host-256-20261117/`. A Python supervisor measured host CPU deltas from `/proc/stat` (excluding guest counters already included in user/nice). Concurrent guest vmstat monitoring was also captured.

The CPU stop rule triggered after roughly four seconds of active workload. Supervisor sent SIGINT; k6 exited with code 105 and explicitly reported interruption, rather than an application-triggered abort. The summary reported 48,515 HTTP requests, 12,261.91 req/s, p95 45.58 ms, zero HTTP failures among recorded requests, and 145,530/145,530 passed checks. It also reported 257 interrupted iterations. This aborted sample is not a completed 30-second comparison or evidence that 256 VUs ran successfully. The direct API health request afterward returned HTTP 200 with expected instance/status. Media-service health remains unverified.

Raw k6 output, summary, supervisor timestamps/stop reason, host idle samples, guest CPU samples, and post-run health response remain in ignored `.local/k6-host-256-20261117/`. Unloaded headroom disappeared quickly under colocated generation. Stop here; no further load increase or repeat was performed. Stage 02 remains pending.
