package trparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"inet.af/netaddr"
)

var (
	lineRegexp = regexp.MustCompile(`^\s?(\d+)\s+([a-zA-z0-9-\.]+)\s+\((\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\)\s+(\d+\.?\d+\s?ms)\s+(\d+\.?\d+\s?ms)\s+(\d+\.?\d+\s?ms)\s?$`)
)

func ParseOutput(output string) ([]*Route, error) {
	lines := strings.Split(output, "\n")
	routes := make([]*Route, 0, len(lines))

	for _, line := range lines {
		route, err := ParseLine(line)
		if err != nil {
			return nil, err
		}

		routes = append(routes, route)
	}

	return routes, nil
}

// ParseLine given a line from the traceroute output, parse the route info
func ParseLine(line string) (*Route, error) {
	matches := lineRegexp.FindAllStringSubmatch(line, -1)

	hop, err := strconv.Atoi(matches[0][1])
	if err != nil {
		return nil, fmt.Errorf("error parsing hop `%s`: %w", matches[0][1], err)
	}

	ip, err := netaddr.ParseIP(matches[0][3])
	if err != nil {
		return nil, fmt.Errorf("error parsing IP address `%s`: %w", matches[0][3], err)
	}

	ms1 := strings.Split(matches[0][4], " ")[0]
	rtt1, err := parseMS(ms1)
	if err != nil {
		return nil, fmt.Errorf("error parsing first RTT `%s`: %w", ms1, err)
	}

	ms2 := strings.Split(matches[0][5], " ")[0]
	rtt2, err := parseMS(ms2)
	if err != nil {
		return nil, err
	}

	println(matches[0][6])
	ms3 := strings.Split(matches[0][6], " ")[0]
	rtt3, err := parseMS(ms3)
	if err != nil {
		return nil, err
	}

	return &Route{
		Hop:      hop,
		Hostname: matches[0][2],
		IP:       ip,
		RTT: []time.Duration{
			rtt1,
			rtt2,
			rtt3,
		},
	}, nil
}

func parseMS(s string) (time.Duration, error) {
	// parts := strings.Split(s, ".")
	// TODO: don't use floats

	// TODO: check for length of parts

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return time.Duration(f * float64(time.Millisecond)), nil
}
