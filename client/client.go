package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kolo/xmlrpc"
)

const bodyTypeXML = "text/xml"

// HTTPClient is ...
type HTTPClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

// Client is ...
type Client struct {
	httpClient HTTPClient
	config     ClientConfig
	Token      string
}

// ClientConfig is the URL of Cobbler plus login credentials.
type ClientConfig struct {
	URL      string
	Username string
	Password string
}

// NewClient returns a new instance of the Cobbler Client.
func NewClient(httpClient HTTPClient, c ClientConfig) Client {
	return Client{
		httpClient: httpClient,
		config:     c,
	}
}

// Call is the general wrapper which makes XML-RPC calls to the Cobbler API possible.
func (c *Client) Call(method string, args ...interface{}) (interface{}, error) {
	var result interface{}

	reqBody, err := xmlrpc.EncodeMethodCall(method, args...)
	if err != nil {
		return nil, err
	}

	r := fmt.Sprintf("%s\n", string(reqBody))
	res, err := c.httpClient.Post(c.config.URL, bodyTypeXML, bytes.NewReader([]byte(r)))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resp := xmlrpc.Response(body)
	if err := resp.Unmarshal(&result); err != nil {
		return nil, err
	}

	if err := resp.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// GenerateAutoinstall is ...
func (c *Client) GenerateAutoinstall(profile string, system string) (string, error) {
	result, err := c.Call("generate_autoinstall", profile, system)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// LastModifiedTime is ...
func (c *Client) LastModifiedTime() (float64, error) {
	result, err := c.Call("last_modified_time")
	if err != nil {
		return 0.0, err
	} else {
		return result.(float64), err
	}
}

// Ping is ...
func (c *Client) Ping() (bool, error) {
	result, err := c.Call("ping")
	if err != nil {
		return false, err
	} else {
		return result.(bool), err
	}
}

// AutoAddRepos is ...
func (c *Client) AutoAddRepos() error {
	_, err := c.Call("auto_add_repos", c.Token)
	return err
}

// GetAutoinstallTemplates is ...
func (c *Client) GetAutoinstallTemplates() error {
	_, err := c.Call("get_autoinstall_templates", c.Token)
	return err
}

// GetAutoinstallSnippets is ...
func (c *Client) GetAutoinstallSnippets() error {
	_, err := c.Call("get_autoinstall_snippets", c.Token)
	return err
}

// IsAutoinstallInUse is ...
func (c *Client) IsAutoinstallInUse(name string) error {
	_, err := c.Call("is_autoinstall_in_use", name, c.Token)
	return err
}

// GenerateIPxe is ...
func (c *Client) GenerateIPxe(profile, image, system string) error {
	_, err := c.Call("generate_ipxe", profile, image, system)
	return err
}

// GenerateBootCfg is ...
func (c *Client) GenerateBootCfg(profile, system string) error {
	_, err := c.Call("generate_bootcfg", profile, system)
	return err
}

// GenerateScript is ...
func (c *Client) GenerateScript(profile, system, name string) error {
	_, err := c.Call("generate_script", profile, system, name)
	return err
}

// GetBlendedData is ...
func (c *Client) GetBlendedData(profile, system string) error {
	_, err := c.Call("get_blended_data", profile, system)
	return err
}

// GetSettings is ...
func (c *Client) GetSettings() error {
	_, err := c.Call("get_settings", c.Token)
	return err
}

// RegisterNewSystem is ...
func (c *Client) RegisterNewSystem(info map[string]interface{}) error {
	_, err := c.Call("register_new_system", info, c.Token)
	return err
}

// RunInstallTriggers is ...
func (c *Client) RunInstallTriggers(mode string, objtype string, name string, ip string) error {
	_, err := c.Call("run_install_triggers", mode, objtype, name, ip, c.Token)
	return err
}

// Version is ...
func (c *Client) Version() (float64, error) {
	res, err := c.Call("version")
	return res.(float64), err
}

// ExtendedVersion is ...
func (c *Client) ExtendedVersion() error {
	_, err := c.Call("extended_version")
	return err
}

// GetReposCompatibleWithProfile is ...
func (c *Client) GetReposCompatibleWithProfile(profile_name string) error {
	_, err := c.Call("get_repos_compatible_with_profile", profile_name, c.Token)
	return err
}

// FindSystemByDnsName is ...
func (c *Client) FindSystemByDnsName(dns_name string) error {
	_, err := c.Call("find_system_by_dns_name", dns_name)
	return err
}

// GetRandomMac is ...
func (c *Client) GetRandomMac() error {
	_, err := c.Call("get_random_mac")
	return err
}

// XmlRpcHacks is ...
func (c *Client) XmlRpcHacks(data interface{}) error {
	_, err := c.Call("xmlrpc_hacks", data)
	return err
}

// GetStatus is ...
func (c *Client) GetStatus(mode string) error {
	_, err := c.Call("get_status", mode, c.Token)
	return err
}

// SyncDhcp is ...
func (c *Client) SyncDhcp() error {
	_, err := c.Call("sync_dhcp", c.Token)
	return err
}

// GetConfigData is ...
func (c *Client) GetConfigData(hostname string) error {
	_, err := c.Call("get_config_data", hostname)
	return err
}
