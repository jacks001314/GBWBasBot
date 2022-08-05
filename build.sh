#!/usr/bin/env bash

export GO111MODULE=on
export GOPROXY=https://goproxy.cn

rm -rf build
mkdir -p build/bin
mkdir -p build/conf/scripts/

echo "build packet scan module-----"
cd pscan/

go build -o pscan cmd/main.go

mv pscan ../build/bin

cd ../

echo "build assets detect module"
cd detect

go build -o detect cmd/main.go

mv detect ../build/bin
cd ../

echo "build attack module"

cd attack 

go build -o attack cmd/attack/main.go
go build -o atarget cmd/source/main.go

mv attack atarget ../build/bin
cd ../

cp -rf attack/scripts/* build/conf/scripts/

echo "build GBWBasBot into build dir is ok.................."
