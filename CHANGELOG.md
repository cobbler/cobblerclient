# Changelog

The main changelog for the library versions can be found at <https://github.com/cobbler/cobblerclient>. This file is
aiming to provide a logical overview about compatibility with the Cobbler server.

## Cobbler 3.3.x support
v0.5.0 of this client introduced support for Cobbler v3.3.0, which was a refactor from runtime-created Python
attributes to Python Properties.  For further details see
[release notes](https://github.com/cobbler/cobbler/releases/tag/v3.3.0).  Breaking changes:
* This client's support for earlier Cobbler versions was dropped.
* next_server attribute is now either next_server_v4 or next_server_v6
* boot_loader string attribute is now boot_loaders list
* The following string attributes are now lists: FetchableFiles, KernelOptions, KernelOptionsPost, TemplateFiles,
  AutoinstallMeta, Repos

## Cobbler <=3.2.x support
Retaining the below notes for the time-being, which only apply to v0.4.2 and earlier clients:

[Cobbler](https://github.com/cobbler/cobbler) (up to version 2.8.x) was written in Python2.
However, Python2 is EOL since January 2020.\
Cobbler 3 has been adapted to use Python3 and so lots of code changed. Sadly this also broke
backward compatability with the original `cobblerclient`.

### XMLRPC API changes

Function `read_or_write_kickstart_template` was replaced with:

- `read_autoinstall_template`
- `write_autoinstall_template`
- `remove_autoinstall_template`

Function `read_or_write_snippet` was replaced with:

- `read_autoinstall_snippet`
- `write_autoinstall_snippet`
- `remove_autoinstall_snippet`

### Other changes

Template names used are now **short names** without a path.\
So `foo.ks` instead of `/var/lib/cobbler/kickstarts/foo.ks`.

#### Renamed

These attributes are renamed in Cobbler 3:

- `kickstart` to `autoinstall`
- `ks_meta` to `autoinstall_meta`, but it is still used as a "legacy field"

These directories have been renamed:

- `/var/www/cobbler/ks_mirror` to `/var/www/cobbler/distro_mirror`
- `/var/lib/cobbler/kickstarts` to `/var/lib/cobbler/templates`

The storage locations for the json files changed from `/var/lib/cobbler/config/{distros,profiles,systems,etc...}.d` to `/var/lib/cobbler/collections/{distros,profiles,systems,etc...}`.

There is being worked on a script to migrate these: `scripts/migrate-data-v2-to-v3.py`.

#### Added

These fields have been added:

- `boot_loader` - must be either `grub`, `pxe`, or `ipxe`

#### Removed

Support for these attributes was dropped in Cobbler 3:

- `ldap_enabled`
- `ldap_type`
- `monit_enabled`
- `redhat_management_server`
