package hosts_updater

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/damonkelley/hostsup/host"
)

func TestAddHostsEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := host.NewHost(hostname, ip)

	h := NewHostsfile("testdata/hosts")
	h.AddHostsEntry(host)

	f, _ := os.Open("testdata/hosts")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "dev.dev") {
			return
		}
	}

	t.Error(fmt.Sprintf("Expected to find %s in testdata/hosts", hostname))
}

func TestRemoveHostsEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := host.NewHost(hostname, ip)

	h := NewHostsfile("testdata/hosts")
	h.RemoveHostsEntry(host)

	f, _ := os.Open("testdata/hosts")
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "dev.dev") {
			return
		}
	}

	t.Error(fmt.Sprintf("Expected to find %s in testdata/hosts", hostname))
}

func TestListHostsEntries(t *testing.T) {
	hostname1, ip1 := "dev1.dev", "192.168.0.1"
	host1 := host.NewHost(hostname1, ip1)

	hostname2, ip2 := "dev2.dev", "192.168.0.2"
	host2 := host.NewHost(hostname2, ip2)

	h := NewHostsfile("testdata/hosts")
	h.AddHostsEntry(host1)
	h.AddHostsEntry(host2)

	hosts := h.ListHostsEntries()

	if len(hosts) != 2 {
		t.Error(fmt.Sprintf("Expected to find 2 host entry. Found %d instead.", len(hosts)))
	}
}
