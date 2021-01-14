package notification

import (
	"errors"
	"fmt"

	. "github.com/gyf841010/pz-infra-new/logging"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	_SES_ACCESS_KEY = "sesAccessKey"
	_SES_SECRET_KEY = "sesSecretKey"
	_SES_REGION     = "sesRegion"
)

type Email struct {
	*Destination
	*Message
	Source *string
}

type Destination struct {
	BccAddresses []*string
	CcAddresses  []*string
	ToAddresses  []*string
}

type Message struct {
	Subject *Content
	Body    *Body
}

// Represents the body of the message. You can specify text, HTML, or both.
// If you use both, then the message should display correctly in the widest
// variety of email clients.
type Body struct {
	Html *Content
	Text *Content
}

// Represents textual data, plus an optional character set specification.
// By default, the text must be 7-bit ASCII, due to the constraints of the
// SMTP protocol. If the text must contain any other characters, then you must
// also specify a character set. Examples include UTF-8, ISO-8859-1, and Shift_JIS.
type Content struct {
	Charset *string
	Data    *string
}

func SendEmail(email *Email) error {
	if email == nil || email.Destination == nil || email.Message == nil || email.Source == nil {
		return errors.New("invalidate parameter")
	}
	if len(email.Destination.ToAddresses) == 0 {
		return errors.New("no ToAddresses")
	}
	if email.Message.Body == nil || (email.Message.Body.Text == nil && email.Message.Body.Html == nil) {
		return errors.New("no body")
	}
	if email.Message.Body.Html != nil && email.Message.Body.Html.Data == nil {
		return errors.New("html body is empty")
	}
	if email.Message.Body.Text != nil && email.Message.Body.Text.Data == nil {
		return errors.New("text body is empty")
	}
	svc, err := createSESClient()
	if err != nil {
		Log.Error("createSESClient failed", WithError(err))
		return err
	}

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{ // Required
			BccAddresses: email.Destination.BccAddresses,
			CcAddresses:  email.Destination.CcAddresses,
			ToAddresses:  email.Destination.ToAddresses,
		},
		Message: &ses.Message{ // Required
			Body: &ses.Body{ // Required
			},
			Subject: &ses.Content{ // Required
				Data:    email.Subject.Data, // Required
				Charset: email.Subject.Charset,
			},
		},
		Source: email.Source,
	}
	if email.Body.Html != nil {
		input.Message.Body.Html = &ses.Content{ // Required
			Data:    email.Body.Html.Data, // Required
			Charset: email.Body.Html.Charset,
		}
	}
	if email.Body.Text != nil {
		input.Message.Body.Text = &ses.Content{ // Required
			Data:    email.Body.Text.Data, // Required
			Charset: email.Body.Text.Charset,
		}
	}
	_, err = cbSendEmail(svc, input)

	if err != nil {
		Log.Error("send email failed", WithError(err))
		return err
	}
	Log.Debug("Send email successfully", With("email", email))
	return nil
}

var cbSendEmail = func(sesClient SESClient, input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return sesClient.SendEmail(input)
}

func createSESClient() (SESClient, error) {
	isMockSes := beego.AppConfig.String("isMockSes")
	if isMockSes == "true" {
		Log.Info("use mock SES client")
		return &MockSESClient{}, nil
	}
	accessKey := beego.AppConfig.String(_SES_ACCESS_KEY)
	secretKey := beego.AppConfig.String(_SES_SECRET_KEY)
	region := beego.AppConfig.String(_SES_REGION)
	if accessKey == "" || secretKey == "" || region == "" {
		return nil, errors.New(fmt.Sprintf("missing one or more conf: %s", []string{_SES_ACCESS_KEY, _SES_SECRET_KEY, _SES_REGION}))
	}
	sesConfig := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")).WithRegion(region)
	return ses.New(session.New(), sesConfig), nil
}
