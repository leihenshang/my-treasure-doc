package pay

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/mall/data/query"
	payReq "fastduck/treasure-doc/service/mall/data/request/pay"
	payResp "fastduck/treasure-doc/service/mall/data/response/pay"
	"fastduck/treasure-doc/service/mall/global"
	orderDao "fastduck/treasure-doc/service/mall/internal/dao/order"
)

func Create(ctx context.Context, params payReq.ParamsPayCreate) (res *payResp.PayCreate, err error) {
	res = new(payResp.PayCreate)
	if params.OrderId <= 0 {
		err = errors.New("订单Id不能为空")
		return
	}
	if params.UserId <= 0 {
		err = errors.New("用户Id不能为空")
		return
	}

	orderF := &orderDao.GetOrderFilter{
		OrderId: params.OrderId,
		UserId:  params.UserId,
	}
	order, orderErr := orderDao.GetOrder(ctx, orderF)
	if orderErr != nil {
		err = errors.New("查询订单失败")
		global.ZapSugar.Errorf("[pay|Create]failed get order info,err:%+v,filter:%+v", orderErr, orderF)
		return
	}
	if order.Status != 1 {
		err = errors.New("订单状态错误")
		global.ZapSugar.Errorf("[pay|Create]failed get order info,status err.orderId:%+v", order.ID)
		return
	}

	//TODO 更新
	_, updateErr := query.Order.WithContext(ctx).Where(query.Order.ID.Eq(params.OrderId)).UpdateColumn(query.Order.Status, 2)
	if updateErr != nil {
		err = errors.New("更新订单状态为已支付失败")
		global.ZapSugar.Errorf("[pay|Create]failed to update order info.orderId:%+v status:%+v", order.ID, 1)
		return
	}

	return
}
