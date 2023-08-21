package item

import (
	"fmt"

	"github.com/cobbler/cobblerclient/client"
	client_internals "github.com/cobbler/cobblerclient/internal/client"
)

// GetImages returns all images in Cobbler.
func GetImages(c client.Client) ([]*Image, error) {
	var images []*Image

	result, err := c.Call("get_images", "-1", c.Token)
	if err != nil {
		return nil, err
	}

	for _, p := range result.([]interface{}) {
		var image Image
		decodedResult, err := client_internals.DecodeCobblerItem(p, &image)
		if err != nil {
			return nil, err
		}
		decodedImage := decodedResult.(*Image)
		images = append(images, decodedImage)
	}

	return images, nil
}

func ListImageNames(c client.Client) ([]string, error) {
	return c.GetItemNames("image")
}

// GetImage returns a single image obtained by its name.
func GetImage(c client.Client, name string) (*Image, error) {
	var image Image

	result, err := c.Call("get_image", name, c.Token)
	if err != nil {
		return &image, err
	}

	if result == "~" {
		return nil, fmt.Errorf("image %s not found", name)
	}

	decodeResult, err := client_internals.DecodeCobblerItem(result, &image)
	if err != nil {
		return nil, err
	}

	s := decodeResult.(*Image)

	return s, nil
}

// FindImage is ...
func FindImage(c client.Client) error {
	_, err := c.Call("find_image")
	return err
}

// GetImagesSince is ...
func GetImagesSince(c client.Client) error {
	_, err := c.Call("get_images_since")
	return err
}
