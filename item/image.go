package item

type Image struct {
	Item `mapstructure:",squash"`

	// Image specific fields
	Arch                 Architecture
	Autoinstall          string
	Breed                string
	File                 string
	ImageType            ImageType
	NetworkCount         int
	OsVersion            string
	BootLoaders          []string
	Menu                 string
	VirtAutoBoot         bool
	VirtBridge           string
	VirtCpus             int
	VirtDiskDriver       VirtDiskDriver
	VirtFileSize         float32
	VirtPath             string
	VirtRam              int
	VirtType             VirtType
	SupportedBootLoaders []string
}
