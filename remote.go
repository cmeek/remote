package main

import (
	"net"
	"net/http"
	"strings"
)

// IP contains simplified IP information that includes first, second, and last IP in the list of forwarded IPs
type IP struct {
	IPCount int8
	IP1     string
	IP2     string
	IPLast  string
}

// GetIP returns the remote IP for the given http.Request. It checks optional CF-Connecting-IP and X-FORWARDED-FOR headers.
func getIP(r *http.Request) *IP {
	info := &IP{}

	cloudflareIP := r.Header.Get("CF-Connecting-IP")
	if len(cloudflareIP) > 0 {
		parseIPs(cloudflareIP, info)
		return info
	}

	xff := r.Header.Get("X-FORWARDED-FOR")
	if len(xff) > 0 {
		parseIPs(xff, info)
		return info
	}

	info.IPCount = 1
	info.IP1, _, _ = net.SplitHostPort(r.RemoteAddr)
	info.IPLast = info.IP1

	return info
}

func parseIPs(s string, info *IP) {
	ips := strings.Split(s, ",")

	info.IPCount = int8(len(ips))
	info.IP1 = strings.TrimSpace(ips[0])
	if len(ips) > 1 {
		info.IP2 = strings.TrimSpace(ips[1])
	}

	info.IPLast = strings.TrimSpace(ips[len(ips)-1])
}
