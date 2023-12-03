package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// @secretKey: JWT 加解密密钥
// @iat: 时间戳
// @seconds: 过期时间，单位秒
// @payload: 数据载体
func GetJwtToken(secretKey string, iat, dur string, payload uint) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat
	claims["iat"] = iat
	claims["payload"] = payload
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
