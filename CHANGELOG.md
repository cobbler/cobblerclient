# Changelog

The main changelog for the library versions can be found at <https://github.com/cobbler/cobblerclient>. This file is
aiming to provide a logical overview about compatibility with the Cobbler server.

## Cobbler 4.0.x support

v1.0.0 of this client targets Cobbler 4.0.0 only. The 4.0.0 release dropped several
item types and reshaped network interface storage; this client follows. Users still on
Cobbler 3.3.x should stay on the v0.5.x line, which continues to receive bug fixes.

### Added

* `NetworkInterface` as a first-class item type with full CRUD (replaces the
  inline interface map on `System` for writes). Nested `IPv4Option`,
  `IPv6Option`, `DNSInterfaceOption` value structs model the per-interface
  layer-3 and DNS configuration.
* `Template` as a first-class item type for autoinstall templates, with
  `URIOption` (file / importlib / environment schemas), `GetTemplateContent`,
  `TemplatesRefreshContent`, and the background variant.
* `DistroGroup`, `ProfileGroup`, `SystemGroup` item types for bulk operations.
* `TransactionBegin`, `TransactionCommit`, `TransactionAbort` for atomic
  multi-step modifications.
* `DumpVars` replaces the deprecated `GetBlendedData`.
* `Input*` helpers: `InputBoolean`, `InputInt`, `InputStringOrList`,
  `InputStringOrListNoInherit`, `InputStringOrDict`, `InputStringOrDictNoInherit`.
* `SetItemResolvedValue`; `GetItemResolvedValue` now takes the attribute path
  as `[]string` and returns the resolved value.
* `GetRandomMacFor(virtType)` overload (the bare `GetRandomMac()` now pins
  `virt_type=kvm` per the 4.0.0 backend default).

### Changed

* **Minimum Cobbler server: 4.0.0.** Cobbler 3.3.x users keep using
  `v0.5.x`.
* **Minimum Go: 1.22** (was 1.18).
* `GetRandomMac()` returns `(string, error)` and the backend default
  `virt_type` flipped from `qemu` to `kvm`.
* The following methods previously returned only `error` and discarded the
  backend's actual result; they now return `(T, error)`:
  `AutoAddRepos`, `IsAutoinstallInUse`, `RegisterNewSystem`,
  `RunInstallTriggers`, `GetReposCompatibleWithProfile`,
  `FindSystemByDnsName`, `GetRandomMac`, `GetItemResolvedValue`.
* `getConcreteItem` no longer branches on `CachedVersion` for the `resolved`
  parameter — the post-3.3.3 wire shape is always used.

### Removed

* `MgmtClass`, `LinuxPackage`, `File`, `Snippet`, `TemplateFile` item types
  (no longer exposed by the 4.0.0 backend).
* `read_autoinstall_template` / `write_autoinstall_template` /
  `remove_autoinstall_template` / `get_autoinstall_templates` /
  `get_autoinstall_snippets` wrappers.
* Client-level interface-mutation shims `ModifyInterface`,
  `DeleteNetworkInterface(systemID, name)`, `RenameNetworkInterface(...)`
  (replaced by top-level `NetworkInterface` CRUD).
* `github.com/fatih/structs` dependency.

### Deprecated

* `GetBlendedData` — use `DumpVars` instead. `get_blended_data` is
  soft-deprecated in the Cobbler backend as of 4.0.0 (still available
  server-side); kept here as a transitional wrapper.

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

* `read_autoinstall_template`
* `write_autoinstall_template`
* `remove_autoinstall_template`

Function `read_or_write_snippet` was replaced with:

* `read_autoinstall_snippet`
* `write_autoinstall_snippet`
* `remove_autoinstall_snippet`

### Other changes

Template names used are now **short names** without a path.\
So `foo.ks` instead of `/var/lib/cobbler/kickstarts/foo.ks`.

#### Renamed

These attributes are renamed in Cobbler 3:

* `kickstart` to `autoinstall`
* `ks_meta` to `autoinstall_meta`, but it is still used as a "legacy field"

These directories have been renamed:

* `/var/www/cobbler/ks_mirror` to `/var/www/cobbler/distro_mirror`
* `/var/lib/cobbler/kickstarts` to `/var/lib/cobbler/templates`

The storage locations for the json files changed from `/var/lib/cobbler/config/{distros,profiles,systems,etc...}.d` to `/var/lib/cobbler/collections/{distros,profiles,systems,etc...}`.

There is being worked on a script to migrate these: `scripts/migrate-data-v2-to-v3.py`.

#### Added

These fields have been added:

* `boot_loader` - must be either `grub`, `pxe`, or `ipxe`

#### Removed

Support for these attributes was dropped in Cobbler 3:

* `ldap_enabled`
* `ldap_type`
* `monit_enabled`
* `redhat_management_server`
