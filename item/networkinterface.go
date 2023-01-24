package item

// Interface is an interface in a system.
type Interface struct {
	CNAMEs             []string `mapstructure:"cnames" structs:"cnames"`
	DHCPTag            string   `mapstructure:"dhcp_tag" structs:"dhcp_tag"`
	DNSName            string   `mapstructure:"dns_name" structs:"dns_name"`
	BondingOpts        string   `mapstructure:"bonding_opts" structs:"bonding_opts"`
	BridgeOpts         string   `mapstructure:"bridge_opts" structs:"bridge_opts"`
	Gateway            string   `mapstructure:"if_gateway" structs:"if_gateway"`
	InterfaceType      string   `mapstructure:"interface_type" structs:"interface_type"`
	InterfaceMaster    string   `mapstructure:"interface_master" structs:"interface_master"`
	IPAddress          string   `mapstructure:"ip_address" structs:"ip_address"`
	IPv6Address        string   `mapstructure:"ipv6_address" structs:"ipv6_address"`
	IPv6Secondaries    []string `mapstructure:"ipv6_secondaries" structs:"ipv6_secondaries"`
	IPv6MTU            string   `mapstructure:"ipv6_mtu" structs:"ipv6_mtu"`
	IPv6StaticRoutes   []string `mapstructure:"ipv6_static_routes" structs:"ipv6_static_routes"`
	IPv6DefaultGateway string   `mapstructure:"ipv6_default_gateway" structs:"ipv6_default_gateway"`
	MACAddress         string   `mapstructure:"mac_address" structs:"mac_address"`
	Management         bool     `mapstructure:"management" structs:"management"`
	Netmask            string   `mapstructure:"netmask" structs:"netmask"`
	Static             bool     `mapstructure:"static" structs:"static"`
	StaticRoutes       []string `mapstructure:"static_routes" structs:"static_routes"`
	VirtBridge         string   `mapstructure:"virt_bridge" structs:"virt_bridge"`
}

// Interfaces is a collection of interfaces in a system.
type Interfaces map[string]Interface
