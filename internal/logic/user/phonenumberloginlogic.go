package user

import (
	"context"
	"net/http"
	"time"

	"carservice/internal/datatypes/user"
	"carservice/internal/pkg/constant"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/x/errors"
)

type PhoneNumberLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPhoneNumberLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PhoneNumberLoginLogic {
	return &PhoneNumberLoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PhoneNumberLoginLogic) PhoneNumberLogin(req *types.PhoneNumberLoginReq) (resp *types.PhoneNumberLoginRep, err error) {
	key := constant.SmsCaptchaPrefix + req.PhoneNumber
	// Check if the phone number is in the cache.
	cmd := l.svcCtx.RDBC.Exists(l.ctx, key)
	exists, err := cmd.Result()
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Redis 数据库查询数据时出现错误")
	}
	if exists == 0 {
		return nil, errors.New(http.StatusBadRequest, "你还未发送短信")
	}
	// If the phone number exists.
	captcha, err := l.svcCtx.RDBC.Get(l.ctx, key).Result()
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Redis 数据库查询数据时出现错误")
	}
	// Check if the captcha is correct.
	if captcha != req.Captcha {
		return nil, errors.New(http.StatusBadRequest, "验证码不正确")
	}
	// // todo: waiting for SMS service to resume.
	// if req.Captcha != "888888" {
	// 	return nil, errors.New(http.StatusBadRequest, "验证码不正确")
	// }
	// Check the user if exsits in the database.
	query := "SELECT 1 FROM `users` WHERE `phone_number` = ?"
	var hasUser int8
	l.svcCtx.DBC.Get(&hasUser, query, req.PhoneNumber)
	nowString := time.Now().String()
	if hasUser == 1 {
		var u user.UserID
		l.svcCtx.DBC.Get(&u, "SELECT id FROM `users` WHERE `phone_number` = ?", req.PhoneNumber)
		// Generate token by jwt util.
		// Payload contains [id].
		token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.SecretKey, nowString, "36000", u.ID)
		if err != nil {
			return nil, errors.New(http.StatusInternalServerError, "Token 颁发时出现错误")
		}
		// Set field token in resp.
		resp = &types.PhoneNumberLoginRep{
			Token: token,
		}
		return resp, nil
	}
	// otherwise the new user because phone number doesn't exsit in database.
	defaultUsername := "新用户"
	query = "INSERT INTO `users`(`phone_number`, `username`) VALUES(?, ?)"
	result, err := l.svcCtx.DBC.Exec(query, req.PhoneNumber, defaultUsername)
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Mysql 数据库创建数据时出现错误")
	}
	// Get last insert id.
	newId, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Mysql 数据库查询数据时出现错误")
	}
	// This newId will be used to generate the jwt token.
	token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.SecretKey, nowString, "36000", uint(newId))
	// Generate token by jwt util.
	// Payload contains [id].
	if err != nil {
		return nil, errors.New(http.StatusInternalServerError, "Mysql 数据库创建数据时出现错误")
	}
	// Set field token in resp.
	resp = &types.PhoneNumberLoginRep{
		Token: token,
	}
	return resp, nil
}
