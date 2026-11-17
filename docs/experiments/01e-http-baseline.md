# Experiment 01e: HTTP baseline through an SSH tunnel

**Status:** baseline checkpoint complete; corrected 1/2/4/8-user runs and sampled resource observations recorded; stop before a separate higher-load experiment

**Architecture stage:** 01

**Source context:** application and concurrency-test checkpoint `347b78e`; load scripts are included with this baseline checkpoint. The executable was not rebuilt for these measurements.

## Concept and research question

How do throughput and response latency change when the number of clients making continuous health requests increases? Begin with small bounded workloads; do not infer maximum API capacity from low-concurrency tests.

## Hypothesis

Learner predicted throughput would not increase indefinitely with clients and that system resources would limit scaling. Instructor explained that virtual users are concurrent request loops, not requests per second. With one outstanding request per user and no deliberate pause, throughput depends on concurrency and iteration duration. A two-user estimate was about 61 requests/second if iteration duration remained near the one-user observation of 32.63 ms.

## Architecture and controls

```text
Windows k6 -> 127.0.0.1:28080 -> SSH to media-server (192.168.2.103)
                            -> guest 192.168.78.136:8080/healthz
```

- k6 installed on the separate Windows desktop using `winget install k6 --source winget`. Learner verified `k6.exe v2.4.0 (commit/dd9395906d, go1.26.11, windows/amd64)` after reopening the terminal.
- Guest `dsl-node-01`: last recorded allocation one vCPU and 2048 MiB RAM, Go 1.24.4, Debian 13. Go/Debian specifications are carried forward. Later process identification confirmed the lab VM allocation of one vCPU and 2048 MiB. Host also runs unrelated services; resource samples were collected later rather than alongside every baseline run.
- Guest NIC MAC `52:54:00:56:8a:8a`, lease `192.168.78.136/24`. Last recorded deployed artifact SHA-256 `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af` was not rechecked during this session.
- Same endpoint, tunnel, response checks, 30-second duration, and no sleep for baseline runs; virtual users changed through the CLI. Each iteration sends one GET and checks HTTP 200, instance `dsl-node-01`, and status `ok`. Request timeout is five seconds.
- SSH forwarding, Windows scheduling, LAN conditions, and shared host activity are part of this measured path and possible influences on latency/throughput.

## Procedure

1. Restore access after the host reboot. Learner found the guest stopped and lab network inactive, with autostart disabled. Learner explicitly approved starting the existing network and guest. Subsequent lease discovery succeeded. Guest unit was static/inactive; instructor supplied `sudo systemctl start dsl-api.service`. Subsequent Optiplex health GET returned 200.
2. Establish Windows connectivity through the correct SSH host. Keep the tunnel open in a separate Windows terminal:

   ```powershell
   ssh -N -o ExitOnForwardFailure=yes -L 127.0.0.1:28080:192.168.78.136:8080 media-server
   ```

3. Edit scripts on the Optiplex using VS Code over SSH, then copy to Windows:

   ```powershell
   scp media-server:/opt/distributed-systems-lab/load-tests/health-smoke.js .
   scp media-server:/opt/distributed-systems-lab/load-tests/health-baseline.js .
   ```

4. Run smoke and baseline workloads on Windows, reviewing each output before proceeding:

   ```powershell
   k6 run -e BASE_URL=http://127.0.0.1:28080 health-smoke.js
   k6 run -e BASE_URL=http://127.0.0.1:28080 health-baseline.js
   k6 run --vus 2 -e BASE_URL=http://127.0.0.1:28080 health-baseline.js
   ```

5. Repeat the two-user run unchanged, then run four users for 30 seconds. Investigate the sharp change with a one-user repeat followed by a four-user repeat. A three-user run and return to two users narrowed the discontinuity; then the SSH setting comparison below isolated a measurement confounder.

## Observations and evidence

All runtime evidence in this experiment is learner-supplied terminal output. Setup occurred on 2026-11-16 Toronto time, with supplied HTTP Date headers on 2026-11-17 UTC. Exact k6 run timestamps were not captured.

Connectivity troubleshooting:

- Direct Windows TCP connection to guest port 8080 failed. Local Windows binding to port 18080 returned permission denied; its cause was not established. Port 28080 could be used.
- Initial tunnel through alias `optiplex` timed out. A remote hostname check returned `optiplex-debian-server`; that alias points to 192.168.2.101. The lab terminal reported `optiplex-media-server`, route source 192.168.2.103. Switching to alias `media-server` resolved access. A separate diagnostic curl failed because its wildcard argument expanded remotely.
- Windows curl through the corrected tunnel returned HTTP 200 with the expected instance/status, response Date 2026-11-17 02:02:41 GMT. No firewall rule, LAN route, or SSH server configuration edit was requested for tunnel access.

Smoke test: one user, five iterations, one-second pause per iteration; five requests, 15/15 checks passed, zero HTTP failures. Average latency 19.21 ms, p95 27.56 ms, throughput 0.980476 requests/second. This validates setup, not capacity; the pause sets the approximate request pace and the sample is tiny.

Baseline observations (each 30 seconds, no deliberate pause):

| Run | VUs | Requests | Requests/second | Average latency | p95 latency | Maximum latency | Checks passed | HTTP failures |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| 1 | 1 | 920 | 30.637581 | 32.61 ms | 46.46 ms | 62.84 ms | 2,760/2,760 | 0/920 |
| 2 | 2 | 1,791 | 59.663101 | 33.49 ms | 47.80 ms | 63.84 ms | 5,373/5,373 | 0/1,791 |
| 3 (repeat) | 2 | 1,853 | 61.766497 | 32.35 ms | 46.28 ms | 63.17 ms | 5,559/5,559 | 0/1,853 |
| 4 | 4 | 187,281 | 6242.676673 | 0.60864 ms | 1.00 ms | 94.13 ms | 561,843/561,843 | 0/187,281 |
| 5 (return to 1 VU) | 1 | 893 | 29.749133 | 33.58 ms | 47.90 ms | 64.04 ms | 2,679/2,679 | 0/893 |
| 6 (4-VU repeat) | 4 | 184,539 | 6151.166622 | 0.61941 ms | 1.03 ms | 5.50 ms | 553,617/553,617 | 0/184,539 |
| 7 | 3 | 133,000 | 4433.269805 | 0.63689 ms | 1.12 ms | 5.38 ms | 399,000/399,000 | 0/133,000 |
| 8 (return to 2 VUs) | 2 | 1,826 | 60.848782 | 32.84 ms | 46.70 ms | 64.02 ms | 5,478/5,478 | 0/1,826 |

All supplied baseline summaries reported zero interrupted iterations. Instructor clarified that p95 means approximately 95% of measured requests completed at or below that latency; it does not establish an application latency target. Learner acknowledged that each request could take similar time while additional users keep more requests in flight.

### Timing-path diagnostics and tunnel-setting comparison

Five separate curl requests from the Optiplex directly to the guest all returned HTTP 200, with total times 0.000926, 0.000643, 0.000712, 0.000745, and 0.000785 seconds. Five Windows curl requests through the original tunnel all returned 200, with totals 0.002149, 0.051849, 0.053212, 0.053706, and 0.053593 seconds. These samples use separate curl connections and differ from k6's connection reuse. They show delay on the Windows/tunnel path outside k6, but do not isolate one component by themselves.

Learner identified the Windows client as `OpenSSH_for_Windows_9.5p1, LibreSSL 3.8.2`, executable `C:\Windows\System32\OpenSSH\ssh.exe`. Instructor found an [upstream report for the same client version](https://github.com/PowerShell/Win32-OpenSSH/issues/2322) describing X11 forwarding delays mitigated by disabling keystroke-timing obfuscation. This related report motivated a hypothesis, not a diagnosis of this HTTP path. The [SSH configuration manual](https://man.openbsd.org/ssh_config.5) describes the feature and its default 20 ms timing interval.

Learner was instructed to stop the original tunnel and restart it with the same destination and forwarding port, adding only a per-invocation option:

```powershell
ssh -N -o ObscureKeystrokeTiming=no -o ExitOnForwardFailure=yes -L 127.0.0.1:28080:192.168.78.136:8080 media-server
k6 run --vus 1 -e BASE_URL=http://127.0.0.1:28080 health-baseline.js
```

The supplied one-user result with the modified tunnel: 55,642 requests, 1,854.715176 requests/second, average HTTP latency 0.45666 ms, p95 1.09 ms, maximum 5.96 ms, average iteration duration 0.53215 ms, 166,926/166,926 checks passed, zero HTTP failures, zero interrupted iterations. This is a separate measurement configuration and must not be pooled with the earlier tunnel baseline.

The reversal used a recreated tunnel with explicit `ObscureKeystrokeTiming=yes` and the same one-user workload. Learner supplied 923 requests, 30.737929 requests/second, average HTTP latency 32.50 ms, p95 46.48 ms, maximum 63.54 ms, average iteration duration 32.53 ms, 2,769/2,769 checks passed, zero HTTP failures and interrupted iterations. Re-enabling the option restored the earlier delay even after another tunnel restart. This disable/re-enable comparison strongly supports SSH timing obfuscation as a cause of the low-concurrency delay in this setup. It does not establish the internal mechanism for the two-to-three-user transition. Restore the per-invocation option to no before further baseline work; retain SSH forwarding in the documented measurement path.

### Baseline with timing obfuscation disabled

After the reversal test, learner was instructed to restore the tunnel with explicit `ObscureKeystrokeTiming=no` and keep it open for the following 30-second runs of the unchanged baseline script. Evidence is supplied k6 summaries; no fresh resource or artifact checks accompanied these runs.

| Run | VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | Checks passed | HTTP failures |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| Corrected baseline 1 | 1 | 46,728 | 1557.592887 | 0.58461 ms | 1.14 ms | 4.82 ms | 140,184/140,184 | 0/46,728 |
| Corrected baseline 2 | 2 | 124,689 | 4156.215754 | 0.44348 ms | 0.77091 ms | 305.20 ms | 374,067/374,067 | 0/124,689 |
| Corrected baseline 3 (repeat) | 2 | 128,180 | 4272.630606 | 0.43060 ms | 0.76450 ms | 5.64 ms | 384,540/384,540 | 0/128,180 |
| Corrected baseline 4 | 4 | 220,788 | 7359.549464 | 0.51670 ms | 1.00 ms | 5.28 ms | 662,364/662,364 | 0/220,788 |
| Corrected baseline 5 | 8 | 352,299 | 11743.155872 | 0.63870 ms | 1.08 ms | 8.61 ms | 1,056,897/1,056,897 | 0/352,299 |

All five completed without interrupted iterations. The one-user summary was pasted twice; the duplicate is not counted as another execution. Two-user throughput differed by about 2.8% relative to the first two-user run, with similar p95 latency. The 305.20 ms maximum did not recur in the repeat; its cause and frequency are unknown. The first corrected one-user throughput was about 16% below the earlier option-comparison run. Thus the observed 2.67-fold increase from one to two users is not established as a stable scaling factor. A small p95 does not exclude rare much slower requests.

Learner supplied four-user and then eight-user summaries. From four to eight users, throughput increased about 59.6%, average HTTP latency increased about 23.6%, and p95 increased from 1.00 to 1.08 ms. All checks passed. These are single runs at each concurrency, without resource observations; reduced scaling does not identify a bottleneck or prove saturation. Subsequent resource samples and an eight-user repeat are recorded below. No load beyond eight users was measured.

### Eight-user repeat and sampled resource observations

An additional supplied eight-user summary reported 366,685 requests, 12,222.684257 requests/second, average HTTP latency 0.60428 ms, p95 1.00 ms, maximum 8.35 ms, average iteration duration 0.64997 ms, 1,100,055/1,100,055 checks passed, zero HTTP failures, and zero interrupted iterations. This was about 4.1% above the earlier corrected eight-user throughput. Later top/Task Manager observations were requested during eight-user repeats, but no matching additional k6 summary or synchronized timestamps were supplied. Do not pair those samples with a specific throughput result or count them as additional measured k6 runs.

Guest top observations supplied by the learner:

- `health-api` PID 638 used 68.0% and 67.3% CPU in separate process samples, with RES 14,008 KiB and 0.7% memory. These are interval observations, not whole-run averages. This freshly observed PID supersedes earlier pre-reboot service PIDs; no fresh InvocationID was collected.
- Guest summary at displayed time 02:41:40: user CPU 36.1%, system 20.6%, software interrupts 10.3%, idle 33.0%, I/O wait 0%, steal 0%. Total memory 1974.4 MiB, available 1744.3 MiB, no swap. This interval did not show guest-wide CPU saturation; it does not rule out other intervals or component bottlenecks.

Host top at displayed time 21:43:53 showed user CPU 39.8%, system 28.9%, software interrupts 4.8%, idle 26.3%, I/O wait 0.2%, steal 0%; memory available 16093.3 MiB, swap used 1.0 MiB. Process samples: qemu PID 1225797 at 112.9% CPU, sshd-session PID 1259990 at 54.5%, and python3 PID 1269761 at 100%. Process CPU percentages use a single-core basis; QEMU includes multiple threads and virtualization overhead. Individual-core graphs were requested but not supplied, so host per-core saturation is unverified.

Learner's later `ps -ww -p 1225797,1259990,1269761 -o pid,ppid,user,pcpu,args` identified PID 1225797 as `guest=dsl-node-01`, UUID `bc9dc19c-d81b-4df5-85bf-0a94ac7b2be3`, `-smp 1`, memory 2147483648 bytes. PID 1259990 was `sshd-session: pior`, but whether it carried the tunnel or another session was not established. The Python PID was absent at that check. The ps CPU figures of 15.2% for QEMU and 13.1% for SSH are lifetime averages and are not comparable to top's short interval samples. No process was stopped as part of this investigation.

On Windows, learner reported total CPU did not exceed 32%, 12 logical processors, and no individual CPU graph staying near 100% during the requested eight-user observation. Individual k6/SSH CPU and memory readings were not supplied. This suggests no apparent sustained Windows CPU saturation during the observation; it does not rule out brief peaks, single-thread constraints, or non-CPU limits.

Learner asked about shortening the measurement path. Instructor discussed routing to the private subnet, LAN bridging, and host-local load generation. Learner declined changing the path; keep the existing SSH tunnel. No route, firewall, bridge, or host-install change was performed for this proposal. Learner agreed to record and save this checkpoint before a separate higher-load experiment.

## Analysis

Corrected eight-user observations reached 11,743-12,223 health requests/second with p95 1.00-1.08 ms and all response checks passing. Disabling/re-enabling SSH timing obfuscation strongly supports it as a cause of the earlier low-concurrency delay. Keep measurements from the original and corrected tunnel configurations separate.

From four to eight users, throughput grew less than the user count while average latency increased. Sampled guest and Windows CPU retained headroom, while the host carried lab VM, SSH, and unrelated work. These observations do not identify a bottleneck or establish API saturation. Resource samples lack synchronization with a supplied k6 summary; they support a limited description of the environment, not causal attribution.

## Trade-offs and limitations

Health GETs are a small read-only workload, not representative task traffic. SSH and network overhead remain included. Repeats are limited and uneven; no statistical confidence interval, dedicated idle-host control, latency target, or maximum-capacity conclusion is supported. Resource observations are sparse learner-supplied samples, not continuous traces; host per-core and Windows per-process measurements remain unavailable. Deployed artifact identity was not freshly checked. The five-second request timeout and default graceful stop mean a 30-second active workload is not an absolute process-runtime bound. The scripts contain checks but no thresholds or automatic load-stop criteria. Summary evidence is transcribed here; raw k6 artifacts were not saved in the repository. No new unrelated-service health measurements were supplied.

## Conclusion and next question

Baseline checkpoint complete: smoke checks, corrected 1/2/4/8-user comparisons, SSH timing-setting reversal, and sampled resource observations are recorded. Eight users achieved roughly 12,000 health requests/second through the fixed SSH path with all checks passing; no API capacity limit or bottleneck was established. Stop here. Stage 01 remains in progress. Before more load, agree a separate bounded experiment with fixed controls, monitoring, stopping conditions, and any explicit approval required for potentially disruptive tests. Retain the existing network and per-invocation `ObscureKeystrokeTiming=no` tunnel.
