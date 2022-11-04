package pay

type ParamsPayCreate struct {
	OrderId int32 `json:"orderId" form:"orderId"`
	UserId  int32 `json:"userId" form:"userId"`
}
