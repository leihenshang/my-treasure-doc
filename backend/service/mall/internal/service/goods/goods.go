package goods

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	goodsReq "fastduck/treasure-doc/service/mall/request/goods"
)

func GoodsList(ctx context.Context, f goodsReq.FilterGoodsList) (res []*model.Good, total int64, err error) {
	q := query.Good.WithContext(ctx)
	if f.GoodsName != "" {
		q = q.Where(query.Good.GoodsName.Like("%" + f.GoodsName + "%"))
	}

	res, total, err = q.FindByPage(int(f.Offset), int(f.Limit))
	return
}

func GoodsDetail(ctx context.Context, f goodsReq.FilterGoodsDetail) (res *model.Good, err error) {
	if f.GoodsId == 0 {
		err = errors.New("商品id不能为空")
		return
	}

	q := query.Good.WithContext(ctx)
	res, err = q.Where(query.Good.ID.Eq(f.GoodsId)).First()
	return
}
