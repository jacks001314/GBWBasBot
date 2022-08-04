package attack

type AttackProcess struct {
	IP      string //攻击目标IP
	Host    string //攻击目标host
	Port    int    //攻击目标端口
	Proto   string //应用协议
	App     string //应用名称
	OS      string //目标操作系统
	Version string //应用程序版本

	IsSSL bool //是否使用ssl

	Id int //攻击模块编号

	Language string //攻击模块语言

	Name string //攻击模块名称

	Type string //攻击利用类型

	CVECode string //CVE编号

	Desc string //攻击描述

	Suggest string //修复建议

	Status int //攻击状态

	Payload string //攻击载荷

	Result string //攻击结果

	Details string //攻击详情
}
