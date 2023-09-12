package controllers

import (
	"context"

	"payment/interfaces"
	"payment/models"
	pro "payment/proto"
)

type RPCServer struct {
	pro.UnimplementedPaymentServiceServer
}

var (
	TransactionService interfaces.Ipayment
)

func (t *RPCServer) CreatePayment(ctx context.Context, req *pro.PaymentDetails) (*pro.PaymentResponse, error) {
	transactions := &models.Paymentscard{
		Cardno:          float32(req.Cardno),
		Brand:           "",
		PanLastFourNo:   "",
		ExpirationMonth: 0,
		ExpirationYear:  0,
		Cvvverified:     int32(req.Cvvverified),
		Balance:         float32(req.Balance),
	}

	newProfile, err := TransactionService.CreatePayment(float32(transactions.Cardno), int(transactions.Cvvverified), transactions.Balance)

	if err != nil {

		return nil, err
	} else {
		responsePayment := &pro.PaymentResponse{
			Status: newProfile,
		}
		return responsePayment, nil
	}

}
