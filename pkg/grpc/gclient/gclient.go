package gclient

import (
	"fmt"
	"log"
	"strings"
	
	"github.com/shinmigo/pb/orderpb"
	"goshop/service-member/pkg/grpc/etcd3"
	"goshop/service-member/pkg/utils"
	
	"github.com/shinmigo/pb/shoppb"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

var (
	ShopAddress shoppb.AreaServiceClient
	OrderClient orderpb.OrderServiceClient
)

func DialGrpcService() {
	shop()
	oms()
}
func shop() {
	r := etcd3.NewResolver(utils.C.Etcd.Host)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://author/"+utils.C.GrpcClient.Name["shop"], grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc没有连接上%s, err: %v \n", utils.C.GrpcClient.Name["shop"], err)
	}
	fmt.Printf("连接成功：%s, host分别为: %s \n", utils.C.GrpcClient.Name["shop"], strings.Join(utils.C.Etcd.Host, ","))
	ShopAddress = shoppb.NewAreaServiceClient(conn)
}

func oms() {
	r := etcd3.NewResolver(utils.C.Etcd.Host)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://author/"+utils.C.GrpcClient.Name["oms"], grpc.WithBalancerName("round_robin"), grpc.WithInsecure())
	if err != nil {
		log.Panicf("grpc没有连接上%s, err: %v \n", utils.C.GrpcClient.Name["oms"], err)
	}
	fmt.Printf("连接成功：%s, host分别为: %s \n", utils.C.GrpcClient.Name["oms"], strings.Join(utils.C.Etcd.Host, ","))
	OrderClient = orderpb.NewOrderServiceClient(conn)
}
