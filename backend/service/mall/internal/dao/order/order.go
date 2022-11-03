package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	reqCommon "fastduck/treasure-doc/service/mall/data/request/common"
)

type OrderListFilter struct {
	UserId int32
	reqCommon.DataSort
}

func OrderList(ctx context.Context, f *OrderListFilter, p *reqCommon.Pagination) (result []*model.Order, total int64, err error) {
	if f == nil {
		err = errors.New("过滤器设置错误,不能为nil")
		return
	}

	if p == nil {
		p = reqCommon.NewPagination()
	}

	q := query.Order.WithContext(ctx)

	if f.UserId > 0 {
		q = q.Where(query.Order.UserID.Eq(f.UserId))
	}

	if f.SortField != "" {
		orderCol, ok := query.Order.GetFieldByName(f.SortField)
		if ok {
			if f.IsDesc {
				q = q.Order(orderCol.Desc())
			} else {
				q = q.Order(orderCol)
			}
		}
	}

	result, total, err = q.FindByPage(int(p.Offset), int(p.Limit))
	return
}
