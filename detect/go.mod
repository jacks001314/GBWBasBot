module detect

go 1.15

replace common => ../common

require (
	common v0.0.0-00010101000000-000000000000
	github.com/d5/tengo v1.24.8
	github.com/d5/tengo/v2 v2.10.0
	github.com/panjf2000/ants/v2 v2.5.0
	github.com/sirupsen/logrus v1.9.0
	github.com/yuin/gopher-lua v0.0.0-20220413183635-c841877397d8
)
