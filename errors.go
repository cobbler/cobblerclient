package cobblerclient

import "fmt"

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
