package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	orderReq "fastduck/treasure-doc/service/mall/request/order"
)

func OrderList(ctx context.Context, f orderReq.FilterOrderList) (res []*model.Order, total int64, err error) {
	q := query.Order.WithContext(ctx)

	res, total, err = q.FindByPage(int(f.Offset), int(f.Limit))
	return
}

func OrderDetail(ctx context.Context, f orderReq.FilterOrderDetail) (res *model.Order, err error) {
	if f.OrderId == 0 {
		err = errors.New("订单id不能为空")
		return
	}

	q := query.Order.WithContext(ctx)
	res, err = q.Where(query.Order.ID.Eq(f.OrderId)).First()
	return
}
