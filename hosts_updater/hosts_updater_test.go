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
