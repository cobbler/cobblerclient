package raw

// Repo is a created repo.
// Get the fields from cobbler/items/repo.py
type Repo struct {
	Item `mapstructure:",squash"`

	// Repo fields
	AptComponents   []string `mapstructure:"apt_components"`
	AptDists        []string `mapstructure:"apt_dists"`
	Arch            string   `mapstructure:"arch"`
	Breed           string   `mapstructure:"breed"`
	CreateRepoFlags string   `mapstructure:"createrepo_flags"`
	Environment     []string `mapstructure:"environment"`
	KeepUpdated     bool     `mapstructure:"keep_updated"`
	Mirror          string   `mapstructure:"mirror"`
	MirrorLocally   bool     `mapstructure:"mirror_locally"`
	Proxy           string   `mapstructure:"proxy" cobbler:"newfield"`
	RpmList         []string `mapstructure:"rpm_list"`
	//YumOpts                map[string]interface{} `mapstructure:"yumopts"`
}
