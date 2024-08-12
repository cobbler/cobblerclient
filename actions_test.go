package cobblerclient

import (
	"reflect"
	"testing"

	"github.com/ContainerSolutions/go-utils"
)

func TestSync(t *testing.T) {
	c := createStubHTTPClient(t, "sync-req.xml", "sync-res.xml")

	err := c.Sync()
	utils.FailOnError(t, err)
}

func TestBackgroundSync(t *testing.T) {
	c := createStubHTTPClient(t, "background-sync-req.xml", "background-sync-res.xml")

	res, err := c.BackgroundSync(BackgroundSyncOptions{Dhcp: false, Dns: false, Verbose: false})
	utils.FailOnError(t, err)
	if res != "2022-09-30_145124_Sync_2cabdc4eddfa4731b45f145d7b625e29" {
		t.Errorf("Problem with event id return")
	}
}

func TestBackgroundSyncSystems(t *testing.T) {
	c := createStubHTTPClient(t, "background-sync-systems-req.xml", "background-sync-systems-res.xml")

	res, err := c.BackgroundSyncSystems(BackgroundSyncSystemsOptions{Systems: []string{"", ""}, Verbose: false})
	utils.FailOnError(t, err)
	if res != "2022-09-30_151856_Syncsystems_76d70bd7f48642f7b4cb5a0b0dcc93a5" {
		t.Errorf("Problem with event id return")
	}
}

func TestCheck(t *testing.T) {
	c := createStubHTTPClient(t, "check-req.xml", "check-res.xml")
	expected := []string{"reposync not installed, install yum-utils"}

	result, err := c.Check()
	var resolvedResult = *result
	utils.FailOnError(t, err)

	if !reflect.DeepEqual(resolvedResult, expected) {
		t.Errorf("%s expected; got %s", expected, result)
	}
}

func TestBackgroundBuildiso(t *testing.T) {
	c := createStubHTTPClient(t, "background-buildiso-req.xml", "background-buildiso-res.xml")

	res, err := c.BackgroundBuildiso(BuildisoOptions{
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
	if res != "2023-01-24_083001_Build Iso_20fa7d4256fc4f61a2b9c2237c80fb41" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundHardlink(t *testing.T) {
	c := createStubHTTPClient(t, "background-hardlink-req.xml", "background-hardlink-res.xml")

	res, err := c.BackgroundHardlink()
	utils.FailOnError(t, err)
	if res != "2022-09-30_203004_Hardlink_800c38f4e0424187aed6a6ffb6553ef8" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestValidateAutoinstallFiles(t *testing.T) {
	c := createStubHTTPClient(
		t,
		"background-validate-autoinstall-files-req.xml",
		"background-validate-autoinstall-files-res.xml",
	)

	res, err := c.BackgroundValidateAutoinstallFiles()
	utils.FailOnError(t, err)
	if res != "2022-09-30_203505_Automated installation files validation_487b1a5d1d914c62834126391ac2b601" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundReplicate(t *testing.T) {
	c := createStubHTTPClient(t, "background-replicate-req.xml", "background-replicate-res.xml")

	res, err := c.BackgroundReplicate(ReplicateOptions{
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
	if res != "2023-01-24_075801_Replicate_ea7a003a81264039b4277ac55664661a" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundAclSetup(t *testing.T) {
	c := createStubHTTPClient(t, "background-aclsetup-req.xml", "background-aclsetup-res.xml")

	res, err := c.BackgroundAclSetup(AclSetupOptions{
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
	c := createStubHTTPClient(t, "background-import-req.xml", "background-import-res.xml")

	res, err := c.BackgroundImport(BackgroundImportOptions{
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
	if res != "2023-01-24_103639_Media import_dd297121f7bc412e9ce4d80f05de4b3f" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundReposync(t *testing.T) {
	c := createStubHTTPClient(t, "background-reposync-req.xml", "background-reposync-res.xml")

	res, err := c.BackgroundReposync(BackgroundReposyncOptions{
		Repos:  nil,
		Only:   "",
		Nofail: false,
		Tries:  0,
	})
	utils.FailOnError(t, err)
	if res != "2023-01-24_103758_Reposync_3478fd19fd5f48bf8b40c728ad247348" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundMkLoaders(t *testing.T) {
	c := createStubHTTPClient(t, "background-mkloaders-req.xml", "background-mkloaders-res.xml")

	res, err := c.BackgroundMkLoaders()
	utils.FailOnError(t, err)
	if res != "2022-09-30_203957_Create bootable bootloader images_9c809af4d6f148e49b071fac84f9a664" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundPowerSystem(t *testing.T) {
	c := createStubHTTPClient(t, "background-power-system-req.xml", "background-power-system-res.xml")

	result, err := c.BackgroundPowerSystem(BackgroundPowerSystemOptions{
		Systems: []string{"testsys1"},
		Power:   "off",
	})
	utils.FailOnError(t, err)
	if result != "2024-08-06_072956_Power management ()_1a44f162efa74806b16d055dfad0fc04" {
		t.Errorf("Event-ID was malformed")
	}
}

func TestPowerSystem(t *testing.T) {
	c := createStubHTTPClient(t, "power-system-req.xml", "power-system-res.xml")

	result, err := c.PowerSystem("system::testsys1", "status")
	utils.FailOnError(t, err)
	if !result {
		t.Errorf("Expected power operation not to fail!")
	}
}
