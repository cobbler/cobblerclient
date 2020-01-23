# cobblerclient

Cobbler Client written in Go.

Original authors:

- [Container Solutions](https://www.container-solutions.com/)
- [Joe Topjian](https://github.com/jtopjian)

## Cobbler 3 support

Adapted by [Devhouse Spindle](https://wearespindle.com/) to support Cobbler 3 XMLRPC API calls.

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

Directories renamed:

- `/var/www/cobbler/ks_mirror` to `/var/www/cobbler/distro_mirror`
- `/var/lib/cobbler/kickstarts` to `/var/lib/cobbler/templates`

Template names used are now short names without a path.
So `foo.ks` instead of `/var/lib/cobbler/kickstarts/foo.ks`.

`ks_meta` was renamed to `autoinstall_meta`.
