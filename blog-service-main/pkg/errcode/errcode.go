package errcode

//编写常用的一些错误处理公共方法，标准化我们的错误输出

import (
	"fmt"
	"net/http"
)


// Error 用于表示错误的响应结果
type Error struct {
	code int `json:"code"`
	msg string `json:"msg"`
	details []string `json:"details"`
}


// codes 作为全局错误码的存储载体，便于查看当前注册情况
var codes = map[int]string{}


// NewError 创建新的Error实例,同时进行排重的校验
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已经存在，请更换一个", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}


func (e *Error) Error() string {
	//返回一个字符串而不带任何输出
	return fmt.Sprintf("错误码：%d, 错误信息：%s", e.Code(), e.Msg())
}


func (e *Error) Code() int {
	return e.code
}


func (e *Error) Msg() string {
	return e.msg
}


func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}


func (e *Error) Details() []string {
	return e.details
}


func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}
	return &newError
}


// StatusCode 针对一些特定错误码进行状态码的转换
func (e *Error) StatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ServerError.Code():
		return http.StatusInternalServerError
	case InvalidParams.Code():
		return http.StatusBadRequest
	case UnauthorizedAuthNotExist.Code():
		fallthrough
	case UnauthorizedTokenError.Code():
		fallthrough
	case UnauthorizedTokenGenerate.Code():
		fallthrough
	case UnauthorizedTokenTimeout.Code():
		return http.StatusUnauthorized
	case TooManyRequests.Code():
		return http.StatusTooManyRequests
	}

	return http.StatusInternalServerError
}



