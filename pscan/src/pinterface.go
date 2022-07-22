package pscan

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"common/routing"
)

var DetectIP = net.ParseIP("8.8.8.8")

type PInterface struct {
	Iface *net.Interface
	IP    net.IP
	GW    net.IP
	GWMac string

	GWMacAddr net.HardwareAddr
}

func GetInterfaceWith(ip string, ifname, gw string) (*PInterface, error) {

	iface, ifaceIP, err := getIface(ip, ifname)

	if err != nil {

		return nil, err
	}

	mac := ""
	if gw != "" {
		mac = GetMac(gw)
	}

	gwmacAddr, err := net.ParseMAC(mac)

	if err != nil {

		return nil, err
	}

	return &PInterface{
		Iface:     iface,
		IP:        ifaceIP,
		GW:        net.ParseIP(gw),
		GWMac:     mac,
		GWMacAddr: gwmacAddr,
	}, nil
}

func GetInterfaceWithGWMac(ip, ifname, gw, gwmac string) (*PInterface, error) {

	iface, ifaceIP, err := getIface(ip, ifname)

	if err != nil {

		return nil, err
	}

	gwmacAddr, err := net.ParseMAC(gwmac)

	if err != nil {

		return nil, err
	}

	return &PInterface{
		Iface:     iface,
		IP:        ifaceIP,
		GW:        net.ParseIP(gw),
		GWMac:     gwmac,
		GWMacAddr: gwmacAddr,
	}, nil

}

func getIpFromAddr(addr net.Addr) (net.IP, net.IPMask) {

	var ip net.IP
	var mask net.IPMask

	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
		mask = v.Mask

	case *net.IPAddr:

		ip = v.IP
		mask = ip.DefaultMask()
	}

	if ip == nil || ip.IsLoopback() {
		return nil, nil
	}

	return ip, mask
}

func getIface(ipstr string, name string) (*net.Interface, net.IP, error) {

	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, nil, err
	}

	for _, iface := range ifaces {

		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {

			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, nil, err
		}

		for _, addr := range addrs {

			ip, mask := getIpFromAddr(addr)

			if ip == nil {
				continue
			}
			prefixlen, _ := mask.Size()

			if prefixlen > 32 {
				continue
			}

			if ipstr != "" && strings.EqualFold(ip.String(), ipstr) {

				return &iface, ip, nil
			}

			if name != "" && strings.EqualFold(name, iface.Name) {
				return &iface, ip, nil
			}
		}
	}

	return nil, nil, fmt.Errorf("not found interface for ip:%s or name:%s", ipstr, name)
}

//only for linux
func GetInterface() (*PInterface, error) {

	if runtime.GOOS != "linux" {

		return nil, fmt.Errorf("only for linux")
	}

	r, err := routing.New()

	if err != nil {
		log.Printf("new err:%v\n", err)
		return nil, err
	}
	//iface *net.Interface, gateway, preferredSrc net.IP, err error
	iface, gw, src, err := r.Route(DetectIP)
	if err != nil {
		log.Printf("route err:%v\n", err)
		return nil, err
	}

	mac := ""
	if gw != nil {
		mac = GetMac(gw.String())
	}

	gwmacAddr, err := net.ParseMAC(mac)

	if err != nil {

		return nil, err
	}

	return &PInterface{
		Iface:     iface,
		IP:        src,
		GW:        gw,
		GWMac:     mac,
		GWMacAddr: gwmacAddr,
	}, nil
}

func GetMac(ip string) string {

	switch runtime.GOOS {

	case "windows":
		return getMacForWindows(ip)
	case "linux":
		return getMacForLinux(ip)

	default:
		return getMacForUnix(ip)

	}

}

func getMacForWindows(ip string) string {

	data, err := exec.Command("arp", "-a").Output()
	if err != nil {
		return ""
	}

	skipNext := false
	for _, line := range strings.Split(string(data), "\n") {
		// skip empty lines
		if len(line) <= 0 {
			continue
		}
		// skip Interface: lines
		if line[0] != ' ' {
			skipNext = true
			continue
		}
		// skip column headers
		if skipNext {
			skipNext = false
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		ipstr := strings.TrimSpace(fields[0])

		if strings.EqualFold(ip, ipstr) {

			return strings.TrimSpace(fields[1])
		}

	}

	return ""
}

func getMacForLinux(ip string) string {

	f, err := os.Open("/proc/net/arp")

	if err != nil {
		return ""
	}

	defer f.Close()

	s := bufio.NewScanner(f)
	s.Scan() // skip the field descriptions

	for s.Scan() {
		line := s.Text()
		fields := strings.Fields(line)

		if len(fields) < 4 {
			continue
		}

		ipstr := strings.TrimSpace(fields[0])

		if strings.EqualFold(ip, ipstr) {

			return strings.TrimSpace(fields[3])
		}

	}

	return ""
}

func getMacForUnix(ip string) string {

	data, err := exec.Command("arp", "-an").Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		// strip brackets around IP
		ip := strings.Replace(fields[1], "(", "", -1)
		ip = strings.Replace(ip, ")", "", -1)

		ipstr := strings.TrimSpace(ip)

		if strings.EqualFold(ip, ipstr) {

			return strings.TrimSpace(fields[3])
		}

	}

	return ""
}
