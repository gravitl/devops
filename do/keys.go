package do

import (
	"context"

	"github.com/digitalocean/godo"
)

func getAllSSHKeys(client *godo.Client) ([]godo.DropletCreateSSHKey, error) {
	ctx := context.Background()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	existingKeys, _, err := client.Keys.List(ctx, opt)
	if err != nil {
		return nil, err
	}
	keys := []godo.DropletCreateSSHKey{}
	key := godo.DropletCreateSSHKey{}
	for _, existingKey := range existingKeys {
		key.Fingerprint = existingKey.Fingerprint
		keys = append(keys, key)
	}
	return keys, nil
}
