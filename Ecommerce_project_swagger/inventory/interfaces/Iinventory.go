package interfaces

import (
	"inventory/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type Inventory interface {
	CreateInventory(in []*models.Inventory) (*mongo.InsertManyResult, error)
	DeleteItems(item string, sku string, quantity float32) (string)
	GetAllItems() ([]models.Inventory, error) 
	GetInventoryItemByItemName(itemName string) (*models.Inventory, error)
}