package machine

import (
	"crypto/md5"
	"fmt"
)

type Machine struct {
	Id, Hostname, IP string
}

func NewMachine(hostname, ip string) Machine {
	id := createMachineId(hostname)

	return Machine{id, hostname, ip}
}

func createMachineId(hostname string) string {
	id := []byte(hostname)

	return fmt.Sprintf("%x", md5.Sum(id))
}
