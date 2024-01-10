package jwt

import (
	"errors"
	"net/http"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/httper/api"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/handler"
)

var (
	errInvalidToken = errors.New("invalid auth token")
	errNoClaims     = errors.New("no auth params")
)

type UserPayload struct {
	UserId uint
}

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GetJwtToken(secretKey string, iat, seconds int64, userId uint) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["user"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func UnauthorizedCallback() handler.UnauthorizedCallback {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		var msg string
		if e, ok := err.(*jwt.ValidationError); ok {
			if e.Is(jwt.ErrTokenSignatureInvalid) {
				api.Response(w, nil, errcode.UserInvalidTokenErr)
				return
			} else if e.Is(jwt.ErrTokenExpired) {
				api.Response(w, nil, errcode.UserTokenExpiredErr)
				return
			} else {
				msg = "其他错误"
				api.Response(w, nil, errcode.UnauthorizedError.SetMsg(msg))
				return
			}
		}
		api.Response(w, nil, errcode.UserUnauthorizedErr.SetMessage("其他错误"))
	}
}

func UnsignedCallback() handler.UnsignedCallback {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler, strict bool, code int) {
	}
}
