package util

import (
	"crypto/md5"
	"encoding/hex"
)


// EncodeMD5 针对上传后的文件名格式化
//将文件名进行MD5操作后再进行写入，防止直接把原始名称就暴露出去
func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	//hex.EncodeToString’函数 将字节切片转换为十六进制字符串，然后再将其作为普通字符串返回
	return hex.EncodeToString(m.Sum(nil))
}



