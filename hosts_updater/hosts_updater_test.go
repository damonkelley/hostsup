package hosts_updater

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/damonkelley/dm-hostsupdater/machine"
)

func TestAddHostsEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	machine := machine.NewMachine(hostname, ip)

	h := NewHostsfile("testdata/hosts")
	h.AddHostsEntry(machine)

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
	machine := machine.NewMachine(hostname, ip)

	h := NewHostsfile("testdata/hosts")
	h.RemoveHostsEntry(machine)

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
