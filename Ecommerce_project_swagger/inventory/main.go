package main

// import (
// 	"context"
// 	"fmt"
// 	"inventory/config"
// 	"inventory/interfaces"
// 	"inventory/models"
// 	"inventory/services"

// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var InventoryService interfaces.Inventory

// func initApp(mclient *mongo.Client) {
// 	ctx := context.Background()
// 	mcoll := mclient.Database("inventory_SKU").Collection("inventory")
// 	iservice := services.InitInventory(mcoll, ctx)
// 	InventoryService = iservice
// }

// // func main() {
// // 	client, err := config.ConnectDataBase()
// // 	defer client.Disconnect(context.Background())
// // 	if err != nil {
// // 		fmt.Println(err.Error())
// // 	}
// // 	initApp(client)
// // 	fmt.Println("Connected to MongoDB!")
// // 	res := InventoryService.DeleteItems("shirts", "SKU001", 10.0)
// // 	fmt.Println(res)
// // }

// func main() {
// 	client, _ := config.ConnectDataBase()
// 	mcoll := client.Database("inventory_SKU").Collection("inventory")
// 	services.InitInventory(mcoll, context.Background())
// 	inventory := []*models.Inventory{
// 		{ID:         1,
// 			Item:       "shirt",
// 			Features:   []string{
// 				"slim fit",
// 				"full sleeves",
// 			},
// 			Categories: []string{
// 				"tshirt",
// 				"formal",
// 			},
// 		},
// 		{ID:         2,
// 			Item:       "pants",
// 			Features:   []string{
// 				"full length",
// 				"slim fit",
// 			},
// 			Categories: []string{
// 				"jeans",
// 				"formal",
// 			},
// 		},
// 		{ID:         3,
// 			Item:       "headphones",
// 			Features:   []string{
// 				"wireless",
// 				"noise cancellation",
// 			},
// 			Categories: []string{
// 				"headphones",
// 				"bluetooth",
// 			},
// 		},
// 	}
	// res, err:= services.CreateInventory(inventory)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(res)

	// res := inventory.DeleteItems("shirts", "SKU001", 10.0)
	// fmt.Println(res)

	// res, err := inventory.GetAllItems()
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// for _, item := range res {
	// 	fmt.Println(item.Item)
	// }

	// res, err :=services.GetInventoryItemByItemName("shirts")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(res)
// }
