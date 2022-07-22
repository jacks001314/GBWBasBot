package detect

//探测到的结果
type DResult struct {
	IP      string //目标IP
	Port    uint16 //目标端口
	App     string //应用名称
	Version string //应用程序版本
	Proto   string //探测使用的协议（http,tcp,udp,....)
	IsSSL   bool   //是否使用SSL
}
