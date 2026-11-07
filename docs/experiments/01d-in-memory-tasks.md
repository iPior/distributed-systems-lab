# Experiment 01d: Application state in one Go process

**Status:** experiment checkpoint complete; tests, executable deployment, restart/state evidence, and reflection recorded

**Architecture stage:** 01

**Starting source revision:** `7a977aa` on `learn/module-01-go-api`

## Resume here

- Checkpoint complete on 2026-11-07. Source checkpoint `7e47027` was transferred to guest `~/dsl-api-7e47027`, formatted, tested, built, copied to the controller artifact, and deployed by the learner. Repository source formatting has not been synchronized back from the guest.
- Ordinary and race-enabled tests passed. Current tests are sequential; no deliberate concurrent store-access test was executed.
- POST created task 1; GET retrieved it. After a controlled service restart with an unchanged executable, health returned 200 and GET of the old ID returned 404 before any new creation. This supports the recorded prediction of process-local state loss.
- Guest: `dsl-node-01`, private IP `192.168.78.136`, Go 1.24.4, Debian 13. Rediscover address and verify identity after a break.
- Last learner-supplied service state: active/running, PID 27439, InvocationID `d48ba32e6aed4185aea429563bc3bc67`.
- Built, controller, and deployed executable SHA-256: `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`. Controller artifact: `.local/artifacts/dsl-api/health-api`; guest executable: `/opt/dsl-api/health-api`.
- Deployment entry point: `infrastructure/ansible/deploy-api.yaml`; unit source: `infrastructure/systemd/dsl-api.service`. Ansible copies the artifact rather than compiling source.
- Stop here. Stage 01 remains in progress; agree the next focused experiment before implementing or advancing.

## Concept and research question

What happens to application data when it lives only in one API process's memory?

The task resource gives the application useful state to investigate later with replicas and PostgreSQL. Stay within Stage 01 for this experiment.

## Hypothesis

Learner predicted tasks would not survive restart because no persistent storage exists. Instructor clarified that process-local memory is discarded on process exit; the proposed old-ID GET should return 404 after restart, before any new task is created. The observed restart experiment supports this prediction.

## Architecture and controls

Proposed architecture:

```text
Optiplex HTTP client -> existing guest:8080 -> one Go API process
                                            shared in-memory task store
```

- Reuse the current VM, systemd service, listener, and deployment playbook.
- Keep `GET /healthz` available with its existing response contract.
- Smallest proposed feature: `POST /tasks` accepts a title and returns a generated ID with HTTP 201; `GET /tasks/{id}` returns that task with HTTP 200 or a missing-task response with HTTP 404.
- Reviewed POST contract: request `{"title":"Learn Go"}`; response HTTP 201 with JSON `{"id":1,"title":"Learn Go"}`. Titles are trimmed; missing/blank titles, failed JSON decoding, and extra JSON/trailing garbage receive HTTP 400 via `http.Error` (plain text). IDs start at 1 and increment within the process. Unknown fields are currently ignored; request media-type enforcement and body-size limits are not implemented.
- Reviewed GET behavior: `strconv.Atoi` parses the path ID; failed parsing (including overflow) or non-positive IDs return HTTP 400 with `invalid id`; missing positive IDs return HTTP 404 with `task not found`. Errors use `http.Error` plain-text responses with a trailing newline. Found tasks return HTTP 200 with JSON ID/title. The same mutex protects lookup and creation; lookup returns a task value, so response encoding happens after releasing the lock.
- Share one in-memory store across handlers and synchronize concurrent access. Go HTTP handlers can execute concurrently; review both the map and ID allocation rather than assuming sequential requests.
- Keep the executable fixed during the restart/state experiment. Record source revision, artifact checksum, deployed service identity, and requests/results. Test the old ID after restart before creating another task, to avoid confusing ID reuse with retained data.
- Persistence and multiple API processes are later architecture stages. This experiment does not require a new VM or host/network changes.

## Procedure

Proposed checkpoints, one at a time:

1. Review scope/contract and collect the restart prediction. Stop before implementation until the design is understood.
2. Learner implements the smallest task representation/store, with reviewed synchronization and ID behavior.
3. Learner adds creation and retrieval handlers incrementally. Verify correctness, invalid input, and missing IDs with focused tests and observable requests.
4. Build the reviewed source in the guest, verify the new artifact, update the ignored controller artifact, and deploy with the existing Ansible playbook. This will exercise an executable update, which Experiment 01c did not test.
5. Verify health, create a task, retrieve the same ID/title, and record evidence. Review the lab-only service-restart procedure before execution, then request the same ID after readiness returns. Record actual service identity and HTTP results.
6. Reflect on process-local state and restart behavior, document limitations, and stop at the checkpoint.

## Observations and evidence

### Historical source-review checkpoint

On 2026-11-06, instructor reviewed learner-written POST source. Creation locks around ID allocation, insertion, and counter increment. Handler validation occurs before insertion; response headers precede HTTP 201 and JSON encoding; encoding failures are logged rather than attempting to replace a committed status. One shared store is wired in `main`. Learner initially attributed unique concurrent IDs to memory allocation; instructor explained mutex serialization and the need to protect future map reads with the same lock. Learner acknowledged the explanation. No formatting command, build, test, guest connection, task creation request, deployment, or restart check has run in this session. Learner requested saving this checkpoint before beginning GET.

### Verification and restart experiment (2026-11-07)

Instructor read-only discovery confirmed the existing libvirt domain running, UUID `bc9dc19c-d81b-4df5-85bf-0a94ac7b2be3`, MAC `52:54:00:56:8a:8a`, and DHCP address `192.168.78.136` on `dsl-tf00-net`. Strict SSH using the existing known-hosts file confirmed hostname `dsl-node-01`, user `lab`, and Go 1.24.4. Remaining command results below are learner-supplied terminal evidence.

- Initial Optiplex test attempt failed with `go: command not found`; compilation belongs in the guest. Initial guest race-test attempt from the home directory failed with `-race requires cgo`. Diagnostics showed `CGO_ENABLED=0`, linux/amd64, expected compiler `gcc`, neither gcc nor clang on PATH, and only the older health-source module under `~/dsl-api-4301092`.
- After source transfer, learner listed `go.mod`, `main.go`, `tasks.go`, and `tasks_test.go` in `~/dsl-api-7e47027`. Ran `gofmt -w main.go tasks.go tasks_test.go` and `go test ./...`; result `ok github.com/iPior/distributed-systems-lab/api 0.003s`.
- Guest package inspection found gcc/libc6-dev absent. Installation simulation proposed 36 new packages, zero upgrades, zero removals. Instructor supplied guest-only installation instructions for `gcc libc6-dev`. Installation output was not supplied, so the exact completed package transaction is unverified. Subsequent `CGO_ENABLED=1 go test -race ./...` passed (`ok`, 1.013s). This checks exercised paths and does not establish concurrent store safety.
- `go build -o health-api .` succeeded. Guest build and transferred controller artifact both reported SHA-256 `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`. Formatting occurred in the guest; no claim of byte-for-byte source equality with the repository or reproducible build is made.
- Before deployment, guest service was active/running with PID 22632, InvocationID `b84a44dd5e274fcdbdc544844a8f2364`, and old executable SHA-256 `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`.
- Learner predicted executable-only deployment would run the restart handler. Ran `ansible-playbook -i infrastructure/ansible/inventory.yaml infrastructure/ansible/deploy-api.yaml` on the Optiplex. Executable copy changed, unit copy unchanged, Restart API ran, no reload handler ran, and HTTP readiness succeeded. Recap: ok=8, changed=2, unreachable=0, failed=0. Python interpreter-discovery warning did not prevent success.

| Observation | Evidence |
| --- | --- |
| After executable deployment | SHA-256 matches new artifact; active/running; PID 27371; InvocationID `14c8aa423c864963870f6e9ac75d1a90` |
| Optiplex POST `/tasks` | HTTP 201, application/json, `{"id":1,"title":"Learn Go"}`; response Date 2026-11-07 23:50:10 GMT |
| Optiplex GET `/tasks/1` | HTTP 200, same ID/title; response Date 23:50:26 GMT |
| After reviewed lab-service restart | Active/running; PID 27439; InvocationID `d48ba32e6aed4185aea429563bc3bc67`; executable SHA-256 unchanged |
| Optiplex GET `/healthz` after restart | HTTP 200, `{"instance":"dsl-node-01","status":"ok"}`; response Date 23:51:06 GMT |
| Optiplex GET `/tasks/1` before new creation | HTTP 404, text/plain; `task not found` with trailing newline; response Date 23:51:13 GMT |

Learner attributed potential ID reuse to the restart. Instructor clarified that a new process initializes an empty map and `nextID = 1`; reuse of an ID is not retention of the original task. No post-restart task creation was requested or observed.

## Analysis

The changed invocation identity confirms a new service invocation, while the unchanged executable checksum controls for a code change during the restart/state test. Successful retrieval before restart, restored health afterward, and missing old ID before new creation support the hypothesis that this API's task data lives only in its process. The earlier executable deployment also exercised copy-triggered restart behavior, with unit content unchanged.

## Trade-offs and limitations

One process and an in-memory store isolate state lifetime without database or replication behavior. Tests passed, including race instrumentation, but are sequential and do not establish concurrent safety. This was one restart observation, with learner-supplied terminal evidence; no restart-duration, load, capacity, multi-instance, crash-recovery, or durable-storage conclusion is supported. Guest formatting has not been copied back to repository source. Build provenance remains manual. No fresh unrelated-service health check was performed in this experiment.

## Conclusion and next question

Checkpoint complete: task creation/retrieval, tests, executable deployment, and controlled restart verified. The old task disappeared after restart with the same executable and a healthy API, supporting the process-local state prediction. Stop here and agree the next focused experiment with the learner; Stage 01 remains in progress.
