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
	"errors"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/kolo/xmlrpc"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strings"
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
	// To allow for version dependant API calls in the client we cache the major, minor and patch version.
	CachedVersion CobblerVersion
}

// ClientConfig is the URL of Cobbler plus login credentials.
type ClientConfig struct {
	URL      string
	Username string
	Password string
}

// NewClient creates a [Client] struct which is ready for usage.
func NewClient(httpClient HTTPClient, c ClientConfig) Client {
	return Client{
		httpClient:    httpClient,
		config:        c,
		CachedVersion: CobblerVersion{},
	}
}

// Call is the generic method for calling an XML-RPC endpoint in Cobbler that has no dedicated method in the client.
// Normally there should be no need to use this if you are just using the client. In case there is an error closing the
// HTTP connection, it hides all errors that occur during the rest of the method.
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

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resp := xmlrpc.Response(body)
	if err = resp.Unmarshal(&result); err != nil {
		return nil, err
	}

	if err = resp.Err(); err != nil {
		return nil, err
	}

	// Return err because the deferred function may set it to non-nil
	return result, err
}

func (c *Client) setCachedVersion() error {
	if c.CachedVersion != (CobblerVersion{}) {
		return nil
	}
	extendedVersion, err := c.ExtendedVersion()
	if err != nil {
		return err
	}
	if len(extendedVersion.VersionTuple) != 3 {
		return errors.New("cobblerclient: invalid length of extended version tuple")
	}
	c.CachedVersion = CobblerVersion{
		Major: extendedVersion.VersionTuple[0],
		Minor: extendedVersion.VersionTuple[1],
		Patch: extendedVersion.VersionTuple[2],
	}
	return nil
}

func (c *Client) invalidateCachedVersion() {
	c.CachedVersion = CobblerVersion{}
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
func (c *Client) GetBlendedData(profile, system string) (map[string]interface{}, error) {
	result, err := c.Call("get_blended_data", profile, system)
	return result.(map[string]interface{}), err
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

// IsValueInherit safely verifies if a given value is set to the magic "<<inherit>>".
func (c *Client) IsValueInherit(value interface{}) bool {
	if value == nil {
		return false
	}
	stringValue, ok := value.(string)
	if !ok {
		return false
	}
	return stringValue == inherit
}

// cobblerDataHacks is a hook for the mapstructure decoder. It's only used by
// decodeCobblerItem and should never be invoked directly.
// It's used to smooth out issues with converting fields and types from Cobbler.
func cobblerDataHacks(fromType, targetType reflect.Kind, data interface{}) (interface{}, error) {
	dataVal := reflect.ValueOf(data)

	// Cobbler uses ~ internally to mean None/nil
	if dataVal.String() == "~" {
		switch targetType {
		case reflect.String:
			return "", nil
		case reflect.Slice:
			return []string{}, nil
		case reflect.Map:
			return map[string]interface{}{}, nil
		case reflect.Int:
			return -1, nil
		case reflect.Interface:
			return nil, nil
		case reflect.Array:
			return []string{}, nil
		case reflect.Struct:
			return Value[interface{}]{RawData: nil}, nil
		default:
			return nil, errors.New("unknown type was nil")
		}
	}

	if fromType == reflect.Int64 && targetType == reflect.Bool {
		// XML-RPC Integer Booleans
		return convertXmlRpcBool(dataVal.Interface())
	}

	if targetType == reflect.Struct {
		// This must be a value that may or may not be inherited or flattened (dual-homed types)

		switch fromType {
		case reflect.String:
			valueStruct := Value[interface{}]{}
			valueStruct.IsInherited = dataVal.String() == inherit
			valueStruct.RawData = data
			return valueStruct, nil
		case reflect.Slice:
			// Slice that may or may not be inherited
			valueStruct := Value[[]interface{}]{}
			valueStruct.RawData = data
			return valueStruct, nil
		case reflect.Map:
			// This can be: Top-level Map, paged search results, page-info struct, network interface or an inherited struct
			mapKeys := dataVal.MapKeys()
			sort.SliceStable(mapKeys, func(i, j int) bool {
				return mapKeys[i].String() < mapKeys[j].String()
			})
			if len(mapKeys) == 2 && mapKeys[0].String() == "items" && mapKeys[1].String() == "pageinfo" {
				// Paged search results
				return data, nil
			}
			if len(mapKeys) == 10 && mapKeys[0].String() == "end_item" {
				// Page-Info struct
				return data, nil
			}
			if len(mapKeys) == 23 && mapKeys[0].String() == "bonding_opts" {
				// Network Interface struct
				return data, nil
			}
			for _, key := range mapKeys {
				// If the uid key is in the map then it is the top level Map
				if key.String() == "uid" {
					return data, nil
				}
			}
			valueStruct := Value[map[string]interface{}]{}
			valueStruct.Data = make(map[string]interface{})
			valueStruct.RawData = data
			return valueStruct, nil
		case reflect.Int64:
			// Int that may or may not be inherited
			valueStruct := Value[int]{}
			integerValue, err := convertToInt(data)
			valueStruct.Data = integerValue
			valueStruct.RawData = data
			if err != nil {
				return Value[int]{}, err
			}
			return valueStruct, nil
		case reflect.Float64:
			// Float that may or may not be inherited
			valueStruct := Value[float64]{}
			floatValue, err := convertToFloat(data)
			valueStruct.Data = floatValue
			valueStruct.RawData = data
			if err != nil {
				return Value[float64]{}, err
			}
			return valueStruct, nil
		case reflect.Bool:
			// Bool that may or may not be inherited
			valueStruct := Value[bool]{}
			valueStruct.Data = data.(bool)
			valueStruct.RawData = data
			return valueStruct, nil
		default:
			return nil, fmt.Errorf("unknown type %s fromType for Inherited or Flattened Value", fromType)
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

	if err = decoder.Decode(raw); err != nil {
		return nil, err
	}

	return result, nil
}

// updateCobblerFields updates all fields in a Cobbler Item structure.
func (c *Client) updateCobblerFields(what string, item reflect.Value, id string) error {
	method := fmt.Sprintf("modify_%s", what)
	typeOfT := item.Type()

	// Update embedded Item struct
	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		fieldType := v.Type().Name()

		if fieldType == "Item" {
			err := c.updateCobblerFields(what, reflect.ValueOf(v.Interface()), id)
			if err != nil {
				return err
			}
			break
		}
	}

	// Fields that can inherit from other items can only be set after the parent is set.
	// Fields that inherit from settings can be modified without this constraint.
	if method == "modify_profile" {
		// In Cobbler v3.3.0, if profile name isn't created first, an empty child gets written to the distro, which
		// causes a ValueError: "calling find with no arguments"
		nameField := item.FieldByName("Name")
		_, err := c.Call(method, id, "name", nameField.String(), c.Token)
		if err != nil {
			return err
		}

		parentField := item.FieldByName("Parent")
		if parentField != (reflect.Value{}) {
			err = c.updateSingleField(method, id, "parent", parentField.String(), "")
			if err != nil {
				return err
			}
		}
		distroField := item.FieldByName("Distro")
		if distroField != (reflect.Value{}) {
			err = c.updateSingleField(method, id, "distro", distroField.String(), "")
			if err != nil {
				return err
			}
		}
	}
	if method == "modify_system" {
		profileField := item.FieldByName("Profile")
		if profileField != (reflect.Value{}) {
			err := c.updateSingleField(method, id, "profile", profileField.String(), "")
			if err != nil {
				return err
			}
		}
		imageField := item.FieldByName("Image")
		if imageField != (reflect.Value{}) {
			err := c.updateSingleField(method, id, "image", imageField.String(), "")
			if err != nil {
				return err
			}
		}
		interfaceField := item.FieldByName("Interfaces")
		if interfaceField != (reflect.Value{}) {
			err := c.updateInterfaces(id, interfaceField.Interface())
			if err != nil {
				return err
			}
		}
	}

	for i := 0; i < item.NumField(); i++ {
		v := item.Field(i)
		tag := typeOfT.Field(i).Tag
		fieldType := v.Type().Name()
		field := tag.Get("mapstructure")
		cobblerTag := tag.Get("cobbler")

		if cobblerTag == "noupdate" || fieldType == "Item" || fieldType == "Meta" {
			continue
		}

		if field == "" || field == "parent" || field == "distro" || field == "profile" || field == "image" || field == "interfaces" {
			// Skip fields that are empty or have been set previously
			continue
		}

		if method == "modify_profile" && field == "name" {
			// Field set above
			continue
		}

		fieldValue := v.Interface()
		if strings.HasPrefix(fieldType, "Value") {
			if v.FieldByName("IsInherited").Bool() {
				fieldValue = inherit
			} else {
				fieldValue = v.FieldByName("Data").Interface()
			}
		}

		err := c.updateSingleField(method, id, field, fieldValue, cobblerTag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) updateSingleField(method, id, field string, fieldValue interface{}, cobblerTag string) error {
	if result, err := c.Call(method, id, field, fieldValue, c.Token); err != nil {
		return err
	} else {
		if !result.(bool) && fieldValue != false {
			// It's possible this is a new field that isn't available on older versions.
			if cobblerTag == "newfield" {
				return nil
			}
			return fmt.Errorf("error updating field \"%s\" to \"%s\"", field, fieldValue)
		}
	}
	return nil
}

// updateInterfaces takes care of pushing interface modifications. Since interfaces don't have unique identifiers in
// Cobbler 3.3.x. As such no reliable tracking of operations can be done when interfaces are renamed. As such this only
// handles modification and creation of interfaces.
func (c *Client) updateInterfaces(systemId string, interfaceData interface{}) error {
	interfaceMap := interfaceData.(Interfaces)
	for name, iface := range interfaceMap {
		res := makeInterfaceOptionsMap(name, iface)
		err := c.ModifyInterface(systemId, res)
		if err != nil {
			return err
		}
	}
	return nil
}
