package email

import (
	"crypto/tls"
	"gopkg.in/gomail.v2" //用于发送电子邮件的简单又高效的第三方开源库，目前只支持使用 SMTP 服务器发送电子邮件
)

type Email struct {
	*SMTPInfo
}

// SMTPInfo 用于传递发送邮箱所必需的信息
type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

func NewEmail(info *SMTPInfo) *Email {
	return &Email{SMTPInfo: info}
}

func (e *Email) SendMail(to []string, subject, body string) error {
	//调用 NewMessage 方法创建一个消息实例
	m := gomail.NewMessage()
	//用于设置邮件的一些必要信息
	m.SetHeader("From", e.From)     //发件人
	m.SetHeader("To", to...)        //收件人
	m.SetHeader("Subject", subject) //邮件主题
	m.SetBody("text/html", body)    //邮件正文

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password) //创建一个新的 SMTP 拨号实例
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}        //设置对应的拨号信息用于连接 SMTP 服务器
	return dialer.DialAndSend(m)                                       //打开与 SMTP 服务器的连接并发送电子邮件
}
