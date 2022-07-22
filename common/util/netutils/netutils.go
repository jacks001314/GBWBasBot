package netutils

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
)

type URLPathCrypt struct {
	TaskId string

	NodeId string

	Fname string

	AttackType string

	AttackIP string

	TargetIP string

	TargetPort int

	DownloadTool string
}

type DNSDomainCrypt struct {
	AttackType string
	AttackIP   string
	TargetIP   string
	TargetPort int
}

func IPv4StrLittle(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip), byte(ip>>8), byte(ip>>16), byte(ip>>24))
}

func IPv4StrBig(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

func IPStrToInt(ip string) uint32 {

	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())

	return uint32(ret.Uint64())
}

func DNSDomainCryptToString(d *DNSDomainCrypt) string {

	s := fmt.Sprintf("%s_%d_%d_%d",
		d.AttackType, IPStrToInt(d.AttackIP), IPStrToInt(d.TargetIP), d.TargetPort)

	return s
}

func DeCryptToDNSDomain(content string) (*DNSDomainCrypt, error) {

	args := strings.Split(content, "_")

	if len(args) != 4 {

		return nil, fmt.Errorf("Invalid dns domain format:%s", content)

	}

	attackIPI, err := strconv.ParseUint(args[1], 10, 32)
	targetIPI, err := strconv.ParseUint(args[2], 10, 32)
	targetPortI, err := strconv.ParseUint(args[3], 10, 32)

	if err != nil {

		return nil, fmt.Errorf("Invalid dns domain format:%s", content)
	}

	return &DNSDomainCrypt{
		AttackType: args[0],
		AttackIP:   IPv4StrBig(uint32(attackIPI)),
		TargetIP:   IPv4StrBig(uint32(targetIPI)),
		TargetPort: int(targetPortI),
	}, nil
}

func URLPathCryptToString(u *URLPathCrypt) string {

	s := fmt.Sprintf("%s,%s,%s,%s,%d,%d,%d,%s", u.TaskId, u.NodeId,
		u.Fname, u.AttackType, IPStrToInt(u.AttackIP), IPStrToInt(u.TargetIP), u.TargetPort, u.DownloadTool)

	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DeCryptToURLPath(content string) (*URLPathCrypt, error) {

	d, err := base64.StdEncoding.DecodeString(content)

	if err != nil {
		return nil, err
	}

	s := string(d)

	if s == "" || strings.Index(s, ",") <= 0 {

		return nil, fmt.Errorf("Invalid url path format:%s", s)

	}

	args := strings.Split(s, ",")

	if len(args) != 8 {
		return nil, fmt.Errorf("Invalid url path format:%s", s)
	}

	attackIPI, err := strconv.ParseUint(args[4], 10, 32)
	targetIPI, err := strconv.ParseUint(args[5], 10, 32)
	targetPortI, err := strconv.ParseUint(args[6], 10, 32)

	if err != nil {

		return nil, fmt.Errorf("Invalid url path format:%s", s)
	}

	return &URLPathCrypt{
		TaskId:       args[0],
		NodeId:       args[1],
		Fname:        args[2],
		AttackType:   args[3],
		AttackIP:     IPv4StrBig(uint32(attackIPI)),
		TargetIP:     IPv4StrBig(uint32(targetIPI)),
		TargetPort:   int(targetPortI),
		DownloadTool: args[7],
	}, nil
}
