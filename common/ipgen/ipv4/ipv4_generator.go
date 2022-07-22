package ipv4

import (
	"fmt"
	"math/big"
)

type IPV4Generator struct {

	wb *WBList
	cycle *Cycle

	first uint64
	last  uint64
	factor uint64
	modulus uint64
	current uint64

	maxIndex uint64
}


func (ipg *IPV4Generator)initIPGen() {

	var numElts = ipg.cycle.Order
	var expBegin = uint64(ipg.cycle.Offset) % numElts
	var expEnd = uint64(ipg.cycle.Offset) % numElts

	// Multiprecision variants of everything above
	genM := new(big.Int).SetUint64(ipg.cycle.Generator)
	expBeginM := new(big.Int).SetUint64(expBegin)
	expEndM := new(big.Int).SetUint64(expEnd)
	primeM := new(big.Int).SetUint64(ipg.cycle.Group.prime)

	startM := new(big.Int)
	stopM := new(big.Int)

	startM = startM.Exp(genM,expBeginM,primeM)
	stopM = stopM.Exp(genM,expEndM,primeM)

	ipg.first = startM.Uint64()
	ipg.last = stopM.Uint64()
	ipg.factor = ipg.cycle.Generator
	ipg.modulus = ipg.cycle.Group.prime

	ipg.current = ipg.first

	ipg.roll2Valid()

}

func NewIPV4Generator(whiteListFName string,blackListFName string,whiteListEntries []string,
	blakListEntries []string,ignoreInvalidHosts bool) (*IPV4Generator,error) {

	var ipg IPV4Generator
	var err error
	var max32Int uint64 = 1<<32

	ipg.wb,err = NewWblist(whiteListFName,blackListFName,whiteListEntries,blakListEntries,ignoreInvalidHosts)

	if err!=nil {

		return nil,err
	}

	numAddr := ipg.wb.WBListCountAllowed()

	group := GetGroup(numAddr)

	if group == nil {
		return nil,fmt.Errorf("Cannot get valid group for numAddr:%d",numAddr)
	}

	if numAddr> max32Int {
		ipg.maxIndex= 0xFFFFFFFF
	} else {
		ipg.maxIndex = numAddr
	}

	ipg.cycle = group.MakeCycle()

	ipg.initIPGen()

	return &ipg,nil
}


func (ipg *IPV4Generator) GetCurIP() uint32 {

	return uint32(ipg.wb.LookupIndex(ipg.current - 1))
}

func (ipg *IPV4Generator) getNextElem() uint32 {

	var max32Int uint64 = 1 << 32

	for {

		ipg.current = ipg.current*ipg.factor
		ipg.current = ipg.current%ipg.modulus

		if ipg.current <max32Int {

			break
		}
	}

	return uint32(ipg.current)
}

func (ipg *IPV4Generator) GetNextIP() uint32 {

	var candidate uint64

	if ipg.current == 0 {
		return 0
	}

	for {
		candidate = uint64(ipg.getNextElem())
		if candidate == ipg.last {

			ipg.current = 0
			return 0
		}

		if candidate -1 < ipg.maxIndex {
			return uint32(ipg.wb.LookupIndex(candidate - 1))
		}
	}
}

func (ipg *IPV4Generator) roll2Valid() uint32 {

	if ipg.current - 1 <ipg.maxIndex {
		return uint32(ipg.current)
	}

	return ipg.GetNextIP()
}


