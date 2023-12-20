package user

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/constant"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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
		return nil, errcode.New(http.StatusInternalServerError, "-", "Redis 数据库查询数据时发生错误")
	}
	if exists == 0 {
		return nil, errcode.New(http.StatusBadRequest, "-", "你还未发送短信")
	}
	// If the phone number exists.
	captcha, err := l.svcCtx.RDBC.Get(l.ctx, key).Result()
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "-", "Redis 数据库查询数据时发生错误")
	}
	// Check if the captcha is correct.
	if captcha != req.Captcha {
		return nil, errcode.New(http.StatusBadRequest, "-", "手机验证码不正确")
	}
	// Check the user if exists in the database.
	nowString := time.Now().Unix()
	if l.svcCtx.Repo.UserRelated().CheckIfUserExistsByPhoneNumber(req.PhoneNumber) {
		var u = l.svcCtx.Repo.UserRelated().GetIdByPhoneNumber(req.PhoneNumber)
		// Generate token by jwt util.
		// Payload contains [id].
		fmt.Println(l.svcCtx.Config.JwtConf.AccessSecret)
		userPayload := jwt.UserPayload{
			UserId: u.ID,
		}
		token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.AccessSecret, nowString, 36000, userPayload)
		if err != nil {
			return nil, errcode.New(http.StatusInternalServerError, "-", "Token 颁发时发生错误")
		}
		// Set field token in resp.
		resp = &types.PhoneNumberLoginRep{
			Token: token,
		}
		return resp, nil
	}
	// otherwise the new user because phone number doesn't exsit in database.
	defaultUsername := "新用户"
	query := "INSERT INTO `members`(`phone_number`, `username`) VALUES(?, ?)"
	result, err := l.svcCtx.DBC.Exec(query, req.PhoneNumber, defaultUsername)
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "-", "Mysql 数据库创建数据时出现错误")
	}
	// Get last insert id.
	newId, err := result.LastInsertId()
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "-", "Mysql 数据库查询数据时出现错误")
	}
	// This newId will be used to generate the jwt token.
	userPayload := jwt.UserPayload{
		UserId: uint(newId),
	}
	token, err := jwt.GetJwtToken(l.svcCtx.Config.JwtConf.AccessSecret, nowString, 36000, userPayload)
	// Generate token by jwt util.
	// Payload contains [id].
	if err != nil {
		return nil, errcode.New(http.StatusInternalServerError, "-", "Mysql 数据库创建数据时出现错误")
	}
	// delete the CAPTCHA in rdb.
	// todo: errors may occur.
	l.svcCtx.RDBC.Del(l.ctx, key)
	// Set field token in resp.
	resp = &types.PhoneNumberLoginRep{
		Token: token,
	}
	return resp, nil
}
