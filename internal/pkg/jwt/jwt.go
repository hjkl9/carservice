package jwt

import (
	"net/http"

	"carservice/internal/pkg/common/errcode"
	stdresponse "carservice/internal/pkg/httper/response"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/rest/handler"
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
		stdresponse.Response(w, nil, errcode.UnauthorizedError.SetMsg(err.Error()))
	}
}
