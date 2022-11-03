package order

type OrderCreate struct {
	OrderId int32  `json:"orderId"`
	OrderNo string `json:"orderNo"`
}

type OrderList struct {
	Total int64          `json:"total"`
	List  []*OrderEntity `json:"list"`
}

type OrderEntity struct {
	ID        int32   `json:"id"`
	OrderNo   string  `json:"orderNo"`
	UserID    int32   `json:"userId"`
	Amount    float64 `json:"amount"`
	Status    int32   `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	DeletedAt string  `json:"deletedAt"`
}
