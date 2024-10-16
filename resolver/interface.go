package resolver

import "net"

func InterfaceToIP(iface *net.Interface) ([]*net.IPNet, error) {
	ips := make([]*net.IPNet, 0)

	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return nil, err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					ips = append(ips, addr)
					continue
				}
				if ip6 := ipnet.IP.To16(); ip6 != nil {
					addr = &net.IPNet{
						IP:   ip6,
						Mask: ipnet.Mask[len(ipnet.Mask)-16:],
					}
					ips = append(ips, addr)
				}
			}
		}
	}
	return ips, nil
}
