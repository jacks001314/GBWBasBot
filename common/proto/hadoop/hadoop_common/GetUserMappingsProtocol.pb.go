// Code generated by protoc-gen-go.
// source: GetUserMappingsProtocol.proto
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

// *
//  Get groups for user request.
type GetGroupsForUserRequestProto struct {
	User             *string `protobuf:"bytes,1,req,name=user" json:"user,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *GetGroupsForUserRequestProto) Reset()         { *m = GetGroupsForUserRequestProto{} }
func (m *GetGroupsForUserRequestProto) String() string { return proto.CompactTextString(m) }
func (*GetGroupsForUserRequestProto) ProtoMessage()    {}

func (m *GetGroupsForUserRequestProto) GetUser() string {
	if m != nil && m.User != nil {
		return *m.User
	}
	return ""
}

// *
// Response for get groups.
type GetGroupsForUserResponseProto struct {
	Groups           []string `protobuf:"bytes,1,rep,name=groups" json:"groups,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *GetGroupsForUserResponseProto) Reset()         { *m = GetGroupsForUserResponseProto{} }
func (m *GetGroupsForUserResponseProto) String() string { return proto.CompactTextString(m) }
func (*GetGroupsForUserResponseProto) ProtoMessage()    {}

func (m *GetGroupsForUserResponseProto) GetGroups() []string {
	if m != nil {
		return m.Groups
	}
	return nil
}

func init() {
}
