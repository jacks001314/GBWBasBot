package pscan

import (
	"log"
	"time"

	"common/util/netutils"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func (ps *PScan) Receive() {

	for {

		if ps.sendDone {

			if time.Since(ps.sendDoneTime) > time.Duration(ps.cfg.WaitRecvTime)*time.Second {

				ps.wg.Done()

				log.Println("Receive ok!")

				break
			}
		}

		pdata, _, err := ps.handler.ReadPacketData()
		if err != nil {
			continue
		}

		packet := gopacket.NewPacket(pdata, layers.LayerTypeEthernet, gopacket.NoCopy)
		iplayer := packet.Layer(layers.LayerTypeIPv4)
		if iplayer == nil {
			continue
		}

		ipl, ok := iplayer.(*layers.IPv4)
		if !ok {
			continue
		}

		srcIP := netutils.IPStrToInt(ipl.DstIP.String())
		dstIP := netutils.IPStrToInt(ipl.SrcIP.String())

		validator := GenValidator(srcIP, dstIP)

		tcplayer := packet.Layer(layers.LayerTypeTCP)
		if tcplayer == nil {
			continue
		}

		tcp, ok := tcplayer.(*layers.TCP)
		if !ok {
			continue
		}
		seq := validator.GetSequence()
		ack := tcp.Ack

		if validator.IsValidDstPort(uint32(tcp.DstPort)) && (ack == seq+1) {

			if tcp.SYN && tcp.ACK {
				//find a live port
				result := &PResult{
					IP:   ipl.SrcIP.String(),
					Port: uint32(tcp.SrcPort),
				}

				ps.PushScanResult(result)
			}
		}

	} //for
}
