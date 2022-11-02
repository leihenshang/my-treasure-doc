package order

type OrderCreate struct {
	OrderId int32  `json:"orderId"`
	OrderNo string `json:"orderNo"`
}
