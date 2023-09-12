package services

import (
	"context"
	"ecommerce_order/order_dal/interfaces"
	"ecommerce_order/order_dal/models"
	"fmt"

	// "strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerService struct {
	client              *mongo.Client
	OrderCollection     *mongo.Collection
	InventoryCollection *mongo.Collection
	ctx                 context.Context
}

func InitCustomerService(client *mongo.Client, collection1 *mongo.Collection, collection2 *mongo.Collection, ctx context.Context) interfaces.IOrder {
	return &CustomerService{client, collection1, collection2, ctx}
}

func (p *CustomerService) CreateOrder(input *models.Orders) (models.Orders, error) {
	var basePrice float32
	var discountvalue float32
	var total float32
	for i, val := range input.Items {
		filter := bson.M{"sku": val.Sku}
		inventoryResult := p.InventoryCollection.FindOne(p.ctx, filter)
		var inventoryDocument bson.M
		if err := inventoryResult.Decode(&inventoryDocument); err != nil {
			return models.Orders{}, err
		}
		price := inventoryDocument["price"].(bson.M)
		quantity := inventoryDocument["quantity"].(float64)
		if quantity < float64(val.Quantity) {
			fmt.Println("Lack of Quantity")
			return models.Orders{}, nil
		}
		price64 := price["base"].(float64)
		discount := price["discount"].(float64)
		basePrice = float32(price64)
		discountvalue = float32(discount)
		fmt.Println(price64, discount, basePrice)
		input.Items[i].Price = val.Quantity * basePrice
		input.Items[i].Discount = discountvalue
		input.Items[i].PreTaxTotal = input.Items[i].Price - ((discountvalue / 100) * 100)
		total = total + input.Items[i].PreTaxTotal
		input.Items[i].Total = input.Items[i].PreTaxTotal
		fmt.Println(val.Price, val.Discount)
	}
	input.TotalAmount = total
	insertResult, err1 := p.OrderCollection.InsertOne(p.ctx, &input)
	if err1 != nil {

	}
	insertedID := insertResult.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": insertedID}
	var result models.Orders
	err2 := p.OrderCollection.FindOne(p.ctx, filter).Decode(&result)
	if err2 != nil {
		return result, err2
	}
	fmt.Println("Success")
	return result, nil
}
func (p *CustomerService) RemoveOrder(Customer_ID string) (string, error) {
	filter := bson.M{"customerid": Customer_ID}
	_, err := p.OrderCollection.DeleteOne(p.ctx, filter)
	if err != nil {
		return "Unable to delete", err
	}
	return "Deleted Successfully", nil
}
func (p *CustomerService) GetAllOrder(CustomerId string) (*models.Orders, error) {
	filter := bson.M{"customerid": CustomerId}
	var res *models.Orders
	result := p.OrderCollection.FindOne(p.ctx, filter)
	err := result.Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (p *CustomerService) UpdateOrder(input *models.UpdateDetailsModel) (string, error) {
	var basePrice float32
	var discountvalue float32
	var total float32
	filter := bson.M{"sku": input.Sku}
	inventoryResult := p.InventoryCollection.FindOne(p.ctx, filter)
	var inventoryDocument bson.M
	if err := inventoryResult.Decode(&inventoryDocument); err != nil {
		return "Empty", err
	}
	price := inventoryDocument["price"].(bson.M)
	quantity := inventoryDocument["quantity"].(float64)
	if quantity < float64(input.Quantity) {
		fmt.Println("Lack of Quantity")
		return "Empty", nil
	}
	price64 := price["base"].(float64)
	discount := price["discount"].(float64)
	discountvalue = float32(discount)
	basePrice = float32(price64)
	fmt.Println("The Price is ", basePrice)
	input.Price = input.Quantity * basePrice
	input.Discount = discountvalue
	input.PreTaxTotal = input.Price - ((discountvalue / 100) * 100)
	total = input.PreTaxTotal
	input.TotalAmount = total
	input.Total = input.PreTaxTotal
	filter1 := bson.M{"customerid": input.Customer_ID}
	update := bson.M{
		"$set": bson.M{
			"totalamount": models.Orders{
				TotalAmount: input.TotalAmount,
			},
			"items": []models.Items{
				{
					Sku:         input.Sku,
					Quantity:    input.Quantity,
					Price:       input.Price,
					Discount:    input.Discount,
					PreTaxTotal: total,
					Total:       input.Total,
				},
			},
		},
	}
	_, err1 := p.OrderCollection.UpdateOne(context.Background(), filter1, update)
	if err1 != nil {
		return "updation failed", err1
	}
	return "Updation Success", nil

}

func (p *CustomerService) AddOrder(input *models.UpdateDetailsModel) (string, error) {
	var basePrice float32
	var discountvalue float32
	var total float32
	filter := bson.M{"sku": input.Sku}
	inventoryResult := p.InventoryCollection.FindOne(p.ctx, filter)
	var inventoryDocument bson.M
	if err := inventoryResult.Decode(&inventoryDocument); err != nil {
		return "Empty", err
	}
	price := inventoryDocument["price"].(bson.M)
	quantity := inventoryDocument["quantity"].(float64)
	if quantity < float64(input.Quantity) {
		fmt.Println("Lack of Quantity")
		return "Empty", nil
	}
	price64 := price["base"].(float64)
	discount := price["discount"].(float64)
	discountvalue = float32(discount)
	basePrice = float32(price64)
	fmt.Println("The Price is ", basePrice)
	fmt.Println("The quantity is ", input.Quantity)
	fmt.Println("The two values are ", input.Quantity, " ", basePrice)
	input.Price = input.Quantity * basePrice
	input.Discount = discountvalue
	input.PreTaxTotal = input.Price - ((discountvalue / 100) * 100)
	total = total + input.PreTaxTotal + input.TotalAmount
	input.TotalAmount = total
	input.Total = input.PreTaxTotal
	itemsToAdd := []models.Items{
		{
			Sku:      input.Sku,
			Quantity: input.Quantity,
			Price:    input.Price,
			Discount: input.Discount,
			Total:    input.Total,
		},
	}
	filter1 := bson.M{"customerid": input.Customer_ID}
	update := bson.M{

		"$push": bson.M{
			"items": bson.M{
				"$each": itemsToAdd,
			},
		},
	}
	_, err1 := p.OrderCollection.UpdateOne(context.Background(), filter1, update)
	if err1 != nil {
		return "updation failed", err1
	}
	return "Updation Success", nil

}