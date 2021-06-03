package trparse

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"inet.af/netaddr"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		line    string
		want    *Route
		wantErr bool
	}{
		{
			" 1  gigabitethernet3-3.exi1.melbourne.telstra.net (203.50.77.49)  0.235 ms  0.285 ms  0.245 ms",
			&Route{
				Hop:      1,
				Hostname: "gigabitethernet3-3.exi1.melbourne.telstra.net",
				IP:       netaddr.IPFrom4([4]byte{203, 50, 77, 49}),
				RTT: []time.Duration{
					235 * time.Microsecond,
					285 * time.Microsecond,
					245 * time.Microsecond,
				},
			},
			false,
		},
		{
			" 2  bundle-ether3-100.exi-core10.melbourne.telstra.net (203.50.80.1)  2.368 ms  1.800 ms  2.119 ms",
			&Route{
				Hop:      2,
				Hostname: "bundle-ether3-100.exi-core10.melbourne.telstra.net",
				IP:       netaddr.IPFrom4([4]byte{203, 50, 80, 1}),
				RTT: []time.Duration{
					2368 * time.Microsecond,
					1800 * time.Microsecond,
					2119 * time.Microsecond,
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got, err := ParseLine(tt.line)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_parseMS(t *testing.T) {
	tests := []struct {
		in      string
		want    time.Duration
		wantErr bool
	}{
		{
			"0.235",
			235 * time.Microsecond,
			false,
		},
		{
			"0.255",
			255 * time.Microsecond,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			got, err := parseMS(tt.in)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

var output1 = ` 1  gigabitethernet3-3.exi1.melbourne.telstra.net (203.50.77.49)  0.235 ms  0.285 ms  0.245 ms
 2  bundle-ether3-100.exi-core10.melbourne.telstra.net (203.50.80.1)  2.368 ms  1.800 ms  2.119 ms
 3  bundle-ether12.chw-core10.sydney.telstra.net (203.50.11.124)  14.988 ms  14.169 ms  14.738 ms`

func TestParseOutput(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Route
		wantErr bool
	}{
		{
			"test case 1",
			args{output1},
			[]*Route{
				{
					Hop:      1,
					Hostname: "gigabitethernet3-3.exi1.melbourne.telstra.net",
					IP:       netaddr.IPFrom4([4]byte{203, 50, 77, 49}),
					RTT: []time.Duration{
						235 * time.Microsecond,
						285 * time.Microsecond,
						245 * time.Microsecond,
					},
				},
				{
					Hop:      2,
					Hostname: "bundle-ether3-100.exi-core10.melbourne.telstra.net",
					IP:       netaddr.IPFrom4([4]byte{203, 50, 80, 1}),
					RTT: []time.Duration{
						2368 * time.Microsecond,
						1800 * time.Microsecond,
						2119 * time.Microsecond,
					},
				},
				{
					Hop:      3,
					Hostname: "bundle-ether12.chw-core10.sydney.telstra.net",
					IP:       netaddr.IPFrom4([4]byte{203, 50, 11, 124}),
					RTT: []time.Duration{
						14988 * time.Microsecond,
						14169 * time.Microsecond,
						14738 * time.Microsecond,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOutput(tt.args.output)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
