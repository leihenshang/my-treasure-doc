package goods

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	goodsReq "fastduck/treasure-doc/service/mall/data/request/goods"
	goodsResp "fastduck/treasure-doc/service/mall/data/response/goods"
)

func GoodsList(ctx context.Context, f goodsReq.FilterGoodsList) (res *goodsResp.GoodsList, err error) {
	res = new(goodsResp.GoodsList)
	q := query.Good.WithContext(ctx)
	if f.GoodsName != "" {
		q = q.Where(query.Good.GoodsName.Like("%" + f.GoodsName + "%"))
	}
	if f.SortField != "" {
		goodsCol, ok := query.Good.GetFieldByName(f.SortField)
		if ok {
			if f.IsDesc {
				q = q.Order(goodsCol.Desc())
			} else {
				q = q.Order(goodsCol)
			}
		}
	}

	result, total, err := q.FindByPage(int(f.Offset), int(f.Limit))

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
