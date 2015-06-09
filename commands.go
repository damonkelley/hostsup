package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/damonkelley/hostsup/hostsfile"
)

var Commands = []cli.Command{
	{
		Name:   "add",
		Usage:  "Add a hosts entry",
		Action: cmdAddEntry,
	},
	{
		Name:   "rm",
		Usage:  "Remove a hosts entry",
		Action: cmdRemoveEntry,
	},
	{
		Name:   "ls",
		Usage:  "List hosts entries",
		Action: cmdListEntries,
	},
	{
		Name:   "clean",
		Usage:  "Remove all hosts entries added by hostsup",
		Action: cmdClean,
	},
}

const fileName string = "/etc/hosts"

// Format strings for the logger.
const addFormat = "Added \"%s %s\" to %s."
const removeFormat = "Removed \"%s %s\" from %s."

// Command to add an entry to the hosts file.
func cmdAddEntry(c *cli.Context) {
	// Argument order is: <IP>, <hostname>.
	ip := c.Args().First()
	hostname := c.Args().Get(1)

	h, err := hostsfile.NewHostsfile(fileName)
	defer h.Close()
	handleHostsfileError(err)

	host := hostsfile.NewHost(ip, hostname)
	err = h.AddEntry(host)

	if err != nil {
		log.Fatal(err)
	}

	log.Infof(addFormat, host.IP, host.Hostname, fileName)
}

// Command to remove an entry from the hosts file.
func cmdRemoveEntry(c *cli.Context) {
	hostname := c.Args().First()

	h, err := hostsfile.NewHostsfile(fileName)
	defer h.Close()
	handleHostsfileError(err)

	entry := h.FindEntry(hostname)

	// If the entry cannot be found, inform the user and exit gracefully.
	// Not providing the queried hostname should not produce a non-zero exit code,
	// but execution should stop here.
	if entry == nil {
		log.Infof("Unable to find a hosts entry with a hostname %s", hostname)
		os.Exit(0)
	}

	err = h.RemoveEntry(entry)

	if err != nil {
		log.Fatal(err)
	}

	log.Infof(removeFormat, entry.IP, entry.Hostname, fileName)
}

// Command to list all entries added by hostsup.
func cmdListEntries(c *cli.Context) {
	h, _ := hostsfile.NewHostsfile(fileName, true)
	defer h.Close()
	entries := h.GetEntries()

	w := tabwriter.NewWriter(os.Stdout, 5, 1, 3, ' ', 0)
	fmt.Fprintln(w, "HOSTNAME\tIP")

	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\n", entry.Hostname, entry.IP)
	}

	w.Flush()
}

// Command to remove all entries added by hostsup.
func cmdClean(c *cli.Context) {
	h, err := hostsfile.NewHostsfile(fileName)
	defer h.Close()
	handleHostsfileError(err)

	entries := h.Clean()

	for _, entry := range entries {
		log.Infof(removeFormat, entry.Hostname, entry.IP, fileName)
	}

	log.Info("Hosts file has been cleaned.")
}

// Report the correct error if the hosts file was not able to be opened.
func handleHostsfileError(err error) {
	if os.IsNotExist(err) {
		log.Fatalf("The file %s does not exists on your system.", fileName)
	}

	if os.IsPermission(err) {
		log.Fatal("You do not have permission to edit this file. Try reissuing the command with sudo.")
	}
}
