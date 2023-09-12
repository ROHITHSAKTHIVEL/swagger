package controller

import (
	"context"
	"ecommerce_order/order_dal/interfaces"
	"ecommerce_order/order_dal/models"

	pro "ecommerce_order/order_proto"
	"fmt"
)

type RPCServer struct {
	pro.UnimplementedOrderServiceServer
}

var (
	OrderService interfaces.IOrder
)

func (s *RPCServer) RemoveOrderCustomer(ctx context.Context, req *pro.RemoveOrderRequest) (*pro.RemoveOrderResponse, error) {
	// fmt.Println(req.CustomerId)
	result, err := OrderService.RemoveOrder(req.CustomerId)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	responseCustomer := &pro.RemoveOrderResponse{
		Status: result,
	}

	return responseCustomer, nil
}

func (s *RPCServer) CreateOrder(ctx context.Context, req *pro.CustomerOrder) (*pro.CustomerResponse, error) {
	dbInsert := &models.Orders{
		CustomerId:    req.CustomerId,
		PaymentId:     req.PaymentId,
		PaymentStatus: req.PaymentStatus,
		Status:        req.Status,
		Currency:      req.Currency,

		Carrier:  req.Carrier,
		Tracking: req.Tracking,
	}

	// var protoItems []models.Items

	// Iterate over the source data (protoItems) and append to the 'items' slice
	for _, protoItem := range req.Items {
		item := models.Items{
			Sku:      protoItem.Sku,
			Quantity: protoItem.Quantity,

			// Discount:    protoItem.Discount,
			// PreTaxTotal: protoItem.PreTaxTotal,
			// Tax:         protoItem.Tax,
			// Total:       protoItem.Total,
		}

		// Append the 'item' to the 'items' slice
		dbInsert.Items = append(dbInsert.Items, item)
	}

	for _, protoShipping := range req.Shipping {
		shipping := models.Shipping{}
		for _, protoAddress := range protoShipping.Address {
			address := models.Address{
				Street1: protoAddress.Street1,
				Street2: protoAddress.Street2,
				City:    protoAddress.City,
				State:   protoAddress.State,
				Country: protoAddress.Country,
				Zip:     protoAddress.Zip,
			}
			shipping.Address = append(shipping.Address, address)
		}
		for _, protoOrigin := range protoShipping.Origin {
			origin := models.Origin{
				Street1: protoOrigin.Street1,
				Street2: protoOrigin.Street2,
				City:    protoOrigin.City,
				State:   protoOrigin.State,
				Country: protoOrigin.Country,
				Zip:     protoOrigin.Zip,
			}
			shipping.Origin = append(shipping.Origin, origin)
		}
		dbInsert.Shipping = append(dbInsert.Shipping, shipping)
	}

	value, err := OrderService.CreateOrder(dbInsert)
	response := &pro.CustomerResponse{
		CustomerId: value.CustomerId,
	}

	if err != nil {
		return nil, err

	}
	return response, nil
}
func (s *RPCServer) GetOrderDetails(ctx context.Context, req *pro.GetOrderRequest) (*pro.CustomerOrder, error) {
	customerID := req.CustomerId
	result, err := OrderService.GetAllOrder(customerID)
	if err != nil {
		return nil, err
	}
	responseCustomer := &pro.CustomerOrder{
		CustomerId:    result.CustomerId,
		PaymentId:     result.PaymentId,
		PaymentStatus: result.PaymentStatus,
		Status:        result.Status,
		Currency:      result.Currency,
		TotalAmount:   result.TotalAmount,
		Carrier:       result.Carrier,
		Tracking:      result.Tracking,
	}
	if len(result.Items) > 0 {
		responseCustomer.Items = []*pro.Items{
			{
				Sku:         result.Items[0].Sku,
				Quantity:    result.Items[0].Quantity,
				Price:       result.Items[0].Price,
				Discount:    result.Items[0].Discount,
				PreTaxTotal: result.Items[0].PreTaxTotal,
				Tax:         result.Items[0].Tax,
				Total:       result.Items[0].Total,
			},
		}
	}
	if len(result.Shipping) > 0 {
		responseCustomer.Shipping = make([]*pro.Shipping, len(result.Shipping))
		for i, shipping := range result.Shipping {
			shippingMsg := &pro.Shipping{}
			shippingMsg.Address = make([]*pro.Address, len(shipping.Address))
			for j, address := range shipping.Address {
				addressMsg := &pro.Address{
					Street1: address.Street1,
					Street2: address.Street2,
					City:    address.City,
					State:   address.State,
					Country: address.Country,
					Zip:     address.Zip,
				}
				shippingMsg.Address[j] = addressMsg
			}
			shippingMsg.Origin = make([]*pro.Origin, len(shipping.Origin))
			for j, origin := range shipping.Origin {
				originMsg := &pro.Origin{
					Street1: origin.Street1,
					Street2: origin.Street2,
					City:    origin.City,
					State:   origin.State,
					Country: origin.Country,
					Zip:     origin.Zip,
				}
				shippingMsg.Origin[j] = originMsg
			}
			responseCustomer.Shipping[i] = shippingMsg
		}
	}
	return responseCustomer, nil
}

func (s *RPCServer) UpdateOrderDetails(ctx context.Context, req *pro.UpdateOrderRequest) (*pro.UpdateOrderResponse, error) {
	dbInsert := &models.UpdateDetailsModel{
		Customer_ID: req.Customer_ID,
		Sku:         req.Sku,
		Quantity:    req.Quantity,
	}

	value, err := OrderService.UpdateOrder(dbInsert)
	response := &pro.UpdateOrderResponse{
		Status: value,
	}

	if err != nil {
		return nil, err

	}
	return response, nil
}

func (s *RPCServer) AddOrderDetails(ctx context.Context, req *pro.UpdateOrderRequest) (*pro.UpdateOrderResponse, error) {
	dbInsert := &models.UpdateDetailsModel{
		Customer_ID: req.Customer_ID,
		Sku:         req.Sku,
		Quantity:    req.Quantity,
	}
	value, err := OrderService.AddOrder(dbInsert)
	response := &pro.UpdateOrderResponse{
		Status: value,
	}

	if err != nil {
		return nil, err

	}
	return response, nil
}
