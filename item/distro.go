package item

import (
	"errors"
	"github.com/cobbler/cobblerclient/client"
	internalitem "github.com/cobbler/cobblerclient/internal/item"
	rawItem "github.com/cobbler/cobblerclient/internal/item/raw"
	resolvedItem "github.com/cobbler/cobblerclient/internal/item/resolved"
	"reflect"
)

// Distro specific fields
type Distro struct {
	Item

	// Internal fields
	raw      *rawItem.Distro
	resolved *resolvedItem.Distro

	// Distro fields
	SourceRepos    internalitem.Property[[]string]
	TreeBuildTime  internalitem.Property[string]
	Arch           internalitem.Property[string]
	BootFiles      internalitem.Property[map[string]string]
	BootLoaders    internalitem.InheritableProperty[[]string]
	Breed          internalitem.Property[string]
	FetchableFiles internalitem.Property[map[string]string]
	Initrd         internalitem.Property[string]
	Kernel         internalitem.Property[string]
	MGMTClasses    internalitem.Property[[]string]
	OSVersion      internalitem.Property[string]
	TemplateFiles  internalitem.Property[map[string]string]
}

func BuildDistro(client client.Client) Distro {
	i := BuildItem(client)
	d := Distro{
		raw: &rawItem.Distro{
			Item: *i.raw,
		},
		resolved: &resolvedItem.Distro{
			Item: *i.resolved,
		},
		Item: i,
	}
	refreshDistroPointers(&d)

	return d
}

func refreshDistroPointers(d *Distro) {
	d.Item.raw = &d.raw.Item
	d.Item.resolved = &d.resolved.Item
	refreshItemPointers(&d.Item)
	d.SourceRepos = SourceRepos{d}
	d.TreeBuildTime = TreeBuildTime{d}
	d.Arch = Arch{d}
	d.BootFiles = BootFiles{&d.Item}
	d.BootLoaders = BootLoaders{d}
	d.Breed = Breed{d}
	d.FetchableFiles = FetchableFiles{&d.Item}
	d.Initrd = Initrd{d}
	d.Kernel = Kernel{d}
	d.MGMTClasses = MGMTClasses{d}
	d.OSVersion = OSVersion{d}
	d.TemplateFiles = TemplateFiles{&d.Item}
}

type SourceRepos struct {
	item *Distro
}

func (s SourceRepos) Get() []string {
	return s.item.raw.SourceRepos
}

func (s SourceRepos) Set(sourceRepos []string) {
	s.item.raw.SourceRepos = sourceRepos
}

type TreeBuildTime struct {
	item *Distro
}

func (t TreeBuildTime) Get() string {
	return t.item.raw.TreeBuildTime
}

func (t TreeBuildTime) Set(treeBuildTime string) {
	t.item.raw.TreeBuildTime = treeBuildTime
}

type Arch struct {
	item *Distro
}

func (a Arch) Get() string {
	return a.item.raw.Arch
}

func (a Arch) Set(arch string) {
	a.item.raw.Arch = arch
}

type BootLoaders struct {
	item *Distro
}

func (b BootLoaders) Get() []string {
	return b.item.raw.SourceRepos
}

func (b BootLoaders) GetRaw() interface{} {
	return b.item.raw.BootLoaders
}

func (b BootLoaders) Set(sourceRepos []string) {
	b.item.raw.SourceRepos = sourceRepos
}

func (b BootLoaders) SetRaw(loaders interface{}) {
	b.item.raw.BootLoaders = loaders
}

func (b BootLoaders) GetCurrentRawType() reflect.Type {
	return reflect.TypeOf(b.item.raw.BootLoaders)
}

func (b BootLoaders) IsInherited() (bool, error) {
	switch b.item.BootLoaders.GetRaw().(type) {
	case []string:
		// Lists are always inherited
		return true, nil
	case string:
		return true, nil
	default:
		return false, errors.New("unexpected type for variable")
	}
}

type Breed struct {
	item *Distro
}

func (s Breed) Get() string {
	return s.item.raw.Breed
}

func (s Breed) Set(breed string) {
	s.item.raw.Breed = breed
}

type Initrd struct {
	item *Distro
}

func (s Initrd) Get() string {
	return s.item.raw.Initrd
}

func (s Initrd) Set(initrd string) {
	s.item.raw.Initrd = initrd
}

type Kernel struct {
	item *Distro
}

func (k Kernel) Get() string {
	return k.item.raw.Kernel
}

func (k Kernel) Set(kernel string) {
	k.item.raw.Kernel = kernel
}

type MGMTClasses struct {
	item *Distro
}

func (m MGMTClasses) Get() []string {
	return m.item.raw.MGMTClasses
}

func (m MGMTClasses) Set(mgmtclasses []string) {
	m.item.raw.MGMTClasses = mgmtclasses
}

type OSVersion struct {
	item *Distro
}

func (o OSVersion) Get() string {
	return o.item.raw.OSVersion
}

func (o OSVersion) Set(osversion string) {
	o.item.raw.OSVersion = osversion
}
