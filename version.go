package cobblerclient

type ExtendedVersion struct {
	Gitdate      string
	Gitstamp     string
	Builddate    string
	Version      string
	VersionTuple []int
}

type CobblerVersion struct {
	Major int
	Minor int
	Patch int
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
