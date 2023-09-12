package interfaces


type Ipayment interface{
	CreatePayment( float32, int,float32)(string, error)
	
}