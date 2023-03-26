package cmd

import (
	"testing"

	"github.com/matryer/is"
)

func TestGetNextIP(t *testing.T) {
	is := is.New(t)
	taken := make(map[string]bool)
	taken["192.168.1.1"] = true
	taken["192.168.1.2"] = true
	taken["192.168.1.3"] = true
	taken["192.168.1.4"] = true
	taken["192.168.1.6"] = true
	taken["192.168.1.7"] = true

	next := getNextIP("192.168.1.2/24", taken)
	is.Equal(next, "192.168.1.5/24")
	next = getNextIP(next, taken)
	is.Equal(next, "192.168.1.8/24")
	next = getNextIP(next, taken)
	is.Equal(next, "192.168.1.9/24")

}
