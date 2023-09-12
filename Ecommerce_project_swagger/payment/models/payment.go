package models

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
	Cardno          float32   `json:"cardno" bson:"cardno"`
	Brand           string  `json:"brand" bson:"brand"`
	PanLastFourNo   string  `json:"panlastfourno" bson:"panlastfourno"`
	ExpirationMonth int32   `json:"expirationmonth" bson:"expirationmonth"`
	ExpirationYear  int32   `json:"expirationyear" bson:"expirationyear"`
	Cvvverified     int32   `json:"cvvverified" bson:"cvvverified"`
	Balance         float32 `json:"balance" bson:"balance"`
}

