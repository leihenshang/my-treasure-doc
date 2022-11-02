package order

import (
	"context"
	"errors"
	"fastduck/treasure-doc/service/admin/global"
	"fastduck/treasure-doc/service/mall/data/model"
	"fastduck/treasure-doc/service/mall/data/query"
	orderReq "fastduck/treasure-doc/service/mall/data/request/order"
	orderResp "fastduck/treasure-doc/service/mall/data/response/order"
	goodsDao "fastduck/treasure-doc/service/mall/internal/dao/goods"
	"fmt"
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

func OrderCreate(ctx context.Context, f orderReq.FilterOrderCreate) (res *orderResp.OrderCreate, err error) {
	if f.Quantity <= 0 {
		err = errors.New("数量不能小于1")
		return
	}
	if f.SkuId <= 0 {
		err = errors.New("skuId 不能小于1")
		return
	}

	if f.UserId <= 0 {
		err = errors.New("userId 不能小于1")
		return
	}

	//查询商品 sku 以及库存
	//扣减库存
	//生成订单
	skuF := &goodsDao.GetSkuFilter{
		SkuId:   f.SkuId,
		Enabled: 1,
	}
	sku, skuErr := goodsDao.GetSku(ctx, skuF)
	if skuErr != nil {
		global.ZAPSUGAR.Errorf("OrderCreate failed to get sku info. params:%+v err:%+v", skuF, skuErr)
		err = errors.New("获取sku信息失败")
		return
	}
	if sku == nil {
		global.ZAPSUGAR.Error("OrderCreate failed to get sku info. err:sku not existed")
		err = errors.New(fmt.Sprintf("sku不存在,skuId:%d", f.SkuId))
		return
	}

	return
}
