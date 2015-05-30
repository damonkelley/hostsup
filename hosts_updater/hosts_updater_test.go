package hosts_updater

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/damonkelley/dm-hostsupdater/machine"
)

func TestAddHostsEntry(t *testing.T) {
	machine := machine.NewMachine("dev.dev", "192.168.0.1")

	h := NewHostsfile("testdata/hosts")
	h.AddHostsEntry(machine)

	f, _ := os.Open("testdata/hosts")
	defer f.Close()

	// Scan the file, line by line.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "dev.dev") {
			return
		}
	}

	t.Error("Expected to find dev.dev in testdata/hosts")
}

func TestRemoveHostsEntry(t *testing.T) {
	machine := machine.NewMachine("dev.dev", "192.168.0.1")
	h := NewHostsfile("testdata/hosts")

	h.RemoveHostsEntry(machine)

	f, _ := os.Open("testdata/hosts")
	defer f.Close()

	// Scan the file, line by line.
	scanner := bufio.NewScanner(f)

	// FIXME
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "dev.dev") {
			return
		}
	}

	t.Error("Expected to not find dev.dev in testdata/hosts")

}
