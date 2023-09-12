package services

import (
	"context"
	"fmt"
	"log"
	"payment/interfaces"
	"payment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var increment = 0

type TransactionService struct {
	client                *mongo.Client
	CustomerCollection    *mongo.Collection
	TransactionCollection *mongo.Collection
	ctx                   context.Context
}

func NewTransactionServiceInit(client *mongo.Client, Customercollection *mongo.Collection, TransactionCollection *mongo.Collection, ctx context.Context) interfaces.Ipayment {
	return &TransactionService{
		client:                client,
		CustomerCollection:    Customercollection,
		TransactionCollection: TransactionCollection,
		ctx:                   ctx,
	}

}

func (a *TransactionService) CreatePayment(cardno float32, cvvverified int, amount float32) (string, error) {

	increment++

	session, err := a.client.StartSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.EndSession(context.Background())

	// filter := bson.M{"Cardno": cardno}

	// var account *models.Paymentscard

	// err1 := a.CustomerCollection.FindOne(context.Background(), filter).Decode(&account)
	// if err1 != nil {
	//     return "error", err1
	// }
	// fmt.Println(cardno)
	filter := bson.M{"cardno": int32(cardno)} // Use float32 for the cardno value

	var account *models.Paymentscard

	err1 := a.CustomerCollection.FindOne(context.Background(), filter).Decode(&account)
	if err1 != nil {
		if err1 == mongo.ErrNoDocuments {
			// Handle the case when no documents match the filter
			fmt.Println("No documents found.")
			return "not found", nil
		}
		// Handle other errors
		fmt.Printf("Error: %v\n", err1)
		return "error", err1
	}

	// Handle the case when a document is found
	fmt.Printf("Account found: %+v\n", account)

	// fmt.Println("its going to in")

	if account.Balance >= amount {
		// fmt.Println("its in...")
		_, err := session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
			_, err := a.CustomerCollection.UpdateOne(context.Background(),
				bson.M{"cardno": int32(cardno)},
				bson.M{"$inc": bson.M{"balance": -amount}})
			if err != nil {
				fmt.Println("Transaction Failed")
				return nil, err
			}

			transactionToInsert := models.Payments{
				Id:          string(rune(increment)),
				Customer_id: "",
				Status:      "success",
				Gateway:     "",
				Type:        "",
				Amount:      amount,
				Card:        models.Paymentscard{},
				Token:       "",
			}

			res, err := a.TransactionCollection.InsertOne(context.Background(), &transactionToInsert)
			if err != nil {
				return nil, err
			}

			var newUser *models.Payments
			query := bson.M{"_id": res.InsertedID}
			err3 := a.TransactionCollection.FindOne(context.Background(), query).Decode(&newUser)
			if err3 != nil {
				return nil, err3
			}
			return newUser, nil
		})

		if err != nil {
			return "failed", err
		}
	} else {
		return "Insufficient balance", nil
	}
	return "success", nil
}
