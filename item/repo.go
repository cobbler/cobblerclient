package item

import (
	"github.com/cobbler/cobblerclient/client"
	internalItem "github.com/cobbler/cobblerclient/internal/item"
	rawItem "github.com/cobbler/cobblerclient/internal/item/raw"
	resolvedItem "github.com/cobbler/cobblerclient/internal/item/resolved"
)

// Repo is a created repo.
// Get the fields from cobbler/items/repo.py
type Repo struct {
	// Internal fields
	raw      *rawItem.Repo
	resolved *resolvedItem.Repo

	Item

	// Repo fields
	AptComponents   internalItem.Property[[]string]
	AptDists        internalItem.Property[[]string]
	RepoArch        internalItem.Property[string]
	RepoBreed       internalItem.Property[string]
	CreateRepoFlags internalItem.Property[string]
	Environment     internalItem.Property[[]string]
	KeepUpdated     internalItem.Property[bool]
	Mirror          internalItem.Property[string]
	MirrorLocally   internalItem.Property[bool]
	Proxy           internalItem.Property[string]
	RpmList         internalItem.Property[[]string]
	//YumOpts                item.Property[map[string]interface{}]
}

func BuildRepo(client client.Client) Repo {
	var i = BuildItem(client)
	var r = Repo{
		raw: &rawItem.Repo{
			Item: *i.raw,
		},
		resolved: &resolvedItem.Repo{
			Item: *i.resolved,
		},
		Item: i,
	}
	refreshRepoPointers(&r)

	return r
}

func refreshRepoPointers(repo *Repo) {
	repo.Item.raw = &repo.raw.Item
	repo.Item.resolved = &repo.resolved.Item
	refreshItemPointers(&repo.Item)
	repo.AptComponents = AptComponents{repo}
	repo.AptDists = AptDists{repo}
	repo.RepoArch = RepoArch{repo}
	repo.RepoBreed = RepoBreed{repo}
	repo.CreateRepoFlags = CreateRepoFlags{repo}
	repo.Environment = Environment{repo}
	repo.KeepUpdated = KeepUpdated{repo}
	repo.Mirror = Mirror{repo}
	repo.MirrorLocally = MirrorLocally{repo}
	repo.Proxy = Proxy{repo}
	repo.RpmList = RpmList{repo}
}

type AptComponents struct {
	repo *Repo
}

func (a AptComponents) Get() []string {
	return a.repo.raw.AptComponents
}

func (a AptComponents) Set(aptComponents []string) {
	a.repo.raw.AptComponents = aptComponents
}

type AptDists struct {
	repo *Repo
}

func (a AptDists) Get() []string {
	return a.repo.raw.AptDists
}

func (a AptDists) Set(aptDists []string) {
	a.repo.raw.AptDists = aptDists
}

type RepoArch struct {
	repo *Repo
}

func (a RepoArch) Get() string {
	return a.repo.raw.Arch
}

func (a RepoArch) Set(arch string) {
	a.repo.raw.Arch = arch
}

type RepoBreed struct {
	repo *Repo
}

func (b RepoBreed) Get() string {
	return b.repo.raw.Breed
}

func (b RepoBreed) Set(breed string) {
	b.repo.raw.Breed = breed
}

type CreateRepoFlags struct {
	repo *Repo
}

func (c CreateRepoFlags) Get() string {
	return c.repo.raw.CreateRepoFlags
}

func (c CreateRepoFlags) Set(createRepoFlags string) {
	c.repo.raw.CreateRepoFlags = createRepoFlags
}

type Environment struct {
	repo *Repo
}

func (e Environment) Get() []string {
	return e.repo.raw.Environment
}

func (e Environment) Set(env []string) {
	e.repo.raw.Environment = env
}

type KeepUpdated struct {
	repo *Repo
}

func (k KeepUpdated) Get() bool {
	return k.repo.raw.KeepUpdated
}

func (k KeepUpdated) Set(keepUpdated bool) {
	k.repo.raw.KeepUpdated = keepUpdated
}

type Mirror struct {
	repo *Repo
}

func (m Mirror) Get() string {
	return m.repo.raw.Mirror
}

func (m Mirror) Set(mirror string) {
	m.repo.raw.Mirror = mirror
}

type MirrorLocally struct {
	repo *Repo
}

func (m MirrorLocally) Get() bool {
	return m.repo.raw.MirrorLocally
}

func (m MirrorLocally) Set(mirrorLocally bool) {
	m.repo.raw.MirrorLocally = mirrorLocally
}

type Proxy struct {
	repo *Repo
}

func (p Proxy) Get() string {
	return p.repo.raw.Proxy
}

func (p Proxy) Set(proxy string) {
	p.repo.raw.Proxy = proxy
}

type RpmList struct {
	repo *Repo
}

func (r RpmList) Get() []string {
	return r.repo.raw.RpmList
}

func (r RpmList) Set(rpmList []string) {
	r.repo.raw.RpmList = rpmList
}
