# Experiment 01h: Direct Windows access through a LAN bridge

**Status:** connectivity checkpoint complete; LAN attachment and host checks passed, direct Windows health access confirmed by learner

**Architecture stage:** existing single-API configuration; Stage 02 remains pending

**Source revision / configuration:** `fe72b64` with existing uncommitted experiment notes and proposed `infrastructure/terraform/vm.tf` NIC edit

## Concept and research question

Can Windows reach the existing API directly on a guest LAN address, removing SSH forwarding while keeping load generation off the shared VM host? This checkpoint tests connectivity, not throughput.

## Hypothesis

Learner predicted that removing the SSH tunnel could both improve throughput and avoid its observed accept/file-resource limit. These performance predictions remain untested. Learner selected a VM LAN address instead of routed access to the existing private subnet.

## Architecture and controls

```text
Windows 192.168.2.18 -> home LAN -> lab VM DHCP address 192.168.2.x:8080
                                         |
                                 existing host bridge br0
```

- Fresh host inspection verified `br0` at `192.168.2.103`, with physical Ethernet `enp0s31f6` and the existing HAOS tap attached. Before activation, the lab VM used private `dsl-tf00-net` at `192.168.78.136`; after activation it uses LAN DHCP address `192.168.2.44`.
- Proposed scope: move the single existing lab NIC from `dsl-tf00-net` to `br0`. Preserve MAC `52:54:00:56:8a:8a`, PCI position, VM UUID `bc9dc19c-d81b-4df5-85bf-0a94ac7b2be3`, disk, CPU, RAM, and existing lab network definition.
- Guest interface `ens3` obtained `192.168.2.44/24` through DHCP, with default route via home router `192.168.2.1`. A router-side reservation has not been configured; rediscover the address after a break.
- API listens on `0.0.0.0:8080`; bridging makes it reachable to LAN peers if firewall rules allow it. Existing SSH service would also join the LAN. No router port forwarding is proposed.
- Windows learner disconnected ProtonVPN. Before disconnection, route selection for the private guest used the VPN; afterward it used Ethernet via home router `192.168.2.1`. VPN state for historical load runs is unverified. Any later performance comparison must record this difference.
- Host firewall and media application health remain unverified due to elevated-access requirements. No host bridge, firewall, or unrelated service reconfiguration is proposed.

## Procedure

### Preparation completed

At the learner's request, edited `vm.tf` to use `bridge = "br0"` and explicitly preserve the MAC. Terraform formatting, validation, and `git diff --check` passed.

Saved a normal Terraform plan in ignored `.local/lan-bridge-review-20261117/bridge.tfplan`. It reports exactly `0 to add, 1 to change, 0 to destroy`, limited to `libvirt_domain.node`. JSON inspection shows addition of `bridge = "br0"` while the computed old network ID/name remain in the proposed state.

**Do not apply that plan to perform the migration.** Inspection of the pinned libvirt provider v0.8.3 `resourceLibvirtDomainUpdate` implementation found no operation to change an existing NIC's network source. The in-place plan is not evidence that the source will actually move. Source: [provider update implementation](https://github.com/dmacvicar/terraform-provider-libvirt/blob/v0.8.3/libvirt/resource_libvirt_domain.go).

Saved the inactive domain XML and generated proposed LAN XML under `.local/lan-bridge-review-20261117/`. A programmatic comparison verified the only semantic differences are the NIC type (`network` to `bridge`) and its source (`dsl-tf00-net` to `br0`). All other XML, including MAC, UUID, disk and PCI location, is retained. `virt-xml-validate ... domain` passed for the proposed XML. Schema validation does not establish runtime connectivity.

Pre-activation API was active, PID 638, InvocationID `63183976a1ee4367ab2a5328a73eaaae`; its unit is static and must be started after reboot. Fresh authenticated guest ED25519 host-key fingerprint: `SHA256:DpTvDviMG8g9QZxziTBpDmPArsFiky0fHFMcNxw5itc`.

### Reviewed activation procedure

Learner explicitly approved the brief lab-VM outage and LAN attachment, including loss of process-local task data, before execution. The reviewed procedure was:

1. Verify current lab identity and inactive XML still match the reviewed baseline; capture available host routes, domain/interface state, and service-preservation evidence. Stop on unexpected changes.
2. Gracefully shut down only `dsl-node-01` and wait for `shut off`. If shutdown fails, stop for review; do not force-stop it automatically.
3. Define the schema-validated proposed domain XML using the same VM name and UUID. Start that existing VM. No disk creation, deletion or host bridge restart is involved.
4. Discover its DHCP LAN address. The home router's lease table can identify MAC `52:54:00:56:8a:8a`; host neighbor discovery may also reveal it. Libvirt's old private-network lease is not evidence of the new LAN address, and the guest agent is currently disconnected. If address discovery or connectivity cannot be established, use the rollback below rather than expanding the change scope.
5. Verify SSH identity using the previously trusted guest key, with strict checking (the old address can be used as `HostKeyAlias` against the existing known-hosts file). Verify hostname, UUID/NIC attachment from libvirt, and guest address/default route. Start only `dsl-api.service` in the guest and verify `/healthz` over direct HTTP. Update Ansible inventory with the verified address.
6. Preview a Terraform refresh-only plan and, after reviewing it, record the observed NIC state. Require a subsequent ordinary plan to show no infrastructure changes. Do not execute the obsolete saved migration plan.
7. Learner tests Windows TCP port 8080 and `GET /healthz` directly. Record Windows selected interface/route and VPN state. Compare available host/unrelated-service evidence with the baseline and record any inspection limitations. Stop before k6 or Stage 02 provisioning.

### Proposed rollback

Gracefully stop only the lab VM, restore `.local/lan-bridge-review-20261117/dsl-node-01-before.xml` with the same UUID, and start it. Rediscover its private DHCP address rather than assuming the old lease. Restore the Terraform NIC reference to `libvirt_network.lab.id`, reconcile inventory/recorded state, start the static API service, and verify health. No disk rollback or broad firewall/network cleanup is required by this proposal.

## Observations and evidence

Activation completed on 2026-11-17 America/Toronto:

- Immediate preflight confirmed the inactive domain still matched the reviewed rollback XML exactly. Graceful ACPI shutdown reached `shut off`; no force-stop was used.
- Defined the schema-validated LAN XML and started the same VM. UUID, MAC, disk and resource allocation were preserved. Live NIC is now attached to `br0` (tap observed as `vnet3`).
- Discovered the new IPv4 address through authenticated SSH to the guest's existing IPv6 link-local address `fe80::5054:ff:fe56:8a8a%br0`, using strict host-key checking against the previously trusted key through `HostKeyAlias=192.168.78.136`. No guest agent or network package was installed.
- Guest reported hostname `dsl-node-01`, DHCP address `192.168.2.44/24` on `ens3`, and default route via `192.168.2.1`.
- Static API unit was inactive after boot, as expected. Started only `dsl-api.service`. New service identity: PID 635, InvocationID `970322a9bd1f4832bbeb6de11f8c27c4`, active/running. Executable SHA-256 remained `22fd6c687a23dad2a980484742c5a01e62999b16a0adc3088e3f2839bdd990af`.
- Bound the new LAN address to the previously authenticated guest public host key in the ignored known-hosts file. Updated Ansible inventory to `192.168.2.44`. Ansible ping succeeded with `changed=false`, using ordinary strict host-key checking; Python interpreter-discovery warning was non-fatal.
- Host direct GET `http://192.168.2.44:8080/healthz` returned HTTP 200 with `{"instance":"dsl-node-01","status":"ok"}`. Saved response Date: 2026-11-18 01:59:24 GMT.
- Saved a state backup and reviewed a refresh-only plan. It contained no infrastructure actions. Recorded drift was limited to the lab NIC, normalization of `cmdline` to an empty list, and computed pool allocation/available-space counters. Applied only this refresh-only plan: zero added, changed or destroyed resources. A subsequent ordinary Terraform plan reported **no changes**, and validation passed. The original saved migration plan was not applied.
- HAOS and devbox remained running, with byte-equivalent parsed persistent XML before/after. Host IPv4 routes differed only by the expected `linkdown` flag on `virbr-tf00`, now lacking an attached guest. The old private network remains active and persistent, with autostart disabled. No host bridge or firewall configuration was edited. Docker/media application health and firewall rules remain unverified due to elevated-access limitations; a guest-side nftables read also found no `nft` command.

Learner was given this Windows PowerShell check and subsequently confirmed success:

```powershell
curl.exe --noproxy "*" --max-time 5 http://192.168.2.44:8080/healthz
```

The expected response supplied with that instruction was `{"instance":"dsl-node-01","status":"ok"}`. This is learner-confirmed success, without pasted raw Windows output or a new route trace. It supports direct Windows HTTP access, not a performance conclusion. No k6 run followed activation. Rollback XML, proposed XML, plans, state backup, route/domain comparisons and guest/API evidence are local ignored artifacts in `.local/lan-bridge-review-20261117/`; they do not travel with a Git checkout.

## Analysis

An existing physical LAN bridge provides a way to give the VM a LAN address without rebuilding host networking. The provider limitation requires a separately reviewed activation step and subsequent Terraform state reconciliation. Performance remains unmeasured on this path.

## Trade-offs and limitations

Bridging changes the lab VM's network isolation and exposes its listening services to LAN peers, subject to firewall policy. Guest DHCP and host direct access are verified; direct Windows health access is learner-confirmed without raw output. A VM restart loses task memory and requires an explicit API service start. No media application health claim is supported by the current evidence.

## Conclusion and next question

Connectivity checkpoint complete: the existing VM is attached to the LAN, host checks pass, and learner confirms direct Windows health access at `192.168.2.44`. Stop before performance testing or HAProxy/replica work. Agree the next focused experiment before proceeding. Historical SSH/host-local measurements remain separate from any later direct-LAN measurements; no direct-LAN load result exists yet.
