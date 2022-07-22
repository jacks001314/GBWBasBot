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

		if pdata, _, err := ps.handler.ReadPacketData(); err == nil {

			//read ok and parse packet data
			// Parse the packet.  We'd use DecodingLayerParser here if we
			// wanted to be really fast.
			packet := gopacket.NewPacket(pdata, layers.LayerTypeEthernet, gopacket.NoCopy)

			if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {

				if ipl, ok := ipLayer.(*layers.IPv4); ok {

					srcIP := netutils.IPStrToInt(ipl.DstIP.String())
					dstIP := netutils.IPStrToInt(ipl.SrcIP.String())

					validator := GenValidator(srcIP, dstIP)

					if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {

						if tcp, ok := tcpLayer.(*layers.TCP); ok {

							seq := validator.GetSequence()
							ack := tcp.Ack

							if validator.IsValidDstPort(uint32(tcp.DstPort)) && (ack == seq+1) {

								if tcp.SYN && tcp.ACK {
									//find a live port
									result := &PResult{IP: ipl.SrcIP.String(),
										Port: uint32(tcp.SrcPort)}

									ps.PushScanResult(result)
								}
							}
						}
					}
				}
			}
		}
	} //for
}
