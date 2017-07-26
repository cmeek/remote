package remote

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// IP contains simplified IP information that includes first, second, and last IP in the list of forwarded IPs
type IP struct {
	IPCount int8
	IP1     string
	IP2     string
	IPLast  string
}

// GetIP returns the remote IP for the given http.Request. It checks optional CF-Connecting-IP and X-FORWARDED-FOR headers.
func GetIP(r *http.Request) (*IP, error) {
	info := &IP{}

	cloudflareIP := r.Header.Get("CF-Connecting-IP")
	if len(cloudflareIP) > 0 {
		if err := parseIPs(cloudflareIP, info); err != nil {
			return nil, errors.Wrapf(err, "GetIP failed to parseIPs %q", cloudflareIP)
		}
		return info, nil
	}

	xff := r.Header.Get("X-FORWARDED-FOR")
	if len(xff) > 0 {
		if err := parseIPs(xff, info); err != nil {
			return nil, errors.Wrapf(err, "GetIP failed to parseIPs %q", xff)
		}
		return info, nil
	}

	info.IPCount = 1
	info.IP1, _, _ = net.SplitHostPort(r.RemoteAddr)
	info.IPLast = info.IP1

	return info, nil
}

func parseIPs(s string, info *IP) error {
	ips := strings.Split(s, ",")

	for _, ip := range ips {
		if net.ParseIP(strings.TrimSpace(ip)) == nil {
			return fmt.Errorf("parseIPs failed to net.ParseIP %q", ip)
		}
	}

	info.IPCount = int8(len(ips))
	info.IP1 = strings.TrimSpace(ips[0])
	if len(ips) > 1 {
		info.IP2 = strings.TrimSpace(ips[1])
	}

	info.IPLast = strings.TrimSpace(ips[len(ips)-1])

	return nil
}
