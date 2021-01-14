package emailUtil

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"

	. "github.com/gyf841010/pz-infra-new/logging"
	"github.com/gyf841010/pz-infra-new/notification"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-gomail/gomail"
)

const (
	CONTENT_TYPE_HTML  = "text/html"
	CONTENT_TYPE_PLAIN = "text/plain"
)

type EmailSender struct {
	User     string //
	Password string
	SmtpUrl  string
	SmtpPort int
	Name     string //显示的发信人名称
}

type EmailReceiver struct {
	Receivers   []string
	ContentType string //default: text/html
	Subject     string
	Attachments []string //附件
}

//Sending email to with HTML Template
func SendEmailAws(to, subject, templateName string, params map[string]string) error {
	content := ReplaceTemplateContent(templateName, params)
	body := &notification.Body{
		Html: &notification.Content{Data: aws.String(content)},
	}
	destination := &notification.Destination{
		ToAddresses: []*string{aws.String(to)},
	}
	message := &notification.Message{
		Subject: &notification.Content{Data: aws.String(subject)},
		Body:    body,
	}
	email := &notification.Email{
		Destination: destination,
		Message:     message,
		Source:      aws.String(beego.AppConfig.String("user.resetpwd.email.from")),
	}
	if err := notification.SendEmail(email); err != nil {
		Log.Error("Failed to Send Email", With("to", to), WithError(err))
		return err
	}
	return nil
}

func ReplaceTemplateContent(templateName string, params map[string]string) string {
	templateContent, _ := ioutil.ReadFile(fmt.Sprintf("conf/email_template/%s.html", templateName))
	content := string(templateContent)
	for key, value := range params {
		key = "{{" + key + "}}"
		content = strings.Replace(content, key, value, -1)
	}
	return content
}

// Legacy Method, no usage now
func SendEmailSmtp(templateName string, content map[string]string, receiver *EmailReceiver) error {
	if receiver.ContentType == "" {
		receiver.ContentType = CONTENT_TYPE_HTML
	}

	sender := getSender()

	m := gomail.NewMessage()
	m.SetHeader("From", sender.Name)
	//m.SetAddressHeader("From", sender.Name, "info")
	m.SetHeader("To", receiver.Receivers...)
	m.SetHeader("Subject", receiver.Subject)
	for _, attach := range receiver.Attachments {
		m.Attach(attach)
	}

	if receiver.ContentType == CONTENT_TYPE_HTML {
		templateContent, _ := ioutil.ReadFile(fmt.Sprintf("conf/email_template/%s.html", templateName))
		tc := string(templateContent)
		for key, value := range content {
			key = "{{" + key + "}}"
			tc = strings.Replace(tc, key, value, -1)
		}
		m.SetBody(CONTENT_TYPE_HTML, tc)
	}

	d := gomail.NewDialer(sender.SmtpUrl, sender.SmtpPort, sender.User, sender.Password)
	if err := d.DialAndSend(m); err != nil {
		Log.Error("failed to send email", WithError(err))
		return err
	}
	return nil
}

func getSender() *EmailSender {
	sender := EmailSender{
		User:     beego.AppConfig.String("email_user"),
		Password: beego.AppConfig.String("email_pwd"),
		SmtpUrl:  beego.AppConfig.String("email_smtp_url"),
		Name:     beego.AppConfig.String("email_name"),
	}
	sender.SmtpPort, _ = beego.AppConfig.Int("email_smtp_port")
	return &sender
}

func GetEmailDisplayName(email string) string {
	if len(email) <= 0 {
		return ""
	}
	seps := strings.Split(email, "@")
	if len(seps) <= 0 {
		return ""
	}
	return seps[0]
}

// Legacy Method, no usage now
func SendEmail(subject, templateName string, params map[string]string, receiver *EmailReceiver) error {
	auth := smtp.PlainAuth("", "258257921@qq.com", "gongyaofei2003", "smtp.qq.com:465")
	to := receiver.Receivers
	user := "258257921@qq.com"
	content := ReplaceTemplateContent(templateName, params)
	body := []byte(content)
	err := smtp.SendMail("smtp.qq.com:465", auth, user, to, body)
	if err != nil {
		Log.Error("Failed to Send Email", With("to", to), WithError(err))
		return err
	}
	return nil
}
