package common

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path"

	"carservice/internal/pkg/common/errcode"
	"carservice/internal/pkg/jwt"
	"carservice/internal/svc"
	"carservice/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImageLogic) UploadImage(f multipart.File, fh *multipart.FileHeader) (resp *types.UploadImageRep, err error) {
	userId := jwt.GetUserId(l.ctx)

	fileName := fh.Filename
	logc.Debugf(l.ctx, "file name: %s\n", fileName)
	// limit 1,500,000 bytes (1.5M).
	fileSize := fh.Size // bytes
	logc.Debugf(l.ctx, "file size: %d\n", fileSize)
	if fileSize > 1000000 {
		return nil, errcode.InvalidParamsError.SetMsgf("图片不能大于 %d bytes", 1000000)
	}

	supportExts := [3]string{
		".jpg",
		".jpeg",
		".png",
	}
	fileExt := path.Ext(fh.Filename)
	logc.Debugf(l.ctx, "file ext: %s\n", fileExt)
	logc.Debugf(l.ctx, "support ext: %v\n", l.checkExt(supportExts[:], fileExt))
	if !l.checkExt(supportExts[:], fileExt) {
		return nil, errcode.InvalidParamsError.SetMsg("不支持该图片格式")
	}

	// 存到本地
	wd, err := os.Getwd()
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("上传图片时出现错误")
	}

	// 随机文件名
	filename := "image_uid" + userId.(string) + "_" + uuid.New().String()

	file, err := os.OpenFile(wd+"/storage"+filename+fileExt, os.O_CREATE, 0755)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("上传图片时出现错误")
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	if err != nil {
		return nil, errcode.InternalServerError.SetMsg("上传图片时出现错误")
	}

	return &types.UploadImageRep{
		Url: filename,
	}, nil
}

func (l *UploadImageLogic) checkExt(supportExts []string, find string) bool {
	for _, v := range supportExts {
		if v == find {
			return true
		}
	}
	return false
}
