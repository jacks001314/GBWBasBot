package pscan

const (
	MAXPORT = 65535
	MINPORT = 1
)

type PScanTargets struct {
	ipTargets *pscanIPTargets

	portTargets *pscanPortTargets
}

type PscanTarget struct {
	IP   string
	Port uint32
}

type pscanIPTargets struct {
	ips []string

	//the iterator index
	curIndex uint32

	num uint32
}

type pscanPortTargets struct {

	//the target start port to scan
	startPort uint32

	//the target end port to scan
	endPort uint32

	//the target ports to scan
	ports []uint32

	//the iterator index
	curIndex uint32
}

func NewPscanTargets(ips []string, startPort, endPort uint32, ports []uint32) *PScanTargets {

	if endPort == 0 {

		endPort = MAXPORT
	}

	if startPort == 0 {
		startPort = MINPORT
	}

	return &PScanTargets{
		ipTargets: &pscanIPTargets{
			ips:      ips,
			curIndex: 0,
			num:      uint32(len(ips)),
		},
		portTargets: &pscanPortTargets{
			startPort: startPort,
			endPort:   endPort,
			ports:     ports,
			curIndex:  0,
		},
	}
}

func (p *PScanTargets) HasNext() bool {

	return p.ipTargets.curIndex < p.ipTargets.num
}

func (p *PScanTargets) initForNext() {

	plen := uint32(len(p.portTargets.ports))

	if plen > 0 {

		if p.portTargets.curIndex >= plen {

			p.portTargets.curIndex = 0
			p.ipTargets.curIndex++
		}
	} else {

		port := p.portTargets.startPort + p.portTargets.curIndex
		if port > p.portTargets.endPort {

			p.portTargets.curIndex = 0
			p.ipTargets.curIndex++
		}
	}
}

func (p *PScanTargets) Next() *PscanTarget {

	ip := p.ipTargets.ips[p.ipTargets.curIndex]

	var port uint32 = 0

	plen := len(p.portTargets.ports)

	if plen > 0 {

		port = p.portTargets.ports[p.portTargets.curIndex]
		p.portTargets.curIndex++

	} else {

		port = p.portTargets.startPort + p.portTargets.curIndex
		p.portTargets.curIndex++

	}

	p.initForNext()

	return &PscanTarget{
		IP:   ip,
		Port: port,
	}
}
