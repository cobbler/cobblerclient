package cobblerclient

import (
	"fmt"
	"reflect"
	"time"
)

type Architecture int64

const (
	aI386 = iota
	aX8664
	aIA64
	aPPC
	aPPC64
	aPPC64LE
	aPPC64EL
	aS390
	aS390X
	aARM
	aAARCH64
)

func (a Architecture) String() string {
	switch a {
	case aI386:
		return "i386"
	case aX8664:
		return "x86_64"
	case aIA64:
		return "ia64"
	case aPPC:
		return "ppc"
	case aPPC64:
		return "ppc64"
	case aPPC64LE:
		return "ppc64le"
	case aPPC64EL:
		return "ppc64el"
	case aS390:
		return "s390"
	case aS390X:
		return "s390x"
	case aARM:
		return "arm"
	case aAARCH64:
		return "aarch64"
	}
	return "unknown"
}

type ImageType int64

const (
	itDIRECT = iota
	itISO
	itMEMDISK
	itVIRTCLONE
)

func (i ImageType) String() string {
	switch i {
	case itDIRECT:
		return "direct"
	case itISO:
		return "iso"
	case itMEMDISK:
		return "memdisk"
	case itVIRTCLONE:
		return "virt-clone"
	}
	return "unknown"
}

type VirtType int64

const (
	itINHERITED = iota
	vtQEMU
	itKVM
	itXENPV
	itXENFV
	itVMWARE
	itVMWAREW
	itOPENVZ
	itAUTO
)

func (v VirtType) String() string {
	switch v {
	case itINHERITED:
		return "<<inherit>>"
	case vtQEMU:
		return "qemu"
	case itKVM:
		return "kvm"
	case itXENPV:
		return "xenpv"
	case itXENFV:
		return "xenfv"
	case itVMWARE:
		return "vmware"
	case itVMWAREW:
		return "vmwarew"
	case itOPENVZ:
		return "openvz"
	case itAUTO:
		return "auto"
	}
	return "unknown"
}

type VirtDiskDriver int64

const (
	vddINHERITED = iota
	vddRAW
	vddQCOW2
	vddQED
	vddVDI
	vddVDMK
)

func (v VirtDiskDriver) String() string {
	switch v {
	case vddINHERITED:
		return "<<inherit>>"
	case vddRAW:
		return "raw"
	case vddQCOW2:
		return "qcow2"
	case vddQED:
		return "qed"
	case vddVDI:
		return "vdi"
	case vddVDMK:
		return "vdmk"
	}
	return "unknown"
}

// Image is a created image.
// Get the fields from cobbler/items/image.py
type Image struct {
	Item `mapstructure:",squash"`

	// Image specific fields
	// Arch                 Architecture
	// ImageType            ImageType      `mapstructure:"image_type"`
	// VirtDiskDriver       VirtDiskDriver `mapstructure:"virt_disk_driver"`
	// VirtType             VirtType `mapstructure:"virt_type"`
	Arch                 string      `mapstructure:"arch"`
	Autoinstall          string      `mapstructure:"autoinstall"`
	Breed                string      `mapstructure:"breed"`
	File                 string      `mapstructure:"file"`
	ImageType            string      `mapstructure:"image_type"`
	NetworkCount         int         `mapstructure:"network_count"`
	OsVersion            string      `mapstructure:"os_version"`
	BootLoaders          []string    `mapstructure:"boot_loaders"`
	Menu                 string      `mapstructure:"menu"`
	VirtAutoBoot         bool        `mapstructure:"virt_auto_boot"`
	VirtBridge           string      `mapstructure:"virt_bridge"`
	VirtCpus             int         `mapstructure:"virt_cpus"`
	VirtDiskDriver       string      `mapstructure:"virt_disk_driver"`
	VirtFileSize         interface{} `mapstructure:"virt_file_size"`
	VirtPath             string      `mapstructure:"virt_path"`
	VirtRam              int         `mapstructure:"virt_ram"`
	VirtType             string      `mapstructure:"virt_type"`
	SupportedBootLoaders []string    `mapstructure:"supported_boot_loaders"`

	Client
}

func convertRawImage(name string, xmlrpcResult interface{}) (*Image, error) {
	var image Image

	if xmlrpcResult == "~" {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	decodeResult, err := decodeCobblerItem(xmlrpcResult, &image)
	if err != nil {
		return nil, err
	}

	return decodeResult.(*Image), nil
}

func convertRawImagesList(xmlrpcResult interface{}) ([]*Image, error) {
	var images []*Image

	for _, d := range xmlrpcResult.([]interface{}) {
		distro, err := convertRawImage("unknown", d)
		if err != nil {
			return nil, err
		}
		images = append(images, distro)
	}

	return images, nil
}

// GetImages returns all images in Cobbler.
func (c *Client) GetImages() ([]*Image, error) {
	result, err := c.Call("get_images", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawImagesList(result)
}

// ListImageNames returns a list of all known image names.
func (c *Client) ListImageNames() ([]string, error) {
	return c.GetItemNames("image")
}

// GetImage returns a single image obtained by its name.
func (c *Client) GetImage(name string) (*Image, error) {
	result, err := c.Call("get_image", name, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawImage(name, result)
}

// CreateImage creates an image.
func (c *Client) CreateImage(image Image) (*Image, error) {
	// To create an image via the Cobbler API, first call new_image to obtain an ID
	result, err := c.Call("new_image", c.Token)
	if err != nil {
		return nil, err
	}
	newID := result.(string)
	// Set the value of all fields
	item := reflect.ValueOf(&image).Elem()
	if err := c.updateCobblerFields("image", item, newID); err != nil {
		return nil, err
	}

	// Save the final image
	if err = c.SaveImage(newID, "new"); err != nil {
		return nil, err
	}

	// Return a clean copy of the image
	return c.GetImage(image.Name)
}

// UpdateImage updates a single image.
func (c *Client) UpdateImage(image *Image) error {
	item := reflect.ValueOf(image).Elem()
	id, err := c.GetItemHandle("image", image.Name)
	if err != nil {
		return err
	}

	if err := c.updateCobblerFields("image", item, id); err != nil {
		return err
	}

	// Save the final image
	if err := c.SaveImage(id, "bypass"); err != nil {
		return err
	}

	return nil
}

// DeleteImage deletes a single Image by its name.
func (c *Client) DeleteImage(name string) error {
	return c.DeleteImageRecursive(name, false)
}

// DeleteImageRecursive deletes a single Image by its name with the option to do so recursively.
func (c *Client) DeleteImageRecursive(name string, recursive bool) error {
	_, err := c.Call("remove_image", name, c.Token, recursive)
	return err
}

// FindImage searches for one or more images by any of its attributes.
func (c *Client) FindImage(criteria map[string]interface{}) ([]*Image, error) {
	result, err := c.Call("find_image", criteria, true, c.Token)
	if err != nil {
		return nil, err
	}

	return convertRawImagesList(result)
}

// FindImageNames searches for one or more distros by any of its attributes.
func (c *Client) FindImageNames(criteria map[string]interface{}) ([]string, error) {
	resultUnmarshalled, err := c.Call("find_image", criteria, false, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}

// GetImageHandle gets the internal ID of a Cobbler item.
func (c *Client) GetImageHandle(name string) (string, error) {
	res, err := c.Call("get_image_handle", name, c.Token)
	return returnString(res, err)
}

// CopyImage duplicates an image on the server with a new name.
func (c *Client) CopyImage(objectId, newName string) error {
	_, err := c.Call("copy_image", objectId, newName, c.Token)
	return err
}

// RenameImage renames an image with a given object id.
func (c *Client) RenameImage(objectId, newName string) error {
	_, err := c.Call("rename_image", objectId, newName, c.Token)
	return err
}

// GetImagesSince returns all images which were created after the specified date.
func (c *Client) GetImagesSince(mtime time.Time) ([]*Image, error) {
	result, err := c.Call("get_images_since", float64(mtime.Unix()))
	if err != nil {
		return nil, err
	}

	return convertRawImagesList(result)
}

// GetImageAsRendered returns the datastructure after it has passed through Cobblers inheritance structure.
func (c *Client) GetImageAsRendered(name string) (map[string]interface{}, error) {
	result, err := c.Call("get_image_as_rendered", name, c.Token)
	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), err
}

// SaveImage saves all changes performed via XML-RPC to disk on the server side.
func (c *Client) SaveImage(objectId, editmode string) error {
	_, err := c.Call("save_image", objectId, c.Token, editmode)
	return err
}

// GetValidImageBootLoaders retrieves the list of bootloaders that can be assigned to an image.
func (c *Client) GetValidImageBootLoaders(imageName string) ([]string, error) {
	resultUnmarshalled, err := c.Call("get_valid_image_boot_loaders", imageName, c.Token)
	return returnStringSlice(resultUnmarshalled, err)
}
