package pscan

type PScanConfig struct {
	SrcIP string `json:"srcIP"`

	Ifname string `json:"ifname"`

	GatewayIP string `json:"gatewayIP"`

	GatewayMac string `json:"gatewayMac"`

	SourcePortFirst uint32 `json:"SourcePortFirst"`

	SourcePortLast uint32 `json:"SourcePortLast"`

	SendRate uint64 `json:"sendRate"`

	SendBandWidth uint64 `json:"sendBandWidth"`

	WaitRecvTime uint64 `json:"waitRecvTime"`
}
