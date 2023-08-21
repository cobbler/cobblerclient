package item

import (
	"github.com/cobbler/cobblerclient/client"
	internalItem "github.com/cobbler/cobblerclient/internal/item"
	rawItem "github.com/cobbler/cobblerclient/internal/item/raw"
	resolvedItem "github.com/cobbler/cobblerclient/internal/item/resolved"
)

// Profile is a created profile.
type Profile struct {
	// Internal fields
	raw      *rawItem.Profile
	resolved *resolvedItem.Profile

	Item

	Autoinstall       internalItem.InheritableProperty[string]
	DHCPTag           internalItem.Property[string]
	Distro            internalItem.Property[string]
	EnableIPXE        internalItem.InheritableProperty[bool]
	EnableMenu        internalItem.InheritableProperty[bool]
	NameServers       internalItem.InheritableProperty[[]string]
	NameServersSearch internalItem.InheritableProperty[[]string]
	NextServerv4      internalItem.InheritableProperty[string]
	NextServerv6      internalItem.InheritableProperty[string]
	Parent            internalItem.Property[string]
	Proxy             internalItem.InheritableProperty[string]
	Repos             internalItem.Property[[]string]
	Server            internalItem.InheritableProperty[string]
	VirtAutoBoot      internalItem.InheritableProperty[string]
	VirtBridge        internalItem.InheritableProperty[string]
	VirtCPUs          internalItem.InheritableProperty[string]
	VirtDiskDriver    internalItem.Property[string]
	VirtFileSize      internalItem.InheritableProperty[string]
	VirtPath          internalItem.Property[string]
	VirtRAM           internalItem.InheritableProperty[string]
	VirtType          internalItem.Property[string]
}

func BuildProfile(client client.Client) Profile {
	var i = BuildItem(client)
	var p = Profile{
		raw: &rawItem.Profile{
			Item: *i.raw,
		},
		resolved: &resolvedItem.Profile{
			Item: *i.resolved,
		},
		Item: i,
	}
	refreshProfilePointers(&p)

	return p
}

func refreshProfilePointers(profile *Profile) {
	profile.Item.raw = &profile.raw.Item
	profile.Item.resolved = &profile.resolved.Item
	refreshItemPointers(&profile.Item)
	// FIXME
}
