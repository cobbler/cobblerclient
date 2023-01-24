package item

import (
	"github.com/cobbler/cobblerclient/client"
)

type PackageFactory struct {
	linuxpackage Package
}

func (p *PackageFactory) SetInstaller(installer string) {
	p.linuxpackage.Installer.Set(installer)
}

func (p *PackageFactory) SetVersion(name string) {
	p.linuxpackage.Version.Set(name)
}

func Build(client client.Client) Package {
	linuxpackage := Package{}
	linuxpackage.Item = BuildItem(client)
	linuxpackage.Installer = Installer{linuxpackage}
	linuxpackage.Version = Version{linuxpackage}
	return linuxpackage
}
