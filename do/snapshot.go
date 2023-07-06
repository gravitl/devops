package do

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/digitalocean/godo"
	"golang.org/x/exp/slog"
)

func ListSnapshot(name, token string) (godo.Snapshot, error) {
	client := godo.NewFromToken(token)
	ctx := context.TODO()
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	snapshots, _, err := client.Snapshots.List(ctx, opt)
	if err != nil {
		return godo.Snapshot{}, fmt.Errorf("list snapshots %v", err)
	}
	for _, snapshot := range snapshots {
		slog.Debug("snapshot", "name", snapshot.Name, "id", snapshot.ID)
		if snapshot.Name == name {
			return snapshot, nil
		}
	}
	return godo.Snapshot{}, errors.New("snapshot not found")
}

func DeleteSnapshot(id, token string) error {
	client := godo.NewFromToken(token)
	ctx := context.TODO()
	_, err := client.Snapshots.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateFromSnapshot(name, region, size, snapshot, token string) (godo.Droplet, error) {
	client := godo.NewFromToken(token)
	ctx := context.TODO()
	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: snapshot,
		},
	}
	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	if err != nil {
		return godo.Droplet{}, err
	}
	return *droplet, nil
}

func (request *Request) ListSnapshot(name string) (godo.Snapshot, error) {
	return ListSnapshot(name, request.Token)
}

func (request *Request) DeleteSnapshot(id string) error {
	return DeleteSnapshot(id, request.Token)
}

func (request *Request) CreateFromSnapshot(snapshot godo.Snapshot) error {
	client := godo.NewFromToken(request.Token)
	// get all ssh keys
	keys, err := getAllSSHKeys(client)
	if err != nil {
		return err
	}
	snapRequest := &godo.DropletCreateRequest{
		Name:    "egress",
		Tags:    request.Tags,
		Region:  snapshot.Regions[0],
		Size:    "s-1vcpu-1gb",
		SSHKeys: keys,
		Image: godo.DropletCreateImage{
			Slug: snapshot.ID,
		},
		IPv6: true,
	}
	_, response, err := client.Droplets.Create(context.Background(), snapRequest)
	if err != nil {
		return err
	}
	//wait for droplet creation to complete
	wg := sync.WaitGroup{}
	for _, action := range response.Links.Actions {
		wg.Add(1)
		go func(action godo.LinkAction) {
			defer wg.Done()
			for {
				status, _, err := client.Actions.Get(ctx, action.ID)
				if err != nil {
					continue
				}
				if status.Status == "completed" {
					break
				}
				log.Print("waiting for droplet creation to complete ... ", status.Status, " ", status.ID)
				time.Sleep(time.Second << 2)
			}
		}(action)
	}
	wg.Wait()

	return nil
}
