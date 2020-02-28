
# Cobbler 3 fields and types

## Distro

 Field| Type |  Remarks |
------|------|--------|
arch |  string |
autoinstall_meta | dict |
breed | string |
boot_files | list |
comment | string |
fetchable_files | list |
kernel | string |
kernel_options | dict |
kernel_options_post | dict |
initrd | string |
mgmt_classes | list |
name | string |
os_version | string |
owners | list |
template_files | list | ⚠️ Conflicts with profile & system

## Profile

 Field| Type |  Remarks |
------|------|--------|
autoinstall | string |
autoinstall_meta | dict |
boot_files | list |
comment | string |
dhcp_tag | string |
distro | string |
enable_gpxe | bool |
enable_menu | bool |
fetchable_files | list |
kernel_options | dict
kernel_options_post | dict |
mgmt_classes | list |
mgmt_parameters | string |
name | string |
name_servers | list |
name_servers_search | list |
next_server | string |
owners | list |
proxy | string |
repos | list |
server | string |
template_files | dict | ⚠️ Conflicts with distro
virt_auto_boot | bool |
virt_bridge | string |
virt_cpus | int |
virt_disk_driver | string |
virt_file_size | int | ⚠️ Conflicts with system
virt_path | string |
virt_ram | int |
virt_type | string |

## System

 Field| Type |  Remarks |
------|------|--------|
autoinstall | string |
autoinstall_meta | dict |
boot_files | list |
comment | string |
enable_gpxe | bool |
fetchable_files | dict |
gateway | string |
hostname | string |
image | string |
ipv6_default_device | string |
kernel_options | dict |
kernel_options_post | dict |
mgmt_classes | list |
mgmt_parameters | string |
name | string |
name_servers | list |
name_servers_search | list |
netboot_enabled | bool |
next_server | string |
owners | list |
power_address | string |
power_id | string |
power_pass | string |
power_type | string |
power_user | string |
profile | string |
proxy | string |
status | string |
template_files | dict | ⚠️ Conflicts with distro
virt_auto_boot | bool |
virt_cpus | int |
virt_disk_driver | string |
virt_file_size | float | ⚠️ Conflicts with profile
virt_path | string |
virt_pxe_boot | bool |
virt_ram | int |
virt_type | string |

## Repo

 Field| Type |  Remarks |
------|------|--------|
apt_components | list |
apt_dists | list |
arch | string |
breed | string |
comment | string |
createrepo_flags | dict |
environment | dict |
keep_updated | bool |
mirror | string |
mirror_locally | bool |
name | string |
owners | list |
proxy | string |
rpm_list | list |

