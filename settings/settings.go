package settings

import "github.com/cobbler/cobblerclient/client"

type SettingsV340 struct {
	// Client
	client client.Client

	// Settings
	autoMigrateSettings                   bool   `mapstructure:"auto_migrate_settings"`
	autoinstallScheme                     string `mapstructure:"autoinstall_scheme"`
	allow_duplicate_hostnames             bool
	allow_duplicate_ips                   bool
	allow_duplicate_macs                  bool
	allow_dynamic_settings                bool
	always_write_dhcp_entries             bool
	anamon_enabled                        bool
	auth_token_expiration                 int64
	authn_pam_service                     string
	autoinstall_snippets_dir              string
	autoinstall_templates_dir             string
	bind_chroot_path                      string
	bind_zonefile_path                    string
	bind_master                           string
	boot_loader_conf_template_dir         string
	bootloaders_dir                       string
	bootloaders_shim_folder               string
	bootloaders_shim_file                 string
	bootloaders_ipxe_folder               string
	bootloaders_formats                   map[string]map[string]string
	bootloaders_modules                   []string
	grubconfig_dir                        string
	build_reporting_enabled               bool
	build_reporting_email                 []string
	build_reporting_ignorelist            []string
	build_reporting_sender                string
	build_reporting_smtp_server           string
	build_reporting_subject               string
	buildisodir                           string
	cheetah_import_whitelist              []string
	client_use_https                      bool
	client_use_localhost                  bool
	cobbler_master                        string
	convert_server_to_ip                  bool
	createrepo_flags                      string
	autoinstall                           string
	default_name_servers                  []string
	default_name_servers_search           []string
	default_ownership                     []string
	default_password_crypted              string
	default_template_type                 string
	default_virt_bridge                   string
	default_virt_disk_driver              string
	default_virt_file_size                float64
	default_virt_ram                      int64
	default_virt_type                     string
	enable_ipxe                           bool
	enable_menu                           bool
	extra_settings_list                   []string
	grub2_mod_dir                         string
	http_port                             int64
	iso_template_dir                      string
	jinja2_includedir                     string
	kernel_options                        map[string]interface{}
	ldap_anonymous_bind                   bool
	ldap_base_dn                          string
	ldap_port                             int64
	ldap_search_bind_dn                   string
	ldap_search_passwd                    string
	ldap_search_prefix                    string
	ldap_server                           string
	ldap_tls                              bool
	ldap_tls_cacertdir                    string
	ldap_tls_cacertfile                   string
	ldap_tls_certfile                     string
	ldap_tls_keyfile                      string
	ldap_tls_reqcert                      string
	ldap_tls_cipher_suite                 string
	bind_manage_ipmi                      bool
	manage_dhcp                           bool
	manage_dhcp_v6                        bool
	manage_dhcp_v4                        bool
	manage_dns                            bool
	manage_forward_zones                  []string
	manage_reverse_zones                  []string
	manage_genders                        bool
	manage_rsync                          bool
	manage_tftpd                          bool
	mgmt_classes                          []string
	mgmt_parameters                       map[string]interface{}
	modules                               map[string]map[string]string
	mongodb                               map[string]interface{}
	next_server_v4                        string
	next_server_v6                        string
	nsupdate_enabled                      bool
	nsupdate_log                          string
	nsupdate_tsig_algorithm               string
	nsupdate_tsig_key                     []string
	power_management_default_type         string
	proxies                               []string
	proxy_url_ext                         string
	proxy_url_int                         string
	puppet_auto_setup                     bool
	puppet_parameterized_classes          bool
	puppet_server                         string
	puppet_version                        int64
	puppetca_path                         string
	pxe_just_once                         bool
	nopxe_with_triggers                   bool
	redhat_management_permissive          bool
	redhat_management_server              string
	redhat_management_key                 string
	register_new_installs                 bool
	remove_old_puppet_certs_automatically bool
	replicate_repo_rsync_options          string
	replicate_rsync_options               string
	reposync_flags                        string
	reposync_rsync_flags                  string
	restart_dhcp                          bool
	restart_dns                           bool
	run_install_triggers                  bool
	scm_track_enabled                     bool
	scm_track_mode                        string
	scm_track_author                      string
	scm_push_script                       string
	serializer_pretty_json                bool
	server                                string
	sign_puppet_certs_automatically       bool
	signature_path                        string
	signature_url                         string
	syslinux_dir                          string
	syslinux_memdisk_folder               string
	syslinux_pxelinux_folder              string
	tftpboot_location                     string
	virt_auto_boot                        bool
	webdir                                string
	webdir_whitelist                      []string
	xmlrpc_port                           int64
	yum_distro_priority                   int64
	yum_post_install_mirror               bool
	yumdownloader_flags                   string
	windows_enabled                       bool
	windows_template_dir                  string
	samba_distro_share                    string
}

func (s *SettingsV340) Update() error {
	_, err := s.client.Call("get_settings", s.client.Token)
	return err
}

func (s *SettingsV340) modifySetting(name string, value interface{}) {
	s.client.Call("modify_setting", name, value, s.client.Token)
}

type AutoMigrateSettings struct {
	settings SettingsV340
}

func (a *AutoMigrateSettings) Get() bool {
	return a.settings.autoMigrateSettings
}

func (a *AutoMigrateSettings) Set(autoMigrateSettings bool) {
	a.settings.modifySetting("auto_migrate_settings", autoMigrateSettings)
	a.settings.autoMigrateSettings = autoMigrateSettings
}
