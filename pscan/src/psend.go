package pscan

import (
	"log"
	"net"
	"time"

	"common/util/netutils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (ps *PScan) Send() {

	InitRSendRate(1, ps.cfg.SendRate, ps.cfg.SendBandWidth)

	psr := NewPSendRatePerThread()

	ps.sendDone = false

	for {

		select {

		case <-ps.done:
			goto out

		case targets := <-ps.targets:

			for targets.HasNext() {

				psr.Sleep()
				t := targets.Next()
				ps.sendPacket(t.IP, t.Port)

			}

		}
	}

out:

	log.Println("Send Packets ok!")

	ps.sendDone = true
	ps.sendDoneTime = time.Now()
	ps.wg.Done()
}

//make a tcp packet and send it
func (ps *PScan) sendPacket(dip string, dstPort uint32) error {

	srcIP := ps.GetInterfaceIP()
	dstIP := net.ParseIP(dip)

	validator := GenValidator(netutils.IPStrToInt(ps.GetInterfaceIPStr()),
		netutils.IPStrToInt(dip))

	srcPort := validator.GetDstPort(1)
	seq := validator.GetSequence()

	// Construct all the network layers we need.
	eth := layers.Ethernet{
		SrcMAC:       ps.GetInterfaceMac(), //the mac address local interface to send and receive packets
		DstMAC:       ps.GetGatewayMac(),   //the mac address of gateway
		EthernetType: layers.EthernetTypeIPv4,
	}

	ip4 := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	tcp := layers.TCP{
		SrcPort: layers.TCPPort(srcPort),
		DstPort: layers.TCPPort(dstPort),
		Seq:     seq,
		SYN:     true,
	}

	tcp.SetNetworkLayerForChecksum(&ip4)

	// to bytes
	if err := gopacket.SerializeLayers(ps.buf, ps.opts, &eth, &ip4, &tcp); err != nil {

		return err
	}

	//to send
	return ps.handler.WritePacketData(ps.buf.Bytes())
}
