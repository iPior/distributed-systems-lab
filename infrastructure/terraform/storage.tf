resource "libvirt_pool" "lab" {
  name = "dsl-tf00"
  type = "dir"
  target {
    path = "/var/lib/libvirt/dsl-tf00"
  }
}

resource "libvirt_volume" "debian_base" {
  name   = "debian13-base.qcow2"
  pool   = libvirt_pool.lab.name
  source = "/opt/distributed-systems-lab/.local/debian13.qcow2"
  format = "qcow2"
}

resource "libvirt_volume" "guest" {
  name           = "dsl-node-01.qcow2"
  pool           = libvirt_pool.lab.name
  format         = "qcow2"
  base_volume_id = libvirt_volume.debian_base.id
  size           = 20 * 1024 * 1024 * 1024
}
