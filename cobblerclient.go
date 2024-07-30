/*
Copyright 2015 Container Solutions

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cobblerclient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/kolo/xmlrpc"
)

const bodyTypeXML = "text/xml"

// HTTPClient is the interface which defines the API required for the [Client] to work correctly. Normally this
// is satisfied by a [http.DefaultClient].
type HTTPClient interface {
	Post(string, string, io.Reader) (*http.Response, error)
}

// Client is the type which all API methods are attached to.
type Client struct {
	httpClient HTTPClient
	config     ClientConfig
	// The longevity of this token is defined server side in the setting "auth_token_duration". Per default no token is
	// retrieved. A token can be obtained via the [Client.Login] method.
	Token string
}

// ClientConfig is the URL of Cobbler plus login credentials.
type ClientConfig struct {
	URL      string
	Username string
	Password string
}

type ExtendedVersion struct {
	Gitdate      string
	Gitstamp     string
	Builddate    string
	Version      string
	VersionTuple []int
}

// NewClient creates a [Client] struct which is ready for usage.
func NewClient(httpClient HTTPClient, c ClientConfig) Client {
	return Client{
		httpClient: httpClient,
		config:     c,
	}
}

// Call is the generic method for calling an XML-RPC endpoint in Cobbler that has no dedicated method in the client.
// Normally there should be no need to use this if you are just using the client.
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

// GenerateAutoinstall generates the autoinstallation file for a given profile or system.
func (c *Client) GenerateAutoinstall(profile string, system string) (string, error) {
	result, err := c.Call("generate_autoinstall", profile, system)
	if err != nil {
		return "", err
	} else {
		return result.(string), err
	}
}

// LastModifiedTime retrieves the timestamp when any object in Cobbler was last modified.
func (c *Client) LastModifiedTime() (float64, error) {
	result, err := c.Call("last_modified_time")
	if err != nil {
		return 0.0, err
	} else {
		return result.(float64), err
	}
}

// Ping is a simple method to check if the XML-RPC API is available.
func (c *Client) Ping() (bool, error) {
	result, err := c.Call("ping")
	if err != nil {
		return false, err
	} else {
		return result.(bool), err
	}
}

// AutoAddRepos automatically imports any repos server side that are known to the daemon. It is the responsitbility
// of the caller to execute [Client.BackgroundReposync].
func (c *Client) AutoAddRepos() error {
	_, err := c.Call("auto_add_repos", c.Token)
	return err
}

// GetAutoinstallTemplates retrieves a list of all templates that are in use by Cobbler.
func (c *Client) GetAutoinstallTemplates() error {
	_, err := c.Call("get_autoinstall_templates", c.Token)
	return err
}

// GetAutoinstallSnippets retrieves a list of all snippets that are in use by Cobbler.
func (c *Client) GetAutoinstallSnippets() error {
	_, err := c.Call("get_autoinstall_snippets", c.Token)
	return err
}

// IsAutoinstallInUse checks if a given system has reported that it is currently installing.
func (c *Client) IsAutoinstallInUse(name string) error {
	_, err := c.Call("is_autoinstall_in_use", name, c.Token)
	return err
}

// GenerateIPxe generates the iPXE (formerly gPXE) configuration data.
func (c *Client) GenerateIPxe(profile, image, system string) error {
	_, err := c.Call("generate_ipxe", profile, image, system)
	return err
}

// GenerateBootCfg generates the bootcfg for a given MS Windows profile or system.
func (c *Client) GenerateBootCfg(profile, system string) error {
	_, err := c.Call("generate_bootcfg", profile, system)
	return err
}

// GenerateScript generates for either a profile or sytem the requested script.
func (c *Client) GenerateScript(profile, system, name string) error {
	_, err := c.Call("generate_script", profile, system, name)
	return err
}

// GetBlendedData passes a profile or system through Cobblers inheritance chain and returns the result.
func (c *Client) GetBlendedData(profile, system string) error {
	_, err := c.Call("get_blended_data", profile, system)
	return err
}

// GetSettings returns the currently active settings.
func (c *Client) GetSettings() error {
	_, err := c.Call("get_settings", c.Token)
	return err
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

// RegisterNewSystem registers a new system without a Cobbler token. This is normally called
// during unattended installation by a script.
func (c *Client) RegisterNewSystem(info map[string]interface{}) error {
	_, err := c.Call("register_new_system", info, c.Token)
	return err
}

// RunInstallTriggers runs installation triggers for a given object. This is normally called during
// unattended installation.
func (c *Client) RunInstallTriggers(mode string, objtype string, name string, ip string) error {
	_, err := c.Call("run_install_triggers", mode, objtype, name, ip, c.Token)
	return err
}

// Version is a shorter and easier version representation. Normally you want to call [Client.ExtendedVersion].
func (c *Client) Version() (float64, error) {
	res, err := c.Call("version")
	return res.(float64), err
}

// ExtendedVersion returns the version information of the server.
func (c *Client) ExtendedVersion() (ExtendedVersion, error) {
	extendedVersion := ExtendedVersion{}
	data, err := c.Call("extended_version")
	if err != nil {
		return extendedVersion, err
	}
	switch data.(type) {
	case map[string]interface{}:
		data := data.(map[string]interface{})
		var versionTuple, err = returnIntSlice(data["version_tuple"], err)
		if err != nil {
			return extendedVersion, err
		}
		extendedVersion.Version = data["version"].(string)
		extendedVersion.VersionTuple = versionTuple
		extendedVersion.Builddate = data["builddate"].(string)
		extendedVersion.Gitdate = data["gitdate"].(string)
		extendedVersion.Gitstamp = data["gitstamp"].(string)
	default:
		return extendedVersion, err
	}
	return extendedVersion, err
}

// GetReposCompatibleWithProfile returns all repositories that can be potentially assigned to a given profile.
func (c *Client) GetReposCompatibleWithProfile(profile_name string) error {
	_, err := c.Call("get_repos_compatible_with_profile", profile_name, c.Token)
	return err
}

// FindSystemByDnsName searches for a system with a given DNS name.
func (c *Client) FindSystemByDnsName(dns_name string) error {
	_, err := c.Call("find_system_by_dns_name", dns_name)
	return err
}

// GetRandomMac generates a random MAC address for use with a virtualized system.
func (c *Client) GetRandomMac() error {
	_, err := c.Call("get_random_mac")
	return err
}

// XmlRpcHacks is an internal endpoint that doesn't make sense to be called externally.
func (c *Client) XmlRpcHacks(data interface{}) error {
	// FIXME: Make private server-side and remove from here.
	_, err := c.Call("xmlrpc_hacks", data)
	return err
}

// GetStatus retrieves the current status of installation progress that has been reported to Cobbler.
func (c *Client) GetStatus(mode string) error {
	_, err := c.Call("get_status", mode, c.Token)
	return err
}

// SyncDhcp updates the DHCP configuration synchronous.
func (c *Client) SyncDhcp() error {
	_, err := c.Call("sync_dhcp", c.Token)
	return err
}

// GetConfigData retrieves configuration data for a given host.
func (c *Client) GetConfigData(hostname string) error {
	_, err := c.Call("get_config_data", hostname)
	return err
}

// cobblerDataHacks is a hook for the mapstructure decoder. It's only used by
// decodeCobblerItem and should never be invoked directly.
// It's used to smooth out issues with converting fields and types from Cobbler.
func cobblerDataHacks(f, t reflect.Kind, data interface{}) (interface{}, error) {
	dataVal := reflect.ValueOf(data)

	// Cobbler uses ~ internally to mean None/nil
	if dataVal.String() == "~" {
		return map[string]interface{}{}, nil
	}

	if f == reflect.Int64 && t == reflect.Bool {
		if dataVal.Int() > 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return data, nil
}

// decodeCobblerItem is a custom mapstructure decoder to handler Cobbler's uniqueness.
func decodeCobblerItem(raw interface{}, result interface{}) (interface{}, error) {
	var metadata mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         &metadata,
		Result:           result,
		WeaklyTypedInput: true,
		DecodeHook:       cobblerDataHacks,
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	return result, nil
}

// updateCobblerFields updates all fields in a Cobbler Item structure.
func (c *Client) updateCobblerFields(what string, item reflect.Value, id string) error {
	method := fmt.Sprintf("modify_%s", what)

	typeOfT := item.Type()
	// In Cobbler v3.3.0, if profile name isn't created first, an empty child gets written to the distro, which causes
	// a ValueError: "calling find with no arguments"  TO-DO: figure a more efficient way of targeting name.
	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		tag := typeOfT.Field(i).Tag
		field := tag.Get("mapstructure")
		if method == "modify_profile" && field == "name" {
			var value interface{}
			switch v.Type().String() {
			case "string", "bool", "int64", "int":
				value = v.Interface()
			case "[]string":
				value = strings.Join(v.Interface().([]string), " ")
			}
			_, err := c.Call(method, id, field, value, c.Token)
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		tag := typeOfT.Field(i).Tag
		field := tag.Get("mapstructure")
		cobblerTag := tag.Get("cobbler")

		if cobblerTag == "noupdate" {
			continue
		}

		if field == "" {
			continue
		}
		var value interface{}
		switch v.Type().String() {
		case "string", "bool", "int64", "int":
			value = v.Interface()
		case "[]string":
			value = strings.Join(v.Interface().([]string), " ")
		}
		if result, err := c.Call(method, id, field, value, c.Token); err != nil {
			return err
		} else {
			if result.(bool) == false && value != false {
				// It's possible this is a new field that isn't available on
				// older versions.
				if cobblerTag == "newfield" {
					continue
				}
				return fmt.Errorf("error updating %s to %s", field, value)
			}
		}
	}
	return nil
}
