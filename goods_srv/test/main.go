package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"srvs/goods_srv/proto/proto"
)

var (
	conn        *grpc.ClientConn
	goodsClient proto.GoodsClient
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.0.101:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	goodsClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()

	//list, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
	//	Mobile:   "15700188888",
	//	NickName: "Cowboy",
	//	PassWord: "123456",
	//})

	list, err := goodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       1,
		PagePerNums: 10,
	})

	//list, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: "15700188888 "})

	//list, err := userClient.GetUserById(context.Background(), &proto.IdRequest{Id: 2})

	//list, err := userClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
	//	Id:       1,
	//	NickName: "gopher",
	//	BirthDay: uint32(time.Now().Unix()),
	//	Gender:   1,
	//})

	//list, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
	//	PassWord:          "212112",
	//	EncryptedPassWord: "pbkdf2-sha512$oVj8oVEEQb062eBA$7ab8439632fbb75e9fd50390d7cb3d2d1affcb5b03a4983e879e0913a28b9156",
	//})

	if err != nil {
		panic(err)
	}
	fmt.Println(1111)
	fmt.Println(list)
	err = conn.Close()
	if err != nil {
		return
	}
}
