# Migration Guide

## v0.5.x → v1.0.0 (Cobbler 4.0.x)

`v1.0.0` is a clean break that targets Cobbler 4.0.x only. The `v0.5.x` line
continues to serve Cobbler 3.3.x; pick the line that matches your server.

| cobblerclient | Cobbler server | Go    |
|---------------|----------------|-------|
| `v1.x`        | `4.0.x`        | 1.22+ |
| `v0.5.x`      | `3.3.x`        | 1.18+ |

### Item types you can no longer use

These types were removed from the Cobbler 4.0.0 server and are gone from the
client as well:

- `MgmtClass`
- `LinuxPackage`
- `File`
- `Snippet`
- `TemplateFile`

If your code referenced any of them, you need to migrate that data into one
of the remaining item types (typically `Template` for the template-shaped
ones, or `KernelOptions` / `BootFiles` / `TemplateFiles` per-item maps for
the file-shaped ones) **before** upgrading.

### Network interfaces are now their own thing

In `v0.5.x`, network interfaces lived inside `System.Interfaces` as a
`map[string]Interface` and were written through ad-hoc `modify_system`
pseudo-attributes (`Client.ModifyInterface`, `System.CreateInterface`,
`System.DeleteInterface`, `System.RenameInterface`).

In `v1.0.0`, `NetworkInterface` is a first-class item type with full CRUD on
`Client`. To replace the old API:

```go
// before (v0.5.x)
iface := cobblerclient.NewInterface()
iface.MACAddress = "aa:bb:cc:dd:ee:ff"
iface.IPAddress = "10.0.0.5"
sys.CreateInterface("eth0", iface)
sys.DeleteInterface("eth1")

// after (v1.0.0)
ni := cobblerclient.NewNetworkInterface()
ni.Name = "eth0@" + sys.Name
ni.MacAddress = "aa:bb:cc:dd:ee:ff"
ni.IPv4.Address = "10.0.0.5"
_, _ = c.CreateNetworkInterface(sys.Uid, ni)
_ = c.DeleteNetworkInterface("eth1@" + sys.Name)
```

`System.Interfaces` still exists as a read-only convenience for inspecting
interfaces attached to a fetched system — it's `map[string]*NetworkInterface`
now (was `map[string]Interface`), populated from the same first-class item
type the CRUD methods above operate on. It's marked `noupdate`: write through
the `NetworkInterface` CRUD methods, not by mutating a fetched `System`.

The per-interface IPv4 / IPv6 / DNS fields are now nested option objects:

```go
ni.IPv4 = cobblerclient.IPv4Option{
    Address: "10.0.0.5",
    Netmask: "255.255.255.0",
}
ni.IfGateway = "10.0.0.1" // gateway lives on NetworkInterface itself, not IPv4Option
ni.IPv6 = cobblerclient.IPv6Option{
    Address: "fe80::1",
    Prefix:  "64",
}
ni.DNS = cobblerclient.DNSInterfaceOption{
    Name:   "server1.example.com",
    CNames: []string{"www"},
}
```

### Autoinstall templates are now items, not files

The file-level wrappers (`CreateTemplateFile`, `GetTemplateFile`,
`DeleteTemplateFile`, `GetAutoinstallTemplates`) are gone — the 4.0.0 backend
removed `read_autoinstall_template` and friends. Use the new `Template`
item type:

```go
tpl := cobblerclient.NewTemplate()
tpl.Name = "preseed-default"
tpl.URI.Schema = cobblerclient.TemplateSchemaFile
tpl.URI.Path = "preseed.j2"
tpl.Content = "..."
_, _ = c.CreateTemplate(tpl)
content, _ := c.GetTemplateContent(tpl.Uid)
```

### Method signatures that changed

Several methods previously returned only `error` and dropped the backend's
result; they now return `(T, error)`. Update call sites accordingly:

| Method                            | Before          | After                                |
|-----------------------------------|-----------------|---------------------------------------|
| `AutoAddRepos`                    | `error`         | `(bool, error)`                      |
| `IsAutoinstallInUse`              | `error`         | `(bool, error)`                      |
| `RegisterNewSystem`               | `error`         | `(int, error)`                       |
| `RunInstallTriggers`              | `error`         | `(bool, error)`                      |
| `GetReposCompatibleWithProfile`   | `error`         | `([]map[string]interface{}, error)`  |
| `FindSystemByDnsName`             | `error`         | `(map[string]interface{}, error)`    |
| `GetRandomMac`                    | `error`         | `(string, error)`                    |
| `GetItemResolvedValue`            | `error`         | `(interface{}, error)`               |

`GetItemResolvedValue` now takes the attribute path as `[]string` (e.g.
`[]string{"ipv4", "address"}`).

`GetRandomMac()` sends `virt_type=kvm` explicitly because the 4.0.0 backend
default flipped from `qemu` to `kvm`. Use `GetRandomMacFor(virtType)` to pin
a different value.

### background_* options now require UIDs (Cobbler 4.0.0b6+)

`BackgroundReposyncOptions.Repos`/`.Only`, `BuildisoOptions.Profiles`/`.Systems`/`.Distro`,
`BackgroundSyncSystemsOptions.Systems`, and `BackgroundPowerSystemOptions.Systems` must now hold
the target items' UIDs, not their names. The corresponding server-side actions
(`background_reposync`, `background_buildiso`, `background_syncsystems`,
`background_power_system`) do a strict uid-keyed lookup; passing a name is not an error you'll
see immediately — the entry is silently skipped and logged as a warning server-side, except for
`BuildisoOptions.Distro`, where an unresolvable value raises a hard error and aborts the task.

Resolve names to UIDs yourself before calling these, e.g. via `Client.GetSystemHandle(name)` /
`Client.GetItemHandle("repo", name)`, the same pattern already required for `System.Profile`,
`System.Image`, `Profile.Distro`, etc.

### More methods now require UIDs (Cobbler 4.0.0b6+)

A second batch of methods, previously overlooked when the above `background_*` batch was fixed,
also requires the target item's UID instead of its name as of Cobbler 4.0.0b6:

- `GetDistroAsRendered`, `GetProfileAsRendered`, `GetSystemAsRendered`, `GetRepoAsRendered`,
  `GetImageAsRendered`, `GetMenuAsRendered`, `GetDistroGroupAsRendered`,
  `GetProfileGroupAsRendered`, `GetSystemGroupAsRendered` — an unresolvable uid silently returns
  an empty map, not an error.
- `GetValidDistroBootLoaders`, `GetValidImageBootLoaders`, `GetValidProfileBootLoaders`,
  `GetValidSystemBootLoaders` — an unresolvable uid returns a single-element error-message
  slice, not an RPC fault.
- `GetBlendedData`, `IsAutoinstallInUse` — an unresolvable uid raises a hard server-side error.
  `IsAutoinstallInUse`'s parameter is, and always was, an autoinstall *template* UID, not a
  system name as its old parameter name implied.
- `GetReposCompatibleWithProfile`, `DisableNetboot` — an unresolvable uid silently returns an
  empty result (`[]`) or `false`, respectively, not an error.
- `GetRepoConfigForProfile`, `GetRepoConfigForSystem`, `GetTemplateFileForProfile`,
  `GetTemplateFileForSystem` — an unresolvable uid returns a `"# object not found: ..."` comment
  string, not an error. These are normally reached anonymously via the `/cblr/svc/...` HTTP
  endpoints, which resolve a name to a uid server-side; direct callers of these Go wrappers must
  resolve the uid themselves.
- `GenerateIPxe`, `GenerateBootCfg`, `GenerateScript` — an unresolvable uid for any of
  `profileUid`/`imageUid`/`systemUid` is treated as if that argument were empty/unset.

As with the `background_*` batch, resolve names to UIDs yourself first, e.g. via
`Client.GetSystemHandle(name)` / `Client.GetItemHandle("<type>", name)`.

### Deprecated but kept

- `GetBlendedData` — use `DumpVars(itemUuid, formattedOutput, removeDicts)`
  instead. `get_blended_data` is soft-deprecated in the Cobbler backend as of
  4.0.0 (still available server-side, not removed); the wrapper survives for
  transitional callers and will be removed in a future version.

### Go version

`go.mod` now requires Go 1.22 (was 1.18). This aligns with the CLI and the
Terraform provider.

### New things worth knowing about

- **Transactions**: wrap atomic batches in `TransactionBegin` /
  `TransactionCommit` / `TransactionAbort`.
- **Groups**: `DistroGroup`, `ProfileGroup`, `SystemGroup` item types with
  full CRUD for bulk-operation patterns.
- **Input helpers**: `InputBoolean`, `InputInt`, the `InputStringOrList` /
  `InputStringOrDict` pair (plus their `NoInherit` variants).
- **Resolved-value writes**: `SetItemResolvedValue(uuid, attrPath, value)`.
