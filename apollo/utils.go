package apollo

import "net"

func GetLocalIP() string {
	ips, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, ip := range ips {
		if ip4 := toIP4(ip); ip4 != nil {
			return ip4.String()
		}
	}
	return ""
}

func toIP4(addr net.Addr) net.IP {
	if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
		return ipNet.IP.To4()
	}
	return nil
}
