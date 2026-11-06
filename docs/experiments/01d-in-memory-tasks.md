# Experiment 01d: Application state in one Go process

**Status:** POST source checkpoint saved; GET implementation and all runtime verification pending

**Architecture stage:** 01

**Starting source revision:** `7a977aa` on `learn/module-01-go-api`

## Resume here

- Completed: health endpoint and binding experiment, guest systemd supervision, and repeatable Ansible artifact deployment. See Experiments [01](01-single-api.md), [01b](01b-service-supervision.md), and [01c](01c-ansible-deployment.md).
- Learner implemented `api/tasks.go`: task representation, initialized map and ID counter, mutex-protected creation, and the POST handler. `api/main.go` creates one shared store and registers `POST /tasks`. Source reviewed only; formatting, build, tests, and deployment remain pending because the learner is working from another machine without guest access.
- Next: begin GET incrementally with a mutex-protected store lookup returning `(Task, bool)`, then review ID parsing, missing-task behavior, response encoding, and route wiring. Do not provide a complete solution unless requested.
- Learner predicted tasks would not survive restart. Instructor clarified process-local RAM is not automatically persisted; requesting the old ID before creating another task should return 404. This remains an untested prediction.
- Last observed guest: `dsl-node-01`, private IP `192.168.78.136`, Go 1.24.4, Debian 13. Rediscover its address and verify identity before reconnecting after the break.
- Last observed service: `dsl-api.service`, active/running, PID 22632, InvocationID `b84a44dd5e274fcdbdc544844a8f2364`, Description `Distributed Systems Lab HTTP API`; Optiplex health GET returned HTTP 200. These are Experiment 01c observations, not fresh pause-time health checks.
- Repository source is in `api/` on the Optiplex. Build tooling is in the guest. Current ignored controller artifact is `.local/artifacts/dsl-api/health-api`, SHA-256 `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`; it still contains only the health API. Ansible deploys that artifact rather than compiling source.
- Unit source: `infrastructure/systemd/dsl-api.service`; deployment entry point: `infrastructure/ansible/deploy-api.yaml`; guest executable: `/opt/dsl-api/health-api`. A source edit requires a verified build, controller-artifact update, and deployment before it affects the running API.
- No live guest checks or service changes were performed during this source-only session. The last deployed artifact still contains only the health API.

## Concept and research question

What happens to application data when it lives only in one API process's memory?

The task resource gives the application useful state to investigate later with replicas and PostgreSQL. Stay within Stage 01 for this experiment.

## Hypothesis

Learner predicted tasks would not survive restart because no persistent storage exists. Instructor clarified that process-local memory is discarded on process exit; the proposed old-ID GET should return 404 after restart, before any new task is created. Not yet tested.

## Architecture and controls

Proposed architecture:

```text
Optiplex HTTP client -> existing guest:8080 -> one Go API process
                                            shared in-memory task store
```

- Reuse the current VM, systemd service, listener, and deployment playbook.
- Keep `GET /healthz` available with its existing response contract.
- Smallest proposed feature: `POST /tasks` accepts a title and returns a generated ID with HTTP 201; `GET /tasks/{id}` returns that task with HTTP 200 or a missing-task response with HTTP 404.
- Reviewed POST contract: request `{"title":"Learn Go"}`; response HTTP 201 with JSON `{"id":1,"title":"Learn Go"}`. Titles are trimmed; missing/blank titles, failed JSON decoding, and extra JSON/trailing garbage receive HTTP 400 via `http.Error` (plain text). IDs start at 1 and increment within the process. Unknown fields are currently ignored; request media-type enforcement and body-size limits are not implemented. GET ID parsing and exact error responses remain to be reviewed.
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

**Runtime behavior not yet measured.** On 2026-11-06, instructor reviewed learner-written POST source. Creation locks around ID allocation, insertion, and counter increment. Handler validation occurs before insertion; response headers precede HTTP 201 and JSON encoding; encoding failures are logged rather than attempting to replace a committed status. One shared store is wired in `main`. Learner initially attributed unique concurrent IDs to memory allocation; instructor explained mutex serialization and the need to protect future map reads with the same lock. Learner acknowledged the explanation. No formatting command, build, test, guest connection, task creation request, deployment, or restart check has run in this session. Learner requested saving this checkpoint before beginning GET.

## Analysis

POST source reviewed; GET implementation and runtime evidence remain pending.

## Trade-offs and limitations

One process and an in-memory store isolate state lifetime without database or replication behavior. ID reuse, concurrent access, invalid input, and service readiness must be handled explicitly. This experiment will not establish durable storage, multi-instance correctness, capacity, or a reproducible source-to-artifact build pipeline.

## Conclusion and next question

POST source checkpoint complete, without build or runtime verification. Next exercise is synchronized GET lookup, followed by incremental HTTP handling. Stage 01 remains in progress.
