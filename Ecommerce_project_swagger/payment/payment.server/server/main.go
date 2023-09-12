package main

import (
	"context"
	"fmt"
	"net"

	"payment/config"
	"payment/constants"
	"payment/payment.server/controllers"
	pro "payment/proto"
	"payment/services"

	"google.golang.org/grpc"

	"go.mongodb.org/mongo-driver/mongo"
)

func initDatabase(client *mongo.Client) {

	CustomerCollection := config.GetCollection(client, "pay", "payments")
	transactionCollection := config.GetCollection(client, "pay", "transactions")

	controllers.TransactionService = services.NewTransactionServiceInit(client, CustomerCollection, transactionCollection, context.Background())
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

	pro.RegisterPaymentServiceServer(s, &controllers.RPCServer{})
	fmt.Println("Server listening on", constants.Port)
	if err := s.Serve(lis); err != nil {
		fmt.Printf("Failed to serve: %v", err)
	}
}
