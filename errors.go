package cobblerclient

import "fmt"

// UnsupportedServerVersionError is returned when the connected Cobbler server reports a
// version older than the minimum supported by this client (4.0.0).
type UnsupportedServerVersionError struct {
	ServerVersion CobblerVersion
}

func (e *UnsupportedServerVersionError) Error() string {
	return fmt.Sprintf(
		"cobbler server version %s is not supported; this client requires Cobbler 4.0.0 or later",
		e.ServerVersion.String(),
	)
}

// InheritanceUnsupportedError is returned when a field is set to IsInherited=true but the
// connected Cobbler server is older than the version that added inheritance support for it.
type InheritanceUnsupportedError struct {
	Field          string
	ServerVersion  CobblerVersion
	MinimumVersion CobblerVersion
}

func (e *InheritanceUnsupportedError) Error() string {
	return fmt.Sprintf(
		"field %q does not support inheritance on Cobbler %s; inheritance requires %s or later — set an explicit value instead",
		e.Field, e.ServerVersion.String(), e.MinimumVersion.String(),
	)
}
