package do

import (
	"fmt"
	"log"
	"net"

	"github.com/digitalocean/godo"
)

//VerifyDNS - verifies that the dns records for each droplet is
//consistent with it's public ip
func (request *Request) VerifyDNS(tag string) {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		log.Fatal("could not retrieve droplets ", err)
	}
	for _, droplet := range droplets {
		fmt.Println("checking dns", droplet.Name)
		var dns string
		done := false
		count := 1
		for {
			count++
			if done || count == 10 {
				break
			}
			if request.SubDomain == "" {
				dns = droplet.Name + ".clustercat.com"
			} else {
				dns = droplet.Name + "." + request.SubDomain + ".clustercat.com"
			}
			IPs, err := net.LookupIP(dns)
			if err != nil {
				log.Println("nslookup failure ", err)
				continue
			}
			publicIP, err := droplet.PublicIPv4()
			if err != nil {
				log.Fatal("droplet ", droplet.Name, " does not have public ip ", err)
			}
			for _, ip := range IPs {
				fmt.Println(droplet.Name, publicIP, ip)
				if ip.String() == publicIP {
					fmt.Println("dnslookup matches public ip")
					done = true
				}
			}
		}
		if !done {
			log.Println("dns lookup is not returning the public ip for node ", dns)
		}
	}
}

// DeleteDNS - deletes dns records specified in the request
func (request *Request) DeleteDNS(tag string) error {
	client := godo.NewFromToken(request.Token)
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, opts)
	if err != nil {
		return fmt.Errorf("err retrieving droplets %w", err)
	}
	if len(droplets) == 0 {
		return fmt.Errorf("no droplets with tag %s", tag)
	}
	records, _, err := client.Domains.Records(ctx, "clustercat.com", opts)
	if err != nil {
		return fmt.Errorf("failed to retrieve domain records %w", err)
	}
	log.Println("deleting dns records")
	for _, droplet := range droplets {
		for _, rec := range records {
			if droplet.Name == rec.Name {
				log.Println("deleting dns record for ", rec.Name)
				_, err = client.Domains.DeleteRecord(ctx, "clustercat.com", rec.ID)
				if err != nil {
					return fmt.Errorf("failure deleting dns record for %s - %w", rec.Name, err)
				}
			}
		}
	}
	return nil
}
