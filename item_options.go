/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

// VirtOptions models the virtualization settings shared by Image, Profile,
// and System items. Get the fields from cobbler/items/options/virt.py.
type VirtOptions struct {
	AutoBoot   Value[bool]    `mapstructure:"auto_boot" json:"auto_boot" yaml:"auto_boot"`
	Cpus       Value[int]     `mapstructure:"cpus" json:"cpus" yaml:"cpus"`
	DiskDriver string         `mapstructure:"disk_driver" json:"disk_driver" yaml:"disk_driver"`
	FileSize   Value[float64] `mapstructure:"file_size" json:"file_size" yaml:"file_size"`
	Path       string         `mapstructure:"path" json:"path" yaml:"path"`
	// PxeBoot is not inheritable in Cobbler (plain property, not InheritableProperty).
	PxeBoot bool       `mapstructure:"pxe_boot" json:"pxe_boot" yaml:"pxe_boot"`
	Ram     Value[int] `mapstructure:"ram" json:"ram" yaml:"ram"`
	Type    string     `mapstructure:"type" json:"type" yaml:"type"`
	// UEFI is not inheritable in Cobbler (plain property, not InheritableProperty). New in Cobbler 4.0.0a8.
	UEFI bool `mapstructure:"uefi" json:"uefi" yaml:"uefi"`
}

func newVirtOptions() VirtOptions {
	return VirtOptions{
		AutoBoot:   Value[bool]{IsInherited: true},
		Cpus:       Value[int]{IsInherited: true},
		DiskDriver: inherit,
		FileSize:   Value[float64]{IsInherited: true},
		// Path can't be "<<inherit>>" for profiles/images (Cobbler rejects it server-side);
		// for systems an empty string is normalized to inherited server-side anyway.
		Path: "",
		Ram:  Value[int]{IsInherited: true},
		Type: inherit,
	}
}

// PowerOptions models the power control settings for a System.
// Get the fields from cobbler/items/options/power.py. None of these fields
// are inheritable.
type PowerOptions struct {
	Type         string `mapstructure:"type" json:"type" yaml:"type"`
	IdentityFile string `mapstructure:"identity_file" json:"identity_file" yaml:"identity_file"`
	Options      string `mapstructure:"options" json:"options" yaml:"options"`
	User         string `mapstructure:"user" json:"user" yaml:"user"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Address      string `mapstructure:"address" json:"address" yaml:"address"`
	ID           string `mapstructure:"id" json:"id" yaml:"id"`
}

// DNSOptions models the DNS settings shared by Profile and System items.
// Get the fields from cobbler/items/options/dns.py.
type DNSOptions struct {
	// NameServers is inheritable, NameServersSearch is not.
	NameServers       Value[[]string] `mapstructure:"name_servers" json:"name_servers" yaml:"name_servers"`
	NameServersSearch []string        `mapstructure:"name_servers_search" json:"name_servers_search" yaml:"name_servers_search"`
}

func newDNSOptions() DNSOptions {
	return DNSOptions{
		NameServers:       Value[[]string]{Data: make([]string, 0)},
		NameServersSearch: make([]string, 0),
	}
}

// TFTPOptions models the TFTP next-server settings shared by Profile and
// System items. Get the fields from cobbler/items/options/tftp.py.
type TFTPOptions struct {
	NextServerV4 string `mapstructure:"next_server_v4" json:"next_server_v4" yaml:"next_server_v4"`
	NextServerV6 string `mapstructure:"next_server_v6" json:"next_server_v6" yaml:"next_server_v6"`
}

func newTFTPOptions() TFTPOptions {
	return TFTPOptions{
		NextServerV4: inherit,
		NextServerV6: inherit,
	}
}

// APTOptions models the APT mirror settings for a Repo item. Get the fields
// from cobbler/items/options/package.py. Neither field is inheritable.
type APTOptions struct {
	Components []string `mapstructure:"components" json:"components" yaml:"components"`
	Dists      []string `mapstructure:"dists" json:"dists" yaml:"dists"`
}

func newAPTOptions() APTOptions {
	return APTOptions{
		Components: make([]string, 0),
		Dists:      make([]string, 0),
	}
}
