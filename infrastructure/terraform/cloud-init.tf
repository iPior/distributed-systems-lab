resource "libvirt_cloudinit_disk" "seed" {
  name           = "dsl-node-01-seed.iso"
  pool           = libvirt_pool.lab.name
  user_data      = file("${path.module}/../cloud-init/user-data.yaml")
  network_config = file("${path.module}/../cloud-init/network-config.yaml")
  meta_data      = file("${path.module}/../cloud-init/meta-data.yaml")
}