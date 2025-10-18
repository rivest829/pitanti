package main

import (
	"flag"
	"pitanti/app"
	"pitanti/conf"
)

func main() {
	port := flag.Int("port", 3100, "the port to listen")
	env := flag.String("env", "dev", "the server type")
	rpcServerPort := flag.Int("rpcsvport", 3200, "the port that grpc server will listen")
	flag.Parse()
	app.Start(conf.Env(*env), *port, string(conf.ServerGateway), *rpcServerPort)
}
