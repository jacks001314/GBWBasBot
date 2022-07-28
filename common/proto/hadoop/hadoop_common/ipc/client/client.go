package ipc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	gohadoop "common/proto/hadoop"
	"common/proto/hadoop/hadoop_common"
	"common/proto/hadoop/hadoop_common/security"

	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	ClientId      *uuid.UUID
	Ugi           *hadoop_common.UserInformationProto
	ServerAddress string
	TCPNoDelay    bool

	Conn *connection
}

type connection struct {
	con net.Conn
}

type connection_id struct {
	user     string
	protocol string
	address  string
}

type call struct {
	callId     int32
	procedure  proto.Message
	request    proto.Message
	response   proto.Message
	err        *error
	retryCount int32
}

func (c *Client) String() string {
	buf := bytes.NewBufferString("")
	fmt.Fprint(buf, "<clientId:", c.ClientId)
	fmt.Fprint(buf, ", server:", c.ServerAddress)
	fmt.Fprint(buf, ">")
	return buf.String()
}

var (
	SASL_RPC_DUMMY_CLIENT_ID     []byte = make([]byte, 0)
	SASL_RPC_CALL_ID             int32  = -33
	SASL_RPC_INVALID_RETRY_COUNT int32  = -1
)

func (c *Client) Close() {

	if c.Conn != nil {
		c.Conn.con.Close()
	}
}

func (c *Client) Call(rpc *hadoop_common.RequestHeaderProto, rpcRequest proto.Message, rpcResponse proto.Message) error {
	// Create connection_id
	connectionId := connection_id{user: *c.Ugi.RealUser, protocol: *rpc.DeclaringClassProtocolName, address: c.ServerAddress}

	// Get connection to server
	//log.Println("Connecting...", c)
	conn, err := getConnection(c, &connectionId)
	if err != nil {
		return err
	}

	// Create call and send request
	rpcCall := call{callId: 0, procedure: rpc, request: rpcRequest, response: rpcResponse}
	err = sendRequest(c, conn, &rpcCall)
	if err != nil {
		return err
	}

	// Read & return response
	err = c.readResponse(conn, &rpcCall)

	return err
}

var connectionPool = struct {
	sync.RWMutex
	connections map[connection_id]*connection
}{connections: make(map[connection_id]*connection)}

func findUsableTokenForService(service string) (*hadoop_common.TokenProto, bool) {
	userTokens := security.GetCurrentUser().GetUserTokens()

	log.Printf("looking for token for service: %s\n", service)

	if len(userTokens) == 0 {
		return nil, false
	}

	token := userTokens[service]
	if token != nil {
		return token, true
	}

	return nil, false
}

func getConnection(c *Client, connectionId *connection_id) (*connection, error) {

	var con *connection
	var err error
	if c.Conn != nil {
		return c.Conn, nil
	}

	con, err = setupConnectionWithTimeout(c)

	if err != nil {
		return nil, err

	}

	c.Conn = con

	var authProtocol gohadoop.AuthProtocol = gohadoop.AUTH_PROTOCOL_NONE

	if _, found := findUsableTokenForService(c.ServerAddress); found {
		log.Printf("found token for service: %s", c.ServerAddress)
		authProtocol = gohadoop.AUTH_PROTOCOL_SASL
	}

	writeConnectionHeader(con, authProtocol)

	if authProtocol == gohadoop.AUTH_PROTOCOL_SASL {

		log.Println("attempting SASL negotiation.")

		if err = negotiateSimpleTokenAuth(c, con); err != nil {
			return nil, err
		}
	} else {
		log.Println("no usable tokens. proceeding without auth.")

	}

	writeConnectionContext(c, con, connectionId, authProtocol)

	return con, nil
}

func setupConnectionWithTimeout(c *Client) (*connection, error) {

	do := &net.Dialer{
		Timeout:   time.Second * 30,
		KeepAlive: time.Minute * 5,
	}

	netConn, err := do.DialContext(context.Background(), "tcp", c.ServerAddress)

	if err != nil {
		return nil, err
	}

	return &connection{con: netConn}, nil
}

func setupConnection(c *Client) (*connection, error) {
	addr, _ := net.ResolveTCPAddr("tcp", c.ServerAddress)
	tcpConn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {

		return nil, err
	} else {

	}

	// TODO: Ping thread

	// Set tcp no-delay
	tcpConn.SetNoDelay(c.TCPNoDelay)

	return &connection{tcpConn}, nil
}

func writeConnectionHeader(conn *connection, authProtocol gohadoop.AuthProtocol) error {
	// RPC_HEADER
	if _, err := conn.con.Write(gohadoop.RPC_HEADER); err != nil {

		return err
	}

	// RPC_VERSION
	if _, err := conn.con.Write(gohadoop.VERSION); err != nil {

		return err
	}

	// RPC_SERVICE_CLASS
	if serviceClass, err := gohadoop.ConvertFixedToBytes(gohadoop.RPC_SERVICE_CLASS); err != nil {

		return err
	} else if _, err := conn.con.Write(serviceClass); err != nil {

		return err
	}

	// AuthProtocol
	if authProtocolBytes, err := gohadoop.ConvertFixedToBytes(authProtocol); err != nil {

		return err
	} else if _, err := conn.con.Write(authProtocolBytes); err != nil {

		return err
	}

	return nil
}

func writeConnectionContext(c *Client, conn *connection, connectionId *connection_id, authProtocol gohadoop.AuthProtocol) error {
	// Create hadoop_common.IpcConnectionContextProto
	ugi, _ := gohadoop.CreateSimpleUGIProto()
	ipcCtxProto := hadoop_common.IpcConnectionContextProto{UserInfo: ugi, Protocol: &connectionId.protocol}

	// Create RpcRequestHeaderProto
	var callId int32 = -3
	var clientId [16]byte = [16]byte(*c.ClientId)

	/*if (authProtocol == gohadoop.AUTH_PROTOCOL_SASL) {
	  callId = SASL_RPC_CALL_ID
	}*/

	rpcReqHeaderProto := hadoop_common.RpcRequestHeaderProto{RpcKind: &gohadoop.RPC_PROTOCOL_BUFFFER, RpcOp: &gohadoop.RPC_FINAL_PACKET, CallId: &callId, ClientId: clientId[0:16], RetryCount: &gohadoop.RPC_DEFAULT_RETRY_COUNT}

	rpcReqHeaderProtoBytes, err := proto.Marshal(&rpcReqHeaderProto)
	if err != nil {

		return err
	}

	ipcCtxProtoBytes, _ := proto.Marshal(&ipcCtxProto)
	if err != nil {

		return err
	}

	totalLength := len(rpcReqHeaderProtoBytes) + sizeVarint(len(rpcReqHeaderProtoBytes)) + len(ipcCtxProtoBytes) + sizeVarint(len(ipcCtxProtoBytes))
	var tLen int32 = int32(totalLength)
	totalLengthBytes, err := gohadoop.ConvertFixedToBytes(tLen)

	if err != nil {

		return err
	} else if _, err := conn.con.Write(totalLengthBytes); err != nil {

		return err
	}

	if err := writeDelimitedBytes(conn, rpcReqHeaderProtoBytes); err != nil {

		return err
	}
	if err := writeDelimitedBytes(conn, ipcCtxProtoBytes); err != nil {

		return err
	}

	return nil
}

func sizeVarint(x int) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func sendRequest(c *Client, conn *connection, rpcCall *call) error {
	//log.Println("About to call RPC: ", rpcCall.procedure)

	// 0. RpcRequestHeaderProto
	var clientId [16]byte = [16]byte(*c.ClientId)
	rpcReqHeaderProto := hadoop_common.RpcRequestHeaderProto{RpcKind: &gohadoop.RPC_PROTOCOL_BUFFFER, RpcOp: &gohadoop.RPC_FINAL_PACKET, CallId: &rpcCall.callId, ClientId: clientId[0:16], RetryCount: &rpcCall.retryCount}
	rpcReqHeaderProtoBytes, err := proto.Marshal(&rpcReqHeaderProto)
	if err != nil {

		return err
	}

	// 1. RequestHeaderProto
	requestHeaderProto := rpcCall.procedure
	requestHeaderProtoBytes, err := proto.Marshal(requestHeaderProto)
	if err != nil {

		return err
	}

	// 2. Param
	paramProto := rpcCall.request
	paramProtoBytes, err := proto.Marshal(paramProto)
	if err != nil {

		return err
	}

	totalLength := len(rpcReqHeaderProtoBytes) + sizeVarint(len(rpcReqHeaderProtoBytes)) + len(requestHeaderProtoBytes) + sizeVarint(len(requestHeaderProtoBytes)) + len(paramProtoBytes) + sizeVarint(len(paramProtoBytes))
	var tLen int32 = int32(totalLength)
	if totalLengthBytes, err := gohadoop.ConvertFixedToBytes(tLen); err != nil {

		return err
	} else {
		if _, err := conn.con.Write(totalLengthBytes); err != nil {

			return err
		}
	}

	if err := writeDelimitedBytes(conn, rpcReqHeaderProtoBytes); err != nil {

		return err
	}
	if err := writeDelimitedBytes(conn, requestHeaderProtoBytes); err != nil {

		return err
	}
	if err := writeDelimitedBytes(conn, paramProtoBytes); err != nil {

		return err
	}

	//log.Println("Succesfully sent request of length: ", totalLength)

	return nil
}

func writeDelimitedTo(conn *connection, msg proto.Message) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {

		return err
	}
	return writeDelimitedBytes(conn, msgBytes)
}

func writeDelimitedBytes(conn *connection, data []byte) error {
	if _, err := conn.con.Write(proto.EncodeVarint(uint64(len(data)))); err != nil {

		return err
	}
	if _, err := conn.con.Write(data); err != nil {

		return err
	}

	return nil
}

func (c *Client) readResponse(conn *connection, rpcCall *call) error {
	// Read first 4 bytes to get total-length
	var totalLength int32 = -1
	var totalLengthBytes [4]byte
	if _, err := conn.con.Read(totalLengthBytes[0:4]); err != nil {

		return err
	}

	if err := gohadoop.ConvertBytesToFixed(totalLengthBytes[0:4], &totalLength); err != nil {

		return err
	}

	var responseBytes []byte = make([]byte, totalLength)
	if _, err := conn.con.Read(responseBytes); err != nil {

		return err
	}

	// Parse RpcResponseHeaderProto
	rpcResponseHeaderProto := hadoop_common.RpcResponseHeaderProto{}
	off, err := readDelimited(responseBytes[0:totalLength], &rpcResponseHeaderProto)
	if err != nil {

		return err
	}
	//log.Println("Received rpcResponseHeaderProto = ", rpcResponseHeaderProto)

	err = c.checkRpcHeader(&rpcResponseHeaderProto)
	if err != nil {

		return err
	}

	if *rpcResponseHeaderProto.Status == hadoop_common.RpcResponseHeaderProto_SUCCESS {
		// Parse RpcResponseWrapper
		_, err = readDelimited(responseBytes[off:], rpcCall.response)
	} else {
		log.Println("RPC failed with status: ", rpcResponseHeaderProto.Status.String())
		errorDetails := [4]string{rpcResponseHeaderProto.Status.String(), "ServerDidNotSetExceptionClassName", "ServerDidNotSetErrorMsg", "ServerDidNotSetErrorDetail"}
		if rpcResponseHeaderProto.ExceptionClassName != nil {
			errorDetails[0] = *rpcResponseHeaderProto.ExceptionClassName
		}
		if rpcResponseHeaderProto.ErrorMsg != nil {
			errorDetails[1] = *rpcResponseHeaderProto.ErrorMsg
		}
		if rpcResponseHeaderProto.ErrorDetail != nil {
			errorDetails[2] = rpcResponseHeaderProto.ErrorDetail.String()
		}
		err = errors.New(strings.Join(errorDetails[:], ":"))
	}
	return err
}

func readDelimited(rawData []byte, msg proto.Message) (int, error) {
	headerLength, off := proto.DecodeVarint(rawData)
	if off == 0 {

		return -1, nil
	}
	err := proto.Unmarshal(rawData[off:off+int(headerLength)], msg)
	if err != nil {

		return -1, err
	}

	return off + int(headerLength), nil
}

func (c *Client) checkRpcHeader(rpcResponseHeaderProto *hadoop_common.RpcResponseHeaderProto) error {
	var callClientId [16]byte = [16]byte(*c.ClientId)
	var headerClientId []byte = []byte(rpcResponseHeaderProto.ClientId)

	if rpcResponseHeaderProto.ClientId != nil && len(headerClientId) >= 16 {
		if !bytes.Equal(callClientId[0:16], headerClientId[0:16]) {

			return errors.New("Incorrect clientId")
		}
	}
	return nil
}

func sendSaslMessage(c *Client, conn *connection, message *hadoop_common.RpcSaslProto) error {
	saslRpcHeaderProto := hadoop_common.RpcRequestHeaderProto{RpcKind: &gohadoop.RPC_PROTOCOL_BUFFFER,
		RpcOp:      &gohadoop.RPC_FINAL_PACKET,
		CallId:     &SASL_RPC_CALL_ID,
		ClientId:   SASL_RPC_DUMMY_CLIENT_ID,
		RetryCount: &SASL_RPC_INVALID_RETRY_COUNT}

	saslRpcHeaderProtoBytes, err := proto.Marshal(&saslRpcHeaderProto)

	if err != nil {

		return err
	}

	saslRpcMessageProtoBytes, err := proto.Marshal(message)

	if err != nil {

		return err
	}

	totalLength := len(saslRpcHeaderProtoBytes) + sizeVarint(len(saslRpcHeaderProtoBytes)) + len(saslRpcMessageProtoBytes) + sizeVarint(len(saslRpcMessageProtoBytes))
	var tLen int32 = int32(totalLength)

	if totalLengthBytes, err := gohadoop.ConvertFixedToBytes(tLen); err != nil {

		return err
	} else {
		if _, err := conn.con.Write(totalLengthBytes); err != nil {

			return err
		}
	}
	if err := writeDelimitedBytes(conn, saslRpcHeaderProtoBytes); err != nil {

		return err
	}
	if err := writeDelimitedBytes(conn, saslRpcMessageProtoBytes); err != nil {

		return err
	}

	return nil
}

func receiveSaslMessage(c *Client, conn *connection) (*hadoop_common.RpcSaslProto, error) {
	// Read first 4 bytes to get total-length
	var totalLength int32 = -1
	var totalLengthBytes [4]byte

	if _, err := conn.con.Read(totalLengthBytes[0:4]); err != nil {

		return nil, err
	}
	if err := gohadoop.ConvertBytesToFixed(totalLengthBytes[0:4], &totalLength); err != nil {

		return nil, err
	}

	var responseBytes []byte = make([]byte, totalLength)

	if _, err := conn.con.Read(responseBytes); err != nil {

		return nil, err
	}

	// Parse RpcResponseHeaderProto
	rpcResponseHeaderProto := hadoop_common.RpcResponseHeaderProto{}
	off, err := readDelimited(responseBytes[0:totalLength], &rpcResponseHeaderProto)
	if err != nil {

		return nil, err
	}

	err = checkSaslRpcHeader(&rpcResponseHeaderProto)
	if err != nil {

		return nil, err
	}

	var saslRpcMessage hadoop_common.RpcSaslProto

	if *rpcResponseHeaderProto.Status == hadoop_common.RpcResponseHeaderProto_SUCCESS {
		// Parse RpcResponseWrapper
		if _, err = readDelimited(responseBytes[off:], &saslRpcMessage); err != nil {

			return nil, err
		} else {
			return &saslRpcMessage, nil
		}
	} else {
		log.Println("RPC failed with status: ", rpcResponseHeaderProto.Status.String())
		errorDetails := [4]string{rpcResponseHeaderProto.Status.String(), "ServerDidNotSetExceptionClassName", "ServerDidNotSetErrorMsg", "ServerDidNotSetErrorDetail"}
		if rpcResponseHeaderProto.ExceptionClassName != nil {
			errorDetails[0] = *rpcResponseHeaderProto.ExceptionClassName
		}
		if rpcResponseHeaderProto.ErrorMsg != nil {
			errorDetails[1] = *rpcResponseHeaderProto.ErrorMsg
		}
		if rpcResponseHeaderProto.ErrorDetail != nil {
			errorDetails[2] = rpcResponseHeaderProto.ErrorDetail.String()
		}
		err = errors.New(strings.Join(errorDetails[:], ":"))
		return nil, err
	}
}

func checkSaslRpcHeader(rpcResponseHeaderProto *hadoop_common.RpcResponseHeaderProto) error {
	var headerClientId []byte = []byte(rpcResponseHeaderProto.ClientId)
	if rpcResponseHeaderProto.ClientId != nil {
		if !bytes.Equal(SASL_RPC_DUMMY_CLIENT_ID, headerClientId) {

			return errors.New("Incorrect clientId")
		}
	}
	return nil
}

func negotiateSimpleTokenAuth(client *Client, con *connection) error {
	var saslNegotiateState hadoop_common.RpcSaslProto_SaslState = hadoop_common.RpcSaslProto_NEGOTIATE
	var saslNegotiateMessage hadoop_common.RpcSaslProto = hadoop_common.RpcSaslProto{State: &saslNegotiateState}
	var saslResponseMessage *hadoop_common.RpcSaslProto
	var err error

	//send a SASL negotiation request
	if err = sendSaslMessage(client, con, &saslNegotiateMessage); err != nil {

		return err
	}

	//get a response with supported mehcanisms/challenge
	if saslResponseMessage, err = receiveSaslMessage(client, con); err != nil {

		return err
	}

	var auths []*hadoop_common.RpcSaslProto_SaslAuth = saslResponseMessage.GetAuths()

	if numAuths := len(auths); numAuths <= 0 {

		return errors.New("No supported auth mechanisms!")
	}

	//for now we only support auth when TOKEN/DIGEST-MD5 is the first/only
	//supported auth mechanism
	var auth *hadoop_common.RpcSaslProto_SaslAuth = auths[0]

	if !(auth.GetMethod() == "TOKEN" && auth.GetMechanism() == "DIGEST-MD5") {

		return errors.New("gohadoop only supports TOKEN/DIGEST-MD5 auth!")
	}

	method := auth.GetMethod()
	mechanism := auth.GetMechanism()
	protocol := auth.GetProtocol()
	serverId := auth.GetServerId()
	challenge := auth.GetChallenge()

	//TODO: token/service mapping + token selection based on type/service
	//we wouldn't have gotten this far if there wasn't at least one available token.
	userToken, _ := findUsableTokenForService(client.ServerAddress)
	response, err := security.GetDigestMD5ChallengeResponse(protocol, serverId, challenge, userToken)

	if err != nil {

		return err
	}

	saslInitiateState := hadoop_common.RpcSaslProto_INITIATE
	authSend := hadoop_common.RpcSaslProto_SaslAuth{Method: &method, Mechanism: &mechanism,
		Protocol: &protocol, ServerId: &serverId}
	authsSendArray := []*hadoop_common.RpcSaslProto_SaslAuth{&authSend}
	saslInitiateMessage := hadoop_common.RpcSaslProto{State: &saslInitiateState,
		Token: []byte(response), Auths: authsSendArray}

	//send a SASL inititate request
	if err = sendSaslMessage(client, con, &saslInitiateMessage); err != nil {

		return err
	}

	//get a response with supported mehcanisms/challenge
	if saslResponseMessage, err = receiveSaslMessage(client, con); err != nil {

		return err
	}

	if saslResponseMessage.GetState() != hadoop_common.RpcSaslProto_SUCCESS {

		return errors.New("expected SASL SUCCESS response!")
	}

	log.Println("Successfully completed SASL negotiation!")

	return nil //errors.New("abort here")
}
