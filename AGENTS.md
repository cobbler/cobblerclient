# AGENTS.md

This file provides guidance to Claude Code (claude.ai/code) and other AI coding agents when working with code in
this repository.

## Project overview

`cobblerclient` is a Go XML-RPC client for [Cobbler](https://github.com/cobbler/cobbler), a Linux provisioning/PXE
server. It's consumed by the [Cobbler CLI](https://github.com/cobbler/cli/) and the
[Terraform provider](https://github.com/cobbler/terraform-provider-cobbler).

This branch (`v1.x`) targets **Cobbler 4.0.x only** and requires **Go 1.22+**. The `v0.5.x` line targets Cobbler
3.3.x and continues to receive bug fixes independently — see `MIGRATION.md` and `CHANGELOG.md` before assuming a
4.0.0-era API applies to older servers.

## Commands

```sh
go build .                 # build the library
go test -v .                # run the full test suite (also: make test)
go test -run TestName -v .  # run a single test
go vet ./...
gofmt -l .                  # list files needing formatting; gofmt -w . to fix
go run ./cmd                # regenerate fixtures against a live Cobbler server (also: make fixture-refresh)
```

CI (`.github/workflows/test.yml`) runs `go test -race -coverprofile="coverage.out" -covermode="atomic" -v ./...`
across Go 1.22–1.26 on Linux/macOS/Windows. There is no linter config beyond `go vet`/`gofmt`.

## Architecture

### Request flow: Item → Client.Call → XML-RPC

Every item type (`Distro`, `Profile`, `System`, `Repo`, `Image`, `Menu`, `Template`, `NetworkInterface`, and the
`DistroGroup`/`ProfileGroup`/`SystemGroup` group types) embeds `Item` (`item.go`), which carries the fields common
to every Cobbler object (`Name`, `Comment`, `KernelOptions`, `Owners`, ...). Concrete per-type files (`distro.go`,
`profile.go`, ...) add their own fields and `Create*`/`Update*`/`Get*`/`Find*` methods, but the actual field-by-field
write path is generic and lives in `cobblerclient.go`:

- `Client.updateCobblerFields` — entry point called by every `Create*`/`Update*`, handles the special ordering
  Cobbler requires (e.g. a profile's `name`/`parent`/`distro` must be set before other fields can inherit from them).
- `Client.updateFields` — the recursive field walker. It reflects over a struct's fields and calls
  `modify_<type>(id, attributePath, value, token)` for each one. `attributePath` is a `[]string`, not a flat string,
  because Cobbler 4.0.0 nests some fields under "option" objects (see below) — a nested field's path is
  `[]string{"virt", "cpus"}`, not `"virt_cpus"`.
- `isOptionStructType` / `item_options.go` — `VirtOptions`, `PowerOptions`, `DNSOptions`, `TFTPOptions`,
  `APTOptions` (plus `IPv4Option`/`IPv6Option`/`DNSInterfaceOption`/`URIOption` in `network_interface_types.go` and
  `template_types.go`) are nested structs that multiple item types embed as a named field (e.g. `Profile.Virt
  VirtOptions`). `updateFields` recurses into them, extending the attribute path by one segment. Adding a new such
  struct requires registering it in `isOptionStructType` *and* adding a matching key-set check to `cobblerDataHacks`
  (see below) — the two are easy to forget independently, and only the second one shows up as a runtime decode bug
  rather than a compile error.
- `Client.Call` wraps `github.com/kolo/xmlrpc`. It sends whatever Go value you give it, and the wire type is chosen
  purely from `reflect.Kind` — an `int`-backed Go type is always encoded as `<int>`, never `<string>`, regardless of
  named constants. This is why every Cobbler string-enum type in this client (`NetworkInterfaceType`, `VirtType`,
  `PowerType`, `TemplateSchema`, ...) is declared as a plain `string` alias rather than an `int`-backed enum with a
  `String()` method — the latter serializes wrong and Cobbler's setters reject it.

### Decoding responses: `decodeCobblerItem` / `cobblerDataHacks`

Incoming XML-RPC responses are decoded with `mapstructure`, driven by a custom `DecodeHookFunc` in
`cobblerDataHacks` (`cobblerclient.go`). Notable non-obvious behavior:

- `Value[T]` (`item.go`) wraps any attribute that can be either a concrete value or the literal string
  `"<<inherit>>"`. `IsInherited` is the sentinel; `Data`/`FlattenedValue` are only populated when `IsInherited` is
  false. Most `Value[T]` fields need a `sanitizeValue*Struct` call in the type's `convertRaw*` function to populate
  `Data` from `RawData` — forgetting this call is a real, easy-to-miss bug (it did happen: see the DNS decode fix
  in git history) since the struct still decodes without error, just with a permanently empty `Data`.
- The hook recognizes a handful of nested-struct shapes by their exact key set (`matchesKeySet`) — this is how it
  tells "an inheritable map value" apart from "an `IPv4Option`" apart from "a `DNSOptions`" apart from "the
  top-level item" (detected by the presence of a `uid` key). Any new nested option struct needs its own key-set
  branch here, independent of the `isOptionStructType` list used on the write path.
- Cobbler's `"~"` is its wire representation of `None`; it gets converted to the appropriate zero value per target
  kind.

### Fixture-based testing

There is no live Cobbler server in tests. `testing.go`'s `StubHTTPClient` replays canned XML-RPC exchanges: each
`createStubHTTPClient(t, []string{"name1", "name2", ...})` call queues up `fixtures/name-req.xml` /
`fixtures/name-res.xml` pairs, one per expected HTTP round-trip, in order. `Post()` asserts the actual outgoing
request matches the `-req.xml` fixture (after sorting struct members, since Go's map iteration order is
randomized) and returns the corresponding `-res.xml` as the response. `createStubHTTPClientSingle` is the
one-fixture-pair shorthand.

This means adding or changing a `Create*`/`Update*` flow requires fixture files matching the *exact* sequence and
content of RPC calls `updateCobblerFields` will make — including one call per non-empty struct field, in Go
struct-declaration order. `cmd/main.go` (see below) is the tool for generating these against a real server; when
that's not available, a throwaway test using a custom `HTTPClient` that logs/dumps requests instead of asserting on
them is the fastest way to get the real call sequence and hand-build fixtures from it.

### `cmd/main.go`: fixture recorder, not shipped code

`cmd/main.go` is a standalone dev tool (not imported by the library) that exercises the entire client against a
real Cobbler server and records every request/response pair into `fixtures/`, normalizing the session token to
`securetoken99` along the way. It is run via `make fixture-refresh` / `go run ./cmd` and is the source of truth for
`fixtures/*.xml` — hand-edited fixtures should stay consistent with what this tool would actually produce against a
live server for the same call sequence, including current-Go-struct field-declaration order.

### Generic item operations

Beyond the per-type `Create*`/`Update*`/`Get*` methods, `item.go` exposes generic operations that work across any
item type by name (`"distro"`, `"profile"`, `"system_group"`, ...): `GetItem`, `FindItems`/`FindItemNames`,
`ModifyItem`, `SaveItem`, `RemoveItem`, `CopyItem`, `RenameItem`, `GetItemHandle`, `GetItemResolvedValue` /
`SetItemResolvedValue`. These are what the group types (`groups.go`) and some per-type bypass paths (e.g.
`CreateDistroGroup` calling `SaveItem` directly instead of a dedicated `SaveDistroGroup`) build on; check whether a
fixture's `methodName` is the generic `save_item`/`modify_item`/... or a type-specific `save_<type>` before assuming
which code path produced it.
