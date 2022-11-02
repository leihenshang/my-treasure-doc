package goods

type GoodsList struct {
	Total int64          `json:"total"`
	List  []*GoodsEntity `json:"list"`
}

type GoodsEntity struct {
	ID        int32  `json:"id"`
	Img       string `json:"img"`
	GoodsName string `json:"goodsName"`
	Quantity  int32  `json:"quantity"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	DeletedAt string `json:"deletedAt"`
}
