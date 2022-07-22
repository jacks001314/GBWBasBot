package pscan

import (
	"net"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type PScan struct {

	//config for scanner
	cfg *PScanConfig

	handler *pcap.Handle

	//ip and ports will been scanned
	targets chan *PScanTargets

	done chan bool

	//the net interface to send and receive packets for scanning
	inf *PInterface

	presults chan *PResult

	buf gopacket.SerializeBuffer

	opts gopacket.SerializeOptions

	wg *sync.WaitGroup

	sendDoneTime time.Time

	sendDone bool
}

func NewPScan(cfg *PScanConfig, targets chan *PScanTargets, presults chan *PResult, done chan bool) (*PScan, error) {

	inf, err := GetInterfaceWithGWMac(cfg.SrcIP, cfg.Ifname, cfg.GatewayIP, cfg.GatewayMac)

	if err != nil {

		return nil, err
	}

	return &PScan{
		cfg:      cfg,
		handler:  nil,
		targets:  targets,
		done:     done,
		inf:      inf,
		presults: presults,
		buf:      gopacket.NewSerializeBuffer(),
		opts: gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		},
	}, nil

}

func (ps *PScan) Start() error {

	InitPValidate(ps.cfg.SourcePortFirst, ps.cfg.SourcePortLast, 1)

	// Open the handle for reading/writing.
	// Note we could very easily add some BPF filtering here to greatly
	// decrease the number of packets we have to look at when getting back
	// scan results.
	handle, err := pcap.OpenLive(ps.inf.Iface.Name, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}

	ps.handler = handle

	ps.wg = &sync.WaitGroup{}

	ps.wg.Add(2)

	go ps.Send()

	go ps.Receive()

	ps.wg.Wait()

	return nil
}

func (ps *PScan) GetInterfaceMac() net.HardwareAddr {
	return ps.inf.Iface.HardwareAddr
}

func (ps *PScan) GetGatewayMac() net.HardwareAddr {
	return ps.inf.GWMacAddr
}

func (ps *PScan) GetInterfaceIP() net.IP {

	return ps.inf.IP
}

func (ps *PScan) GetInterfaceIPStr() string {

	return ps.inf.IP.String()
}

func (ps *PScan) PushScanResult(result *PResult) {
	ps.presults <- result
}

func (ps *PScan) PushScanTargets(t *PScanTargets) {

	ps.targets <- t
}

func (ps *PScan) SetDone() {

	ps.done <- true
}

func (ps *PScan) SubScanResults() chan *PResult {

	return ps.presults
}
