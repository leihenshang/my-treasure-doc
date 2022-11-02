package goods

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	goodsReq "fastduck/treasure-doc/service/mall/data/request/goods"
	goodsResp "fastduck/treasure-doc/service/mall/data/response/goods"
	"fastduck/treasure-doc/service/mall/global"
	goodsDao "fastduck/treasure-doc/service/mall/internal/dao/goods"
)

func GoodsList(ctx context.Context, f goodsReq.FilterGoodsList) (res *goodsResp.GoodsList, err error) {
	res = new(goodsResp.GoodsList)
	filter := &goodsDao.GoodsListFilter{
		GoodsName: f.GoodsName,
		DataSort:  f.DataSort,
	}
	pagination := &f.Pagination

	result, total, qErr := goodsDao.GoodsList(ctx, filter, pagination)
	if qErr != nil {
		global.ZAPSUGAR.Errorf("GoodsList|goodsDao.GoodsList failed to get goods list.filter:%+v,pagination:%+v,err:%+v", filter, pagination, qErr)
		err = errors.New("获取商品列表失败")
		return
	}

	for _, v := range result {
		res.List = append(res.List, &goodsResp.GoodsEntity{
			ID:        v.ID,
			Img:       v.Img,
			GoodsName: v.GoodsName,
			CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			DeletedAt: v.DeletedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	res.Total = total

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
