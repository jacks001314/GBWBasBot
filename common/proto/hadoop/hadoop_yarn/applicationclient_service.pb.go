// Code generated by protoc-gen-go.
// source: applicationclient_protocol.proto
// DO NOT EDIT!

package hadoop_yarn

import (
	gohadoop "common/proto/hadoop"

	proto "github.com/golang/protobuf/proto"

	uuid "github.com/nu7hatch/gouuid"

	json "encoding/json"

	math "math"

	hadoop_common "common/proto/hadoop/hadoop_common"

	hadoop_ipc_client "common/proto/hadoop/hadoop_common/ipc/client"

	yarn_conf "common/proto/hadoop/hadoop_yarn/conf"
)

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf
var APPLICATION_CLIENT_PROTOCOL = "org.apache.hadoop.yarn.api.ApplicationClientProtocolPB"

func init() {
}

type ApplicationClientProtocolService interface {
	GetNewApplication(in *GetNewApplicationRequestProto, out *GetNewApplicationResponseProto) error
	GetApplicationReport(in *GetApplicationReportRequestProto, out *GetApplicationReportResponseProto) error
	GetApplicationAttemptReport(in *GetApplicationAttemptReportRequestProto, out *GetApplicationAttemptReportResponseProto) error
	SubmitApplication(in *SubmitApplicationRequestProto, out *SubmitApplicationResponseProto) error
	ForceKillApplication(in *KillApplicationRequestProto, out *KillApplicationResponseProto) error
	GetClusterMetrics(in *GetClusterMetricsRequestProto, out *GetClusterMetricsResponseProto) error
	GetApplications(in *GetApplicationsRequestProto, out *GetApplicationsResponseProto) error
	GetClusterNodes(in *GetClusterNodesRequestProto, out *GetClusterNodesResponseProto) error
	GetQueueInfo(in *GetQueueInfoRequestProto, out *GetQueueInfoResponseProto) error
	GetQueueUserAcls(in *GetQueueUserAclsInfoRequestProto, out *GetQueueUserAclsInfoResponseProto) error
	GetDelegationToken(in *hadoop_common.GetDelegationTokenRequestProto, out *hadoop_common.GetDelegationTokenResponseProto) error
	RenewDelegationToken(in *hadoop_common.RenewDelegationTokenRequestProto, out *hadoop_common.RenewDelegationTokenResponseProto) error
	CancelDelegationToken(in *hadoop_common.CancelDelegationTokenRequestProto, out *hadoop_common.CancelDelegationTokenResponseProto) error
	Close()
}

type ApplicationClientProtocolServiceClient struct {
	*hadoop_ipc_client.Client
}

func (c *ApplicationClientProtocolServiceClient) GetNewApplication(in *GetNewApplicationRequestProto, out *GetNewApplicationResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetApplicationReport(in *GetApplicationReportRequestProto, out *GetApplicationReportResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetApplicationAttemptReport(in *GetApplicationAttemptReportRequestProto, out *GetApplicationAttemptReportResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) SubmitApplication(in *SubmitApplicationRequestProto, out *SubmitApplicationResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) ForceKillApplication(in *KillApplicationRequestProto, out *KillApplicationResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetClusterMetrics(in *GetClusterMetricsRequestProto, out *GetClusterMetricsResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetApplications(in *GetApplicationsRequestProto, out *GetApplicationsResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetClusterNodes(in *GetClusterNodesRequestProto, out *GetClusterNodesResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetQueueInfo(in *GetQueueInfoRequestProto, out *GetQueueInfoResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetQueueUserAcls(in *GetQueueUserAclsInfoRequestProto, out *GetQueueUserAclsInfoResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) GetDelegationToken(in *hadoop_common.GetDelegationTokenRequestProto, out *hadoop_common.GetDelegationTokenResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) RenewDelegationToken(in *hadoop_common.RenewDelegationTokenRequestProto, out *hadoop_common.RenewDelegationTokenResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}
func (c *ApplicationClientProtocolServiceClient) CancelDelegationToken(in *hadoop_common.CancelDelegationTokenRequestProto, out *hadoop_common.CancelDelegationTokenResponseProto) error {
	return c.Call(gohadoop.GetCalleeRPCRequestHeaderProto(&APPLICATION_CLIENT_PROTOCOL), in, out)
}

// DialApplicationClientProtocolService connects to an ApplicationClientProtocolService at the specified network address.
func DialApplicationClientProtocolService(conf yarn_conf.YarnConfiguration) (ApplicationClientProtocolService, error) {
	clientId, _ := uuid.NewV4()
	ugi, _ := gohadoop.CreateSimpleUGIProto()
	serverAddress, _ := conf.GetRMAddress()
	c := &hadoop_ipc_client.Client{ClientId: clientId, Ugi: ugi, ServerAddress: serverAddress, Conn: nil}

	return &ApplicationClientProtocolServiceClient{c}, nil
}

/*
// DialApplicationClientProtocolServiceTimeout connects to an ApplicationClientProtocolService at the specified network address.
func DialApplicationClientProtocolServiceTimeout(network, addr string,
	timeout time.Duration) (*ApplicationClientProtocolServiceClient, *rpc.Client, error) {
	c, err := protorpc.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, nil, err
	}
	return &ApplicationClientProtocolServiceClient{c}, c, nil
}*/