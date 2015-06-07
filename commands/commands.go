package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/damonkelley/hostsup/hosts_updater"
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
		Action: cmdListEntry,
	},
	{
		Name:   "clean",
		Usage:  "Remove all hosts entries added by hostsup",
		Action: cmdClean,
	},
}

var log = logrus.New()

const HOSTSFILE string = "/etc/hosts"

func cmdAddEntry(c *cli.Context) {
	hostname := c.Args().Get(1)
	ip := c.Args().First()

	h, err := hosts_updater.NewHostsfile(HOSTSFILE)
	handleHostsfileError(err)

	host := hosts_updater.NewHost(hostname, ip)
	h.AddEntry(host)
}

func cmdRemoveEntry(c *cli.Context) {
	hostname := c.Args().First()

	h, err := hosts_updater.NewHostsfile(HOSTSFILE)
	handleHostsfileError(err)

	host := hosts_updater.NewHost(hostname, "")

	// TODO: Add a FindHost struct method to find an create a Host by the hostname.
	// Then use the returned host as a parameter to RemoveEntry.
	h.RemoveEntry(host)

}

func cmdListEntry(c *cli.Context) {
	h, _ := hosts_updater.NewHostsfile(HOSTSFILE, true)
	entries := h.ListEntries()

	w := tabwriter.NewWriter(os.Stdout, 5, 1, 3, ' ', 0)
	fmt.Fprintln(w, "HOSTNAME\tIP")

	for _, entry := range entries {
		fmt.Fprintf(w, "%s\t%s\n", entry.Hostname, entry.IP)
	}

	w.Flush()
}

func cmdClean(c *cli.Context) {
	h, err := hosts_updater.NewHostsfile(HOSTSFILE)
	handleHostsfileError(err)

	entries := h.Clean()

	for _, entry := range entries {
		log.Infof("Removed %s\t%s", entry.Hostname, entry.IP)
	}

	log.Info("Hosts file has been cleaned.")
}

func handleHostsfileError(err error) {
	if os.IsNotExist(err) {
		log.Fatalf("The file %s does not exists on your system.", HOSTSFILE)
	}

	if os.IsPermission(err) {
		log.Fatal("You do not have permission to edit this file. Try reissueing the command with sudo.")
	}
}
