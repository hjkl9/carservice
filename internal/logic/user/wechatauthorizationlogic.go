package user

import (
	"context"
	"database/sql"
	"time"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/pkg/generator/username"
	"carservice/internal/pkg/wechat"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type WechatAuthorizationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWechatAuthorizationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatAuthorizationLogic {
	return &WechatAuthorizationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// WechatAuthorization
// todo: test this api.
func (l *WechatAuthorizationLogic) WechatAuthorization(req *types.WechatAuthorizationReq) (resp *types.WechatAuthorizationRep, err error) {
	if len(req.Code) == 0 {
		return nil, errcode.InvalidWechatCodeErr
	}
	// call code2session from api of official of wechat.
	mp := wechat.NewWechatProvider(l.svcCtx.Config.WechatConf).MiniProgram()
	code2session, err := mp.Code2session(req.Code)
	if err != nil {
		logc.Errorf(l.ctx, "微信获取 Session 时发生错误 (1), err: %s\n", err.Error())
		return nil, errcode.WechatCode2SessionErr
	}
	if code2session.Errcode != 0 {
		if code2session.Errcode == 40163 {
			return nil, errcode.WechatCode2SessionErr.SetMessage("code 已被使用")
		}
		logc.Errorf(l.ctx, "微信获取 Session 时发生错误 (2), err: %s\n", code2session.Errmsg)
		return nil, errcode.WechatCode2SessionErr
	}
	openid := code2session.Openid
	unionid := code2session.Unionid
	sessionKey := code2session.SessionKey
	// check if openid of user exists in the member table.
	var exist int8
	query := "SELECT (COUNT(1) = 1) AS `exist` FROM `member_binds` WHERE `open_id` = ? LIMIT 1"
	stmt, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
	if err != nil {
		logc.Error(l.ctx, "预处理查询用户 open_id 是否存在语句时发生错误, err:"+err.Error())
		return nil, errcode.DatabasePrepareErr
	}
	if err = stmt.GetContext(l.ctx, &exist, openid); err != nil {
		logc.Error(l.ctx, "查询用户 open_id 是否存在语句时发生错误, err:"+err.Error())
		return nil, errcode.DatabaseGetErr
	}
	// get now time
	nowString := time.Now().Unix()
	if exist == 1 {
		// get user id and make jwt token.
		query = "SELECT `m`.`id` AS `userId` FROM `members` AS `m` JOIN `member_binds` AS `mb` ON `mb`.`user_id` = `m`.`id` WHERE `mb`.`open_id` = ? LIMIT 1"
		var userId uint
		stmtx, err := l.svcCtx.DBC.PreparexContext(l.ctx, query)
		if err != nil {
			logc.Error(l.ctx, "预处理查询登录用户的 open_id 语句时发生错误, err:"+err.Error())
			return nil, errcode.DatabasePrepareErr
		}
		if err = stmtx.GetContext(l.ctx, &userId, openid); err != nil {
			logc.Error(l.ctx, "查询登录用户的 open_id 语句时发生错误, err:"+err.Error())
			return nil, errcode.DatabaseGetErr
		}
		// make token
		token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.AccessSecret, nowString, 36000, uint(userId))
		if err != nil {
			logc.Errorf(l.ctx, "生成已存在用户的 Token 时发生错误, RealError: %s\n", err.Error())
			return nil, errcode.GenTokenErr
		}
		return &types.WechatAuthorizationRep{
			Token: token,
		}, nil
	} else {
		// create a new user.
		// get max increment id.
		var maxIncrementId uint
		query := "SELECT IFNULL(MAX(`id`), 0) AS `maxIncrementId` FROM `members` LIMIT 1"
		if err = l.svcCtx.DBC.Get(&maxIncrementId, query); err != nil {
			logc.Error(l.ctx, "获取最新的用户 ID 语句时发生错误, err:"+err.Error())
			return nil, errcode.DatabaseGetErr
		}
		if maxIncrementId == 0 {
			maxIncrementId = 1
		}
		// begin transaction.
		txx, err := l.svcCtx.DBC.BeginTxx(l.ctx, &sql.TxOptions{})
		if err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			return nil, errcode.DatabasePrepareErr
		}
		// create a new user.
		query = "INSERT INTO `members`(`username`, `phone_number`) VALUES(?, ?)"
		var newUserId int64
		var newUsername = "新用户 " + username.GenerateHexById(maxIncrementId+1)
		// default unbound phone number.
		stmtx, err := txx.PreparexContext(l.ctx, query)
		if err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			logc.Error(l.ctx, "预处理创建用户语句时发生错误, err:"+err.Error())
			return nil, errcode.DatabasePrepareErr
		}
		result, err := stmtx.ExecContext(l.ctx, newUsername, "")
		if err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			logc.Errorf(l.ctx, "创建用户时发生错误, RealError: %s\n", err.Error())
			return nil, errcode.DatabaseCreateErr
		}
		newUserId, _ = result.LastInsertId()
		// create a user binding record.
		query = "INSERT INTO `member_binds`(`user_id`, `app_id`, `open_id`, `union_id`, `session_key`, `access_token`) VALUES(?, ?, ?, ?, ?, ?)"
		// _, err = l.svcCtx.DBC.Exec(query, newUserId, l.svcCtx.Config.WechatConf.MiniProgram.AppId, openid, unionid, sessionKey, "none")
		// don't save appid/unionid/access_token.
		stmtx, err = txx.PreparexContext(l.ctx, query)
		if err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			logc.Error(l.ctx, "预处理创建用户授权表语句时发生错误, err:"+err.Error())
			return nil, errcode.DatabaseCreateErr
		}

		_, err = stmtx.ExecContext(l.ctx, newUserId, "", openid, unionid, sessionKey, "")
		if err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			logc.Errorf(l.ctx, "执行创建用户授权表时发生错误, RealError: %s\n", err.Error())
			return nil, errcode.DatabaseCreateErr
		}
		if err = txx.Commit(); err != nil {
			if rberr := txx.Rollback(); rberr != nil {
				logc.Error(l.ctx, "数据库回滚时发生错误, err:"+err.Error())
				return nil, errcode.DatabaseRollbackErr
			}
			logc.Error(l.ctx, "数据库提交时时发生错误, err:"+err.Error())
			return nil, errcode.DatabaseCommitErr
		}
		// make jwt token.
		token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.AccessSecret, nowString, 36000, uint(newUserId))
		if err != nil {
			logc.Errorf(l.ctx, "创建新用户的 Token 时发生错误, RealError: %s\n", err.Error())
			return nil, errcode.GenTokenErr
		}
		return &types.WechatAuthorizationRep{
			Token: token,
		}, nil
	}
}
