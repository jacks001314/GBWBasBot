package ipv4

import (
	"fmt"

	stengo "common/scripts/tengo"
	"common/util/netutils"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

const (
	TengoMethodCurIP  = "curIP"
	TengoMethodNextIP = "nextIP"
)

type IPV4Tengo struct {
	stengo.TengoObj

	ipgen *IPV4Generator
}

func newTengoIPGenFromFile(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	wlistFname, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "wlistFname",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	blistFname, ok := objects.ToString(args[1])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "blistFname",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	ipgen, err := NewIPV4Generator(wlistFname, blistFname, []string{}, []string{}, true)

	if err != nil {
		return nil, err
	}

	return &IPV4Tengo{
		TengoObj: stengo.TengoObj{Name: "IPV4_Gen_Tengo"},
		ipgen:    ipgen,
	}, nil
}

func newTengoIPGenFromArray(args ...objects.Object) (objects.Object, error) {

	if len(args) != 2 {

		return nil, tengo.ErrWrongNumArguments
	}

	wlistArr, ok := objects.ToInterface(args[0]).([]interface{})

	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "wlistArray",
			Expected: "[]string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	blistArr, ok := objects.ToInterface(args[1]).([]interface{})

	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "blistArray",
			Expected: "[]string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	wlists := make([]string, 0)
	blists := make([]string, 0)

	for _, w := range wlistArr {

		wlists = append(wlists, w.(string))
	}

	for _, b := range blistArr {

		blists = append(blists, b.(string))
	}

	ipgen, err := NewIPV4Generator("", "", wlists, blists, true)

	if err != nil {
		return nil, err
	}

	return &IPV4Tengo{
		TengoObj: stengo.TengoObj{Name: "IPV4_Gen_Tengo"},
		ipgen:    ipgen,
	}, nil

}

func (ipt *IPV4Tengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	switch key {

	case TengoMethodCurIP:

		return &IPV4TengoMethod{
			TengoObj: stengo.TengoObj{Name: TengoMethodCurIP},
			ipt:      ipt,
		}, nil

	case TengoMethodNextIP:

		return &IPV4TengoMethod{
			TengoObj: stengo.TengoObj{Name: TengoMethodNextIP},
			ipt:      ipt,
		}, nil

	}

	return nil, fmt.Errorf("Unknown ipgen method:%s", key)
}

type IPV4TengoMethod struct {
	stengo.TengoObj

	ipt *IPV4Tengo
}

func (iptm *IPV4TengoMethod) getCurIP(args ...objects.Object) (objects.Object, error) {

	ip := iptm.ipt.ipgen.GetCurIP()

	if ip == 0 {

		return objects.FromInterface("")
	}

	return objects.FromInterface(netutils.IPv4StrBig(ip))
}

func (iptm *IPV4TengoMethod) getNextIP(args ...objects.Object) (objects.Object, error) {

	ip := iptm.ipt.ipgen.GetNextIP()

	if ip == 0 {

		return objects.FromInterface("")
	}

	return objects.FromInterface(netutils.IPv4StrBig(ip))
}

func (iptm *IPV4TengoMethod) Call(args ...objects.Object) (objects.Object, error) {

	switch iptm.Name {

	case TengoMethodCurIP:
		return iptm.getCurIP(args...)

	case TengoMethodNextIP:
		return iptm.getNextIP(args...)

	default:
		return nil, fmt.Errorf("unknown ipv4 generator method:%s", iptm.Name)
	}
}

var moduleMap objects.Object = &objects.ImmutableMap{

	Value: map[string]objects.Object{

		"newIPGenFromFile": &objects.UserFunction{
			Name:  "new_ipgen_fromfile",
			Value: newTengoIPGenFromFile,
		},

		"newIPGenFromArray": &objects.UserFunction{
			Name:  "new_ipgen_from_array",
			Value: newTengoIPGenFromArray,
		},
	},
}

func (IPV4Tengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case "ip4gen":
		return moduleMap, nil
	default:
		return nil, fmt.Errorf("undefined module %q", moduleName)
	}
}
