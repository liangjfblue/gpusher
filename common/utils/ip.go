/**
 *
 * @author liangjf
 * @create on 2020/6/1
 * @version 1.0
 */
package utils

import (
	"errors"
	"net"
)

// ExternalIP 获取公网IP , 假设本机有独立公网IP
func ExternalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1", err
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && isExternalIP(ip.IP) {
			return ip.IP.To4().String(), nil
		}
	}

	return "127.0.0.1", errors.New("can not find external IP")
}

//isExternalIP 判断是否是公网ip
func isExternalIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}
