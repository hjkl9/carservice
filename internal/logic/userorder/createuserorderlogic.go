package userorder

import (
	"context"
	"fmt"
	"time"

	"carservice/internal/data/tables"
	"carservice/internal/enum/payment"
	"carservice/internal/enum/userorder"
	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/generator/order"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateUserOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserOrderLogic {
	return &CreateUserOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateUserOrderLogic) CreateUserOrder(req *types.CreateUserOrderReq) error {
	// ! 注意 user_id 和 member_id 都是用户的 ID
	userId := jwt.GetUserId(l.ctx)

	// 创建事务
	tx, err := l.svcCtx.DBC.Beginx()
	if err != nil {
		logc.Errorf(l.ctx, "创建订单->创建事务时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建事务时发生错误").SetDetails(err.Error())
	}

	// 创建车主信息
	query := "INSERT INTO `%s`(`name`, `user_id`, `phone_number`, `multilevel_address`, `full_address`, `longitude`, `latitude`) VALUES(?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.ExecContext(
		l.ctx,
		fmt.Sprintf(query, tables.CarOwnerInfo),
		req.CarOwnerName, userId,
		req.CarOwnerPhoneNumber,
		req.CarOwnerMultilevelAddress,
		req.CarOwnerFullAddress,
		req.CarOwnerLongitude,
		req.CarOwnerLatitude,
	)
	if err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("回滚失败").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->创建车主信息时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建用户车主信息时发生错误").SetDetails(err.Error())
	}
	// 新创建的车主信息 ID
	newCarOwnerInfoId, _ := result.LastInsertId()

	// 品牌和车系
	// 查询是否存在品牌和车系
	query = "SELECT count(1) AS `count` FROM `%s` AS `cb` JOIN `%s` AS `cbs` ON `cb`.`brand_id` = `cbs`.`brand_id` WHERE `cb`.`brand_id` = %d AND `cbs`.`series_id` = %d"
	var count int
	if err = tx.Get(
		&count,
		fmt.Sprintf(query, tables.CarBrand, tables.CarBrandSeries, req.CarBrandId, req.CarBrandSeriesId),
	); err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->查询车数据时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("数据库查询时发生错误").SetDetails(err.Error())
	}
	// 如果不存在
	if count == 0 {
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		return errcode.InvalidParamsError.SetMsg("无效的汽车品牌和系列")
	}
	// 安装需求
	// req.Requirements
	// 订单默认状态
	orderStatus := userorder.Pending
	// 订单默认支付方式
	paymentMethod := payment.Unknown
	// 生成订单号
	orderNumber := order.GenerateNumber(time.Now())
	// 创建时间和更新时间
	nowUnix := time.Now()
	var createdAt = nowUnix
	var updatedAt = nowUnix
	// 创建订单
	query = "INSERT INTO `%s`(`member_id`, `car_brand_id`, `car_brand_series_id`, `car_info_id`, `car_owner_info_id`, `partner_store_id`, `order_number`, `order_status`, `comment`, `est_amount`, `act_amount`, `payment_method`, `created_at`, `updated_at`) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	rs, err := tx.ExecContext(
		l.ctx,                                // 上下文
		fmt.Sprintf(query, tables.UserOrder), // 创建语句
		userId,                               // 用户
		req.CarBrandId,                       // 汽车品牌
		req.CarBrandSeriesId,                 // 汽车品牌系列
		0,                                    // ! car_info_id 应该删除该字段
		newCarOwnerInfoId,                    // 用户车主信息
		req.PartnerStoreId,                   // ! partner_store_id 合作门店 ID
		orderNumber,                          // 订单号
		orderStatus,                          // 订单状态
		req.Requirements,                     // 需求
		0.00,                                 // ! unused 默认的预估服务费
		0.00,                                 // ! unused 默认的实际服务费
		paymentMethod,                        // 默认的支付方式
		createdAt,                            // 创建时间
		updatedAt,                            // 更新时间
	)
	if err != nil {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单时发生错误, err: %s\n", err.Error())
		return errcode.DatabaseError.SetMsg("创建订单时发生错误").SetDetails(err.Error())
	}

	rows, _ := rs.RowsAffected()
	if rows == 0 {
		// 回滚
		if err1 := tx.Rollback(); err1 != nil {
			return errcode.DatabaseError.SetMsg("数据库回滚时发生错误").SetDetails(err1.Error())
		}
		logc.Errorf(l.ctx, "创建订单->创建订单失败, err: %s\n", err.Error())
		return errcode.InternalServerError.SetMsg("订单创建失败")
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return errcode.DatabaseError.SetMsg("提交事务事务时发生错误").SetDetails(err.Error())
	}

	return nil
}
