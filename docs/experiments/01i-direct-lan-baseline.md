# Experiment 01i: Windows k6 directly over the LAN

**Status:** closed with Stage 01 at learner's request; two 7500-VU passes, p95 failures at 7750 and 8000 VUs; limiting component unidentified

**Architecture stage:** existing single-API configuration, supplementary experiment after Stage 01 closure

**Source revision / configuration:** `fe72b64` application source context; unchanged `load-tests/health-baseline.js`; LAN attachment from Experiment 01h

## Concept and research question

Does direct Windows-to-guest HTTP access improve the eight-VU health workload compared with the historical Windows SSH tunnel path?

## Hypothesis

Learner predicted both improved throughput and avoidance of the Windows SSH accept/file-resource limit after removing the tunnel. The initial eight-VU checkpoint tested low-concurrency path performance. Learner subsequently supplied the 128/256-VU observations recorded below.

### Learner-defined pass/fail criteria

After clarifying the single-server/scaling comparison objective, learner specified a run fails if **any** of these are true:

- HTTP request duration p95 > 500 ms.
- HTTP request duration p99 > 1000 ms.
- HTTP error rate > 1%.

Equality at a boundary passes. Use `http_req_duration` for latency and `http_req_failed` for the HTTP error-rate criterion; retain existing status/JSON checks as separate correctness evidence. These are learner-selected lab criteria, not a production service target. The earlier direct-LAN summaries through 5000 VUs meet them: reported p95 values are below 500 ms, HTTP failure rates are zero, and reported maximum HTTP durations are below 1000 ms, which bounds p99 below 1000 ms even though it was not displayed. The later 7750/8000-VU summaries fail the p95 criterion, as recorded below. These are assessments of supplied metrics, not automatic threshold verdicts from those scripts.

The repository baseline script still has no thresholds and its default summary does not display p99. Next script exercise is to encode these exact inclusive limits and display p99, following the [k6 threshold documentation](https://grafana.com/docs/k6/latest/using-k6/thresholds/) and [summary statistics options](https://grafana.com/docs/k6/latest/using-k6/k6-options/reference/#summary-trend-stats). No script edit or additional load was executed when the criteria were recorded.

## Architecture and controls

```text
Windows k6 192.168.2.18 -> LAN -> guest 192.168.2.44:8080/healthz
                                      existing host bridge br0
```

- k6 runs on the separate Windows desktop, not the shared VM host. Last verified Windows version is k6 2.4.0; it was not rechecked for this run.
- Existing `dsl-node-01` on the Debian 13 Optiplex, i5-6500, 32 GB RAM. VM allocation remains one vCPU and 2048 MiB. Same VM UUID and MAC as the previous experiments.
- Guest NIC now attaches to `br0`; current DHCP address is `192.168.2.44`. No SSH forwarding carries the HTTP workload. Separate SSH used only for CPU monitoring.
- Executable checksum was verified unchanged during LAN activation: `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`. No application edit, build, deployment or service restart was requested for this workload. No separate checksum check accompanied the pasted k6 summary.
- Same baseline script, eight continuous VUs, 30-second duration, five-second HTTP timeout and three checks per iteration. Client location and concurrency match the historical eight-VU comparison, but topology, run time, and background work differ.
- Learner reported disconnecting ProtonVPN during connectivity setup. Instructor requested it remain disconnected for this run; no fresh route/VPN output accompanied the performance summary. VPN state during historical tunnel runs is unverified.

## Procedure

Instructor verified direct guest health, started three minutes of host/guest CPU logging, and asked the learner to run this from the Windows folder containing the script:

```powershell
Get-Date -Format o
k6 run --vus 8 --duration 30s --summary-export direct-lan-8vu.json -e BASE_URL=http://192.168.2.44:8080 health-baseline.js
Get-Date -Format o
```

Learner supplied the eight-VU summary below, without the requested timestamps. The requested Windows JSON export has not been transferred or inspected. No higher-concurrency run was performed by the instructor.

## Observations and evidence

Learner-supplied k6 summary, received on 2026-11-17 America/Toronto:

| Metric | Direct-LAN eight-VU result |
| --- | --- |
| Requests | 515,117 |
| Requests/second | 17,170.106393 |
| Average HTTP latency | 0.43149 ms |
| Median HTTP latency | 0.52750 ms |
| p95 HTTP latency | 0.65683 ms |
| Maximum HTTP latency | 6.05 ms |
| HTTP failures | 0/515,117 |
| Checks passed | 1,545,351/1,545,351 |
| Interrupted iterations | 0 |
| Average iteration duration | 0.46317 ms |

Both CPU loggers completed: 181 vmstat rows each, including the initial since-boot row. Collection window was 21:04:50-21:07:50 America/Toronto (guest timestamps 02:04:50-02:07:50 UTC on November 18). Excluding the first row, host idle ranged 27-97% (mean 76.8%), and guest idle ranged 6-100% (mean 84.6%). Guest steal ranged 0-2%; host steal was zero. These whole-window statistics include unloaded time and must not be described as workload averages. Host and guest samples share timestamps, but the missing k6 run timestamps prevent exact association with the 30-second summary. No per-process CPU attribution was collected.

Host direct health after receiving the result returned HTTP 200 with expected instance/status (response Date 2026-11-18 02:07:56 GMT). No media-service application health check was performed. Local CPU logs are in ignored `.local/windows-lan-baseline-8vu-20261117/`; they do not travel with a Git checkout.

## Analysis

| Eight-VU path | Requests/second | p95 HTTP latency |
| --- | --- | --- |
| Historical corrected Windows SSH run 1 | 11,743.155872 | 1.08 ms |
| Historical corrected Windows SSH run 2 | 12,222.684257 | 1.00 ms |
| Current Windows direct LAN | 17,170.106393 | 0.65683 ms |

This direct-LAN run showed 40.5-46.2% higher throughput and 34.3-39.2% lower p95 than the two historical corrected eight-VU observations in Experiment 01e. All checks passed. This supports the learner's expected improvement at this concurrency.

The SSH process is absent from the HTTP path, so its local-forwarding accept error cannot be emitted as part of this direct path. The initial eight-VU run alone did not establish behavior at higher concurrency; see the follow-up below. It also does not establish API maximum capacity or attribute the full measured improvement solely to SSH removal.

## Trade-offs and limitations

One direct run compared with historical observations, no same-session alternating comparison, no statistical confidence interval, missing run timestamps, and no fresh Windows k6 version or VPN confirmation. Guest health is a small read-only workload rather than task traffic. The VM was restarted for the preceding network migration, and background host/LAN activity can differ. CPU logs cover both unloaded and potentially loaded time; exact overlap with the summary is not established. No media-service health conclusion is supported.

## Conclusion and next question

Eight-VU direct-LAN checkpoint complete: about 17,170 health req/s with p95 0.657 ms, no HTTP failures and all checks passing. Stop before another load level. Agree the next bounded experiment before proceeding. Stage 02 remains pending; replica correctness must precede replica performance testing.

## Follow-up: learner-supplied 128 and 256 VUs

After reviewing the eight-VU checkpoint, learner supplied the following additional summaries. Instructor did not execute additional load or start new monitoring for these runs. The message repeats identical copies of both summaries; each distinct result is counted once, not as repeated executions. Exact timestamps were not supplied. The 128-VU command was not pasted; its summary identifies 128 VUs and 30 seconds. The displayed 256-VU command explicitly uses the same baseline script and direct guest URL:

```powershell
k6 run --vus 256 --duration 30s --summary-export direct-lan-8vu.json -e BASE_URL=http://192.168.2.44:8080 health-baseline.js
```

| VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | HTTP failures | Checks passed |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 128 | 682,489 | 22,746.768150 | 5.56 ms | 10 ms | 305.75 ms | 0 | 2,047,467/2,047,467 |
| 256 | 686,908 | 22,889.963873 | 11.1 ms | 17 ms | 41.05 ms | 0 | 2,060,724/2,060,724 |

Both summaries show zero interrupted iterations. The 256-VU command reused `direct-lan-8vu.json`; if run from the same directory, this overwrites the earlier exported summary. No Windows JSON artifact was inspected or transferred. Future exports should have distinct filenames per concurrency and execution.

Doubling concurrency from 128 to 256 produced only 0.63% more throughput, with average HTTP latency nearly doubling and p95 increasing 70%. These observations show a throughput plateau around 22.7k-22.9k health req/s on this direct path; they do not identify the component responsible. The 305.75 ms maximum in the 128-VU summary shows that a low p95 does not exclude rarer longer requests. Its cause and recurrence are unknown.

Direct-LAN 128-VU throughput is within 0.7% of the historical SSH 128-VU observation (22,899.97 req/s; p95 9.36 ms). Direct-LAN 256 VUs completed without reported failures, unlike the historical tunnel 256-VU result (19,981.76 req/s, p95 21.23 ms, 123 HTTP failures). Thus bypassing SSH avoided the earlier tunnel failure during this sample, while the observed throughput plateau remains close to the historical 23k region. These are different runs under incompletely controlled background conditions, not proof of a sole bottleneck or stable maximum.

The earlier three-minute CPU logs must not be assigned to these follow-up summaries: exact run timestamps and monitoring overlap are not established. No fresh per-run resource samples or Windows process observations accompanied them. A host health request after receipt returned HTTP 200 with expected instance/status (response Date 2026-11-18 02:10:15 GMT); it is not media-service health evidence.

Stop increasing load at this checkpoint. API capacity remains unestablished. Agree a separately monitored experiment before seeking the limiting component, or resume the planned Stage 02 topology discussion. No replicas or HAProxy have been provisioned.

## Follow-up: learner-supplied 512 and 1024 VUs

Learner then supplied two additional 30-second summaries. The leading `256 VUs` completion line belongs to the previous output; the first new summary identifies 512 VUs in both its metrics and completion line. It is counted as a 512-VU observation, not another 256-VU run. The 512-VU command was not supplied. The displayed 1024-VU command uses the same script and direct guest URL, again reusing `direct-lan-8vu.json` for its export.

| VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | HTTP failures | Checks passed |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 512 | 694,388 | 23,131.681756 | 22.02 ms | 30.10 ms | 79.53 ms | 0 | 2,083,164/2,083,164 |
| 1024 | 639,302 | 21,279.193748 | 47.95 ms | 58.99 ms | 134.13 ms | 0 | 1,917,906/1,917,906 |

Both summaries report zero interrupted iterations. No exact timestamps, fresh software/route/VPN confirmation, synchronized resource data, or media-service health evidence accompanied them. Instructor executed no additional load. Host direct API health after receipt returned HTTP 200 with expected instance/status (response Date 2026-11-18 02:12:27 GMT).

From 256 to 512 VUs, throughput increased only 1.06%. At 1024 VUs, throughput fell 8.01% relative to 512, while p95 rose from 30.10 to 58.99 ms (about 96%). Across the 128/256/512 observations, throughput remained around 22.7k-23.1k req/s while latency rose. These are workload/path observations; they do not establish a stable maximum, the optimal concurrency, or whether the VM, client, host, or another component is limiting throughput. Zero failures did not prevent performance degradation.

Learner clarified the objective: establish the capacity of the existing single server, then compare how horizontal scaling changes throughput. For the current baseline, the server is the existing **one-vCPU, 2048-MiB API VM**, not the entire Optiplex. Highest observed health throughput is a provisional path result. Usable capacity needs a latency/error criterion, repeated measurements, and evidence that the generator/path is not the limiting component.

Next proposed focused experiment: repeat one fixed workload in the observed plateau region while collecting matching k6 timestamps, Windows k6/CPU observations, guest/API CPU data, and host CPU data. Predict the limiting component before running it. Do not automatically increase VUs further. For the later three-replica comparison, record total allocated resources and shared-host contention; three one-vCPU VMs add CPU allocation as well as instances. Agree the Stage 02 topology and test independent-store task correctness before replica throughput work. No new benchmark, resource change, or replica provisioning has been performed for this clarified objective.

## Follow-up: learner-supplied 3000 and 5000 VUs

After stating the latency/error criteria, learner supplied two more summaries and their commands, using the direct guest URL, unchanged baseline script, 30-second duration, and 3000 then 5000 VUs. Both commands again reuse `direct-lan-8vu.json`. The repeated filename is not evidence that the export still contains eight-VU results. Instructor executed no additional load, started no new per-run monitoring, and has not inspected the Windows JSON file.

| VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | HTTP failures | Checks passed |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 3000 | 658,112 | 21,835.063730 | 136.14 ms | 158.02 ms | 250.05 ms | 0 | 1,974,336/1,974,336 |
| 5000 | 606,284 | 20,064.107691 | 246.37 ms | 286.49 ms | 543.47 ms | 0 | 1,818,852/1,818,852 |

Both report zero interrupted iterations. The displayed final elapsed active times are about 30.1 and 30.2 seconds respectively; each command requested 30 seconds. No exact timestamps, matched host/guest/client resource observations, fresh VPN/route confirmation or media-service health measurements accompanied these summaries. A host direct API health request after receipt returned HTTP 200 with expected instance/status (response Date 2026-11-18 02:17:23 GMT).

**Both runs pass the declared criteria.** Their p95 values are below 500 ms, HTTP failure rates are zero, and p99 cannot exceed their reported maxima (250.05 and 543.47 ms), both below 1000 ms. Actual p99 values remain unreported. The 5000-VU maximum exceeding 500 ms is not a p95 failure: the learner's 500 ms bound applies to the percentile, not every request.

The observed 5000-VU throughput is 13.26% below the highest supplied sample, 23,131.68 req/s at 512 VUs. The existing workload is closed-loop: each VU waits for its request/checks before starting the next iteration. Higher VU counts therefore increase outstanding work and latency without necessarily increasing requests per second. These summaries support a 30-second observation of 5000 continuous VUs within the chosen limits, not 5000 ordinary users, a stable maximum capacity, or a first-failing concurrency boundary. Highest observed passing throughput remains approximately 23.1k health req/s; repeated, instrumented observations are still needed for a defensible single-VM baseline and later replica comparison.

## Follow-up: observed pass/fail range at 7000-8000 VUs

Learner supplied a sequence of commands/summaries at 7000, 8000, 7750, 7500 and 7500 VUs and explicitly said the 7500-VU test was run twice. The entire pasted sequence appears twice in the message. Count each distinct result once: five executions, including two distinct 7500-VU runs, rather than ten executions. All commands specify the same direct URL, baseline script, 30-second duration and reused summary filename. Exact timestamps and matched resource observations were not supplied. Instructor executed no further load.

Table sorted by concurrency; the actual pasted order placed 8000 and 7750 before the two 7500-VU runs:

| VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | Checks passed | SLO result |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 7000 | 599,272 | 19,800.875498 | 347.70 ms | 397.49 ms | 847.84 ms | 1,797,816/1,797,816 | Pass |
| 7500, run 1 | 605,880 | 20,011.591120 | 368.30 ms | 420.60 ms | 547.97 ms | 1,817,640/1,817,640 | Pass |
| 7500, run 2 | 593,774 | 19,568.182035 | 376.42 ms | 430.63 ms | 788.92 ms | 1,781,322/1,781,322 | Pass |
| 7750 | 483,815 | 15,946.072332 | 477.28 ms | 549.62 ms | 576.73 ms | 1,451,445/1,451,445 | Fail: p95 >500 ms |
| 8000 | 562,676 | 18,580.209367 | 422.79 ms | 550.04 ms | 865.10 ms | 1,688,028/1,688,028 | Fail: p95 >500 ms |

All five summaries report zero HTTP failures, all checks passing and zero interrupted iterations. Their reported HTTP maxima are below 1000 ms, so all satisfy the p99 criterion by an upper bound, although exact p99 values remain unreported. Some iteration-duration maxima exceed one second; those are not HTTP request-duration p99 and do not establish failure of the declared HTTP latency criterion. The green completion/check output does not encode the learner's latency SLOs: the repository baseline script still has no threshold definitions.

**7500 VUs passes in both supplied repeats; 7750 is the lowest tested failing concurrency.** Chronologically, 8000 was the first supplied failure in this sequence. The two 7500-VU throughputs differ by about 2.22%, with p95 differing by about 2.38%. This gives an observed passing/failing range between 7500 and 7750 VUs for these short continuous-loop runs, not an exact or stable maximum concurrency. The lower throughput at 7750 than at 8000 shows run variability; background activity and generator/path effects were not controlled or measured alongside these summaries.

Highest tested passing concurrency is therefore 7500 VUs, with about 19.6k-20.0k health req/s in the two supplied repeats. Highest observed passing throughput remains 23,131.68 req/s at 512 VUs. Keep these two findings separate when comparing replicas. They do not establish ordinary-user capacity, the limiting resource, or scaling across independent physical hosts.

Host direct health after receipt returned HTTP 200 with expected instance/status (response Date 2026-11-18 02:21:51 GMT). Media-service health remains unverified. Exploratory concurrency-boundary checkpoint complete: stop further increases. A claim about sustainable single-VM capacity still requires longer/repeated observations with matched generator, guest/API and host data; later replica performance must follow Stage 02 topology agreement and independent-store correctness checks.

## Final Stage 01 closure

Learner explicitly closed Stage 01 including these follow-ups and requested a separate agent for Stage 02. Preserve the observed results and their limits; do not automatically resume Stage 01 load exploration. Continue with the [Stage 02 handoff](../STAGE_02_HANDOFF.md): agree the topology for three API instances behind HAProxy, obtain a prediction about independent in-memory stores, and investigate task correctness before performance. No replicas, HAProxy installation, new network change, commit or push occurred as part of this closure.
