package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"srvs/inventory_srv/proto/proto"
	"sync"
)

var (
	conn            *grpc.ClientConn
	inventoryClient proto.InventoryClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.0.101:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	inventoryClient = proto.NewInventoryClient(conn)
}

func TestSetInv(GoodsId int32) {
	_, err := inventoryClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: GoodsId,
		Num:     100,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(GoodsId)
}

func TestInvDetail() {
	res, err := inventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
		Num:     100,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}

func TestSell(wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := inventoryClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInvInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     1,
			},
		},
		GoodsSn: "",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(111)
}

func TestReBack() {
	res, err := inventoryClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInvInfo: []*proto.GoodsInvInfo{
			{
				GoodsId: 421,
				Num:     10,
			},
			{
				GoodsId: 422,
				Num:     20,
			},
		},
		GoodsSn: "",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
func main() {
	Init()

	//TestSetInv(421)

	//TestInvDetail()

	//TestSell()

	//TestReBack()

	wg := sync.WaitGroup{}

	wg.Add(30)

	for i := 0; i < 30; i++ {
		go TestSell(&wg)
	}

	wg.Wait()

	err := conn.Close()
	if err != nil {
		return
	}
}
