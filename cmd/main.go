package main

import (
	"fmt"
	"github.com/cobbler/cobblerclient/item"
	"net/http"
	"os"

	cobblerClient "github.com/cobbler/cobblerclient/client"
)

var config = cobblerClient.ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

func main() {
	c := cobblerClient.NewClient(http.DefaultClient, config)
	_, err := c.Login()
	if err != nil {
		fmt.Printf("Error logging in: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Token: %s\n", c.Token)

	fmt.Println("Creating a repo")
	r := item.BuildRepo(c)
	r.Name.Set("myrepo")

	newRepo, err := r.Create()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", newRepo)

	fmt.Println("Listing all repos")
	repos, err := item.GetRepos(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, repo := range repos {
		fmt.Printf("%+v\n", repo)
	}

	fmt.Println("Deleting a repo")
	err = r.Delete()
	if err != nil {
		fmt.Println(err)
	}

	d := item.BuildDistro(c)
	d.Name.Set("testdistro")
	d.Arch.Set("x86_64")
	d.Breed.Set("Ubuntu")
	d.OSVersion.Set("focal")
	d.Kernel.Set("/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/linux")
	d.Initrd.Set("/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/initrd.gz")

	fmt.Println("Listing all distros")
	distros, err := item.GetDistros(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Distros: %#v\n", distros)
	fmt.Println("Creating a Distro")
	newDistro, err := d.Create()
	if err != nil {
		fmt.Printf("Error creating distro: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("New Distro: %+v\n", newDistro)

	if newDistro.Name.Get() != "Test" {
		fmt.Println("Distro name does not match.")
	}

	fmt.Println("Updating Distro")
	newDistro.Comment.Set("Update Test")
	if err := newDistro.Update(); err != nil {
		fmt.Printf("Error creating distro: %s\n", err)
	}

	fmt.Println("Listing all profiles")
	profiles, err := item.GetProfiles(c)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Profiles: %#v\n", profiles)
	fmt.Println("Creating a Profile")
	p := item.BuildProfile(c)
	p.Name.Set("testprofile")
	p.Distro.Set("Ubuntu-20.04-x86_64")
	p.Parent.Set("Ubuntu-20.04-x86_64")
	p.Autoinstall.Set("sample.seed")
	p.VirtDiskDriver.Set("raw")

	newProfile, err := p.Create()
	if err != nil {
		fmt.Printf("Error creating profile: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("New Profile: %+v\n", newProfile)

	fmt.Println("Creating a System")
	s := item.System{
		Item:        item.BuildItem(c),
		Profile:     "testprofile",
		NameServers: []string{"8.8.8.8", "1.1.1.1"},
		PowerID:     "foo",
	}
	s.Name.Set("testsystem")
	s.Comment.Set("I'd like to teach the world to sing")

	newSystem, err := s.Create()
	if err != nil {
		fmt.Printf("Error creating system: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("New System: %+v\n", newSystem)

	eth0 := item.Interface{
		MACAddress:    "aa:bb:cc:dd:ee:ff",
		Static:        true,
		InterfaceType: "bridge",
	}

	eth1 := item.Interface{
		MACAddress:    "aa:bb:cc:dd:ee:fa",
		Static:        true,
		Management:    true,
		InterfaceType: "na",
	}

	fmt.Println("Adding NIC to System")
	if err := eth0.CreateInterface(*newSystem, "eth0"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Adding second NIC to System")
	if err := eth1.CreateInterface(*newSystem, "eth1"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//
	fmt.Println("Syncing the cobbler server")
	if err := c.Sync(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Getting system")
	s2, err := item.GetSystem(c, "testsystem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("s2: %#v\n", s2)

	fmt.Println("Verifying NIC data")
	interfaces, err := item.GetInterfaces(*s2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n\n", interfaces)

	iface, err := item.GetInterface(*s2, "eth0")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("eth0:\n%+v\n\n", iface)

	iface, err = item.GetInterface(*s2, "eth1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("eth1:\n%+v\n\n", iface)

	fmt.Println("Deleting Interface")
	err = item.DeleteInterface(*s2, "eth0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deleting Interface")
	err = item.DeleteInterface(*s2, "eth1")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s2, err = item.GetSystem(c, "testsystem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	interfaces, err = item.GetInterfaces(*s2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(interfaces) != 1 {
		fmt.Println("Error deleting interface eth1")
		fmt.Printf("%+v\n", interfaces)
	}

	sysToDelete, err := item.GetSystem(c, "testsystem")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Deleting System")
	err = sysToDelete.Delete()
	if err != nil {
		fmt.Printf("Error deleting system: %s\n", err)
	}

	fmt.Println("Deleting Profile")
	err = p.Delete()
	if err != nil {
		fmt.Printf("Error deleting profile: %s\n", err)
		os.Exit(1)
	}
	//
	fmt.Println("Deleting Distro")
	err = d.Delete()
	if err != nil {
		fmt.Printf("Error deleting distro: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Creating a Snippet")
	snippet := cobblerClient.Snippet{
		Name: "testsnippet",
		Body: "sample content",
	}

	err = c.CreateSnippet(snippet)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Deleting a Snippet")
	if err := c.DeleteSnippet("testsnippet"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Creating a Template")
	ks := cobblerClient.TemplateFile{
		Name: "testtemplate.ks",
		Body: "sample content",
	}

	err = c.CreateTemplateFile(ks)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Deleting a Template")
	if err := c.DeleteTemplateFile("testtemplate.ks"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
