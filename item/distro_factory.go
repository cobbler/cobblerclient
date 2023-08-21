package item

import (
	"github.com/cobbler/cobblerclient/client"
)

func DistroFactory(client client.Client) Distro {
	var item = Distro{
		Item: Item{
			Client: client,
		},
	}
	return item
}
