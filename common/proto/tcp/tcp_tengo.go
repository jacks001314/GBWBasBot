package tcp

import (
	"errors"
	"fmt"

	"time"

	stengo "common/scripts/tengo"

	"github.com/d5/tengo/objects"
	"github.com/d5/tengo/v2"
)

var tcpMethodMaps = map[string]*TCPMethodTengo{

	closeMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: closeMethod},
	},

	flushMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: flushMethod},
	},

	writeBytesMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: writeBytesMethod},
	},

	writeHexMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: writeHexMethod},
	},

	writeStringMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: writeStringMethod},
	},

	readLineMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: readLineMethod},
	},

	readBytesMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: readBytesMethod},
	},

	readStringMethod: &TCPMethodTengo{
		TengoObj: stengo.TengoObj{Name: readStringMethod},
	},
}

type TCPTengo struct {
	stengo.TengoObj
	conn *Connection
}

type TCPMethodTengo struct {
	stengo.TengoObj

	tcp *TCPTengo
}

/*
*
  args[0] ---network
  args[1] ---host
  args[2] ---port
  args[3] ---isSSL
  args[4] ---connectionTimeout
  args[5] ---readTimeout
  args[6] ---writeTimeout
*/

func newConnection(args ...objects.Object) (objects.Object, error) {

	var conn *Connection
	var err error

	if len(args) != 7 {

		return nil, tengo.ErrWrongNumArguments
	}

	network, ok := objects.ToString(args[0])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "network",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	host, ok := objects.ToString(args[1])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "host",
			Expected: "string(compatible)",
			Found:    args[1].TypeName(),
		}
	}

	port, ok := objects.ToInt(args[2])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "port",
			Expected: "int(compatible)",
			Found:    args[2].TypeName(),
		}
	}

	isSSL, ok := objects.ToBool(args[3])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "isSSL",
			Expected: "bool(compatible)",
			Found:    args[3].TypeName(),
		}
	}

	connTimeout, ok := objects.ToInt64(args[4])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "connTimeout",
			Expected: "int64(compatible)",
			Found:    args[4].TypeName(),
		}
	}

	readTimeout, ok := objects.ToInt64(args[5])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "readTimeout",
			Expected: "int64(compatible)",
			Found:    args[5].TypeName(),
		}
	}

	writeTimeout, ok := objects.ToInt64(args[6])
	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "writeTimeout",
			Expected: "int64(compatible)",
			Found:    args[6].TypeName(),
		}
	}

	if isSSL {

		conn, err = Dial(network, fmt.Sprintf("%s:%d", host, port),
			DialConnectTimeout(time.Duration(connTimeout)*time.Millisecond),
			DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
			DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond),
			DialTLSSkipVerify(true),
			DialUseTLS(true))

	} else {

		conn, err = Dial(network, fmt.Sprintf("%s:%d", host, port),
			DialConnectTimeout(time.Duration(connTimeout)*time.Millisecond),
			DialReadTimeout(time.Duration(readTimeout)*time.Millisecond),
			DialWriteTimeout(time.Duration(writeTimeout)*time.Millisecond))
	}

	if err != nil {
		return nil, err
	}

	return &TCPTengo{
		TengoObj: stengo.TengoObj{Name: "tcp"},
		conn:     conn,
	}, nil

}

func (t *TCPTengo) IndexGet(index objects.Object) (value objects.Object, err error) {

	key, ok := objects.ToString(index)

	if !ok {
		return nil, tengo.ErrInvalidArgumentType{
			Name:     "index",
			Expected: "string(compatible)",
			Found:    index.TypeName(),
		}
	}

	if m, ok := tcpMethodMaps[key]; ok {
		m.tcp = t

		return m, nil
	}

	return nil, errors.New("Cannot support method:" + key)
}

func (tm *TCPMethodTengo) Call(args ...objects.Object) (objects.Object, error) {

	switch tm.Name {
	case closeMethod:
		return tm.TClose(args...)

	case flushMethod:
		return tm.TFlush(args...)

	case writeBytesMethod:
		return tm.TWriteBytes(args...)

	case writeHexMethod:
		return tm.TWriteHex(args...)

	case writeStringMethod:
		return tm.TWriteString(args...)

	case readLineMethod:
		return tm.TReadLine(args...)

	case readBytesMethod:
		return tm.TReadBytes(args...)

	case readStringMethod:
		return tm.TReadString(args...)

	default:
		return nil, errors.New("unknown http client method:" + tm.Name)

	}

}

func (tm *TCPMethodTengo) TClose(args ...objects.Object) (objects.Object, error) {

	tm.tcp.conn.Close()

	return nil, nil
}

func (tm *TCPMethodTengo) TFlush(args ...objects.Object) (objects.Object, error) {

	tm.tcp.conn.Flush()

	return nil, nil
}

func (tm *TCPMethodTengo) TWriteBytes(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	data, ok := objects.ToByteSlice(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "[]byte(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil, tm.tcp.conn.WriteBytes(data)
}

func (tm *TCPMethodTengo) TWriteHex(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	data, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil, tm.tcp.conn.WriteHex(data)
}

func (tm *TCPMethodTengo) TWriteString(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	data, ok := objects.ToString(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "data",
			Expected: "string(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	return nil, tm.tcp.conn.WriteString(data)
}

func (tm *TCPMethodTengo) TReadLine(args ...objects.Object) (objects.Object, error) {

	data, err := tm.tcp.conn.ReadLine()
	if err != nil {

		return nil, err
	}

	return objects.FromInterface(data)
}

func (tm *TCPMethodTengo) TReadBytes(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	n, ok := objects.ToInt(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "n",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	data, err := tm.tcp.conn.ReadBytes(n)
	if err != nil {

		return nil, err
	}

	return objects.FromInterface(data)
}

func (tm *TCPMethodTengo) TReadString(args ...objects.Object) (objects.Object, error) {

	if len(args) != 1 {

		return nil, tengo.ErrWrongNumArguments
	}

	n, ok := objects.ToInt(args[0])
	if !ok {

		return nil, tengo.ErrInvalidArgumentType{
			Name:     "n",
			Expected: "int(compatible)",
			Found:    args[0].TypeName(),
		}
	}

	data, err := tm.tcp.conn.ReadBytes(n)
	if err != nil {

		return nil, err
	}

	return objects.FromInterface(string(data))
}

var moduleMap objects.Object = &objects.ImmutableMap{
	Value: map[string]objects.Object{
		newConnectionMethod: &objects.UserFunction{
			Name:  newConnectionMethod,
			Value: newConnection,
		},
	},
}

func (TCPTengo) Import(moduleName string) (interface{}, error) {

	switch moduleName {
	case tcpModName:
		return moduleMap, nil
	default:
		return nil, errors.New("undefined module:" + moduleName)
	}
}
