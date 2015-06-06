package hosts_updater

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func createTestHostsFile() (*os.File, error) {
	f, err := ioutil.TempFile("testdata", "hosts-")

	if err != nil {
		return nil, err
	}

	contents, _ := ioutil.ReadFile("testdata/hosts")

	f.Write(contents)
	f.Seek(0, 0)

	return f, nil
}

func remove(f *os.File) error {
	return os.Remove(f.Name())
}

func TestAddHostsEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(hostname, ip)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	h.AddEntry(host)

	// Reset the offset after AddHostsEntry changes it.
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
	host := NewHost(hostname, ip)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	h.RemoveEntry(host)

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
	host1 := NewHost(hostname1, ip1)

	hostname2, ip2 := "dev2.dev", "192.168.0.2"
	host2 := NewHost(hostname2, ip2)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.AddEntry(host1)
	h.AddEntry(host2)

	hosts := h.ListEntries()

	if len(hosts) != 2 {
		t.Error(fmt.Sprintf("Expected to find 2 host entry. Found %d instead.", len(hosts)))
	}
}
