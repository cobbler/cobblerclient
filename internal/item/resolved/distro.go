package resolved

type Distro struct {
	Item

	// These are internal fields and cannot be modified.
	SourceRepos    []string
	TreeBuildTime  string
	Arch           string
	BootFiles      []string
	BootLoaders    []string
	Breed          string
	FetchableFiles []string
	Initrd         string
	Kernel         string
	MGMTClasses    []string
	OSVersion      string
	TemplateFiles  []string
}
