package attack

type AttackTarget struct {
	IP string //攻击目标IP

	Host string //攻击目标主机

	Port int //攻击目标端口

	IsSSL bool //是否需要ssl链接

	Version string //应用程序版本

	Proto string //使用的协议（http,tcp,udp,....)

	App string //攻击目标应用名称

	Types map[string]struct{} //攻击类型列表
}
