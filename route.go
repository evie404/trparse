package trparse

import (
	"time"

	"inet.af/netaddr"
)

type Route struct {
	Hop      int
	RTT      []time.Duration
	IP       netaddr.IP
	Hostname string
}
