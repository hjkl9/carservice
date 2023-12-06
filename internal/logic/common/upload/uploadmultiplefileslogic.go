package upload

import (
	"context"

	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadMultipleFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadMultipleFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadMultipleFilesLogic {
	return &UploadMultipleFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadMultipleFilesLogic) UploadMultipleFiles() (resp *types.UploadMultipleFilesRep, err error) {
	// todo: add your logic here and delete this line

	return
}
