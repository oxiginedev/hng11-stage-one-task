package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// GetIpAddress parses the ip for a given request
func GetIpAddress(r *http.Request) (string, error) {
	rip := r.Header.Get("X-Real-IP")

	netRip := net.ParseIP(rip)
	if netRip != nil {
		return netRip.String(), nil
	}

	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIp := net.ParseIP(splitIps[0])
		if netIp != nil {
			return netIp.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err != nil {
		return "", err
	}

	netIp := net.ParseIP(ip)
	if netIp != nil {
		ip := netIp.String()
		if ip == "::1" {
			return "127.0.0.1", nil
		}
		return ip, nil
	}

	return "", errors.New("could not get ip address from request")
}
