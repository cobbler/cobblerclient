package cobblerclient

import (
	"errors"
	"github.com/go-viper/mapstructure/v2"
	"reflect"
)

type OsVersion struct {
	Signatures          []string            `mapstructure:"signatures" json:"signatures" yaml:"signatures"`
	VersionFile         string              `mapstructure:"version_file" json:"version_file" yaml:"version_file"`
	VersionFileRegex    string              `mapstructure:"version_file_regex" json:"version_file_regex" yaml:"version_file_regex"`
	KernelArch          string              `mapstructure:"kernel_arch" json:"kernel_arch" yaml:"kernel_arch"`
	KernelArchRegex     string              `mapstructure:"kernel_arch_regex" json:"kernel_arch_regex" yaml:"kernel_arch_regex"`
	SupportedArches     []string            `mapstructure:"supported_arches" json:"supported_arches" yaml:"supported_arches"`
	SupportedRepoBreeds []string            `mapstructure:"supported_repo_breeds" json:"supported_repo_breeds" yaml:"supported_repo_breeds"`
	KernelFile          string              `mapstructure:"kernel_file" json:"kernel_file" yaml:"kernel_file"`
	InitrdFile          string              `mapstructure:"initrd_file" json:"initrd_file" yaml:"initrd_file"`
	IsolinuxOk          bool                `mapstructure:"isolinux_ok" json:"isolinux_ok" yaml:"isolinux_ok"`
	DefaultAutoinstall  string              `mapstructure:"default_autoinstall" json:"default_autoinstall" yaml:"default_autoinstall"`
	KernelOptions       string              `mapstructure:"kernel_options" json:"kernel_options" yaml:"kernel_options"`
	KernelOptionsPost   string              `mapstructure:"kernel_options_post" json:"kernel_options_post" yaml:"kernel_options_post"`
	TemplateFiles       string              `mapstructure:"template_files" json:"template_files" yaml:"template_files"`
	BootFiles           []string            `mapstructure:"boot_files" json:"boot_files" yaml:"boot_files"`
	BootLoaders         map[string][]string `mapstructure:"boot_loaders" json:"boot_loaders" yaml:"boot_loaders"`
}

type DistroSignatures struct {
	Breeds map[string]map[string]OsVersion `mapstructure:"breeds" json:"breeds" yaml:"breeds"`
}

// cobblerSignatureHacks is a hook for the mapstructure decoder. It's only used by
// decodeCobblerSignatures and should never be invoked directly.
// It's used to smooth out issues with converting fields and types from Cobbler.
func cobblerSignatureHacks(sourceType, targetType reflect.Kind, data interface{}) (interface{}, error) {
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
		default:
			return nil, errors.New("unknown type was nil")
		}

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

// decodeCobblerSignatures is a custom mapstructure decoder to handler Cobbler's uniqueness.
func decodeCobblerSignatures(raw interface{}, result interface{}) (interface{}, error) {
	var metadata mapstructure.Metadata
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         &metadata,
		Result:           result,
		WeaklyTypedInput: true,
		DecodeHook:       cobblerSignatureHacks,
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(raw); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSignatures retrieves the complete signatures that are loaded by Cobbler.
func (c *Client) GetSignatures() (*DistroSignatures, error) {
	var distroSignatures DistroSignatures
	rawSignatures, err := c.Call("get_signatures", c.Token)
	if err != nil {
		return &distroSignatures, err
	}
	decodedResult, err := decodeCobblerSignatures(rawSignatures, &distroSignatures)
	if err != nil {
		return &distroSignatures, err
	}
	return decodedResult.(*DistroSignatures), err
}

// GetValidBreeds retrieves all valid OS breeds that a distro can have.
func (c *Client) GetValidBreeds() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_breeds", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidOsVersionsForBreed retrieves all valid OS versions for a given breed.
func (c *Client) GetValidOsVersionsForBreed(breed string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_os_versions_for_breed", breed, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidOsVersions retrieves all valid OS versions that a distro can have.
func (c *Client) GetValidOsVersions() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_os_versions", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetValidArchs retrieves all valid architectures that Cobbler is offering.
func (c *Client) GetValidArchs() ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_archs", c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// BackgroundSignatureUpdate runs a signatures update in the background on the server.
func (c *Client) BackgroundSignatureUpdate() (string, error) {
	res, err := c.Call("background_signature_update", map[string]string{}, c.Token)
	return returnString(res, err)
}
