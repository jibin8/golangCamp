/*
 * IP相关的操作
 */
package components

import (
	"bytes"
	"errors"
	"net"
)

func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			if isPrivateIP(ip.String()) {
				return ip.String(), nil
			}
		}
	}
	return "", errors.New("are you connected to the network?")
}

func isPrivateIP(ip string) bool {
	// 根据tcp/ip协议，如下IP地址都是内网地址:
	// 192.168.0.0 - 192.168.255.255
	// 172.16.0.0 - 172.31.255.255
	// 10.0.0.0 - 10.255.255.255

	ret := true
	trial := net.ParseIP(ip)
	if trial.To4() == nil {
		ret = false
	}
	if bytes.Compare(trial, net.ParseIP("192.168.0.0")) >= 0 && bytes.Compare(trial, net.ParseIP("192.168.255.255")) <= 0 {
		ret = true
	} else if bytes.Compare(trial, net.ParseIP("172.16.0.0")) >= 0 && bytes.Compare(trial, net.ParseIP("172.31.255.255")) <= 0 {
		ret = true
	} else if bytes.Compare(trial, net.ParseIP("10.0.0.0")) >= 0 && bytes.Compare(trial, net.ParseIP("10.255.255.255")) <= 0 {
		ret = true
	}
	return ret
}
