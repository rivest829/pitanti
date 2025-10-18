package main

import (
	"flag"
	"pitanti/app"
	"pitanti/conf"
)

func main() {
	serverId := flag.Int("serverId", 1, "server id")
	env := flag.String("env", "dev", "the server type")
	rpcServerPort := flag.Int("rpcsvport", *serverId+3200, "the port that grpc server will listen")

	flag.Parse()
	app.Start(conf.Env(*env), 0, string(conf.ServerGame), *rpcServerPort)
}
