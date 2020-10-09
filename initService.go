package main

import "goshop/service-member/pkg/grpc/gclient"

func initService() {
	go gclient.DialGrpcService()
	//go user.Hello()
}
