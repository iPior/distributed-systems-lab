resource "libvirt_network" "lab" {
  name      = "dsl-tf00-net"
  mode      = "nat"
  bridge    = "virbr-tf00"
  addresses = ["192.168.78.0/24"]
  autostart = false
  dhcp {
    enabled = true
  }
  dns {
    enabled = true
  }
}