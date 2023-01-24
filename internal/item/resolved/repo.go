package resolved

type Repo struct {
	Item

	AptComponents   []string
	AptDists        []string
	Arch            string
	Breed           string
	CreateRepoFlags string
	Environment     []string
	KeepUpdated     bool
	Mirror          string
	MirrorLocally   bool
	Proxy           string
	RpmList         []string
	//YumOpts                map[string]interface{}
}
