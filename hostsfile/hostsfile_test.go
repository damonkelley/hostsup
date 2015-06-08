package hostsfile

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

func TestAddEntryAddsEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

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

func TestRemoveEntryRemovesEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

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

func TestFindEntryFindsEntry(t *testing.T) {
	hostname, ip := "dev1.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.AddEntry(host)

	entry := h.FindEntry(hostname)

	if entry == nil {
		t.Error(fmt.Sprintf("Unable to find host entry %s.", hostname))
	}
}

func TestListEntriesReturnsEntries(t *testing.T) {
	hostname1, ip1 := "dev1.dev", "192.168.0.1"
	host1 := NewHost(ip1, hostname1)

	hostname2, ip2 := "dev2.dev", "192.168.0.2"
	host2 := NewHost(ip2, hostname2)

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

func TestListEntriesReturnsNoEntries(t *testing.T) {

}

func TestClean(t *testing.T) {
	hostname1, ip1 := "dev1.dev", "192.168.0.1"
	host1 := NewHost(ip1, hostname1)

	hostname2, ip2 := "dev2.dev", "192.168.0.2"
	host2 := NewHost(ip2, hostname2)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.AddEntry(host1)
	h.AddEntry(host2)

	h.Clean()
	hosts := h.ListEntries()

	if len(hosts) != 0 {
		t.Error(fmt.Sprintf("Expected to find 2 host entry. Found %d instead.", len(hosts)))
	}
}
