package cobblerclient

import (
	"reflect"
	"testing"
)

func TestSync(t *testing.T) {
	c := createStubHTTPClientSingle(t, "sync")

	err := c.Sync()
	FailOnError(t, err)
}

func TestBackgroundSync(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-sync")

	res, err := c.BackgroundSync(BackgroundSyncOptions{Dhcp: false, Dns: false, Verbose: false})
	FailOnError(t, err)
	if res != "2000-01-01_000000_Sync_00000000000000000000000000000001" {
		t.Errorf("Problem with event id return")
	}
}

func TestBackgroundSyncSystems(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-sync-systems")

	res, err := c.BackgroundSyncSystems(BackgroundSyncSystemsOptions{Systems: []string{"", ""}, Verbose: false})
	FailOnError(t, err)
	if res != "2000-01-01_000000_Syncsystems_00000000000000000000000000000002" {
		t.Errorf("Problem with event id return")
	}
}

func TestCheck(t *testing.T) {
	c := createStubHTTPClientSingle(t, "check")
	expected := []string{
		"The 'server' field in /etc/cobbler/settings.yaml must be set to something other than localhost, or automatic installation features will not work.  This should be a resolvable hostname or IP for the boot server as reachable by all machines that will use it.",
		"For PXE to be functional, the 'next_server_v4' field in /etc/cobbler/settings.yaml must be set to something other than 127.0.0.1, and should match the IP of the boot server on the PXE network.",
		"For PXE to be functional, the 'next_server_v6' field in /etc/cobbler/settings.yaml must be set to something other than ::1, and should match the IP of the boot server on the PXE network.",
		"some network boot-loaders are missing from /var/lib/cobbler/loaders. If you only want to handle x86/x86_64 netbooting, you may ensure that you have installed a *recent* version of the syslinux package installed and can ignore this message entirely. Files in this directory, should you want to support all architectures, should include pxelinux.0, andmenu.c32.",
		"enable and start rsyncd.service with systemctl",
		"reposync is not installed, install yum-utils or dnf-plugins-core",
		"yumdownloader is not installed, install yum-utils or dnf-plugins-core",
		"The default password used by the sample templates for newly installed machines (default_password_crypted in /etc/cobbler/settings.yaml) is still set to 'cobbler' and should be changed, try: \"openssl passwd -1 -salt 'random-phrase-here' 'your-password-here'\" to generate new one",
	}

	result, err := c.Check()
	var resolvedResult = *result
	FailOnError(t, err)

	if !reflect.DeepEqual(resolvedResult, expected) {
		t.Errorf("%s expected; got %s", expected, result)
	}
}

func TestBackgroundBuildiso(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-buildiso")

	res, err := c.BackgroundBuildiso(BuildisoOptions{
		Iso:            "",
		Profiles:       nil,
		Systems:        nil,
		BuildisoDir:    "",
		Distro:         "",
		Standalone:     false,
		Airgapped:      false,
		Source:         "",
		ExcludeDns:     false,
		ExcludeSystems: false,
		XorrisofsOpts:  "",
		Esp:            "",
	})
	FailOnError(t, err)
	if res != "2000-01-01_000000_Build Iso_00000000000000000000000000000003" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundHardlink(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-hardlink")

	res, err := c.BackgroundHardlink()
	FailOnError(t, err)
	if res != "2000-01-01_000000_Hardlink_00000000000000000000000000000004" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestValidateAutoinstallFiles(t *testing.T) {
	c := createStubHTTPClientSingle(
		t,
		"background-validate-autoinstall-files",
	)

	res, err := c.BackgroundValidateAutoinstallFiles()
	FailOnError(t, err)
	if res != "2000-01-01_000000_Automated installation files validation_00000000000000000000000000000005" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundReplicate(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-replicate")

	res, err := c.BackgroundReplicate(ReplicateOptions{
		Master:          "",
		Port:            "",
		DistroPatterns:  "",
		ProfilePatterns: "",
		SystemPatterns:  "",
		RepoPatterns:    "",
		Imagepatterns:   "",
		Prune:           false,
		OmitData:        false,
		SyncAll:         false,
		UseSsl:          false,
	})
	FailOnError(t, err)
	if res != "2000-01-01_000000_Replicate_00000000000000000000000000000006" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundAclSetup(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-aclsetup")

	res, err := c.BackgroundAclSetup(AclSetupOptions{
		AddUser:     "testing",
		AddGroup:    "",
		RemoveUser:  "",
		RemoveGroup: "",
	})
	if res != "2000-01-01_000000_(CLI) ACL Configuration_00000000000000000000000000000007" {
		t.Errorf("Event-ID was malformed")
	}
	FailOnError(t, err)
}

func TestBackgroundImport(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-import")

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
	FailOnError(t, err)
	if res != "2000-01-01_000000_Media import_00000000000000000000000000000008" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundReposync(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-reposync")

	res, err := c.BackgroundReposync(BackgroundReposyncOptions{
		Repos:  nil,
		Only:   "",
		Nofail: false,
		Tries:  0,
	})
	FailOnError(t, err)
	if res != "2000-01-01_000000_Reposync_00000000000000000000000000000009" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundMkLoaders(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-mkloaders")

	res, err := c.BackgroundMkLoaders()
	FailOnError(t, err)
	if res != "2000-01-01_000000_Create bootable bootloader images_0000000000000000000000000000000a" {
		t.Fatalf("Expected a different Event-ID!")
	}
}

func TestBackgroundPowerSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "background-power-system")

	result, err := c.BackgroundPowerSystem(BackgroundPowerSystemOptions{
		Systems: []string{"testsys1"},
		Power:   "off",
	})
	FailOnError(t, err)
	if result != "2000-01-01_000000_Power management ()_0000000000000000000000000000000b" {
		t.Errorf("Event-ID was malformed")
	}
}

func TestPowerSystem(t *testing.T) {
	c := createStubHTTPClientSingle(t, "power-system")
	expected := `Fault(1): <class 'cobbler.cexceptions.CX'>:'command failed (rc=1), please validate the physical setup and cobbler config'`

	result, err := c.PowerSystem("0000000000000000000000000000001b", "status")
	if result {
		t.Errorf("Expected power operation to fail!")
	}
	if err == nil || err.Error() != expected {
		t.Errorf("%s expected; got %s", expected, err)
	}
}
