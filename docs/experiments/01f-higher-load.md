# Experiment 01f: Behavior as concurrency increases

**Status:** overload observation checkpoint complete; Windows SSH tunnel emitted accept/file-limit errors; API capacity not established

**Architecture stage:** 01

**Source context:** baseline checkpoint `fe72b64`; existing `load-tests/health-baseline.js`. No application change was requested.

## Concept and research question

What happens to throughput, latency, and failures as more active clients compete for the existing API path? Learner prefers studying overload behavior and moving toward HAProxy/replicas rather than fully diagnosing the limiting component.

## Hypothesis

Learner predicted roughly 24,000 requests/second at sixteen users, assuming linear scaling from approximately 12,000 at eight. Instructor explained that additional waiting can reduce scaling and that throughput can flatten while latency rises.

## Architecture and controls

Use the Windows k6 client, SSH tunnel through media-server (192.168.2.103), local port 28080, and guest health endpoint from [Experiment 01e](01e-http-baseline.md). The intended tunnel setting is per-invocation `ObscureKeystrokeTiming=no`; no fresh confirmation or artifact checksum accompanied these summaries. Each baseline run specifies concurrent users and a 30-second duration with one continuous GET per iteration, five-second request timeout, and three response checks. This measures active k6 users through the entire path, not a real-user population or isolated server capacity.

## Procedure

Instructor proposed agreeing on monitoring and stopping conditions before higher load. Learner supplied a sixteen-user summary before those conditions were agreed. Subsequently learner explicitly approved a bounded 32-user run to execute themselves, with exploratory cutoffs of greater than 1% HTTP/check failures or p95 over 100 ms, plus immediate manual stop if unrelated services were affected. The script has no automatic threshold abort. Instructor did not authorize or execute further increases; learner subsequently supplied 64/128/256-user summaries. No 32-user output was supplied. Exact run chronology for the pasted 64/128/256 results is not established; table rows below are sorted by concurrency.

## Observations and evidence

All measurements are learner-supplied summaries; no fresh instructor runtime checks or synchronized resource observations were performed. Exact run timestamps were not captured.

| VUs | Requests | Requests/second | Average HTTP latency | p95 HTTP latency | Maximum HTTP latency | HTTP failures | Failed checks |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 16 | 494,149 | 16468.818922 | 0.92123 ms | 1.99 ms | 17.06 ms | 0 | 0/1,482,447 |
| 64 | 674,639 | 22486.122954 | 2.79 ms | 4.83 ms | 21.90 ms | 0 | 0/2,023,917 |
| 128 | 687,138 | 22899.968405 | 5.52 ms | 9.36 ms | 31.22 ms | 0 | 0/2,061,414 |
| 256 | 601,031 | 19981.757962 | 12.70 ms | 21.23 ms | 5 seconds | 123/601,031 | 246/1,802,970 |

All supplied summaries reported zero interrupted iterations. Learner subsequently supplied repeated `accept: Too many open files` messages in the Windows PowerShell window running `ssh -N -o ObscureKeystrokeTiming=no -o ExitOnForwardFailure=yes -L 127.0.0.1:28080:192.168.78.136:8080 media-server`. This identifies the Windows SSH tunnel process as the emitter, with failure during accept of local forwarded connections. It supports a tunnel-side socket/file-resource limit, not an API-emitted file-limit error. The exact underlying Windows/OpenSSH limit was not inspected. The maximum duration of five seconds matches the configured request timeout, but this alone does not establish why those requests failed. At 256 users, 123 status checks and 123 instance checks failed; the JSON status check showed no failures. Total checks are fewer than three times requests, so do not assume every response reached all three checks; the supplied SSH error does not explain exactly how each failed response was processed by k6.

## Analysis

Sixteen users reached about 16,469 requests/second rather than the predicted 24,000. From 64 to 128 users, throughput rose only about 1.8% while average latency nearly doubled. At 256 users, throughput fell about 12.7% relative to 128 users, latency increased, and failures appeared. This is evidence of a plateau and degraded behavior on the measured path, not identification of an API CPU limit.

HTTP failures at 256 users were approximately 0.0205%, and failed checks approximately 0.0136%; p95 was 21.23 ms. These remain below the proposed 1%/100 ms cutoffs. Therefore no last-passing/first-failing range under those cutoffs has been established. The first supplied run with errors was 256 users; 128 users passed without reported HTTP or check failures. The later exact error output attributes the observed accept/file-limit failures to the Windows SSH tunnel. This is a measurement-path limit, not evidence of the API's intrinsic connection capacity. Other constraints may also contribute to the plateau; no exclusive bottleneck claim is supported.

## Trade-offs and limitations

Single short observations at each supplied load, no synchronized resource monitoring or post-load unrelated-service health evidence, no exact timestamps or retained raw logs beyond transcribed summaries/error text, and no fresh deployed-artifact verification. SSH forwarding remains in the path. Do not generalize these numbers to normal users, task operations, sustained reliability, or intrinsic server capacity. Do not increase file limits or rerun higher load before reviewing the actual error and agreeing scope.

## Conclusion and next question

Overload observation checkpoint complete: the measured path plateaued around 22,000-23,000 health requests/second at 64-128 active users; 256 users produced lower throughput and some failures. Exact error output identified the Windows SSH tunnel as emitting `accept: Too many open files`. Thus the first observed failures expose a tunnel-side limit, not an established API failure boundary. Stop increasing load. Learner explicitly closed Stage 01 and requested resuming with HAProxy next time. Stage 01 is complete with the documented limits; it does not claim an API maximum-capacity result. Stage 02 remains unstarted. Next session: agree the topology for three API instances behind HAProxy, predict what happens when task creation and retrieval reach different in-memory stores, and choose one small learner exercise. Retain the current network unless separately approved. The tunnel may still cap future performance comparisons. No resource-limit, network, service, or replica provisioning change was performed or authorized by this closure.
