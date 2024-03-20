package service

import (
	"errors"
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/pkg/upload"
	"mime/multipart"
	"os"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	//获取文件所需的基本信息
	fileName := upload.GetFileName(fileHeader.Filename)

	//对其进行业务所需的文件检查(文件大小是否符合需求、文件后缀是否达到要求)
	if !upload.CheckContainExt(fileType, fileName) { //文件后缀是否达到要求
		return nil, errors.New("file suffix is not supported.")
	}
	if upload.CheckMaxSize(fileType, file) { //文件大小是否符合需求
		return nil, errors.New("exceeded maximum file limit.")
	}

	//判断在写入文件前是否具备必要的写入条件（目录是否存在、权限是否足够）
	uploadSavePath := upload.GetSavePath()    //获取文件保存地址
	if upload.CheckSavePath(uploadSavePath) { //目录是否存在
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory.")
		}
	}
	if upload.CheckPermission(uploadSavePath) { //权限是否足够
		return nil, errors.New("insufficient file permissions.")
	}

	//真正的写入文件
	dst := uploadSavePath + "/" + fileName
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
