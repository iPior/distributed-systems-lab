# Experiment 01: A single Go HTTP API

**Status:** paused before formatting, build, and HTTP verification, 2026-11-05

**Architecture stage:** 01

**Source revision / configuration:** planning begins at `ef28b24` on `learn/rebuild-experiment-00`; module file and unfinished handler scaffold added during this session

## Resume here

- Module 00 is complete. Stage 01 was requested on 2026-10-28.
- Read repository guidance, README, roadmap, Module 00 overview, prior learning-run evidence, and experiment template. Available curriculum covers only Module 00; no Stage 01 curriculum exists yet.
- Continued on `learn/module-01-go-api`, branched from `927360b`; `learn/rebuild-experiment-00` is retained as the Module 00 legacy branch. Paused on 2026-11-05 because the learner cannot currently access the guest from their personal device.
- Learner implemented the JSON health response, encoding-error logging, `GET /healthz` route, and a loopback listener at `127.0.0.1:8080` with server-error logging in `api/main.go`. Source review completed; formatting, Go validation/build, deployment, and HTTP checks have not run.
- Resume by ensuring the guest has the latest `api/` source, then run `gofmt -w main.go` and `go build -o health-api .` from its API directory. Review the outcome before starting the server. The learner correctly predicted that a guest-loopback listener cannot accept an Optiplex request addressed to the guest's private IP. Verify that prediction with bounded requests after a successful build; do not change the binding yet.
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

## Analysis

Pending real evidence.

## Trade-offs and limitations

A stateless endpoint isolates HTTP behavior and network reachability. It cannot demonstrate persistence, transactions, sessions, replica behavior, resilience, or capacity. Foreground execution does not establish service supervision or reboot recovery.

## Conclusion and next question

Initial scope agreed on 2026-11-04. Guest discovery and learner-reported Go availability are recorded. Handler and loopback server source were reviewed on 2026-11-05; formatting, build, HTTP verification, and final reflection remain pending. No Stage 01 API runtime behavior has been measured.
