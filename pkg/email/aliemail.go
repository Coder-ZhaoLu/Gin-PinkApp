package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"pinkacg/settings"
)

func SendEmail(title, sendEmail, content string) (err error) {
	email := settings.Conf.EmailConfig.User
	username := settings.Conf.EmailConfig.Username
	password := settings.Conf.EmailConfig.Password
	host := settings.Conf.EmailConfig.Host
	port := 465
	header := make(map[string]string)
	header["From"] = username + "<" + email + ">"
	header["To"] = sendEmail
	header["Subject"] = title
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := "<html><body><h3>" + title + "</h3><p>" + content + "</p></body></html>"
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		email,
		password,
		host,
	)
	err = SendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{sendEmail},
		[]byte(message),
	)
	if err != nil {
		panic(err)
	} else {
	}
	return
}

//return a smtp client
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
