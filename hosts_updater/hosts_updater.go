package hosts_updater

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/damonkelley/dm-hostsupdater/machine"
)

const entryTemplate string = "\n%s\t%s\t# %s"

type Hostsfile struct {
	Filename string
	File     *os.File
}

func handleError(err error) {
	panic(err)
}

func NewHostsfile(filename string) Hostsfile {
	f, err := os.OpenFile(filename, os.O_RDWR, 0666)

	if err != nil {
		handleError(err)
	}

	return Hostsfile{filename, f}
}

func (h Hostsfile) Close() error {
	return h.File.Close()
}

func (h Hostsfile) AddHostsEntry(machine machine.Machine) {
	// Go the end of the file to append the new host entry.
	h.File.Seek(0, 2)

	entry := fmt.Sprintf(entryTemplate, machine.IP, machine.Hostname, machine.Id)

	if _, err := h.File.WriteString(entry); err != nil {
		handleError(err)
	}
}

func (h Hostsfile) RemoveHostsEntry(machine machine.Machine) {
	h.File.Seek(0, 0)

	f, err := ioutil.ReadAll(h.File)

	if err != nil {
		handleError(err)
	}

	lines := strings.Split(string(f), "\n")
	updatedLines := []string{}

	for _, line := range lines {
		if !strings.Contains(line, machine.Id) {
			updatedLines = append(updatedLines, line)
		}
	}

	output := strings.Join(updatedLines, "\n")

	err = ioutil.WriteFile(h.Filename, []byte(output), 0666)

	if err != nil {
		handleError(err)
	}
}
