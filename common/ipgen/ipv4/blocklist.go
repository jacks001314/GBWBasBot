package ipv4

import (
	"common/util/fileutils"
	"common/util/netutils"
	"fmt"

	"strconv"
	"strings"
)

const (
	ADDR_DISALLOWED = 0
	ADDR_ALLOWED    = 1
)

type bl_cidr_node struct {
	address   uint32
	prefixLen int
}

type WBList struct {
	con       *Constraint
	blacklist []*bl_cidr_node
	whitelist []*bl_cidr_node
}

func (wb *WBList) addNode(address uint32, prefixLen int, value uint32) {

	nd := &bl_cidr_node{
		address:   address,
		prefixLen: prefixLen,
	}

	if value == ADDR_ALLOWED {

		wb.whitelist = append(wb.whitelist, nd)
	} else {

		wb.blacklist = append(wb.blacklist, nd)
	}

}

func (wb *WBList) LookupIndex(index uint64) uint64 {

	return wb.con.LookupIndex(index, ADDR_ALLOWED)

}

// check whether a single IP address is allowed to be scanned.

func (wb *WBList) IsAllowed(addr uint32) bool {

	return wb.con.LookupIP(addr) == ADDR_ALLOWED
}

func (wb *WBList) addConstraint(addr uint32, prefixLen int, value uint32) {

	wb.con.Set(addr, prefixLen, value)

	wb.addNode(addr, prefixLen, value)

}

func (wb *WBList) InitFromString(ips string, value uint32) error {

	var prefixLen int = 32
	var ip string = ips
	if strings.Contains(ips, "/") {

		arr := strings.Split(ips, "/")
		if len(arr) != 2 {
			return fmt.Errorf("Invalid ip format:%s", ips)
		}

		ip = arr[0]

		v, err := strconv.ParseInt(arr[1], 10, 32)
		if err != nil || v > 32 || v < 0 {

			return fmt.Errorf("Invalid ip format:%s,err:%v", ips, err)

		}

		prefixLen = int(v)
	}

	wb.addConstraint(netutils.IPStrToInt(ip), prefixLen, value)

	return nil
}

func (wb *WBList) InitFromFile(fname string, value uint32, ignoreInvalidHosts bool) error {

	lines, err := fileutils.ReadAllLines(fname)

	if err != nil {

		return err
	}

	for _, line := range lines {

		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "#") {

			err = wb.InitFromString(line, value)

			if err != nil && !ignoreInvalidHosts {

				return err
			}
		}
	}

	return nil

}

func (wb *WBList) InitFromArray(arr []string, value uint32, ignoreInvalidHosts bool) error {

	for _, line := range arr {

		line = strings.TrimSpace(line)

		if !strings.HasPrefix(line, "#") {

			err := wb.InitFromString(line, value)

			if err != nil && !ignoreInvalidHosts {

				return err
			}
		}
	}

	return nil
}

func (wb *WBList) WBListCountAllowed() uint64 {

	return wb.con.CountIPS(ADDR_ALLOWED)
}

func (wb *WBList) WBListCountNotAllowed() uint64 {
	return wb.con.CountIPS(ADDR_DISALLOWED)
}

// network order
func (wb *WBList) WBListIpToIndex(ip uint32) uint32 {
	return wb.con.LookupIP(ip)
}

// Initialize address constraints from allowlist and blocklist files.
// Either can be set to NULL to omit.

func NewWblist(whiteListFName string, blackListFName string, whiteListEntries []string,
	blakListEntries []string, ignoreInvalidHosts bool) (*WBList, error) {

	var wblist WBList
	var err error

	wblist.blacklist = make([]*bl_cidr_node, 0)
	wblist.whitelist = make([]*bl_cidr_node, 0)

	if whiteListFName != "" || len(whiteListEntries) > 0 {
		// using a allowlist, so default to allowing nothing
		wblist.con = NewConstraint(ADDR_DISALLOWED)

		if whiteListFName != "" {

			err = wblist.InitFromFile(whiteListFName, ADDR_ALLOWED, ignoreInvalidHosts)
			if err != nil {

				return nil, err
			}
		}

		if len(whiteListEntries) > 0 {

			err = wblist.InitFromArray(whiteListEntries, ADDR_ALLOWED, ignoreInvalidHosts)
			if err != nil {

				return nil, err
			}
		}
	} else {
		// no allowlist, so default to allowing everything
		wblist.con = NewConstraint(ADDR_ALLOWED)
	}

	if blackListFName != "" {

		err = wblist.InitFromFile(blackListFName, ADDR_DISALLOWED, ignoreInvalidHosts)
		if err != nil {

			return nil, err
		}
	}

	if len(blakListEntries) > 0 {

		err = wblist.InitFromArray(blakListEntries, ADDR_DISALLOWED, ignoreInvalidHosts)
		if err != nil {

			return nil, err
		}

	}

	wblist.InitFromString("0.0.0.0", ADDR_DISALLOWED)
	wblist.con.PaintValue(ADDR_ALLOWED)

	allowed := wblist.WBListCountAllowed()

	if allowed == 0 {

		return nil, fmt.Errorf("WBList no addresses are eligible to be scanned in the "+
			"current configuration. This may be because the "+
			"blocklist being used by scan (%s) prevents "+
			"any addresses from receiving probe packets.",
			whiteListFName)

	}

	return &wblist, nil
}
