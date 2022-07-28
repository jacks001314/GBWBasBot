// Code generated by protoc-gen-go.
// source: ZKFCProtocol.proto
// DO NOT EDIT!

package hadoop_common

import (
	json "encoding/json"

	proto "github.com/golang/protobuf/proto"

	math "math"
)

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type CedeActiveRequestProto struct {
	MillisToCede     *uint32 `protobuf:"varint,1,req,name=millisToCede" json:"millisToCede,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *CedeActiveRequestProto) Reset()         { *m = CedeActiveRequestProto{} }
func (m *CedeActiveRequestProto) String() string { return proto.CompactTextString(m) }
func (*CedeActiveRequestProto) ProtoMessage()    {}

func (m *CedeActiveRequestProto) GetMillisToCede() uint32 {
	if m != nil && m.MillisToCede != nil {
		return *m.MillisToCede
	}
	return 0
}

type CedeActiveResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *CedeActiveResponseProto) Reset()         { *m = CedeActiveResponseProto{} }
func (m *CedeActiveResponseProto) String() string { return proto.CompactTextString(m) }
func (*CedeActiveResponseProto) ProtoMessage()    {}

type GracefulFailoverRequestProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *GracefulFailoverRequestProto) Reset()         { *m = GracefulFailoverRequestProto{} }
func (m *GracefulFailoverRequestProto) String() string { return proto.CompactTextString(m) }
func (*GracefulFailoverRequestProto) ProtoMessage()    {}

type GracefulFailoverResponseProto struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *GracefulFailoverResponseProto) Reset()         { *m = GracefulFailoverResponseProto{} }
func (m *GracefulFailoverResponseProto) String() string { return proto.CompactTextString(m) }
func (*GracefulFailoverResponseProto) ProtoMessage()    {}

func init() {
}
