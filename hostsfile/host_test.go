package hostsfile

import (
	"testing"
)

func TestNewHost(t *testing.T) {
	host := NewHost("IP", "hostname")

	id := createHostId("hostname")
	stub := Host{"IP", "hostname", id}

	if *host != stub {
		t.Error("Ordering of Host fields do not match.")
	}
}
