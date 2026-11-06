# Experiment 01d: Application state in one Go process

**Status:** paused at learner request; proposed scope recorded, prediction pending

**Architecture stage:** 01

**Starting source revision:** `7a977aa` on `learn/module-01-go-api`

## Resume here

- Completed: health endpoint and binding experiment, guest systemd supervision, and repeatable Ansible artifact deployment. See Experiments [01](01-single-api.md), [01b](01b-service-supervision.md), and [01c](01c-ansible-deployment.md).
- Next proposed feature: create and retrieve tasks using process-local memory. No task API implementation has begun.
- Resume by reviewing this research question and architecture, collecting the unanswered prediction below, and agreeing the precise HTTP contract. Then give the learner one small implementation exercise. Do not provide a complete solution unless requested.
- Pending prediction: if we create a task, restart the API service, and request the same ID, what should happen, and why?
- Last observed guest: `dsl-node-01`, private IP `192.168.78.136`, Go 1.24.4, Debian 13. Rediscover its address and verify identity before reconnecting after the break.
- Last observed service: `dsl-api.service`, active/running, PID 22632, InvocationID `b84a44dd5e274fcdbdc544844a8f2364`, Description `Distributed Systems Lab HTTP API`; Optiplex health GET returned HTTP 200. These are Experiment 01c observations, not fresh pause-time health checks.
- Repository source is in `api/` on the Optiplex. Build tooling is in the guest. Current ignored controller artifact is `.local/artifacts/dsl-api/health-api`, SHA-256 `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`; it still contains only the health API. Ansible deploys that artifact rather than compiling source.
- Unit source: `infrastructure/systemd/dsl-api.service`; deployment entry point: `infrastructure/ansible/deploy-api.yaml`; guest executable: `/opt/dsl-api/health-api`. A source edit requires a verified build, controller-artifact update, and deployment before it affects the running API.
- Service is left in its last observed running state. No task implementation, guest restart, new deployment, or runtime test is included in this pause/documentation request.

## Concept and research question

What happens to application data when it lives only in one API process's memory?

The task resource gives the application useful state to investigate later with replicas and PostgreSQL. Stay within Stage 01 for this experiment.

## Hypothesis

**Pending learner prediction.** Ask what a GET for a previously created task should return after an API process restart. Record the learner's answer before testing.

## Architecture and controls

Proposed architecture:

```text
Optiplex HTTP client -> existing guest:8080 -> one Go API process
                                            shared in-memory task store
```

- Reuse the current VM, systemd service, listener, and deployment playbook.
- Keep `GET /healthz` available with its existing response contract.
- Smallest proposed feature: `POST /tasks` accepts a title and returns a generated ID with HTTP 201; `GET /tasks/{id}` returns that task with HTTP 200 or a missing-task response with HTTP 404.
- Review exact request/response JSON, invalid-input behavior, and ID-generation/lifetime rules before coding. These details are not yet agreed. Expected validation failures should have explicit HTTP responses.
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

**Not yet measured.** Only this proposed scope and resume notes have been recorded. No task has been created and no restart prediction has been supplied.

## Analysis

Pending implementation and real evidence.

## Trade-offs and limitations

One process and an in-memory store isolate state lifetime without database or replication behavior. ID reuse, concurrent access, invalid input, and service readiness must be handled explicitly. This experiment will not establish durable storage, multi-instance correctness, capacity, or a reproducible source-to-artifact build pipeline.

## Conclusion and next question

Paused before implementation. Resume with the pending prediction and reviewed HTTP contract. Stage 01 remains in progress.
