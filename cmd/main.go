package main

import (
	"go_awd/conf"
	"go_awd/routers"
)

// main
// @Description: 开启所有服务
func main() {
	conf.Init()
	r := routers.NewRouter()
	if err := r.Run(conf.HttpPort); err != nil {
		panic(err)
	}
}
