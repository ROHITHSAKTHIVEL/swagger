package main

import (
	"context"
	"fmt"
	_ "inventory/client/docs"
	h "inventory/proto"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type Inventory struct {
	ID         int64           `json:"id" bson:"_id"`
	Item       string          `json:"item" bson:"item"`
	Features   []string        `json:"features" bson:"features"`
	Categories []string        `json:"categories" bson:"categories"`
	Skus       []Inventory_SKU `json:"skus" bson:"skus"`
}
type ItemToUpdate struct{
	 Item string `json:"item" bson:"item"`
	 Sku string `json:"sku" bson:"sku"`
	 Quantity float64 `json:"quantity" bson:"quantity"` 
}
	
type InventoryResponse struct {
	ID         int64           `json:"id" bson:"_id"`
	Item       string          `json:"item" bson:"item"`
	Features   []string        `json:"features" bson:"features"`
	Categories []string        `json:"categories" bson:"categories"`
	Skus       []Inventory_SKU `json:"skus" bson:"skus"`
}
type Inventory_SKU struct {
	// Id 	 primitive.ObjectID `json:"_id" bson:"_id"`
	Sku      string       `json:"sku" bson:"sku"`
	Price    Price_type   `json:"price" bson:"price"`
	Quantity float32      `json:"quantity" bson:"quantity,truncate"`
	Options  Options_type `json:"options" bson:"options"`
}
type UpdateResponse struct{
	Message string `json:"message" bson:"message"`
}
type CreateResponse struct{
	Id primitive.ObjectID `json:"id"`
}
type Price_type struct {
	Base     float32 `json:"base" bson:"base,truncate"`
	Currency string  `json:"currency" bson:"currency"`
	Discount float32 `json:"discount" bson:"discount, truncate"`
}
type Options_type struct {
	Size     Size_type `json:"size" bson:"size"`
	Features []string  `json:"features" bson:"features"`
	Colors   []string  `json:"colors" bson:"colors"`
	Ruling   string    `json:"ruling" bson:"ruling"`
	Image    string    `json:"image" bson:"image"`
}
type GetInventoryItemByItemName struct {
    ItemName string `json:"item_name" bson:"item_name"`
}

type Size_type struct {
	H float32 `json:"h" bson:"h,truncate"`
	L float32 `json:"l" bson:"l, truncate"`
	W float32 `json:"w" bson:"w,truncate"`
}

// @title Documenting API (INVENTORY)
// @version 1
// @Description Buy Anything In Our Webstite

// @contact.name INVENTORY
// @contact.url http://localhost:8080/swagger/index.html#/Payment/post_payment_createpayment
// @contact.name INVENTORY
// @contact.url http://localhost:8080/swagger/index.html#/Payment/post_payment_createpayment
// @contact.email rohith.s@netxd.com
// @contact.phone 1234567787



// @host localhost:8083
// @BasePath /api/v1
func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	user := v1.Group("/inventory")
	{
		user.GET("/getitems", GetAllItems)
		user.POST("/getitem", GetItemsByID)
		user.POST("/create", CreateItems)
		user.POST("/update", UpdateItems)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.Run(":8083")
	}
}

// Getitems
// @Summary getitems 
// @Description get all items
// @Tags Inventory
// @Success 200 {object} Inventory
// @Router /inventory/getitems [get]
func GetAllItems(c *gin.Context) {
	conn, err := grpc.Dial("localhost:2023", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect: ", err)
	}
	defer conn.Close()
	client := h.NewInventoryServiceClient(conn)
	response, err := client.GetAllItems(context.Background(), &h.Empty{})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println("Response: ", response)
	c.JSON(http.StatusOK, gin.H{"value": response})
}

// Getitems By id
// @Summary getitems by id
// @Description getitemsbyid
// @Tags Inventory
// @Accept json
// @Produce json
// @Param user body GetInventoryItemByItemName true "user"
// @Success 200 {object} Inventory
// @Router /inventory/getitem [post]
func GetItemsByID(c *gin.Context) {
	conn, err := grpc.Dial("localhost:2023", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect: ", err)
	}
	defer conn.Close()
	client := h.NewInventoryServiceClient(conn)
	var request h.ItemName
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := client.GetInventoryItemByItemName(context.Background(), &h.ItemName{
		ItemName: request.ItemName,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println("Response: ", response)
	c.JSON(http.StatusOK, gin.H{"value": response})
}

// CreateItems
// @Summary Create and return it
// @Description CreateItems
// @Tags Inventory
// @Accept json
// @Produce json
// @Param user body []Inventory true "user"
// @Success 200 {object} CreateResponse
// @Router /inventory/create [post]
func CreateItems(c *gin.Context) {
	conn, err := grpc.Dial("localhost:2023", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect: ", err)
	}
	defer conn.Close()
	client := h.NewInventoryServiceClient(conn)
	var request []*h.InventoryItem
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.CreateInventory(context.Background(), &h.AllInventoryItems{
		Items: request,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println("Response: ", response)
	c.JSON(http.StatusOK, gin.H{"value": response})
}
// UpdateItems
// @Summary Update Items in Inventory
// @Description UpdateItems
// @Tags Inventory
// @Accept json
// @Produce json
// @Param user body ItemToUpdate true "user"
// @Success 200 {object} UpdateResponse
// @Router /inventory/update [post]
func UpdateItems(c *gin.Context) {
	conn, err := grpc.Dial("localhost:2023", grpc.WithInsecure())
	if err != nil {
		fmt.Println("Failed to connect: ", err)
	}
	defer conn.Close()
	client := h.NewInventoryServiceClient(conn)
	var request h.ItemToDelete
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.UpdateInventory(context.Background(), &h.ItemToDelete{
		Item:     request.Item,
		Sku:      request.Sku,
		Quantity: request.Quantity,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println(response.Msg)
	Response := UpdateResponse{Message: response.Msg}
	c.JSON(http.StatusOK, gin.H{"value": Response})
}