resource "libvirt_domain" "node" {
  name      = "dsl-node-01"
  vcpu      = 1
  memory    = 2048
  running   = true
  autostart = false
  cloudinit = libvirt_cloudinit_disk.seed.id

  disk {
    volume_id = libvirt_volume.guest.id
  }

  network_interface {
    network_id = libvirt_network.lab.id
  }

  console {
    type        = "pty"
    target_type = "serial"
    target_port = "0"
  }

  video {
    type = "vga"
  }

  graphics {
    type           = "vnc"
    listen_type    = "address"
    listen_address = "127.0.0.1"
    autoport       = true
  }

  xml {
    xslt = templatefile("${path.module}/disk-chain.xsl.tftpl", {
      guest_disk_path = libvirt_volume.guest.id
      guest_disk_name = libvirt_volume.guest.name
      pool_name       = libvirt_pool.lab.name
      base_disk_path  = libvirt_volume.debian_base.id
    })
  }
}
