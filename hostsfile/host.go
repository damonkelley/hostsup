package hostsfile

import (
	"crypto/md5"
	"fmt"
)

type Host struct {
	IP, Hostname, Id string
}

func NewHost(ip, hostname string) *Host {
	id := createHostId(hostname)

	return &Host{ip, hostname, id}
}

func createHostId(hostname string) string {
	id := []byte(hostname)

	return fmt.Sprintf("%x", md5.Sum(id))
}
