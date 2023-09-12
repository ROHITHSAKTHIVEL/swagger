package main

import (
	"context"
	"fmt"

	// "fmt"

	"log"
	"net/http"
   _ "ecommerce_order/order_client/client/docs"
	pb "ecommerce_order/order_proto"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)
type Orders struct {
	CustomerId    string     `json:"customerid" bson:"customerid"`
	PaymentId     string     `json:"payment_id" bson:"payment_id" `
	PaymentStatus string     `json:"paymentstatus" bson:"paymentstatus"`
	Status        string     `json:"status" bson:"status"`
	Currency      string     `json:"currency" bson:"currency"`
	Items         []Items    `json:"items" bson:"items"`
	TotalAmount   float32    `json:"totalamount" bson:"totalamount"`
	Shipping      []Shipping `json:"shipping" bson:"shipping"`
	Carrier       string     `json:"carrier" bson:"carrier"`
	Tracking      string     `json:"tracking" bson:"tracking"`
}

type Shipping struct {
	Address []Address `json:"address" bson:"address"`
	Origin  []Origin  `json:"origin" bson:"origin"`
}
type Items struct {
	Sku         string  `json:"sku" bson:"sku"`
	Quantity    float32 `json:"quantity" bson:"quantity"`
	Price       float32 `json:"price" bson:"price"`
	Discount    float32 `json:"discount" bson:"discount"`
	PreTaxTotal float32 `json:"pretaxtotal" bson:"pretaxtotal"`
	Tax         float32 `json:"tax" bson:"tax"`
	Total       float32 `json:"total" bson:"total"`
}

type Address struct {
	Street1 string `json:"street1" bson:"street1"`
	Street2 string `json:"street2" bson:"street2"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
	Zip     string `json:"zip" bson:"zip"`
}
type Origin struct {
	Street1 string `json:"street1" bson:"street1"`
	Street2 string `json:"street2" bson:"street2"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
	Zip     string `json:"zip" bson:"zip"`
}
type UpdateDetails struct {
	Sku         string  `json:"sku" bson:"sku"`
	Quantity    float32 `json:"quantity" bson:"quantity"`
}
type UpdateOrderResponse struct{
	Status string `json:"status" bson:"status"`
}
type Inventory_SKU struct {
	Sku      string       `json:"sku" bson:"sku"`
	Price    Price_type   `json:"price" bson:"price"`
	Quantity float64      `json:"quantity" bson:"quantity"`
	Options  Options_type `json:"options" bson:"options"`
}

type Price_type struct {
	Base     float64 `json:"base" bson:"base"`
	Currency string  `json:"currency" bson:"currency"`
	Discount float64 `json:"discount" bson:"discount"`
}
type Options_type struct {
	Size     Size_type `json:"size" bson:"size"`
	Features []string  `json:"features" bson:"features"`
	Colors   []string  `json:"colors" bson:"colors"`
	Ruling   string    `json:"ruling" bson:"ruling"`
	Image    string    `json:"image" bson:"image"`
}

type Size_type struct {
	H float64 `json:"h" bson:"h"`
	L float64 `json:"l" bson:"l"`
	W float64 `json:"w" bson:"w"`
}


var (
	mongoclient *mongo.Client
	ctx         context.Context
	server      *gin.Engine
)

// @title Documenting API (ORDER)
// @version 1
// @Description Buy Anything In Our Webstite

// @contact.name Order
// @contact.url http://localhost:8081/swagger/index.html
// @contact.name INVENTORY
// @contact.url http://localhost:8083/swagger/index.html#/
// @contact.email rohith.s@netxd.com
// @contact.phone 1234567787

// @securityDefinitions.apikey bearer
// @in header
// @name Authorization

// @host localhost:8089
// @BasePath /api/v1
func main() {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	user := v1.Group("/order")
	{
	user.POST("/createorder", CreateOrder)
	user.POST("/updateorder", UpdateOrder)
	user.POST("/addorderdetails", AddOrderDetails)
	user.GET("/deletecustomer", DeleteCustomer)
	user.GET("/getbyid", GetbyId)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8089")
	}
}
// CreateOrder create new order
// @Summary return created order
// @Description createorder
// @Tags Orders
// @Accept json
// @Produce json
// @Param user body Orders true "user"
// @Success 200 {object} Orders
// @Router /order/createorder [post]
func CreateOrder(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	var request pb.CustomerOrder
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.CreateOrder(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": response})
}
// UpdateOrder Update the created order
// @Summary return updated order
// @Description UpdateOrder
// @Security bearerToken
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user body UpdateDetails true "user"
// @Success 200 {object} UpdateOrderResponse
// @Router /order/updateorder [post]
func UpdateOrder(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	var request pb.UpdateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.UpdateOrderDetails(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	Response := UpdateOrderResponse{Status: response.Status}
	c.JSON(http.StatusOK, gin.H{"value": Response})
}
// AddordersDetails 
// @Summary return Added Details
// @Description AddOrderDetails
// @Security bearerToken
// @Tags Orders
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param user body UpdateDetails true "user"
// @Success 200 {object} UpdateOrderResponse
// @Router /order/addorderdetails [post]
func AddOrderDetails(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	var request pb.UpdateOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.UpdateOrderDetails(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"value": response})
}
// DeleteCustomer Delete Customer order
// @Summary return Deleted Meassage
// @Description DeleteCustomer Order
// @Security bearerToken
// @Tags Orders
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} UpdateOrderResponse
// @Router /order/deletecustomer [get]
func DeleteCustomer(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	var user pb.RemoveOrderRequest

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	client.RemoveOrderCustomer(c.Request.Context(), &pb.RemoveOrderRequest{CustomerId: user.CustomerId})

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})

}
// Getby Id
// @Summary return GetId
// @Description GetById
// @Security bearerToken
// @Tags Orders
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} Orders
// @Router /order/getbyid [get]
func GetbyId(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrderServiceClient(conn)
	var res pb.GetOrderRequest
	if err := c.ShouldBindJSON(&res); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(res.CustomerId)
	result, err := client.GetOrderDetails(c.Request.Context(), &pb.GetOrderRequest{CustomerId: res.CustomerId})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": result})
}
