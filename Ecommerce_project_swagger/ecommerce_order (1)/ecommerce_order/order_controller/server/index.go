package main

import (
	"context"
	"fmt"
	"net"

	"ecommerce_order/order_config/config"
	"ecommerce_order/order_config/constants"
	"ecommerce_order/order_controller/controller"
	"ecommerce_order/order_dal/services"
	pro "ecommerce_order/order_proto"


	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func initDatabase(client *mongo.Client) {
	OrderCollection := config.GetCollection(client, "DataBase", "Order")
	InventoryCollection:=config.GetCollection(client,"inventory_SKU","items")
	controller.OrderService = services.InitCustomerService(client, OrderCollection,InventoryCollection, context.Background())
}

func main() {
	mongoclient, err := config.ConnectDataBase()
	defer mongoclient.Disconnect(context.TODO())
	if err != nil {
		panic(err)
	}
	initDatabase(mongoclient)
	lis, err := net.Listen("tcp", constants.Port)
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		return
	}
	s := grpc.NewServer() 
	pro.RegisterOrderServiceServer(s, &controller.RPCServer{})

	fmt.Println("Server listening on", constants.Port)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
