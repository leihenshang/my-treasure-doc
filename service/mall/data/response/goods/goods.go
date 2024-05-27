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

type GoodsDetail struct {
	GoodsEntity
	// sku
	Sku []*GoodsSkuEntity
}

type GoodsSkuEntity struct {
	ID           int32       `json:"id"`
	GoodsID      int32       `json:"goodsId"`
	GoodsSpecIds string      `json:"goodsSpecIds"`
	GoodsSpec    []*SpecDesc `json:"goodsSpec"`
	Price        float64     `json:"price"`
	Stock        int32       `json:"stock"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt"`
	DeletedAt    string      `json:"deletedAt"`
}

type SpecDesc struct {
	SpecId  int32  `json:"specId"`
	Spec    string `json:"spec"`
	SpecVal string `json:"specVal"`
	Units   string `json:"units"`
}
