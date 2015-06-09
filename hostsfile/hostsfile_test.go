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
	defer f.Seek(0, 0)

	if err != nil {
		return nil, err
	}

	contents, _ := ioutil.ReadFile("testdata/hosts")
	f.Write(contents)

	return f, nil
}

func convertToReadOnlyHostsFile(f *os.File) (*os.File, error) {
	f.Close()
	f, err := os.OpenFile(f.Name(), os.O_RDONLY, 0666)

	if err != nil {
		return nil, err
	}

	return f, err
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
			return // Return as soon as we match the hostname.
		}
	}

	t.Error(fmt.Sprintf("Expected to find %s in testdata/hosts", hostname))
}

func TestAddEntryReturnsError(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

	f, _ := createTestHostsFile()
	f, _ = convertToReadOnlyHostsFile(f)
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	err := h.AddEntry(host)

	if err == nil {
		t.Error("Expected an error to be returned from AddEntry.")
	}
}

func TestRemoveEntryRemovesEntry(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	h.AddEntry(host)
	h.RemoveEntry(host)

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "dev.dev") {
			return
		}
	}

	t.Error(fmt.Sprintf("Expected to find %s in testdata/hosts", hostname))
}

func TestRemoveEntryReturnsError(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	h.AddEntry(host)

	f, _ = convertToReadOnlyHostsFile(f)
	err := h.RemoveEntry(host)

	if err == nil {
		t.Error("Expected an error to be returned from RemoveEntry.")
	}
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

func TestFindEntryDoesNotFindsEntry(t *testing.T) {
	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	entry := h.FindEntry("dev.dev")

	if entry != nil {
		t.Error(fmt.Sprintf("Unexpectedly found a host entry."))
	}
}

func TestGetEntriesReturnsEntries(t *testing.T) {
	hostname1, ip1 := "dev1.dev", "192.168.0.1"
	host1 := NewHost(ip1, hostname1)

	hostname2, ip2 := "dev2.dev", "192.168.0.2"
	host2 := NewHost(ip2, hostname2)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.AddEntry(host1)
	h.AddEntry(host2)

	hosts := h.GetEntries()

	if numHosts := len(hosts); numHosts != 2 {
		t.Error(fmt.Sprintf("Expected to find 2 host entries. Found %d instead.", numHosts))
	}
}

func TestGetEntriesReturnsDuplicateEntries(t *testing.T) {
	hostname, ip := "dev.dev", "192.168.0.1"
	host := NewHost(ip, hostname)

	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.AddEntry(host)
	h.AddEntry(host)
	h.AddEntry(host)

	hosts := h.GetEntries()

	if numHosts := len(hosts); numHosts != 3 {
		t.Error(fmt.Sprintf("Expected to find 3 host entries. Found %d instead.", numHosts))
	}
}

func TestListEntriesReturnsNoEntries(t *testing.T) {
	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}
	hosts := h.GetEntries()

	if numHosts := len(hosts); numHosts != 0 {
		t.Error(fmt.Sprintf("Expected to find 0 host entries. Found %d instead.", numHosts))
	}

}

func TestCleanReturnsAllEntries(t *testing.T) {
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
	hosts := h.GetEntries()

	if len(hosts) != 0 {
		t.Error(fmt.Sprintf("Expected to find 2 host entry. Found %d instead.", len(hosts)))
	}
}

func TestCleanWithNoEntries(t *testing.T) {
	f, _ := createTestHostsFile()
	defer remove(f)

	h := Hostsfile{f.Name(), f}

	h.Clean()
	hosts := h.GetEntries()

	if numHosts := len(hosts); numHosts != 0 {
		t.Error(fmt.Sprintf("Expected to find 2 host entry. Found %d instead.", numHosts))
	}
}
