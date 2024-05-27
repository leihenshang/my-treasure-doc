package goods

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	reqCommon "fastduck/treasure-doc/service/mall/data/request/common"
	"fmt"
)

type GetGoodsFilter struct {
	GoodId  int32
	Enabled int32
}

func GetGoods(ctx context.Context, f *GetGoodsFilter) (res *model.Good, err error) {
	q := query.Good.WithContext(ctx)
	if f.GoodId > 0 {
		q = q.Where(query.Good.ID.Eq(f.GoodId))
	}
	if f.Enabled > 0 {
		q = q.Where(query.Good.Enabled.Eq(f.Enabled))
	}

	res, err = q.First()

	return
}

type GoodsListFilter struct {
	GoodsName string
	reqCommon.DataSort
}

// 获取商品列表
func GoodsList(
	ctx context.Context,
	f *GoodsListFilter,
	p *reqCommon.Pagination,
) (result []*model.Good, total int64, err error) {
	if f == nil {
		err = errors.New("过滤器设置错误,不能为nil")
		return
	}

	if p == nil {
		p = reqCommon.NewPagination()
	}

	q := query.Good.WithContext(ctx)
	if f.GoodsName != "" {
		q = q.Where(query.Good.GoodsName.Like(fmt.Sprintf("%%%s%%", f.GoodsName)))
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

	result, total, err = q.FindByPage(int(p.Offset), int(p.Limit))
	return
}
