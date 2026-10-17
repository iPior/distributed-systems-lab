# Lesson 03 — Guest networking

## Learn
With libvirt NAT, the host provides a private virtual subnet and routes outbound guest traffic; LAN devices cannot normally initiate connections directly to the guest. A bridged guest may be reachable like another LAN device. DHCP assigns the address, while SSH provides remote access. Network changes can affect host firewall rules.

## Predict
Will a Windows k6 client on your LAN be able to reach a private NAT guest directly? What does NAT change and what does it not change?

## Explore
Read current libvirt network definitions, routes, Docker networks and host firewall state *before* changing anything. Draw expected guest ↔ host ↔ LAN paths. Choose an isolated subnet that does not overlap existing routes.

## Verify
After separately approved network setup, record guest address, default route, DNS resolution, host-to-guest SSH, and any intentionally inaccessible path.

## Checkpoint
Explain why a later API might need port forwarding, bridging, or routing for Windows k6. Stop before Terraform.
