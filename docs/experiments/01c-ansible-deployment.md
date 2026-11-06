# Experiment 01c: Repeatable API deployment with Ansible

**Status:** experiment checkpoint complete; evidence and learner reflection recorded

**Architecture stage:** 01

**Source revision / configuration:** start from service-supervision checkpoint `8232876`, branch `learn/module-01-go-api`

## Resume here

- Learner requested beginning repeatable Ansible deployment after completing service supervision.
- Repository was clean at initial inspection. Existing `infrastructure/ansible/playbook.yaml` only asserts Debian and manages curl. Inventory targets the lab guest with strict SSH settings.
- Current concept: a changed task can notify a named handler; repeated notifications to the same handler are coalesced within a handler run. Unchanged copy tasks should not request a restart.
- Checkpoint complete: first deployment changed2/restarted, unchanged repeat changed0/no restart, Description-only update changed2/restarted; all post-run host health checks returned HTTP 200 with expected JSON. Learner correctly concluded that an ordinary restarted task would restart even on an unchanged run. Last observed service active/running, PID 22632, InvocationID `b84a44dd5e274fcdbdc544844a8f2364`, updated Description. Stop here; agree the next focused experiment before implementation.
- Guest is `dsl-node-01`, last verified private IP `192.168.78.136`. Artifact remains ignored under controller `.local/artifacts/dsl-api/health-api`. Raw baseline/run1/run2/run3 evidence is under ignored `.local/stage01-ansible-checks/`. Rediscover address/identity before reconnecting after a break. Learner requested committing and pushing the playbook, this log, and the unit Description change.

## Concept and research question

Can Ansible deploy the existing API, repeat the deployment without changes or an unnecessary restart, and apply one controlled update with a restart and successful HTTP verification?

## Hypothesis

Learner predictions: unchanged executable/unit copies should not run the restart handler; if both copies change and notify the same handler before it runs, it should execute once. These are correct expectations to verify. Record expected first-run changes based on actual initial state; do not assume that the already configured guest needs changes.

## Architecture and controls

Proposed deployment path:

```text
Optiplex Ansible controller
  ignored executable artifact + tracked service unit
                 |
                 v
guest /opt/dsl-api/health-api + /etc/systemd/system/dsl-api.service
                 |
        systemd runs API as lab
```

- New learner-authored entry point: `infrastructure/ansible/deploy-api.yaml`, using the existing single-guest inventory. Preserve the Module 00 playbook.
- Initially use the existing compiled executable as a deployment artifact. Copy it to ignored controller `.local/artifacts/dsl-api/health-api` and record matching source/destination SHA-256 checksums before use. Track provenance as the current manually built artifact; this is not yet automated compilation from Git.
- Tracked unit source remains `infrastructure/systemd/dsl-api.service`. Record artifact checksum, unit configuration, and playbook revision for each run.
- Manage guest directory and file ownership/modes, executable copy, service-unit copy, and desired service state. Continue using `lab`, port 8080, the same API contract, and the same VM allocation.
- Unit changes notify manager reload and service restart. Executable changes notify restart. Define reload before restart and ensure pending handlers run before final service/HTTP verification.
- Request a running service through desired-state management rather than unconditionally restarting on every playbook run. Multiple changed files should request one restart-handler execution per handler run.
- Check bounded guest HTTP readiness through Ansible and verify the Optiplex HTTP path separately. Preserve boot-enablement policy; enabling boot startup is not part of this experiment.
- No host installation or networking change is needed for the proposed artifact-deployment design. Guest source compilation remains a separate step for now.

## Procedure

Proposed checkpoints, one at a time:

1. Explain handlers and record learner predictions. Inspect installed module documentation and confirm the controller/guest prerequisites.
2. Prepare and checksum the existing executable artifact in the ignored controller directory; record current service PID/invocation identity and HTTP baseline.
3. Learner builds the deployment playbook incrementally: target/assertion, directory/files, handler notifications, desired service state, and bounded HTTP readiness. Review and syntax-check before execution.
4. Run the reviewed deployment once. Record task/handler output and actual service/HTTP state, including failures. Initial manual deployment may already satisfy some tasks.
5. Run an immediate unchanged second deployment. Compare changed counts, service invocation identity/PID, journal, and HTTP. Do not infer lack of a restart solely from task output or a coincidentally matching PID.
6. Make one reviewed, meaningful unit update (for example, a description change), predict its effects, deploy once, and verify reload/restart behavior plus HTTP restoration. Do not claim executable-update handling was tested if only the unit changed.
7. Reflect and document limitations. Stop before the next experiment.

## Observations and evidence

Instructor checked controller Ansible core 2.19.11 and installed documentation for `ansible.builtin.copy` and `ansible.builtin.systemd_service`. Copy documentation describes local source to remote destination and content comparison; service documentation distinguishes idempotent `started` from unconditional `restarted` and supports manager reload.

Fresh read-only guest check through the existing strict SSH settings returned expected hostname/user. SHA-256 of `/opt/dsl-api/health-api` is `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`. Manager inspection returned active, MainPID 21433, InvocationID `b871202d0475460c8267b0c0fb1646ff`. Controller artifact destination does not exist and `git check-ignore` confirms it will be ignored. No artifact copy, deployment playbook execution, controlled update, or service change has occurred yet.

Learner supplied first real deployment output: facts/assertion/directory/executable copy ok; service definition copy changed; manager reload handler ran (ok); restart handler ran (changed); desired started state and URI readiness ok. Recap: ok=9, changed=2, unreachable=0, failed=0, skipped/rescued/ignored=0. Interpreter-discovery warning references `/usr/bin/python3.13`, but execution succeeded. First-run prediction was requested but not supplied; record observed behavior only.

Instructor read-only post-run verification found active/running service PID 21788, InvocationID `b682f4a167ce4c1d87bdb8df81efe340`, start timestamp `Fri 2026-11-06 01:58:56 UTC`. Journal records one stop/deactivate/stopped/start sequence at that timestamp. Executable checksum remains `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`; installed unit checksum now matches controller `c90223183447ba3c1521c095201dc98a0debf8e1b963018559043fe17f48906c`. Optiplex health GET returns HTTP 200 and expected JSON (Date `Fri, 06 Nov 2026 01:59:29 GMT`). Raw verification is saved in ignored `.local/stage01-ansible-checks/run1/`. These establish first-run deployment/restart, not unchanged-run idempotency yet.

Learner supplied the unchanged second run: all ordinary tasks ok, no handler output, recap ok=7, changed=0, unreachable=0, failed=0, skipped/rescued/ignored=0. Second-run prediction was requested but not supplied. Instructor post-run checks confirmed active/running, MainPID 21788, unchanged InvocationID `b682f4a167ce4c1d87bdb8df81efe340` and ExecMainStartTimestamp `Fri 2026-11-06 01:58:56 UTC`; the ten-entry unit journal excerpt has no additions relative to run 1. Host health returned HTTP 200 and expected JSON (Date `Fri, 06 Nov 2026 02:01:14 GMT`). Raw output is in ignored `.local/stage01-ansible-checks/run2/`. These observations support unchanged-run idempotency and no service restart in this repetition, beyond Ansible's changed count alone. Next controlled update is a reviewed service Description change; executable bytes remain constant.

For the proposed Description-only update, learner predicted "neither and none" for copy changes and handler execution. Instructor clarified that this applies to unchanged input, whereas modifying the unit Description changes file content: expected executable copy unchanged, unit copy changed, manager reload then one restart-handler execution. The copy module detects file content differences rather than deciding whether a particular directive changes application behavior. Controlled update remains to be reviewed and executed.

After an initial review found the old saved Description, learner saved the intended edit. Instructor confirmed the exact one-line tracked unit diff, with no executable or playbook change, and reran whitespace/syntax validation successfully. Controlled update is ready for learner execution; no third deployment or changed service result has been observed yet.

Learner supplied the controlled Description-only deployment: executable copy unchanged, unit copy changed, manager reload handler ran (ok), restart handler ran (changed), desired started state/readiness ok. Recap ok=9, changed=2, unreachable=0, failed=0, skipped/rescued/ignored=0. Instructor post-run checks show Description `Distributed Systems Lab HTTP API`, active/running, MainPID 22632, InvocationID `b84a44dd5e274fcdbdc544844a8f2364`, start timestamp `Fri 2026-11-06 02:04:25 UTC`. Journal records one stop/deactivate/stopped/start sequence at that timestamp with the new Description. Executable SHA-256 is unchanged. Installed/controller unit SHA-256 now both equal `8d00b28caf30e36d03bcfef99d90285950a8100e0fbc1b64439cd044f15a771c`. Playbook SHA-256 is `35c750f2c4823e4893cbb840eb71fbed5a6b8c92fe087013a9895999d1c5c4f4`. Optiplex health returns HTTP 200 and expected JSON (Date `Fri, 06 Nov 2026 02:05:10 GMT`). Raw verification is in ignored `run3/` evidence.

| Run | Recap changed | Reload/restart handlers | Post-run identity | Host HTTP |
| --- | --- | --- | --- | --- |
| First deployment (unit final newline) | 2 | Both ran; one restart | PID 21788, new invocation | 200 |
| Unchanged second deployment | 0 | Neither ran | Same PID/invocation/start time | 200 |
| Description-only update | 2 | Both ran; one restart | PID 22632, new invocation | 200 |

Final reflection: learner initially described a restart as needing to be handled. Instructor clarified that `started` is an idempotent desired state while `restarted` requests a new lifecycle action each time its task executes. Asked what an ordinary restarted task would do on an unchanged second deployment, learner correctly answered "it would restart". This completes the guided handler-versus-normal-task reflection.

## Drafting and baseline notes (historical)

At learner request to begin baseline collection, instructor rechecked the guest MAC/lease and connected with strict SSH settings. Baseline guest identity is dsl-node-01/lab; service active/running, MainPID 21433, InvocationID `b871202d0475460c8267b0c0fb1646ff`, ExecMainStartTimestamp `Fri 2026-11-06 00:44:26 UTC`, NRestarts 0. Deployed executable SHA-256 still matches the controller artifact. Directory and executable are root:root 0755; unit is root:root 0644. Optiplex health GET returned HTTP 200 and expected JSON, Date `Fri, 06 Nov 2026 01:57:56 GMT`.

Guest unit SHA-256 is `23d33d596d7f45d4b0ac859b7131b338129fe89192b3d43edd8afd6791cc41f7`; controller unit SHA-256 is `c90223183447ba3c1521c095201dc98a0debf8e1b963018559043fe17f48906c`. Retrieved guest unit and inspected the exact diff: the only difference is the missing guest final newline, already normalized in repository source after Experiment 01b. This is a content difference even though service semantics are unchanged. Raw guest state, HTTP response, and unit snapshot are saved in ignored `.local/stage01-ansible-checks/baseline.VcKiCr/`. No deployment, manager reload, or restart was performed during baseline collection. Requested first-run prediction based on this evidence before learner execution.

Learner explicitly requested instructor implementation of the readiness task. Instructor added it after desired started state, with task-level retry controls and no changes to the existing handler logic. Rechecked the full playbook order and ran syntax validation successfully; whitespace check also passes. This is offline validation, not an executed HTTP check or API deployment. Source-level tests were not added for this small configuration change; the planned real deployment/idempotency/update observations will verify behavior.

Instructor reviewed the standalone flush correction and reran syntax validation successfully. The playbook now describes directory/files, conditional reload/restart, and desired service state in the intended order. The normal task named `Monitor state` actually enforces started state; a descriptive ensure-running name is clearer. Installed URI module documentation confirms URL, expected status codes, and socket timeout parameters. Next: add bounded application-level readiness so process state alone is not treated as proof of HTTP success. No guest tasks ran.

Instructor reviewed the next learner edit: normal systemd task uses the intended service and `state: started`. Both copy tasks also contain `ansible.builtin.meta: flush_handlers` as a second action. Syntax check fails (exit 4) with conflicting action statements `ansible.builtin.copy` and `ansible.builtin.meta`. Explained that each task has one action, so one standalone flush task belongs after both copies. Two separate flush points between changed copies could defeat the intended single combined restart. No guest tasks ran.

Asked whether unchanged copies would notify handlers or start an already stopped service, learner correctly answered no to notifications but suggested "only restart" for service action. Instructor clarified that neither handler runs without notification; therefore the existing playbook leaves the stopped service stopped. Introduced normal `state: started` desired-state management and explicit flushing of pending handlers before that task. Flush processes only already-notified handlers and does not unconditionally restart. No guest tasks ran.

Instructor reviewed learner's unit copy with correct source/destination, root ownership, mode 0644, and both notifications. Handler order is correct: manager reload precedes service restart. Syntax validation passes. A harmless task-label typo (`defition`) does not affect execution. Next concept: notifications alone do not enforce service availability when files are unchanged; add desired running-state management after pending handlers have been processed. No guest tasks ran.

Instructor confirmed the corrected notification placement and handlers hierarchy. `Restart API` uses `ansible.builtin.systemd_service` with the intended unit name and `state: restarted`; it is no longer an unconditional ordinary task. Syntax check passes. Next exercise adds the tracked unit copy and reload-before-restart ordering; service-state enforcement and readiness checks remain later reviewed steps.

After initial saved-file rechecks showed no edits, instructor inspected the newly saved notification/restart draft. Copy task now has a descriptive name, but `notify` is indented as a copy-module option, and `Restart API` remains a normal task rather than a handler. Offline syntax validation still passes; it does not execute module argument checks or prove the desired conditional behavior. Explained two structural corrections and provided a partial hierarchy scaffold while preserving learner-authored task contents. No guest tasks ran.

Instructor reviewed learner's new copy task: controller source is anchored to `playbook_dir`, guest destination is `/opt/dsl-api/health-api`, owner/group root, mode quoted 0755. Syntax validation passes. Task label currently says `Syntax validation passes`, which does not describe its action; requested a descriptive deployment label. Next exercise introduces a task-level notification and a separately defined play-level handler, distinguishing handler name from service-module name. No guest tasks ran.

Instructor confirmed the learner added the Debian assertion before the directory task and reran syntax validation successfully. The playbook now has the reviewed target, privilege escalation, precondition, and desired directory state. No tasks have been executed. Next copy-task exercise distinguishes controller `src` from guest `dest`; use `playbook_dir` to anchor the ignored artifact path independently of the invoking shell's working directory.

Instructor rechecked the renamed `deploy-api.yaml`: `tasks` now contains a properly named list item with the correct file module parameters. Syntax check passes (exit 0), but the file still contains only the directory task; the requested Debian assertion is missing. Explained that syntax validation does not verify all experiment requirements and requested the assertion before any configuration changes. No guest tasks ran.

Instructor attempted to read the expected deployment filename, which was absent, then used repository file search to find learner draft `infrastructure/ansible/deply-api.yaml`. The draft targets lab with become true and correct directory attributes, but has `ansible.builtin.file` directly under `tasks` without a list item and omits the Debian assertion. Syntax-checking the actual draft failed (exit 4): malformed task block, expected a list rather than a mapping. No Ansible tasks or guest changes ran. Explained task list structure and requested learner corrections before adding deployment behavior.

After the instructed guest-to-controller SCP, learner supplied controller `sha256sum` output for `.local/artifacts/dsl-api/health-api`: `8f9549d822a09b3afadf6a24ed8d59d0a662fdda16c1e3bfb2bd5951fbe830a5`. This matches the independently observed guest executable checksum. It establishes matching artifact content at those checks; it does not establish a reproducible build from source. No playbook has yet deployed this artifact.

## Analysis

The observed deployments distinguish convergence from unnecessary restart. First deployment reconciled a real unit-content difference; the unchanged repeat made no changes and preserved service invocation identity, start timestamp, and journal. A controlled unit-content update ran reload before restart once, left executable bytes unchanged, and restored HTTP 200. This supports the focused question on this prepared guest. Ansible's changed count was checked against independent process/journal/HTTP evidence.

`state: started` enforces availability without restarting an already-running service. Conditional handlers apply changed deployed content, and one flush after both copies places reload/restart before final state/readiness checks. A changed unit may legitimately trigger the configured restart policy even when the change is administrative metadata.

## Trade-offs and limitations

This first scope deploys an existing executable and manages an already provisioned Debian guest. It does not establish a source-to-artifact build pipeline, clean-VM reproducibility, boot recovery, or application capacity. Only unit-content changes were exercised; executable-update notifications and simultaneous executable/unit-change coalescing were discussed but not tested. Repair of a deliberately stopped service through the deployment playbook was not exercised. URI readiness checks status 200; instructor host probes also inspected expected JSON. No controlled downtime or restart-duration measurement was performed.

## Conclusion and next question

Observed first deployment, unchanged repeat, and controlled unit update behaved as intended, with matching file/process/journal/HTTP evidence. Learner reflection is complete. This experiment checkpoint is complete; learner requested committing and pushing it. Service is left active/running with the updated Description. Stage 01 continues; agree the next focused question before proceeding.
