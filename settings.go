package cobblerclient

import (
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"reflect"
)

type BootloaderFormatSettings struct {
	BinaryName      string   `json:"binary_name" mapstructure:"binary_name"`
	ExtraModules    []string `json:"extra_modules" mapstructure:"extra_modules"`
	ModuleDirectory string   `json:"mod_dir" mapstructure:"mod_dir"`
}

type Settings struct {
	AllowDuplicateHostnames           bool                                  `json:"allow_duplicate_hostnames" mapstructure:"allow_duplicate_hostnames"`
	AllowDuplicateIPs                 bool                                  `json:"allow_duplicate_ips" mapstructure:"allow_duplicate_ips"`
	AllowDuplicateMACs                bool                                  `json:"allow_duplicate_macs" mapstructure:"allow_duplicate_macs"`
	AllowDynamicSettings              bool                                  `json:"allow_dynamic_settings" mapstructure:"allow_dynamic_settings"`
	AlwaysWriteDhcpEntries            bool                                  `json:"always_write_dhcp_entries" mapstructure:"always_write_dhcp_entries"`
	AnamonEnabled                     bool                                  `json:"anamon_enabled" mapstructure:"anamon_enabled"`
	AuthTokenExpiration               int                                   `json:"auth_token_expiration" mapstructure:"auth_token_expiration"`
	AuthnPamService                   string                                `json:"authn_pam_service" mapstructure:"authn_pam_service"`
	AutoMigrateSettings               bool                                  `json:"auto_migrate_settings" mapstructure:"auto_migrate_settings"`
	Autoinstall                       string                                `json:"autoinstall" mapstructure:"autoinstall"`
	AutoinstallSnippetsDirectory      string                                `json:"autoinstall_snippets_dir" mapstructure:"autoinstall_snippets_dir"`
	AutoinstallTemplatesDirectory     string                                `json:"autoinstall_templates_dir" mapstructure:"autoinstall_templates_dir"`
	BUildReportingIgnorelist          []string                              `json:"build_reporting_ignorelist" mapstructure:"build_reporting_ignorelist"`
	BindChrootPath                    string                                `json:"bind_chroot_path" mapstructure:"bind_chroot_path"`
	BindManageIpmi                    bool                                  `json:"bind_manage_ipmi" mapstructure:"bind_manage_ipmi"`
	BindMaster                        string                                `json:"bind_master" mapstructure:"bind_master"`
	BindZonefilePath                  string                                `json:"bind_zonefile_path" mapstructure:"bind_zonefile_path"`
	BootloaderConfTemplateDirectory   string                                `json:"boot_loader_conf_template_dir" mapstructure:"boot_loader_conf_template_dir"`
	BootloaderDirectory               string                                `json:"bootloaders_dir" mapstructure:"bootloaders_dir"`
	BootloaderFormats                 map[string][]BootloaderFormatSettings `json:"bootloaders_formats" mapstructure:"bootloaders_formats"`
	BootloaderIPXEFolder              string                                `json:"bootloaders_ipxe_folder" mapstructure:"bootloaders_ipxe_folder"`
	BootloaderModules                 []string                              `json:"bootloaders_modules" mapstructure:"bootloaders_modules"`
	BootloaderShimDirectory           string                                `json:"bootloaders_shim_folder" mapstructure:"bootloaders_shim_folder"`
	BootloaderShimFile                string                                `json:"bootloaders_shim_file" mapstructure:"bootloaders_shim_file"`
	BuildReportingEmail               []string                              `json:"build_reporting_email" mapstructure:"build_reporting_email"`
	BuildReportingEnabled             bool                                  `json:"build_reporting_enabled" mapstructure:"build_reporting_enabled"`
	BuildReportingSender              string                                `json:"build_reporting_sender" mapstructure:"build_reporting_sender"`
	BuildReportingSmtpServer          string                                `json:"build_reporting_smtp_server" mapstructure:"build_reporting_smtp_server"`
	BuildReportingSubject             string                                `json:"build_reporting_subject" mapstructure:"build_reporting_subject"`
	BuildisoDirectory                 string                                `json:"buildisodir" mapstructure:"buildisodir"`
	CacheEnabled                      bool                                  `json:"cache_enabled" mapstructure:"cache_enabled"`
	CheetahImportWhitelist            []string                              `json:"cheetah_import_whitelist" mapstructure:"cheetah_import_whitelist"`
	ClientUseHttps                    bool                                  `json:"client_use_https" mapstructure:"client_use_https"`
	ClientUseLocalhost                bool                                  `json:"client_use_localhost" mapstructure:"client_use_localhost"`
	CobblerMaster                     string                                `json:"cobbler_master" mapstructure:"cobbler_master"`
	ConvertServerToIp                 bool                                  `json:"convert_server_to_ip" mapstructure:"convert_server_to_ip"`
	CreaterepoFlags                   string                                `json:"createrepo_flags" mapstructure:"createrepo_flags"`
	DefaultNameServers                []string                              `json:"default_name_servers" mapstructure:"default_name_servers"`
	DefaultNameServersSearch          []string                              `json:"default_name_servers_search" mapstructure:"default_name_servers_search"`
	DefaultOwnership                  []string                              `json:"default_ownership" mapstructure:"default_ownership"`
	DefaultPasswordCrypted            string                                `json:"default_password_crypted" mapstructure:"default_password_crypted"`
	DefaultTemplateType               string                                `json:"default_template_type" mapstructure:"default_template_type"`
	DefaultVirtBridge                 string                                `json:"default_virt_bridge" mapstructure:"default_virt_bridge"`
	DefaultVirtDiskDriver             string                                `json:"default_virt_disk_driver" mapstructure:"default_virt_disk_driver"`
	DefaultVirtFileSize               float64                               `json:"default_virt_file_size" mapstructure:"default_virt_file_size"`
	DefaultVirtRam                    int                                   `json:"default_virt_ram" mapstructure:"default_virt_ram"`
	DefaultVirtType                   string                                `json:"default_virt_type" mapstructure:"default_virt_type"`
	EnableIpxe                        bool                                  `json:"enable_ipxe" mapstructure:"enable_ipxe"`
	EnableMenu                        bool                                  `json:"enable_menu" mapstructure:"enable_menu"`
	ExtraSettingsList                 []string                              `json:"extra_settings_list" mapstructure:"extra_settings_list"`
	Grub2ModDirectory                 string                                `json:"grub2_mod_dir" mapstructure:"grub2_mod_dir"`
	GrubconfigDirectory               string                                `json:"grubconfig_dir" mapstructure:"grubconfig_dir"`
	HttpPort                          int                                   `json:"http_port" mapstructure:"http_port"`
	Include                           []string                              `json:"include" mapstructure:"include"`
	IsoTemplateDirectory              string                                `json:"iso_template_dir" mapstructure:"iso_template_dir"`
	Jinja2IncludeDirectory            string                                `json:"jinja2_includedir" mapstructure:"jinja2_includedir"`
	KernelOptions                     map[string]interface{}                `json:"kernel_options" mapstructure:"kernel_options"`
	LazyStart                         bool                                  `json:"lazy_start" mapstructure:"lazy_start"`
	LdapAnonymousBind                 bool                                  `json:"ldap_anonymous_bind" mapstructure:"ldap_anonymous_bind"`
	LdapBaseDn                        string                                `json:"ldap_base_dn" mapstructure:"ldap_base_dn"`
	LdapPort                          int                                   `json:"ldap_port" mapstructure:"ldap_port"`
	LdapSearchBindDn                  string                                `json:"ldap_search_bind_dn" mapstructure:"ldap_search_bind_dn"`
	LdapSearchPasswd                  string                                `json:"ldap_search_passwd" mapstructure:"ldap_search_passwd"`
	LdapSearchPrefix                  string                                `json:"ldap_search_prefix" mapstructure:"ldap_search_prefix"`
	LdapServer                        string                                `json:"ldap_server" mapstructure:"ldap_server"`
	LdapTls                           bool                                  `json:"ldap_tls" mapstructure:"ldap_tls"`
	LdapTlsCaCertDirectory            string                                `json:"ldap_tls_cacertdir" mapstructure:"ldap_tls_cacertdir"`
	LdapTlsCaCertFile                 string                                `json:"ldap_tls_cacertfile" mapstructure:"ldap_tls_cacertfile"`
	LdapTlsCertFile                   string                                `json:"ldap_tls_certfile" mapstructure:"ldap_tls_certfile"`
	LdapTlsCipherSuite                string                                `json:"ldap_tls_cipher_suite" mapstructure:"ldap_tls_cipher_suite"`
	LdapTlsKeyFile                    string                                `json:"ldap_tls_keyfile" mapstructure:"ldap_tls_keyfile"`
	LdapTlsReqcert                    string                                `json:"ldap_tls_reqcert" mapstructure:"ldap_tls_reqcert"`
	ManageDhcp                        bool                                  `json:"manage_dhcp" mapstructure:"manage_dhcp"`
	ManageDhcpV4                      bool                                  `json:"manage_dhcp_v4" mapstructure:"manage_dhcp_v4"`
	ManageDhcpV6                      bool                                  `json:"manage_dhcp_v6" mapstructure:"manage_dhcp_v6"`
	ManageDns                         bool                                  `json:"manage_dns" mapstructure:"manage_dns"`
	ManageForwardZones                []string                              `json:"manage_forward_zones" mapstructure:"manage_forward_zones"`
	ManageGenders                     bool                                  `json:"manage_genders" mapstructure:"manage_genders"`
	ManageReverseZones                []string                              `json:"manage_reverse_zones" mapstructure:"manage_reverse_zones"`
	ManageRsync                       bool                                  `json:"manage_rsync" mapstructure:"manage_rsync"`
	ManageTftpd                       bool                                  `json:"manage_tftpd" mapstructure:"manage_tftpd"`
	MgmtClasses                       []string                              `json:"mgmt_classes" mapstructure:"mgmt_classes"`
	MgmtParameters                    map[string]interface{}                `json:"mgmt_parameters" mapstructure:"mgmt_parameters"`
	NextServerV4                      string                                `json:"next_server_v4" mapstructure:"next_server_v4"`
	NextServerV6                      string                                `json:"next_server_v6" mapstructure:"next_server_v6"`
	NoPxeWithTriggers                 bool                                  `json:"nopxe_with_triggers" mapstructure:"nopxe_with_triggers"`
	NsupdateEnabled                   bool                                  `json:"nsupdate_enabled" mapstructure:"nsupdate_enabled"`
	NsupdateLog                       string                                `json:"nsupdate_log" mapstructure:"nsupdate_log"`
	NsupdateTsigAlorithm              string                                `json:"nsupdate_tsig_algorithm" mapstructure:"nsupdate_tsig_algorithm"`
	NsupdateTsigKey                   []string                              `json:"nsupdate_tsig_key" mapstructure:"nsupdate_tsig_key"`
	PowerManagementDefaultType        string                                `json:"power_management_default_type" mapstructure:"power_management_default_type"`
	Proxies                           []string                              `json:"proxies" mapstructure:"proxies"`
	ProxyUrlExternal                  string                                `json:"proxy_url_ext" mapstructure:"proxy_url_ext"`
	ProxyUrlInternal                  string                                `json:"proxy_url_int" mapstructure:"proxy_url_int"`
	PuppetAutoSetup                   bool                                  `json:"puppet_auto_setup" mapstructure:"puppet_auto_setup"`
	PuppetCaPath                      string                                `json:"puppetca_path" mapstructure:"puppetca_path"`
	PuppetParametrizedClasses         bool                                  `json:"puppet_parameterized_classes" mapstructure:"puppet_parameterized_classes"`
	PuppetServer                      string                                `json:"puppet_server" mapstructure:"puppet_server"`
	PuppetVersion                     int                                   `json:"puppet_version" mapstructure:"puppet_version"`
	PxeJustOnce                       bool                                  `json:"pxe_just_once" mapstructure:"pxe_just_once"`
	RedhatManagementKey               string                                `json:"redhat_management_key" mapstructure:"redhat_management_key"`
	RedhatManagementPermissive        bool                                  `json:"redhat_management_permissive" mapstructure:"redhat_management_permissive"`
	RedhatManagementServer            string                                `json:"redhat_management_server" mapstructure:"redhat_management_server"`
	RegisterNewInstalls               bool                                  `json:"register_new_installs" mapstructure:"register_new_installs"`
	RemoveOldPuppetCertsAutomatically bool                                  `json:"remove_old_puppet_certs_automatically" mapstructure:"remove_old_puppet_certs_automatically"`
	ReplicateRepoRsyncOptions         string                                `json:"replicate_repo_rsync_options" mapstructure:"replicate_repo_rsync_options"`
	ReplicateRsyncOptions             string                                `json:"replicate_rsync_options" mapstructure:"replicate_rsync_options"`
	ReposyncFlags                     string                                `json:"reposync_flags" mapstructure:"reposync_flags"`
	ReposyncRsyncFlags                string                                `json:"reposync_rsync_flags" mapstructure:"reposync_rsync_flags"`
	RestartDhcp                       bool                                  `json:"restart_dhcp" mapstructure:"restart_dhcp"`
	RestartDns                        bool                                  `json:"restart_dns" mapstructure:"restart_dns"`
	RunInstallTriggers                bool                                  `json:"run_install_triggers" mapstructure:"run_install_triggers"`
	SambaDistroShare                  string                                `json:"samba_distro_share" mapstructure:"samba_distro_share"`
	ScmPushScript                     string                                `json:"scm_push_script" mapstructure:"scm_push_script"`
	ScmTrackAuthor                    string                                `json:"scm_track_author" mapstructure:"scm_track_author"`
	ScmTrackEnabled                   bool                                  `json:"scm_track_enabled" mapstructure:"scm_track_enabled"`
	ScmTrackMode                      string                                `json:"scm_track_mode" mapstructure:"scm_track_mode"`
	SerializerPrettyJson              bool                                  `json:"serializer_pretty_json" mapstructure:"serializer_pretty_json"`
	Server                            string                                `json:"server" mapstructure:"server"`
	SignPuppetCertsAutomatically      bool                                  `json:"sign_puppet_certs_automatically" mapstructure:"sign_puppet_certs_automatically"`
	SignaturePath                     string                                `json:"signature_path" mapstructure:"signature_path"`
	SignatureUrl                      string                                `json:"signature_url" mapstructure:"signature_url"`
	SyslinuxDirectory                 string                                `json:"syslinux_dir" mapstructure:"syslinux_dir"`
	SyslinuxMemdiskDirectory          string                                `json:"syslinux_memdisk_folder" mapstructure:"syslinux_memdisk_folder"`
	SyslinuxPxelinuxDirectory         string                                `json:"syslinux_pxelinux_folder" mapstructure:"syslinux_pxelinux_folder"`
	TftpbootLocation                  string                                `json:"tftpboot_location" mapstructure:"tftpboot_location"`
	VirtAutoBoot                      bool                                  `json:"virt_auto_boot" mapstructure:"virt_auto_boot"`
	WebDirectory                      string                                `json:"webdir" mapstructure:"webdir"`
	WebDirectoryWhitelist             []string                              `json:"webdir_whitelist" mapstructure:"webdir_whitelist"`
	WindowsEnabled                    bool                                  `json:"windows_enabled" mapstructure:"windows_enabled"`
	WindowsTemplateDirectory          string                                `json:"windows_template_dir" mapstructure:"windows_template_dir"`
	XmlrpcPort                        int                                   `json:"xmlrpc_port" mapstructure:"xmlrpc_port"`
	YumDistroPriority                 int                                   `json:"yum_distro_priority" mapstructure:"yum_distro_priority"`
	YumDownloaderFlags                string                                `json:"yumdownloader_flags" mapstructure:"yumdownloader_flags"`
	YumPostInstallMirror              bool                                  `json:"yum_post_install_mirror" mapstructure:"yum_post_install_mirror"`
}

// cobblerDataHacks is a hook for the mapstructure decoder. It's only used by
// decodeCobblerItem and should never be invoked directly.
// It's used to smooth out issues with converting fields and types from Cobbler.
func cobblerSettingsHacks(sourceType, targetType reflect.Kind, data interface{}) (interface{}, error) {
	dataVal := reflect.ValueOf(data)

	// Cobbler uses ~ internally to mean None/nil
	if dataVal.String() == "~" {
		return map[string]interface{}{}, nil
	}

	if sourceType == reflect.Int64 && targetType == reflect.Bool {
		if dataVal.Int() > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return data, nil
}

// decodeCobblerItem is a custom mapstructure decoder to handler Cobbler's uniqueness.
func decodeCobblerSettings(raw interface{}, result interface{}) (interface{}, error) {
	var metadata mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         &metadata,
		Result:           result,
		WeaklyTypedInput: true,
		DecodeHook:       cobblerSettingsHacks,
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSettings returns the currently active settings.
func (c *Client) GetSettings() (*Settings, error) {
	var settings Settings
	resultUnmarshalled, err := c.Call("get_settings", c.Token)

	if resultUnmarshalled == "~" {
		return nil, fmt.Errorf("settings not found")
	}

	decodeResult, err := decodeCobblerSettings(resultUnmarshalled, &settings)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Settings), nil
}

// ModifySetting modifies a settings if "allow_dynamic_settings" is turned on server side.
func (c *Client) ModifySetting(name string, value interface{}) (int, error) {
	result, err := c.Call("modify_setting", name, value, c.Token)
	if err != nil {
		return -1, err
	} else {
		convertedInteger, err := convertToInt(result)
		if err != nil {
			return -1, err
		}
		return convertedInteger, err
	}
}
