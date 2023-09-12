package main

import (
	"context"
	"fmt"
	"inventory/config"
	rpc "inventory/controllers"
	h "inventory/proto"
	"inventory/services"
	"net"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)


func initApp(mclient *mongo.Client) {
	ctx := context.Background()
	mcoll := mclient.Database("inventory_SKU").Collection("inventory")
	iservice := services.InitInventory(mcoll, ctx)
	rpc.InventoryService = iservice
}

func main(){
	fmt.Println("Starting server...")
	client, err := config.ConnectDataBase()
	defer client.Disconnect(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	initApp(client)
	fmt.Println("Connected to MongoDB!")
	lis,err := net.Listen("tcp",":2023")
	if err!=nil{
		fmt.Println("Failed to listen: ",err)
		return
	}
	s := grpc.NewServer()
	h.RegisterInventoryServiceServer(s,&rpc.RPCServer{})

	fmt.Println("Server listening on port: 2023")
	if err := s.Serve(lis); err!=nil{
		fmt.Println("Failed to serve: ",err)
		return
	}
}