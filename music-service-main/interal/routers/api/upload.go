package api

import (
	"github.com/gin-gonic/gin"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/interal/service"
	"github.com/travel_study/blog-service/pkg/app"
	"github.com/travel_study/blog-service/pkg/convert"
	"github.com/travel_study/blog-service/pkg/errcode"
	"github.com/travel_study/blog-service/pkg/upload"
)

type Upload struct{}

func NewUpload() Upload {
	return Upload{}
}

func (u Upload) UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file") //读取入参 file 字段的上传文件信息
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	//把入参 type 字段作为所上传文件类型的确立依据 （也可以通过解析上传文件后缀来确定文件类型）
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if fileHeader == nil || fileType <= 0 {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	//通过入参检查后进行 Service 的调用
	svc := service.New(c.Request.Context())
	//完成上传和文件保存
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		global.Logger.Errorf(c, "svc.UploadFile err: %v", err)
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	//返回文件的展示地址
	response.ToResponse(gin.H{
		"file_access_url": fileInfo.AccessUrl,
	})
}
