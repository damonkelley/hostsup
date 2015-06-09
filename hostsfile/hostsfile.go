package hostsfile

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"
)

const entryTemplate string = "\n%s\t%s\t# HOSTSUP %s"
const entryTag string = "HOSTSUP"

type Hostsfile struct {
	Filename string
	File     *os.File
}

// Instantiate a new Hostsfile.
func NewHostsfile(filename string, ro ...bool) (*Hostsfile, error) {
	var f *os.File
	var err error

	// Determine if the hosts file should opened as read only or not.
	if len(ro) > 0 && ro[0] == true {
		f, err = os.OpenFile(filename, os.O_RDONLY, 0666)

	} else {
		f, err = os.OpenFile(filename, os.O_RDWR, 0666)
	}

	if err != nil {
		return nil, err
	}

	return &Hostsfile{filename, f}, nil
}

// Close the hostsfile.
func (h *Hostsfile) Close() error {
	return h.File.Close()
}

// Add a new entry to the hosts file. The entry resembles:
// <IP Address>	<hostname>	# HOSTSUP <md5 sum>
func (h *Hostsfile) AddEntry(host *Host) error {
	defer h.File.Seek(0, 0)

	// Go the end of the file to append the new host entry.
	h.File.Seek(0, 2)

	entry := fmt.Sprintf(entryTemplate, host.IP, host.Hostname, host.Id)

	if _, err := h.File.WriteString(entry); err != nil {
		return errors.New("Unable to write entry to hosts file.")
	}

	return nil
}

// Remove an entry from the hosts file.
func (h *Hostsfile) RemoveEntry(host *Host) error {
	defer h.File.Seek(0, 0)

	// Read the contents of the hosts file.
	f, err := ioutil.ReadAll(h.File)

	if err != nil {
		return errors.New("Unable to read hosts file.")
	}

	// Explode the contents into a slice by line. Create an empty slice for
	// the modified file contents.
	lines := strings.Split(string(f), "\n")
	updatedLines := []string{}

	// Filter out any lines of the file that contain the Host ID.
	for _, line := range lines {
		if !strings.Contains(line, host.Id) {
			updatedLines = append(updatedLines, line)
		}
	}

	// Implode the lines and write them back to the file.
	output := strings.Join(updatedLines, "\n")
	err = ioutil.WriteFile(h.Filename, []byte(output), 0666)

	if err != nil {
		return errors.New("Unable to remove entry from hosts file.")
	}

	return nil
}

// Find the first host entry that matches a hostname.
func (h *Hostsfile) FindEntry(hostname string) *Host {
	entries := h.GetEntries()

	for _, entry := range entries {
		if entry.Hostname == hostname {
			return entry
		}
	}

	return nil
}

// Get a list of all the host entries that were added. This will exclude any
// entry that was not added by AddEntry().
func (h *Hostsfile) GetEntries() []*Host {
	defer h.File.Seek(0, 0)

	var (
		ipIndex       = 0
		hostnameIndex = 1
		idIndex       = 2
	)

	// Set up the Reader for a tab delimited file.
	reader := csv.NewReader(h.File)
	reader.Comma, _ = utf8.DecodeRuneInString("\t")
	reader.Comment, _ = utf8.DecodeRuneInString("#")
	reader.FieldsPerRecord = -1

	// Get all lines. Sets lines to [][]string.
	lines, _ := reader.ReadAll()
	hosts := make([]*Host, 0)

	for _, line := range lines {
		// Verify that the line contains the entryTag. Hostsup entries will
		// always have 3 columns.
		if len(line) >= 3 && strings.Contains(line[idIndex], entryTag) {
			host := NewHost(line[ipIndex], line[hostnameIndex])
			hosts = append(hosts, host)
		}
	}

	return hosts
}

// Remove all entries added by AddEntries.
func (h *Hostsfile) Clean() []*Host {
	entries := h.GetEntries()

	for _, entry := range entries {
		h.RemoveEntry(entry)
	}

	return entries
}
