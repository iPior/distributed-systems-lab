# Experiment 01: A single Go HTTP API

**Status:** first experiment checkpoint complete; evidence and reflection recorded

**Architecture stage:** 01

**Source revision / configuration:** planning began at `ef28b24`; module/scaffold checkpoint `927360b`; handler and loopback server checkpoint `4301092` on `learn/module-01-go-api`

## Resume here

- Module 00 is complete. Stage 01 was requested on 2026-10-28.
- Read repository guidance, README, roadmap, Module 00 overview, prior learning-run evidence, and experiment template. Available curriculum covers only Module 00; no Stage 01 curriculum exists yet.
- Continued on `learn/module-01-go-api`, branched from `927360b`; `learn/rebuild-experiment-00` is retained as the Module 00 legacy branch. Resumed on the Optiplex after the other-laptop checkpoint `4301092`; initial working tree clean.
- Learner implemented the JSON health response, encoding-error logging, `GET /healthz` route, and a loopback listener at `127.0.0.1:8080` with server-error logging in `api/main.go`. Source review completed. Learner copied the source, formatted it, and successfully built the executable in the guest; process startup and HTTP checks remain pending.
- Current checkpoint complete: guest-local contract checks passed; wildcard-bound process PID 21246 answered guest-loopback and Optiplex-private-IP requests with HTTP 200 and expected JSON. Learner correctly explained the copy/build/restart sequence. Instructor synchronized the formatted guest source back to repository `api/main.go` after verifying only indentation, import order, blank-line whitespace, and final-newline differences. Learner requested committing/pushing this checkpoint, then beginning the service-supervision experiment with a focused question before implementation.
- Authoritative source lives in repository `api/`; the guest will receive a working copy for builds. Recheck the guest's state and DHCP address before reconnecting after the break. Do not assume that the Go installation or unfinished source has been automated through Ansible.
- Discussed progression: this endpoint belongs to the API we will grow. Later proposed experiments cover service supervision, a small application feature, and baseline measurement before replicas. The application domain and those later experiments remain undecided; they are not implementation authorization.
- Repository inspection confirmed the expected branch, latest commit `ef28b24`, and an initially clean working tree. No live infrastructure check or guest connection was performed.
- On 2026-11-04, learner-supplied `domiflist` shows `dsl-node-01` on `dsl-tf00-net` through `vnet13`, virtio MAC `52:54:00:56:8a:8a`. The network's DHCP lease matches that MAC to `192.168.78.136/24`, hostname `dsl-node-01`. This rechecks the address before connecting; it does not yet establish SSH authentication or guest identity.

## Concept and research question

First concept: an HTTP contract describes application behavior; a listening address determines which network interfaces can accept connections to that process.

Proposed question: can a single Go process in the existing guest satisfy a defined HTTP contract when called from the guest and from the Optiplex?

This first experiment establishes correctness and reachability. Latency, resource usage, and later Windows-based k6 capacity measurements require separate focused experiments; they are not completed by this checkpoint.

## Hypothesis

Original prediction prompt before execution:

- If the process listens on `127.0.0.1:8080` inside the guest, what happens to a guest-local request and an Optiplex request to the guest's private address?

On 2026-11-04, the learner predicted that a request inside the VM would work because the API listens locally, while an outside request would fail, describing the address as the wrong IP.

Instructor clarification: guest-local access to `127.0.0.1:8080` is expected to work. The guest's private IP is a valid address for reaching that VM, but a listener bound only to guest loopback does not accept connections addressed to the private interface. Each machine has its own loopback; `127.0.0.1` on the Optiplex refers to the Optiplex. Even a guest-local request addressed to the guest's private IP would not reach this loopback-only listener. These are expected outcomes, not executed checks.

## Architecture and controls

Proposed request path:

```text
Optiplex HTTP client -> existing private guest network
                    -> dsl-node-01:8080 -> one Go HTTP process
```

- Reuse `dsl-node-01`, recorded as Debian 13, 1 vCPU, 2 GiB RAM, 20 GiB disk, on `dsl-tf00-net`. These are prior observations, not fresh runtime verification.
- Use Go's standard HTTP library, one process, and no mutable application state or external dependencies.
- Proposed endpoint: `GET /healthz`, HTTP 200, JSON content type, body containing `status: "ok"` and a stable configured instance identifier `instance: "dsl-node-01"`.
- Contract also includes HTTP 404 for an unknown path and HTTP 405 for an unsupported method on `/healthz`. Exact response headers/body details will be reviewed before implementation.
- Proposed final listener is guest IPv4 `0.0.0.0:8080`; explain and review the binding choice before running it.
- Begin with a foreground process owned by an unprivileged guest user. Go 1.24.4 runs in the guest per learner-supplied output; source lives in repository `api/` and will be copied to the guest for building.
- Keep infrastructure source in `infrastructure/`. The guest Go installation is a manual learning step, not yet declared in Ansible.
- Keep the existing VM allocation and networking. Windows access is a later architecture decision. No host installation, firewall change, VM creation, disruptive test, or teardown is authorized by this proposal.
- The physical host also runs unrelated services, so future performance results must describe that shared-host limitation.

## Procedure

Proposed checkpoints, each requiring a learner attempt/review and a stop:

1. Agree on this experiment and answer the listener-address prediction. No implementation at this checkpoint.
2. Recheck guest addressing/identity, inspect tooling, and agree on source/build location and the precise HTTP contract.
3. Learner writes the minimal API; review and debug it before deployment.
4. Run the reviewed process and verify the contract from the guest and Optiplex with bounded requests. Record response status, headers, JSON fields, and failures. Investigate binding with a small controlled check if needed.
5. Reflect on what the observations establish and record the checkpoint. Starting a later experiment requires learner agreement.

Verification proposal: check the successful GET, an unknown route, an unsupported method, and repeated GET responses with the same instance identifier. Use a listener inspection to connect process binding to reachability. This endpoint reports process responsiveness and identity; it does not establish dependency health or service capacity.

## Observations and evidence

Learner supplied the NIC and matching DHCP lease listed above. Instructor checked existing inventory SSH settings: user `lab`, dedicated lab key, `IdentitiesOnly=yes`, strict host-key checking, and the existing dedicated known-hosts file. No private key contents were read. No Stage 01 API, HTTP probe, guest connection, deployment, or benchmark has been performed.

Subsequently, the learner supplied guest-session output on 2026-11-04: `hostname` returned `dsl-node-01`; `whoami` returned `lab`; `command -v go` returned no path; `go version` returned `-bash: go: command not found`. This establishes the reported guest identity and no Go executable on that session's PATH, not proof that no Go installation exists anywhere on disk. Instructor also ran `command -v go` on the Optiplex; no path was returned (exit 1). No application deployment or HTTP test has occurred.

Proposed tooling arrangement: keep authoritative application source under `api/` in this repository, copy a working source snapshot to the guest, and build/run there with guest-installed Go. Next learner check: `apt-cache policy golang-go` inside the guest, before choosing a package or performing an installation. No host install or package change has occurred.

Learner supplied guest `apt-cache policy golang-go`: installed none, candidate `2:1.24~2` from trixie/main at priority 500; backports offers `2:1.27~1~bpo13+1` at priority 100. These are the guest's cached package metadata, not a freshly verified upstream release list or the installed compiler version. Next learner exercise: predict the machine affected by a guest package installation and preview the dependency transaction with `apt-get -s --no-install-recommends install golang-go`. No installation has been performed.

Learner correctly predicted that installing from the guest session affects the guest rather than the Optiplex. Supplied simulation proposes four new packages (`golang-1.24-go`, `golang-1.24-src`, `golang-go`, `golang-src`), 0 upgrades, 0 removals, and 4 not upgraded. Versioned compiler/source packages are listed as `1.24.4-1`; metapackages as `2:1.24~2`. Recommended C/C++ build packages are excluded by `--no-install-recommends`. Next: learner performs this scoped guest installation and reports `go version`. Simulation is not evidence of installation success.

Learner subsequently supplied `go version go1.24.4 linux/amd64` from `lab@dsl-node-01`. This verifies that Go runs in the learner's guest session. Full installation transaction output was not supplied; do not claim an independently inspected package transaction. Next source exercise: learner creates `api/go.mod` with module path `github.com/iPior/distributed-systems-lab/api` and language version `1.24.0`, then requests review. This file can be authored on the Optiplex without installing Go there. No application build or HTTP behavior has been verified.

Instructor reviewed the learner-created `api/go.mod`: correct intended module path and language version, but both remain enclosed in literal `<...>` placeholder delimiters. Explained that the delimiters must be removed. Learner correction is pending; instructor did not edit the module file or run a Go validation/build.

Learner corrected the module file; instructor directly confirmed `module github.com/iPior/distributed-systems-lab/api` and `go 1.24.0`. The file lacks a final newline, a formatting detail to fix in the next learner edit. Next exercise introduces `http.ResponseWriter` and `*http.Request`: author only `healthz`, returning JSON `status: "ok"` and `instance: "dsl-node-01"` with content type `application/json` and HTTP 200. Routing, method behavior, server startup, and runtime verification remain later reviewed steps. No handler code was supplied or written by the instructor.

At the pause, instructor inspected the learner's unfinished `api/main.go` and confirmed the scaffold only. The learner also added the final newline to `go.mod`. The learner requested committing and pushing this learning checkpoint. No handler implementation, infrastructure mutation, or guest shutdown is included in that request.

On 2026-11-05, instructor reviewed the learner's completed handler and server wiring in repository source. The handler sets the JSON content type before committing HTTP 200, encodes the two string fields, and logs encoding errors. `main` registers `GET /healthz` on a dedicated mux and listens on guest loopback, logging a fatal server error if the listener returns. Go's GET pattern also accepts HEAD; unknown-path and unsupported-method behavior remains unverified. Learner correctly explained why the private guest IP would not reach a loopback-only listener. No guest connection, formatting command, build, process startup, or HTTP probe was performed in this session. Learner requested committing and pushing this checkpoint.

On resuming from the other laptop's checkpoint, instructor reviewed the two-file diff since `927360b`. The handler and startup code match the small experiment; source requires gofmt formatting and a final newline. Instructor freshly matched MAC `52:54:00:56:8a:8a` to DHCP address `192.168.78.136`, then used existing strict SSH settings with batch mode and a five-second connect timeout. Guest returned hostname `dsl-node-01`, user `lab`, and `go version go1.24.4 linux/amd64`; destination `~/dsl-api-4301092` was absent. These were read-only checks. No source copy, formatter, build, or API startup has yet run.

Learner subsequently supplied successful SCP output for `go.mod` and `main.go`, followed by guest commands in `~/dsl-api-4301092`: `gofmt -w main.go` and `go build -o health-api . && ls -lh health-api`. The conditional listing ran and showed an executable owned by `lab:lab`, mode `-rwxrwxr-x`, rounded size `8.1M`. This supports successful compilation of the formatted guest copy. No build-duration measurement, process startup, listener inspection, or HTTP response evidence was supplied. The build/start prediction prompt was not answered; explain the distinction before running the binary.

Learner correctly predicted that the server logs only errors. Their attempted probe shows a guest prompt (`lab@dsl-node-01:~/dsl-api-4301092`) running the SSH command intended for the Optiplex. SSH reported missing identity file `.local/lesson03_ed25519` and an unknown ED25519 host key under strict checking, then exited with host-key verification failure. This is an SSH client-location/path failure before HTTP; it does not establish API startup or failure. Correction: run bounded curl directly in the second guest session, or run the original SSH command from the Optiplex. No key copying, host-key bypass, or trust-store edit is needed.

Learner subsequently reported that `./health-api` remained running. After the instructed second-terminal SSH-mediated guest-local curl, learner supplied `Content-Type: application/json`, `Date: Fri, 06 Nov 2026 00:04:40 GMT`, `Content-Length: 41`, and body `{"instance":"dsl-node-01","status":"ok"}`. The supplied excerpt omits the HTTP status line, so do not claim HTTP 200 was observed yet. Header/body evidence supports guest-local process responsiveness and the intended JSON fields; it does not establish host-to-guest HTTP reachability. Response Date is a server-reported timestamp, not a measured duration.

Learner subsequently supplied explicit guest-local status checks using bounded curl with discarded bodies: `GET /healthz: 200`, `GET /missing: 404`, and `POST /healthz: 405`. These confirm the three tested status outcomes. Predictions for the latter two were requested but not supplied; record the observed results only. Together with the earlier response headers/body, these complete the initial guest-local HTTP-contract checkpoint. HEAD, other routes/methods, and host-to-guest HTTP access have not been tested.

Learner predicted the private-IP request would fail and supplied `curl -i --max-time 5 http://192.168.78.136:8080/healthz` from prompt `lab@dsl-node-01:~`. Curl returned error 7 after 0 ms with no HTTP response. This tests the guest addressing its own private interface, not an Optiplex-origin request. It agrees with the loopback-only binding expectation. Instructor clarified the client location and directed the learner to exit only the request session back to the Optiplex prompt before repeating; keep the original server session running.

Learner then supplied the same bounded private-IP curl from the confirmed Optiplex prompt `pior@optiplex-media-server:/opt/distributed-systems-lab`. It returned curl error 7 after 0 ms and no HTTP response. This establishes the failed host-to-guest HTTP attempt under the current loopback-only binding. Together with successful guest-loopback HTTP checks and the prior guest-private-IP failure, it completes the loopback reachability checkpoint. No firewall/network changes were made, and no latency or throughput conclusion follows from curl's error timing.

Instructor subsequently inspected repository `api/main.go` after the learner's edit. The diff changes only `http.ListenAndServe`'s address to `0.0.0.0:8080`; response/routing behavior is unchanged. This is a source review, not evidence of an updated running process. Next exercise distinguishes source, compiled executable, and running process before copying, rebuilding, and restarting only the lab API.

After the updated copy/format/build/restart instructions, learner reported "okay im restarted". Record this as learner-reported restart, not independently inspected deployment or build output. Listener inspection and both HTTP paths are the next checks; no successful all-interface response has yet been supplied.

Learner supplied replacement verification: guest `ss` shows `LISTEN` on `127.0.0.1:8080`, owned by `health-api` PID 21150, fd 3. Guest-loopback curl returns explicit `HTTP/1.1 200 OK`, JSON content type, length 41, and the expected two fields (response Date `Fri, 06 Nov 2026 00:13:02 GMT`). The Optiplex private-IP curl again returns error 7 after 0 ms. This directly contradicts the expected all-interface listener after replacement; the actual process is still loopback-bound. Diagnose the guest source, executable, and running-process path rather than changing firewall rules or assuming the reviewed repository edit was deployed.

Learner supplied guest source/executable diagnosis: `grep` reports line 26 still calling `http.ListenAndServe("127.0.0.1:8080", mux)` in `~/dsl-api-4301092/main.go`; `readlink /proc/21150/exe` returns `/home/lab/dsl-api-4301092/health-api`. Instructor rechecked repository source, which declares `0.0.0.0:8080`. The guest source copy has not picked up the reviewed listener edit. Next: explicitly copy and inspect that one source file before rebuilding; the old process can remain running during the copy.

Learner's next supplied SCP command targets `/home/lab/dsl-api-4301092/main.g` (missing final `o`), with a successful 639-byte transfer. It therefore did not replace the build's `main.go`. Subsequent grep ran at the Optiplex prompt, expanding `~` to `/home/pior`, and failed because that local path does not exist. Correction: exact `main.go` destination followed by SSH-executed grep of its absolute guest pathname. No rebuild or restart has occurred in this attempt.

Learner supplied the corrected SCP to `/home/lab/dsl-api-4301092/main.go` (639 bytes), followed by SSH-executed grep returning `26: ... http.ListenAndServe("0.0.0.0:8080", mux)`. This establishes that the guest source now contains the all-interface binding. It does not change the existing compiled executable or running process; learner rebuild/restart and runtime verification remain next.

After the verified source copy and rebuild/restart instructions, learner again reported the server running. Supplied guest `ss` now shows `LISTEN` at wildcard `*:8080`, peer `*:*`, owned by `health-api` PID 21246 (fd 3), replacing the previously observed PID 21150. Guest-loopback curl returned `HTTP/1.1 200 OK`, JSON content type, length 41, and the intended fields (response Date `Fri, 06 Nov 2026 00:17:20 GMT`). The Optiplex private-IP curl returned the same HTTP status, content type, length, and JSON (Date `Fri, 06 Nov 2026 00:17:29 GMT`). These establish success on both tested IPv4 paths after the listener change. The abbreviated wildcard display is not a separate IPv6 reachability test. The response timestamps are not latency measurements.

Learner reflection: "it needed to be copied over and rebuilt and then reran". This correctly distinguishes source transfer, compilation, and starting the replacement process. Instructor retrieved only guest `main.go` into an ignored temporary file under `.local/stage01-checks/`, reviewed its diff against repository source, and confirmed only gofmt formatting differences (tabs, import order, whitespace, final newline). Copied the reviewed formatted source back to repository `api/main.go`. No executable was copied or added to Git. This closes the source synchronization and guided reflection checkpoint; no commit/push or next experiment was requested.

## Analysis

The initial guest-loopback connection failure was followed by an empty listener inspection. Starting the executable resolved guest-loopback HTTP access. This demonstrates why a successful build alone does not establish a running service.

With the loopback listener, guest-local contract checks passed but private-IP requests from both guest and Optiplex failed. After the all-interface source change reached the guest, rebuilding and restarting produced a wildcard listener and HTTP 200 on both tested paths. Holding handler, port, VM, and client paths constant supports listener binding as the explanation for the changed reachability.

The first attempted replacement still used loopback-bound guest source. Inspecting source and the running executable path revealed the missed copy; a later filename typo and client-location mistake were also recorded. The successful sequence required copying the correct file, building it, starting the replacement process, and checking actual listener/response behavior rather than relying on a reported restart.

Guest-local checks established 200 for the health path, 404 for an unknown path, and 405 for POST to the health path. The wider-binding checks established HTTP 200 and the expected JSON fields on guest loopback and Optiplex-to-private-IP paths. These are bounded correctness/reachability observations, not capacity or reliability measurements.

## Trade-offs and limitations

A stateless endpoint isolates HTTP behavior and network reachability. It cannot demonstrate persistence, transactions, sessions, replica behavior, resilience, or capacity. Foreground execution does not establish service supervision or reboot recovery.

## Conclusion and next question

Initial scope agreed on 2026-11-04. Guest discovery and Go availability are recorded. Handler and loopback server source were reviewed, and the learner formatted and built the guest copy. Guest-local JSON and status outcomes 200/404/405 were observed. Private-IP requests failed with the loopback-bound process. After correcting a stale guest source copy, rebuilding, and restarting with the all-interface binding, both guest-loopback and Optiplex-to-private-IP health requests returned HTTP 200 with the intended JSON. This supports the research question for the tested paths and contract checks. Formatted source is synchronized back to the repository, and the learner completed the guided copy/build/restart reflection. This first experiment checkpoint is complete; Stage 01 is not complete. The API was last observed running in the foreground, and no shutdown was requested. A proposed next question is how service supervision changes process startup and lifecycle; agree its scope before proceeding. No latency, throughput, resource-usage, reboot recovery, or Windows-access conclusion is included.
