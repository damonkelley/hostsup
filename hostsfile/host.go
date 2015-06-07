package hostsfile

import (
	"crypto/md5"
	"fmt"
)

type Host struct {
	Id, Hostname, IP string
}

func NewHost(hostname, ip string) *Host {
	id := createHostId(hostname)

	return &Host{id, hostname, ip}
}

func createHostId(hostname string) string {
	id := []byte(hostname)

	return fmt.Sprintf("%x", md5.Sum(id))
}
