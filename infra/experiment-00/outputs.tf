output "vm_ip" {
  value = local.vm_ip
}

output "vm_name" {
  value = libvirt_domain.node.name
}
