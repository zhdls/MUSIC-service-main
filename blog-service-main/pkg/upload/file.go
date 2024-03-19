package upload

import (
	"github.com/travel_study/blog-service/global"
	"github.com/travel_study/blog-service/pkg/util"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)


type FileType int //定义了 FileType 为 int 的类型别名
const TypeImage FileType = iota + 1 //利用 FileType 作为类别标识的基础类型，并用iota作为它的初始值



// GetFileName 获取文件名称
func GetFileName(name string) string {
	//通过获取文件后缀并筛出原始文件名进行 MD5 加密
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)
	//返回经过加密处理后的文件名
	return fileName + ext
}

// GetFileExt 获取文件后缀
func GetFileExt(name string) string {
	//调用 path.Ext 方法进行循环查找”.“符号，
	//通过切片索引返回对应的文件后缀名称
	return path.Ext(name)
}

// GetSavePath 获取文件保存地址
func GetSavePath() string {
	//这里直接返回配置中的文件保存目录即可，也便于后续的调整
	return global.AppSetting.UploadSavePath
}



// CheckSavePath 检查保存目录是否存在
func CheckSavePath(dst string) bool {
	//利用 os.Stat 方法所返回的 error 值与系统中所定义的 oserror.ErrNotExist 进行判断，以此达到校验效果

	//调用 os.Stat 方法获取文件的描述信息 FileInfo
	_, err := os.Stat(dst)
	//调用 os.IsNotExist 方法进行判断
	return os.IsNotExist(err)
}

// CheckContainExt 检查文件后缀是否包含在约定的后缀配置项中
func CheckContainExt(t FileType, name string) bool {
	//获取文件后缀
	ext := GetFileExt(name)
	//所上传的文件的后缀有可能是大写、小写、大小写等，因此我们需要调用 strings.ToUpper 方法统一转为大写（固定的格式）来进行匹配
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

// CheckMaxSize 检查文件大小是否超出最大大小限制
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := ioutil.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}

// CheckPermission 检查文件权限是否足够
func CheckPermission(dst string) bool {
	//与 CheckSavePath 方法原理一致，是利用 oserror.ErrPermission 进行判断

	_, err := os.Stat(dst)
	return os.IsPermission(err)
}




// CreateSavePath 创建在上传文件时所使用的保存目录
func CreateSavePath(dst string, perm os.FileMode) error {
	//os.MkdirAll 方法，将会以传入的 os.FileMode 权限位去递归创建所需的所有目录结构
	//若涉及的目录均已存在，则不会进行任何操作，直接返回 nil
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}
	return nil
}

// SaveFile 保存所上传的文件
func SaveFile(file *multipart.FileHeader, dst string) error {
	//打开源地址的文件
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//创建目标地址的文件
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	//两者之间的文件内容拷贝
	_, err = io.Copy(out, src)
	return err
}




