package hosts_updater

import (
	"encoding/csv"
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

func handleError(err error) {
	panic(err)
}

func NewHostsfile(filename string, ro ...bool) (*Hostsfile, error) {
	var f *os.File
	var err error

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

func (h *Hostsfile) Close() error {
	return h.File.Close()
}

func (h *Hostsfile) AddEntry(host *Host) {
	defer h.File.Seek(0, 0)

	// Go the end of the file to append the new host entry.
	h.File.Seek(0, 2)

	entry := fmt.Sprintf(entryTemplate, host.IP, host.Hostname, host.Id)

	if _, err := h.File.WriteString(entry); err != nil {
		handleError(err)
	}
}

func (h *Hostsfile) RemoveEntry(host *Host) {
	defer h.File.Seek(0, 0)

	f, err := ioutil.ReadAll(h.File)

	if err != nil {
		handleError(err)
	}

	lines := strings.Split(string(f), "\n")
	updatedLines := []string{}

	for _, line := range lines {
		if !strings.Contains(line, host.Id) {
			updatedLines = append(updatedLines, line)
		}
	}

	output := strings.Join(updatedLines, "\n")

	err = ioutil.WriteFile(h.Filename, []byte(output), 0666)

	if err != nil {
		handleError(err)
	}
}

func (h *Hostsfile) ListEntries() []*Host {
	defer h.File.Seek(0, 0)

	reader := csv.NewReader(h.File)
	tab, _ := utf8.DecodeRuneInString("\t")
	comment, _ := utf8.DecodeRuneInString("#")

	reader.Comma = tab
	reader.Comment = comment
	reader.FieldsPerRecord = -1

	lines, _ := reader.ReadAll()

	hosts := []*Host{}

	for _, line := range lines {
		// TODO: Add a check to determine if the entry was added by hostsup
		if len(line) >= 3 {
			// TODO: See if we can unpack the list to create the Host
			// host.NewHost(line...)
			host := NewHost(line[1], line[0])
			hosts = append(hosts, host)
		}
	}

	return hosts
}

func (h *Hostsfile) Clean() []*Host {
	entries := h.ListEntries()

	for _, entry := range entries {
		h.RemoveEntry(entry)
	}

	return entries
}