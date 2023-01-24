package client_test

import (
	"reflect"
	"testing"

	"github.com/ContainerSolutions/go-utils"

	"github.com/cobbler/cobblerclient/client"
	cobblerTesting "github.com/cobbler/cobblerclient/internal/testing"
)

func TestSync(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "sync-req.xml", "sync-res.xml")

	err := c.Sync()
	utils.FailOnError(t, err)
}

func TestBackgroundSync(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-sync-req.xml", "background-sync-res.xml")

	res, err := c.BackgroundSync(client.BackgroundSyncOptions{Dhcp: false, Dns: false, Verbose: false})
	utils.FailOnError(t, err)
	if res != "2022-09-30_145124_Sync_2cabdc4eddfa4731b45f145d7b625e29" {
		t.Errorf("Problem with event id return")
	}
}

func TestBackgroundSyncSystems(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-sync-systems-req.xml", "background-sync-systems-res.xml")

	res, err := c.BackgroundSyncSystems(client.BackgroundSyncSystemsOptions{Systems: []string{"", ""}, Verbose: false})
	utils.FailOnError(t, err)
	if res != "2022-09-30_151856_Syncsystems_76d70bd7f48642f7b4cb5a0b0dcc93a5" {
		t.Errorf("Problem with event id return")
	}
}

func TestCheck(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "check-req.xml", "check-res.xml")
	expected := []string{"reposync not installed, install yum-utils"}

	result, err := c.Check()
	var resolvedResult = *result
	utils.FailOnError(t, err)

	if !reflect.DeepEqual(resolvedResult, expected) {
		t.Errorf("%s expected; got %s", expected, result)
	}
}

func TestBackgroundBuildiso(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-buildiso-req.xml", "background-buildiso-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundBuildiso(client.BuildisoOptions{
		Iso:           "",
		Profiles:      nil,
		Systems:       nil,
		BuildisoDir:   "",
		Distro:        "",
		Standalone:    false,
		Airgapped:     false,
		Source:        "",
		ExcludeDns:    false,
		XorrisofsOpts: "",
	})
	utils.FailOnError(t, err)
}

func TestBackgroundHardlink(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-hardlink-req.xml", "background-hardlink-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundHardlink()
	utils.FailOnError(t, err)
}

func TestValidateAutoinstallFiles(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(
		t,
		"background-validate-autoinstall-files-req.xml",
		"background-validate-autoinstall-files-res.xml",
	)

	// FIXME: Test event-id return
	_, err := c.BackgroundValidateAutoinstallFiles()
	utils.FailOnError(t, err)
}

func TestBackgroundReplicate(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-replicate-req.xml", "background-replicate-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundReplicate(client.ReplicateOptions{
		Master:            "",
		Port:              "",
		DistroPatterns:    "",
		ProfilePatterns:   "",
		SystemPatterns:    "",
		RepoPatterns:      "",
		Imagepatterns:     "",
		MgmtclassPatterns: "",
		PackagePatterns:   "",
		FilePatterns:      "",
		Prune:             false,
		OmitData:          false,
		SyncAll:           false,
		UseSsl:            false,
	})
	utils.FailOnError(t, err)
}

func TestBackgroundAclSetup(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-aclsetup-req.xml", "background-aclsetup-res.xml")

	res, err := c.BackgroundAclSetup(client.AclSetupOptions{
		AddUser:     "testing",
		AddGroup:    "",
		RemoveUser:  "",
		RemoveGroup: "",
	})
	if res != "2023-01-24_083137_(CLI) ACL Configuration_334327920d2946fda3ac95dbf457e76d" {
		t.Errorf("Event-ID was malformed")
	}
	utils.FailOnError(t, err)
}

func TestBackgroundImport(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-import-req.xml", "background-import-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundImport(client.BackgroundImportOptions{
		Path:            "",
		Name:            "",
		AvailableAs:     "",
		AutoinstallFile: "",
		RsyncFlags:      "",
		Arch:            "",
		Breed:           "",
		OsVersion:       "",
	})
	utils.FailOnError(t, err)
}

func TestBackgroundReposync(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-reposync-req.xml", "background-reposync-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundReposync(client.BackgroundReposyncOptions{
		Repos:  nil,
		Only:   "",
		Nofail: false,
		Tries:  0,
	})
	utils.FailOnError(t, err)
}

func TestBackgroundMkLoaders(t *testing.T) {
	c := cobblerTesting.CreateStubHTTPClient(t, "background-mkloaders-req.xml", "background-mkloaders-res.xml")

	// FIXME: Test event-id return
	_, err := c.BackgroundMkLoaders()
	utils.FailOnError(t, err)
}
