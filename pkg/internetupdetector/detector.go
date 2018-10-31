package internetupdetector

import (
	"context"
	"net"
	"time"
)

// rationale for bypassing LAN DNS servers: to bypass caching

var (
	primaryDetector   = newDetector("208.67.222.222:53") // documented at https://use.opendns.com/
	secondaryDetector = newDetector("208.67.220.220:53")
)

var defaultDialer = net.Dialer{
	Timeout: 2 * time.Second,
}

type detector struct {
	resolver *net.Resolver
}

func (d *detector) detect(ctx context.Context) bool {
	// our IP address would be in the return value
	_, err := d.resolver.LookupHost(ctx, "myip.opendns.com")

	return err == nil
}

func newDetector(overrideAddr string) *detector {
	return &detector{
		resolver: &net.Resolver{
			Dial: func(ctx context.Context, network string, address string) (net.Conn, error) {
				return defaultDialer.DialContext(ctx, network, overrideAddr)
			},
		},
	}
}

func IsInternetUp(ctx context.Context) bool {
	if ok := primaryDetector.detect(ctx); ok {
		return true
	}

	// primary detector failed

	return secondaryDetector.detect(ctx)
}
