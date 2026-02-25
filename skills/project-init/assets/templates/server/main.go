package main

import (
	"flag"
	"server/core"
	"server/entrance"
)

func main() {
	// 解析命令行参数
	port := flag.Int("port", 0, "覆盖配置文件中的端口号")
	flag.Parse()

	err := core.InitAll("local.yaml")
	if err != nil {
		panic(err)
	}
	srv := entrance.RunServer(*port)
	entrance.WaitForShutdown()
	entrance.ShutdownServer(srv)
}
