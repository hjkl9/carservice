package jwt

import (
	"context"
)

func GetUserId(ctx context.Context) interface{} {
	return ctx.Value("user")
}
