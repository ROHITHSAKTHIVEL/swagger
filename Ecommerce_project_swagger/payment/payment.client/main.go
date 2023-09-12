package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "payment/payment.client/docs"
	pb "payment/proto"
	model "payment/models"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var (
	mongoclient *mongo.Client
	ctx         context.Context
	server      *gin.Engine
)

type Payments struct {
	Id          string       `json:"id" bson:"id"`
	Customer_id string       `json:"customer_id" bson:"customer_id"`
	Status      string       `json:"status" bson:"status"`
	Gateway     string       `json:"gateway" bson:"gateway"`
	Type        string       `json:"type" bson:"type"`
	Amount      float32      `json:"amount" bson:"amount"`
	Card        Paymentscard `json:"card" bson:"card"`
	Token       string       `json:"token" bson:"token"`
}

type Paymentscard struct {
	Cardno          float32 `json:"cardno" bson:"cardno"`
	Brand           string  `json:"brand" bson:"brand"`
	PanLastFourNo   string  `json:"panlastfourno" bson:"panlastfourno"`
	ExpirationMonth int32   `json:"expirationmonth" bson:"expirationmonth"`
	ExpirationYear  int32   `json:"expirationyear" bson:"expirationyear"`
	Cvvverified     int32   `json:"cvvverified" bson:"cvvverified"`
	Balance         float32 `json:"balance" bson:"balance"`
}


type PaymentResponse struct {
	Message string `json:"message"`
}

// @title Documenting API (PAYMENT)

// @Description Buy Anything In Our Webstite
// @contact.name NetXD
// @contact.url https://netxd.com/
// @contact.email rohith.s@netxd.com
// @contact.phone 1234567787

// @securityDefinitions.apikey bearer
// @in header
// @name Authorization

// @host localhost:8081
// @BasePath /api/v1
func main() {

	r := gin.Default()
m:=model.Paymentscard{Cardno : 12345}
fmt.Println(m)
	v1 := r.Group("/api/v1")
	user := v1.Group("/payment")
	{
		user.POST("/createpayment", CreatePayment)
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		r.Run(":8081")
	}

}
// CreatePayment
// @Summary return created payment
// @Description createpayment
// @Tags Payment
// @Accept json
// @Produce json
// @Param user body model.Paymentscard true "user"
// @Success 200 {object} PaymentResponse
// @Router /payment/createpayment [post]
func CreatePayment(c *gin.Context) {
	conn, err := grpc.Dial("localhost:5002", grpc.WithInsecure())
	if err != nil {
		fmt.Println("111")
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)

	var request pb.PaymentDetails
	if err := c.ShouldBindJSON(&request); err != nil {
		fmt.Println("12")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	response, err := client.CreatePayment(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(response.Status)
	Response := PaymentResponse{Message: response.Status}
	c.JSON(http.StatusOK, gin.H{"value": Response})
}


