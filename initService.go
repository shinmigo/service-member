package main

import (
	grpcserver "goshop/service-member/pkg/grpc/server"
)

func initService() {
	go grpcserver.Run()
	//go user.Hello()
}
