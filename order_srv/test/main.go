package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop-srvs/order_srv/proto/proto"
)

var (
	conn            *grpc.ClientConn
	inventoryClient proto.OrderClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.0.101:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	inventoryClient = proto.NewOrderClient(conn)
}

func TestCreateCartItem() {
	_, err := inventoryClient.CreateCartItem(context.Background(), &proto.CartItemRequest{})

	if err != nil {
		panic(err)
	}

	fmt.Println("Success")
}

func main() {
	Init()

	err := conn.Close()
	if err != nil {
		return
	}
}
