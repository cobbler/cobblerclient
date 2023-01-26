package cobblerclient

// BuildisoOptions is a struct which describes the options one can set for the buildiso action of Cobbler.
type BuildisoOptions struct {
	Iso           string   `mapstructure:"iso"`
	Profiles      []string `mapstructure:"profiles"`
	Systems       []string `mapstructure:"systems"`
	BuildisoDir   string   `mapstructure:"buildisodir"`
	Distro        string   `mapstructure:"distro"`
	Standalone    bool     `mapstructure:"standalone"`
	Airgapped     bool     `mapstructure:"airgapped"`
	Source        string   `mapstructure:"source"`
	ExcludeDns    bool     `mapstructure:"exclude_dns"`
	XorrisofsOpts string   `mapstructure:"xorrisofs_opts"`
}

// AclSetupOptions is a struct which describes the options one can set for the actlsetup action of Cobbler.
type AclSetupOptions struct {
	AddUser     string `mapstructure:"adduser"`
	AddGroup    string `mapstructure:"addgroup"`
	RemoveUser  string `mapstructure:"removeuser"`
	RemoveGroup string `mapstructure:"removegroup"`
}

// ReplicateOptions is a struct which descibres the options one can set for the replicate action of Cobbler.
type ReplicateOptions struct {
	Master            string `mapstructure:"master"`
	Port              string `mapstructure:"port"`
	DistroPatterns    string `mapstructure:"distro_patterns"`
	ProfilePatterns   string `mapstructure:"profile_patterns"`
	SystemPatterns    string `mapstructure:"system_patterns"`
	RepoPatterns      string `mapstructure:"repo_patterns"`
	Imagepatterns     string `mapstructure:"image_patterns"`
	MgmtclassPatterns string `mapstructure:"mgmtclass_patterns"`
	PackagePatterns   string `mapstructure:"package_patterns"`
	FilePatterns      string `mapstructure:"file_patterns"`
	Prune             bool   `mapstructure:"prune"`
	OmitData          bool   `mapstructure:"omit_data"`
	SyncAll           bool   `mapstructure:"sync_all"`
	UseSsl            bool   `mapstructure:"use_ssl"`
}

type BackgroundSyncOptions struct {
	Dhcp    bool `mapstructure:"dhcp"`
	Dns     bool `mapstructure:"dns"`
	Verbose bool `mapstructure:"verbose"`
}

type BackgroundSyncSystemsOptions struct {
	Systems []string `mapstructure:"systems"`
	Verbose bool     `mapstructure:"verbose"`
}

type BackgroundImportOptions struct {
	Path            string `mapstructure:"path"`
	Name            string `mapstructure:"name"`
	AvailableAs     string `mapstructure:"available_as"`
	AutoinstallFile string `mapstructure:"autoinstall_file"`
	RsyncFlags      string `mapstructure:"rsync_flags"`
	Arch            string `mapstructure:"arch"`
	Breed           string `mapstructure:"breed"`
	OsVersion       string `mapstructure:"os_version"`
}

type BackgroundReposyncOptions struct {
	Repos  []string `mapstructure:"repos"`
	Only   string   `mapstructure:"only"`
	Nofail bool     `mapstructure:"nofail"`
	Tries  int      `mapstructure:"tries"`
}
