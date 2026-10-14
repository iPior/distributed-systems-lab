locals {
  vm_name = "dsl-node-01"
  vm_ip   = "192.168.100.10"
}

# This directory must contain only this experiment's managed volumes.
resource "libvirt_pool" "lab" {
  name = "dsl-exp00"
  type = "dir"
  target {
    path = "/var/lib/libvirt/images/dsl-exp00"
  }
}

resource "libvirt_network" "lab" {
  name      = "dsl-exp00"
  mode      = "nat"
  bridge    = "virbr-dsl00"
  addresses = ["192.168.100.0/24"]
  autostart = false

  dhcp {
    enabled = true
  }
  dns {
    enabled = true
  }
}

resource "libvirt_volume" "root" {
  name   = "dsl-node-01-root.qcow2"
  pool   = libvirt_pool.lab.name
  format = "qcow2"
  # Preparation verifies the original Debian image and expands a standalone
  # sparse copy to 20 GiB. The provider takes size from this source image.
  source = abspath("${path.module}/../../.local/debian13-root.qcow2")
}

resource "libvirt_cloudinit_disk" "seed" {
  name = "dsl-node-01-seed.iso"
  pool = libvirt_pool.lab.name
  user_data = "#cloud-config\n${yamlencode({
    hostname         = local.vm_name
    manage_etc_hosts = true
    ssh_pwauth       = false
    disable_root     = true
    users = [{
      name                = "lab"
      shell               = "/bin/bash"
      lock_passwd         = true
      sudo                = "ALL=(ALL) NOPASSWD:ALL"
      ssh_authorized_keys = [trimspace(file("${path.module}/../../.local/id_ed25519.pub"))]
    }]
    growpart      = { mode = "auto", devices = ["/"], ignore_growroot_disabled = false }
    resize_rootfs = true
  })}"
  meta_data = yamlencode({
    "instance-id"    = "dsl-node-01-exp00"
    "local-hostname" = local.vm_name
  })
  network_config = yamlencode({
    version = 2
    ethernets = {
      lab = {
        match = { macaddress = "52:54:00:00:00:01" }
        dhcp4 = true
        dhcp6 = false
      }
    }
  })
}

resource "libvirt_domain" "node" {
  name      = local.vm_name
  memory    = 2048
  vcpu      = 1
  autostart = false
  cloudinit = libvirt_cloudinit_disk.seed.id

  # Keep the provider's canonical volume ID in state, but generate an explicit
  # qcow2 file source for libvirt's AppArmor helper. No host security changes.
  xml {
    xslt = templatefile("${path.module}/disk-file.xsl.tftpl", {
      root_disk_path = libvirt_volume.root.id
    })
  }

  # Override the legacy provider's SPICE default with built-in VNC. No graphics
  # listener is exposed; access is through SSH and the serial console.
  graphics {
    type        = "vnc"
    listen_type = "none"
  }

  disk {
    volume_id = libvirt_volume.root.id
  }
  network_interface {
    network_id     = libvirt_network.lab.id
    hostname       = local.vm_name
    addresses      = [local.vm_ip]
    mac            = "52:54:00:00:00:01"
    wait_for_lease = true
  }
  console {
    type        = "pty"
    target_type = "serial"
    target_port = "0"
  }
}
