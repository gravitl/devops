package netmaker

import (
	"github.com/guumaster/hostctl/pkg/file"
	"github.com/guumaster/hostctl/pkg/types"
)

// AddDNSEntry
func AddDNSEntry(ip, host, profile string, windows bool) error {
	etchosts := "/etc/hosts"
	if windows {
		etchosts = "c:\\windows\\system32\\drivers\\etc\\hosts"
	}
	route := types.NewRoute(ip, host)
	hosts, err := file.NewFile(etchosts)
	if err != nil {
		return err
	}
	if err := hosts.AddRoute(profile, route); err != nil {
		return err
	}
	if err := hosts.Flush(); err != nil {
		return err
	}
	return nil
}

func DeleteDNSEntry(host, profile string, windows bool) error {
	var todelete []string
	etchosts := "/etc/hosts"
	if windows {
		etchosts = "c:\\windows\\system32\\drivers\\etc\\hosts"
	}
	hosts, err := file.NewFile(etchosts)
	if err != nil {
		return err
	}
	section, err := hosts.GetProfile(profile)
	if err != nil {
		return err
	}
	todelete = append(todelete, host)
	section.RemoveHostnames(todelete)
	if err := hosts.ReplaceProfile(section); err != nil {
		return err
	}
	hosts.Flush()
	return nil
}
