package main

import (
	"fmt"
	"net/http"

	cobbler "github.com/cobbler/cobblerclient"
)

var config = cobbler.ClientConfig{
	URL:      "http://localhost:8081/cobbler_api",
	Username: "cobbler",
	Password: "cobbler",
}

func main() {
	c := cobbler.NewClient(http.DefaultClient, config)
	_, err := c.Login()
	if err != nil {
		fmt.Printf("Error logging in: %s\n", err)
	}

	fmt.Printf("Token: %s\n", c.Token)

	res, err := c.BackgroundSync(cobbler.BackgroundSyncOptions{Dhcp: true, Dns: true, Verbose: true})
	if err != nil {
		fmt.Println(err)
	}
	eventLog, err := c.GetEventLog(res)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Print event log for %s\n", res)
	fmt.Println(eventLog)

	events, err := c.GetEvents("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Listing all events")
	for _, event := range events {
		fmt.Printf("%#v\n", event)
	}

	fmt.Println("Creating a repo")
	r := cobbler.Repo{
		Item: cobbler.Item{
			Name: "myrepo",
		},
		Arch:          "x86_64",
		Breed:         "yum",
		Mirror:        "http://repo/homeawayel7/",
		MirrorLocally: false,
	}

	newRepo, err := c.CreateRepo(r)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v\n", newRepo)

	fmt.Println("Listing all repos")
	repos, err := c.GetRepos()
	if err != nil {
		fmt.Println(err)
	}

	for _, repo := range repos {
		fmt.Printf("%+v\n", repo)
	}

	fmt.Println("Deleting a repo")
	err = c.DeleteRepo("myrepo")
	if err != nil {
		fmt.Println(err)
	}

	d := cobbler.Distro{
		Item: cobbler.Item{
			Name: "mydistro",
		},
		Arch:      "x86_64",
		Breed:     "Ubuntu",
		Initrd:    "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/initrd.gz",
		Kernel:    "/srv/www/cobbler/distro_mirror/Ubuntu-20.04/install/netboot/ubuntu-installer/amd64/linux",
		OSVersion: "focal",
	}

	fmt.Println("Listing all distros")
	distros, err := c.GetDistros()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Distros: %#v\n", distros)
	fmt.Println("Creating a Distro")
	newDistro, err := c.CreateDistro(d)
	if err != nil {
		fmt.Printf("Error creating distro: %s\n", err)
	}

	fmt.Printf("New Distro: %+v\n", newDistro)

	if newDistro.Name != "Test" {
		fmt.Println("Distro name does not match.")
	}

	fmt.Println("Updating Distro")
	newDistro.Comment = "Update Test"
	if err := c.UpdateDistro(newDistro); err != nil {
		fmt.Printf("Error creating distro: %s\n", err)
	}

	fmt.Println("Listing all profiles")
	profiles, err := c.GetProfiles()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Profiles: %#v\n", profiles)
	fmt.Println("Creating a Profile")
	p := cobbler.Profile{
		Item: cobbler.Item{
			Name:   "testprofile",
			Parent: "Ubuntu-20.04-x86_64",
		},
		Autoinstall:    "sample.seed",
		Distro:         "Ubuntu-20.04-x86_64",
		VirtDiskDriver: "raw", // For some reason the virt_disk_driver must be set in Cobbler 3...
	}

	newProfile, err := c.CreateProfile(p)
	if err != nil {
		fmt.Printf("Error creating profile: %s\n", err)
	}

	fmt.Printf("New Profile: %+v\n", newProfile)

	fmt.Println("Creating a System")
	s := cobbler.System{
		Item: cobbler.Item{
			Comment: "I'd like to teach the world to sing",
			Name:    "testsystem",
		},
		NameServers: []string{"8.8.8.8", "1.1.1.1"},
		PowerID:     "foo",
		Profile:     "testprofile",
	}

	newSystem, err := c.CreateSystem(s)
	if err != nil {
		fmt.Printf("Error creating system: %s\n", err)
	}

	fmt.Printf("New System: %+v\n", newSystem)

	eth0 := cobbler.Interface{
		MACAddress:    "aa:bb:cc:dd:ee:ff",
		Static:        true,
		InterfaceType: "bridge",
	}

	eth1 := cobbler.Interface{
		MACAddress:    "aa:bb:cc:dd:ee:fa",
		Static:        true,
		Management:    true,
		InterfaceType: "na",
	}

	fmt.Println("Adding NIC to System")
	if err := newSystem.CreateInterface("eth0", eth0); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Adding second NIC to System")
	if err := newSystem.CreateInterface("eth1", eth1); err != nil {
		fmt.Println(err)
	}
	//
	fmt.Println("Syncing the cobbler server")
	if err := c.Sync(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Getting system")
	s2, err := c.GetSystem("testsystem", false, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("s2: %#v\n", s2)

	fmt.Println("Verifying NIC data")
	interfaces, err := s2.GetInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n\n", interfaces)

	iface, err := s2.GetInterface("eth0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("eth0:\n%+v\n\n", iface)

	iface, err = s2.GetInterface("eth1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("eth1:\n%+v\n\n", iface)

	fmt.Println("Deleting Interface")
	err = s2.DeleteInterface("eth0")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Deleting Interface")
	err = s2.DeleteInterface("eth1")
	if err != nil {
		fmt.Println(err)
	}

	s2, err = c.GetSystem("testsystem", false, false)
	if err != nil {
		fmt.Println(err)
	}
	interfaces, err = s2.GetInterfaces()
	if err != nil {
		fmt.Println(err)
	}
	if len(interfaces) != 1 {
		fmt.Println("Error deleting interface eth1")
		fmt.Printf("%+v\n", interfaces)
	}

	_, err = c.GetSystem("testsystem", false, false)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Deleting System")
	err = c.DeleteSystem("testsystem")
	if err != nil {
		fmt.Printf("Error deleting system: %s\n", err)
	}

	fmt.Println("Deleting Profile")
	err = c.DeleteProfile("testprofile")
	if err != nil {
		fmt.Printf("Error deleting profile: %s\n", err)
	}
	//
	fmt.Println("Deleting Distro")
	err = c.DeleteDistro("testdistro")
	if err != nil {
		fmt.Printf("Error deleting distro: %s\n", err)
	}

	fmt.Println("Creating a Snippet")
	snippet := cobbler.Snippet{
		Name: "testsnippet",
		Body: "sample content",
	}

	err = c.CreateSnippet(snippet)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Deleting a Snippet")
	if err := c.DeleteSnippet("testsnippet"); err != nil {
		fmt.Println(err)
	}

	fmt.Println("Creating a Template")
	ks := cobbler.TemplateFile{
		Name: "testtemplate.ks",
		Body: "sample content",
	}

	err = c.CreateTemplateFile(ks)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Deleting a Template")
	if err := c.DeleteTemplateFile("testtemplate.ks"); err != nil {
		fmt.Println(err)
	}

}
