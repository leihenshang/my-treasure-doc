package goods

import (
	"context"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	goodsReq "fastduck/treasure-doc/service/mall/request/goods"
)

func GoodsList(ctx context.Context, f goodsReq.FilterGoodsList) (res []*model.Good, total int64, err error) {
	q := query.Good
	if f.GoodsName != "" {
		q.Where(query.Good.GoodsName.Like(f.GoodsName))
	}

	res, total, err = q.FindByPage(int(f.Offset), int(f.Limit))
	return
}

func GoodsDetail(f goodsReq.FilterGoodsDetail) (res *model.Good, err error) {

	return
}
