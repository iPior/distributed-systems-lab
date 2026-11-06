# Experiment 01b: Run the API as a guest service

**Status:** experiment checkpoint complete; evidence and learner reflection recorded

**Architecture stage:** 01

**Source revision / configuration:** API checkpoint `4b93b9d` on `learn/module-01-go-api`

## Resume here

- Learner requested committing/pushing the completed health-endpoint experiment and beginning service supervision. Commit `4b93b9d` was pushed and local/remote branch tips matched before this proposal was created.
- Current concept: the systemd system service manager manages a service independently of the interactive SSH shell used to administer it.
- Learner correctly predicted that disconnecting/reconnecting SSH does not affect the API because systemd controls it; before/after HTTP checks supported that prediction.
- Checkpoint complete: process identity, SSH independence, explicit stop/start, journal lifecycle, and restored host HTTP were verified. Learner explained that a running process can still have an incorrect endpoint. Service is last observed active/running with MainPID 21433. Learner requested committing and pushing the completed unit source and this log. Stop here; boot enablement, reboot, automatic recovery, and Ansible deployment remain later decisions.

## Concept and research question

Can the guest manage the existing API independently of an SSH session, with observable service state, logs, and explicit start/stop behavior?

## Hypothesis

Learner prediction: "since systemd is contorlling the api, nothing happens when we disconnect and reconnect". Instructor confirmed the expected independence of a system service from its administering SSH session. Explicit service stop, process failure, and VM shutdown are different lifecycle events; they are not part of an ordinary SSH disconnect/reconnect.

## Architecture and controls

Proposed arrangement:

```text
Optiplex client -> guest private IP:8080 -> same Go API
                                              ^
                                  guest systemd manages process
```

- Reuse the current VM, API response contract, all-interface listener, and port 8080.
- Keep the reviewed service-unit source in `infrastructure/systemd/dsl-api.service`.
- Install the compiled executable at guest `/opt/dsl-api/health-api`, separate from the temporary build directory.
- Run the system service as existing unprivileged guest user `lab`; guest administration uses sudo to install the executable and unit.
- Manage only the named lab service. Observe service lifecycle and captured stdout/stderr through the guest journal.
- Inspect guest service-manager availability and existing paths/unit names before installation. Recheck guest identity/address before reconnecting after a break.
- First scope is explicit start/stop and independence from SSH. Boot enablement, reboot testing, restart-on-failure policy, graceful request draining, and Ansible deployment are later decisions.

## Procedure

Proposed checkpoints, one at a time:

1. Review the research question and architecture; learner predicts SSH-disconnection behavior. Inspect the guest service manager using read-only commands.
2. Explain unit sections and have the learner draft the smallest service unit. Review source and validate before guest installation.
3. Inspect target paths, install the reviewed binary/unit in the guest, and deliberately stop the old foreground API before starting the service to avoid a port conflict.
4. Verify the service state, process user, listening socket, and bounded Optiplex health response. End only the administrative SSH session and repeat the Optiplex health check.
5. Exercise an explicit service stop/start and inspect journal lifecycle entries. Reflect on what supervision changes and record actual evidence. Stop at the checkpoint.

## Observations and evidence

Learner supplied guest `ps -p 1 -o pid,comm,args`: PID 1, command `systemd`, args `/sbin/init`. `systemctl --version` reports systemd 257 (`257.13-1~deb13u1`). `systemctl show dsl-api.service -p LoadState -p FragmentPath` reports `LoadState=not-found` and empty `FragmentPath`, so the proposed unit is not currently loaded/known to the manager. This is read-only prerequisite evidence, not service installation or runtime verification. The API was last observed running in the foreground.

Learner supplied guest path checks: both `/opt/dsl-api` and `/etc/systemd/system/dsl-api.service` report no such file or directory. Existing `/home/lab/dsl-api-4301092/health-api` is executable, owned by `lab:lab`, mode `-rwxrwxr-x`, size 8,411,319 bytes, with displayed timestamp `Nov 6 00:16`. No deployment destination was found to overwrite. Next proposed commands create only the guest lab application directory and install a root-owned 0755 executable copy; they do not start a process.

Learner predicted installing the executable would not start the API and that an explicit start would be needed. Instructor confirmed the separation between deployment and execution. Learner then supplied `ls -l /opt/dsl-api/health-api`: `-rwxr-xr-x`, owner/group `root root`, size 8,411,319 bytes, displayed timestamp `Nov 6 00:32`. This confirms the intended deployed file metadata; it is not proof of a systemd-managed process. The original foreground API has not yet been deliberately stopped in this experiment.

Learner supplied `Unit validation passed` after the instructed unit transfer and SSH-executed `systemd-analyze verify ... && echo ...` command. This supports successful validation of the staged guest unit. No warning text was supplied. This does not establish installation, manager reload, or service activation. Next steps distinguish installing a definition, loading that definition into the manager, and explicitly starting the service.

Learner initially predicted the pre-start state as "closed" and post-start as active. Instructor clarified the systemd term `inactive`. After reporting execution of the install/reload commands, learner supplied `LoadState=loaded`, `ActiveState=inactive`, and `FragmentPath=/etc/systemd/system/dsl-api.service`. This verifies manager loading from the intended guest unit path without activation. The original foreground API is separate from this inactive service. Next prediction: starting another listener on the same port while the foreground API remains bound would cause a bind conflict; stop the original process before the planned service start.

Learner supplied guest `sudo systemctl start dsl-api.service` followed by status: loaded from `/etc/systemd/system/dsl-api.service`, static; active (running) since `Fri 2026-11-06 00:38:31 UTC`; Main PID 21388 (`health-api`); cgroup `/system.slice/dsl-api.service`; executable `/opt/dsl-api/health-api`. Journal excerpt records systemd[1] Started the service. This establishes the reported manager-owned startup state. The 38ms status age and reported resource counters are not controlled performance measurements. Process user, HTTP behavior, continued operation after SSH disconnect, and explicit service stop/start remain to verify. The port-conflict prediction was requested but not answered; no deliberate conflict test was performed.

Learner supplied guest process inspection: PID 21388, PPID 1, user `lab`, command `/opt/dsl-api/health-api`. Guest-loopback health GET returned explicit HTTP 200, JSON content type, length 41, and expected two fields (Date `Fri, 06 Nov 2026 00:39:40 GMT`). Learner exited that SSH session, supplied the logout/connection-closed message, and then issued the Optiplex-to-guest private-IP health GET from the confirmed host prompt. That returned HTTP 200 with the same expected content (Date `Fri, 06 Nov 2026 00:39:52 GMT`). This supports the learner's prediction that the system service continues independently of the administrative SSH session. No post-disconnect PID observation or response-time measurement was supplied.

Learner predicted an inactive service and no HTTP response after explicit stop. Supplied guest commands show `sudo systemctl stop dsl-api.service` followed by `ActiveState=inactive`, `SubState=dead`. Learner then ran the private-IP curl at the guest prompt `lab@dsl-node-01`, receiving curl error 7 after 0 ms with no HTTP response. This confirms the stopped state and a failed guest-to-own-private-IP connection; do not describe that request as an Optiplex-origin probe. Next: restore the service, inspect lifecycle journal, and collect host HTTP restoration evidence.

Learner restored the service and supplied `MainPID=21433`, `ActiveState=active`, `SubState=running`. Unit journal excerpt records Started at `Nov 06 00:38:31`, Stopping/Deactivated successfully/Stopped at `00:41:05`, and Started at `00:44:26`, on guest `dsl-node-01`. These entries are from systemd[1], not successful-request logging added to the API. The subsequent confirmed Optiplex-origin health GET returned `HTTP/1.1 200 OK`, JSON content type, length 41, and expected fields (response Date `Fri, 06 Nov 2026 00:44:40 GMT`). This verifies explicit service restoration and host HTTP behavior. Journal timestamps are displayed lifecycle events, not controlled recovery-time or uptime measurements.

Learner final reflection: "there could have been something wrong with the endpoint, application running doesnt mean its bug free". Instructor confirmed that manager process state and application-level behavior are separate checks; a successful health request verifies only the tested endpoint/path, not the absence of all application bugs. Instructor added the missing final newline to repository unit source as a formatting-only change; no guest configuration or runtime behavior was changed by that normalization.

## Analysis

The reviewed minimal unit declares the description, `User=lab`, and the agreed absolute executable path. The staged unit passed validation, the executable had the intended ownership/permissions, and the manager loaded the installed unit from the intended path before activation.

The guest manager started the deployed API under the specified unprivileged user. Successful HTTP responses before and after the administrative SSH disconnection establish continued service on the tested path, consistent with system-service ownership rather than shell ownership. Explicit stop produced inactive/dead and a failed guest-private-IP request. A subsequent explicit start produced a new MainPID, active/running state, a lifecycle journal entry, and restored Optiplex HTTP 200. These observations support the research question for the exercised lifecycle. The learner completed the process-state versus endpoint-behavior reflection.

## Trade-offs and limitations

Service state alone does not establish HTTP correctness; verify both. Separating process lifetime from SSH does not provide redundancy or protect against VM/host failure. The Go source, executable, and deployed service configuration still require explicit synchronization.

## Conclusion and next question

The guest systemd service ran the same API as unprivileged `lab`, remained responsive after the administrative SSH session closed, stopped on an explicit stop command, and resumed serving HTTP 200 after an explicit start. Unit journal entries made the observed lifecycle visible. Learner reflection is complete, and the service is left running. This experiment checkpoint is complete; the learner requested committing and pushing it. This is one manual deployment/lifecycle experiment, not boot persistence, automatic crash recovery, reproducible Ansible deployment, or full Stage 01 completion. Agree on the next focused question before proceeding.
