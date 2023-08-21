package resolved

// Item general fields
type Item struct {
	// Internal XML-RPC handle
	Handle string

	// Item fields
	Parent            string
	Depth             int
	Children          []string
	CTime             float64
	MTime             float64
	Uid               string
	Name              string
	Comment           string
	KernelOptions     map[string]string
	KernelOptionsPost map[string]string
	AutoinstallMeta   map[string]string
	FetchableFiles    map[string]string
	BootFiles         map[string]string
	TemplateFiles     map[string]string
	Owners            []string
	MgmtClasses       []string
	MgmtParameters    map[string]interface{}
}
