package cobblerclient

import (
	"fmt"
	"strconv"
	"strings"
)

// parseCobblerVersion parses a "MAJOR.MINOR.PATCH" version string into a CobblerVersion.
func parseCobblerVersion(s string) (*CobblerVersion, error) {
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("expected MAJOR.MINOR.PATCH, got %q", s)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version in %q: %w", s, err)
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version in %q: %w", s, err)
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version in %q: %w", s, err)
	}
	return &CobblerVersion{Major: major, Minor: minor, Patch: patch}, nil
}

type ExtendedVersion struct {
	Gitdate      string `json:"gitdate" yaml:"gitdate"`
	Gitstamp     string `json:"gitstamp" yaml:"gitstamp"`
	Builddate    string `json:"builddate" yaml:"builddate"`
	Version      string `json:"version" yaml:"version"`
	VersionTuple []int  `json:"version_tuple" yaml:"version_tuple"`
}

type CobblerVersion struct {
	Major int `json:"major" yaml:"major"`
	Minor int `json:"minor" yaml:"minor"`
	Patch int `json:"patch" yaml:"patch"`
}

func (cv *CobblerVersion) GreaterThan(otherVersion *CobblerVersion) bool {
	if cv.Equal(otherVersion) {
		return false
	}
	if cv.Major > otherVersion.Major {
		return true
	}
	if cv.Major == otherVersion.Major && cv.Minor > otherVersion.Minor {
		return true
	}
	if cv.Major == otherVersion.Major && cv.Minor == otherVersion.Minor && cv.Patch > otherVersion.Patch {
		return true
	}
	return false
}

func (cv *CobblerVersion) LessThan(otherVersion *CobblerVersion) bool {
	if cv.Equal(otherVersion) {
		return false
	}
	return !cv.GreaterThan(otherVersion)
}

func (cv *CobblerVersion) Equal(otherVersion *CobblerVersion) bool {
	return cv.Major == otherVersion.Major && cv.Minor == otherVersion.Minor && cv.Patch == otherVersion.Patch
}

func (cv *CobblerVersion) NotEqual(otherVersion *CobblerVersion) bool {
	return !cv.Equal(otherVersion)
}

func (cv *CobblerVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", cv.Major, cv.Minor, cv.Patch)
}

// Version is a shorter and easier version representation. Normally you want to call [Client.ExtendedVersion].
func (c *Client) Version() (float64, error) {
	res, err := c.Call("version")
	if err != nil {
		return 0, err
	}
	return res.(float64), err
}

// ExtendedVersion returns the version information of the server.
func (c *Client) ExtendedVersion() (ExtendedVersion, error) {
	extendedVersion := ExtendedVersion{}
	data, err := c.Call("extended_version")
	if err != nil {
		return extendedVersion, err
	}
	switch data.(type) {
	case map[string]interface{}:
		data := data.(map[string]interface{})
		var versionTuple, err = returnIntSlice(data["version_tuple"], err)
		if err != nil {
			return extendedVersion, err
		}
		extendedVersion.Version = data["version"].(string)
		extendedVersion.VersionTuple = versionTuple
		extendedVersion.Builddate = data["builddate"].(string)
		extendedVersion.Gitdate = data["gitdate"].(string)
		extendedVersion.Gitstamp = data["gitstamp"].(string)
	default:
		return extendedVersion, err
	}
	return extendedVersion, err
}
