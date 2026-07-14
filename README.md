[![Go Report Card](https://goreportcard.com/badge/github.com/cobbler/cobblerclient)](https://goreportcard.com/report/github.com/cobbler/cobblerclient)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/2bd1196e7ee7427b894689ca47d4e170)](https://app.codacy.com/gh/cobbler/cobblerclient/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/2bd1196e7ee7427b894689ca47d4e170)](https://app.codacy.com/gh/cobbler/cobblerclient/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
[![Go Reference](https://pkg.go.dev/badge/github.com/cobbler/cobblerclient.svg)](https://pkg.go.dev/github.com/cobbler/cobblerclient)

# cobblerclient

Cobbler Client written in Go. Used by the [CLI](https://github.com/cobbler/cli/) and by the
[Terraform Provider](https://github.com/cobbler/terraform-provider-cobbler).

## Compatibility

| cobblerclient | Cobbler server | Go    |
|---------------|----------------|-------|
| `v1.x`        | `4.0.x`        | 1.22+ |
| `v0.5.x`      | `3.3.x`        | 1.18+ |

Upgrading from `v0.5.x` to `v1.0.0`? See [MIGRATION.md](./MIGRATION.md) — the
release is a clean break that drops support for Cobbler 3.x and reshapes
network interface and autoinstall template handling.

For more details please see:

- the [AUTHORS](./AUTHORS.md) file,
- the [Changelog](./CHANGELOG.md) for the Server compatibility or
- the [changelog](https://github.com/cobbler/cobblerclient/releases) for this library.
