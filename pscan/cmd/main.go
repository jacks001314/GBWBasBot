package main

import (
	"flag"
	"fmt"
	"log"
	pscan "pscan/src"
	"strconv"
	"strings"
	"time"

	"common/ipgen/ipv4"
	"common/util/netutils"
)

func main() {

	ifaceIP := flag.String("ifaceIP", "", "the interface ip to send/receive packets")

	ifaceName := flag.String("ifaceName", "", "the interface name to send/receive packets")

	gwIP := flag.String("gwIP", "", "the interface gateway ip  to send/receive packets")

	gwMac := flag.String("gwMac", "", "the interface gateway mac address")

	sportFirst := flag.Uint("sportFirst", 2222, "source port range fist")

	sportLast := flag.Uint("sportLast", 6666, "source port range last")

	sendRate := flag.Uint64("sendRate", 128, "send packets rate,packets/second")

	sendBandWidth := flag.Uint64("sendBandWidth", 0, "send packets bandwidth")

	waitReceiveTime := flag.Uint64("waitReceiveTime", 60, "when send packets over,the time out to wait ack packets")

	scanPortStart := flag.Uint("scanPortStart", 0, "the scan port range start")

	scanPortEnd := flag.Uint("scanPortEnd", 0, "the scan port range end")

	scanPorts := flag.String("scanPorts", "", "the scan ports:22,33,44...")

	scanIPWhiteList := flag.String("scanIPWhiteList", "", "scan ip whitelist")

	scanIPBlackList := flag.String("scanIPBlackList", "", "scan ip blacklist")

	flag.Parse()

	cfg := &pscan.PScanConfig{

		SrcIP:           *ifaceIP,
		Ifname:          *ifaceName,
		GatewayIP:       *gwIP,
		GatewayMac:      *gwMac,
		SourcePortFirst: uint32(*sportFirst),
		SourcePortLast:  uint32(*sportLast),
		SendRate:        *sendRate,
		SendBandWidth:   *sendBandWidth,
		WaitRecvTime:    *waitReceiveTime,
	}

	wlist := strings.Split(*scanIPWhiteList, ",")
	blist := strings.Split(*scanIPBlackList, ",")

	ipgen, err := ipv4.NewIPV4Generator("", "", wlist, blist, true)

	if err != nil {

		log.Fatal("create ipv4 generate failed:%v", err)

	}

	ports := make([]uint32, 0)

	if *scanPorts != "" {

		parr := strings.Split(*scanPorts, ",")

		for _, p := range parr {

			if v, err := strconv.ParseUint(p, 10, 32); err == nil {

				ports = append(ports, uint32(v))
			}
		}
	}

	targets := make(chan *pscan.PScanTargets)
	done := make(chan bool)
	results := make(chan *pscan.PResult)

	scan, err := pscan.NewPScan(cfg, targets, results, done)
	if err != nil {

		log.Fatal("create a scanner failed:%v", err)
	}

	go func() {
		i := 0
		ips := make([]string, 0)

		for ip := ipgen.GetCurIP(); ip != 0; ip = ipgen.GetNextIP() {

			ipStr := netutils.IPv4StrBig(ip)

			if i >= 128 {
				i = 0
				t := pscan.NewPscanTargets(ips, uint32(*scanPortStart), uint32(*scanPortEnd), ports)

				scan.PushScanTargets(t)
				ips = make([]string, 0)
			}

			ips = append(ips, ipStr)

			i++
		}

		if len(ips) > 0 {

			t := pscan.NewPscanTargets(ips, uint32(*scanPortStart), uint32(*scanPortEnd), ports)
			scan.PushScanTargets(t)
		}

		//wait 1miniute
		time.Sleep(1 * time.Minute)

		scan.SetDone()
	}()

	go func() {

		results := scan.SubScanResults()

		for {

			select {

			case r := <-results:
				fmt.Printf("{ip:%s,port:%d}\n", r.IP, r.Port)
			}

		}
	}()

	scan.Start()
}
