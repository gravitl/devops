package do

import (
	"context"
	"errors"

	"github.com/digitalocean/godo"
)

func Name(name, tag, token string) (godo.Droplet, error) {
	client := godo.NewFromToken(token)
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opt)
	if err != nil {
		return godo.Droplet{}, err
	}
	for _, droplet := range droplets {
		if droplet.Name == name {
			return droplet, nil
		}
	}
	return godo.Droplet{}, errors.New("droplet not found")
}
