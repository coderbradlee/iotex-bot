package mailutil

import (
	"fmt"
	"os"

	"path/filepath"
	"strings"
)

type SendEmailConf struct {
	SmtpAddr    string   `json:"smtp_addr"`
	SmtpPort    int      `json:"smtp_port"`
	SmtpAccount string   `json:"smtp_account"`
	SmtpPwd     string   `json:"smtp_pwd"`
	From        string   `json:"from"`
	To          []string `json:"to"`
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

type EmailConf struct {
	To []string
}
type Email struct {
}

func (e *Email) Send(msg string) error {
	fmt.Println("mail:", msg)
	//emailCfg := &EmailConf{}
	//if !LoadConfig("email.toml", emailCfg) {
	//	log.Println("local email.toml err")
	//	return false
	//}
	//
	//date := time.Now().Format("2006-01-02")
	//m := gomail.NewMessage()
	//
	//m.SetHeader("From", smtpFrom)
	//m.SetHeader("To", emailCfg.To...)
	//
	//m.SetHeader("Subject", "91pool-etf-data-"+date)
	//m.SetBody("text/html", "91pool-etf-data-"+date)
	//
	//if len(att_files) > 0 {
	//	m.Attach(att_files[0])
	//}
	//
	//d := gomail.NewDialer(smtpHost, smtpPort, smtpAccount, smtpPwd)
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	//
	//if err := d.DialAndSend(m); err != nil {
	//	log.Println(err)
	//	return false
	//}
	//log.Println("send email ok!")
	//return true
	return nil
}
