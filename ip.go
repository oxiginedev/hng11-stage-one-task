package main

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// Get the IP address for a given request
func GetIpAddress(r *http.Request) (string, error) {
	ips := r.Header.Get("X-Forwarded-For")
	splitIps := strings.Split(ips, ",")

	if len(splitIps) > 0 {
		netIp := net.ParseIP(splitIps[len(splitIps)-1])
		if netIp != nil {
			return netIp.String(), nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
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
