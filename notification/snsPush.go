package notification

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/yiGmMk/pz-infra-new/logging"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"

	"encoding/json"

	"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

//  Client Token State
type PushClientState int

const (
	PushClientStateActive PushClientState = iota
	PushClientStateInactive
)

//  Client Type , 0: unknown, 1: IOS 2: android
type PushClientType int

const (
	PushClientTypeUnknown PushClientType = iota
	PushClientTypeApple
	PushClientTypeAndroid
)

const (
	PUSH_PLATFORM_APNS         = "APNS"
	PUSH_PLATFORM_APNS_SANDBOX = "APNS_SANDBOX"
	PUSH_PLATFORM_GCM          = "GCM"
)

const (
	_SNS_ACCESS_KEY                    = "snsAccessKey"
	_SNS_SECRET_KEY                    = "snsSecretKey"
	_SNS_REGION                        = "snsRegion"
	_SNS_APNS_PLATFORM_APPLICATION_ARN = "snsApnsPlatformApplicationArn"
	_SNS_GCM_PLATFORM_APPLICATION_ARN  = "snsGcmPlatformApplicationArn"
)

var _PLATFORM_ARN_MAP = map[string]string{
	PUSH_PLATFORM_APNS: _SNS_APNS_PLATFORM_APPLICATION_ARN,
	PUSH_PLATFORM_GCM:  _SNS_GCM_PLATFORM_APPLICATION_ARN,
}

var (
	EndPointNotFound = errors.New("NotFound: Endpoint does not exist")
	MessageTooLong   = errors.New("Message too long")
)

type apnsNotificationBody struct {
	Aps         apnsAps         `json:"aps"`
	MessageMeta messageMetaBody `json:"data"`
}

type apnsAps struct {
	Alert apnsAlert `json:"alert"`
	Badge int       `json:"badge"`
}

type apnsAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type gcmNotificationBody struct {
	Notification gcmNotificationPayload `json:"notification"`
	Data         messageMetaBody        `json:"data"`
}

type gcmNotificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Badge int    `json:"badge"`
}

//type gcmDataPayload struct {
//	MessageMeta messageMetaBody `json:"data"`
//}

type messageMetaBody struct {
	MessageType      int         `json:"msg_type"`
	MessageArguments interface{} `json:"msg_argus"`
}

type PopupPromptsMessageArguments struct {
	Url string `json:"url"`
}

type UserClient struct {
	UserId    string
	Platform  string
	PushToken string
	SnsToken  string
}

type NotificationMessage struct {
	Title            string
	Body             string
	MessageType      int
	MessageArguments map[string]string
}

var (
	cbSnsPublish          = snsPublish
	cbSendSnsNotification = sendSnsNotification
)

const (
	ENDPOINT_ATTRIBUTE_KEY_ENABLED = "Enabled"
	ENDPOINT_ATTRIBUTE_KEY_TOKEN   = "Token"
)

func IsSupportPlatform(platform string) bool {
	_, isSupport := _PLATFORM_ARN_MAP[platform]
	return isSupport
}

func SendNotification(tokens []*UserClient, messageCount int, message *NotificationMessage) error {
	if len(tokens) <= 0 {
		return nil
	}
	platformMessageMap := make(map[string]string)
	errorMessages := []string{}
	for _, token := range tokens {
		body, err := getMessage(token.Platform, messageCount, message, platformMessageMap)
		if err != nil {
			if err == MessageTooLong {
				return err
			}
			errorMessages = append(errorMessages, err.Error())
			continue
		}
		Log.Debug("####notification message ", With("body", body))
		if err := cbSendSnsNotification(token.SnsToken, body); err != nil {
			Log.Error("###send notification failed.", With("token", token), WithError(err))
			errorMessages = append(errorMessages, err.Error())
		}
		Log.Debug("successfully send notification, userClient=> ", With("token", token))
	}
	if len(errorMessages) > 0 {
		err := errors.New(strings.Join(errorMessages, ","))
		Log.Error("send notification failed.", WithError(err))
		return err
	}
	return nil
}

func sendSnsNotification(snsEndPoint, snsBody string) error {
	if snsEndPoint == "" || snsBody == "" {
		return Log.Error("invalid sns message")
	}

	snsClient, err := createSnsClient()
	if err != nil {
		return Log.Error("createSnsClient failed", WithError(err))
	}
	_, err = cbSnsPublish(snsClient, &sns.PublishInput{
		Message:          aws.String(snsBody),
		MessageStructure: aws.String("json"),
		TargetArn:        aws.String(snsEndPoint),
	})
	return err
}

func getMessage(platform string, messageCount int, message *NotificationMessage, platformMessageMap map[string]string) (string, error) {
	if message, isExist := platformMessageMap[platform]; isExist {
		return message, nil
	}
	m, err := generateNotificationBody(platform, messageCount, message)
	if err != nil {
		return "", err
	}
	// iOS increase this to 2048 bytes after version 8.0
	if len(m) > 2048 {
		return "", MessageTooLong
	}
	platformMessageMap[platform] = m
	return m, nil
}

func generateNotificationBody(platform string, messageCount int, message *NotificationMessage) (string, error) {
	if platform == PUSH_PLATFORM_APNS {
		return generateAPNSNotificationBody(messageCount, message)
	} else if platform == PUSH_PLATFORM_GCM {
		return generateGCMNotificationBody(messageCount, message)
	} else {
		return "", errors.New(fmt.Sprintf("platform not supported: %s", platform))
	}
}

func generateAPNSNotificationBody(messageCount int, message *NotificationMessage) (string, error) {
	body := apnsNotificationBody{
		Aps: apnsAps{
			Alert: apnsAlert{
				Title: message.Title,
				Body:  message.Body,
			},
			Badge: messageCount,
		},
		MessageMeta: messageMetaBody{
			MessageType:      message.MessageType,
			MessageArguments: message.MessageArguments,
		},
	}
	bodyBytes, err := json.Marshal(&body)
	if err != nil {
		return "", err
	}
	m := make(map[string]string)
	if "dev" == beego.AppConfig.String("runmode") {
		m[PUSH_PLATFORM_APNS_SANDBOX] = string(bodyBytes)
	} else {
		m[PUSH_PLATFORM_APNS] = string(bodyBytes)
	}
	mStr, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(mStr), nil
}

func generateGCMNotificationBody(messageCount int, message *NotificationMessage) (string, error) {
	body := gcmNotificationBody{
		Notification: gcmNotificationPayload{
			Title: message.Title,
			Body:  message.Body,
			Badge: messageCount,
		},
		Data: messageMetaBody{
			MessageType:      message.MessageType,
			MessageArguments: message.MessageArguments,
		},
	}
	result, err := json.Marshal(&body)
	if err != nil {
		return "", err
	}
	m := make(map[string]string)
	m[PUSH_PLATFORM_GCM] = string(result)
	mStr, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(mStr), nil
}

func snsPublish(client SnsClient, input *sns.PublishInput) (*sns.PublishOutput, error) {
	return client.Publish(input)
}

func GetSnsEndpointArn(platform, platformPushToken string) (string, error) {
	platformArn := beego.AppConfig.String(_PLATFORM_ARN_MAP[platform])
	if platformArn == "" {
		return "", errors.New(fmt.Sprintf("missing config; %s", _PLATFORM_ARN_MAP[platform]))
	}
	resp, err := createPlatformEndpoint(&sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(platformArn),
		Token:                  aws.String(platformPushToken),
	})
	if err != nil {
		Log.Error("create platform endpoint failed.", With("platform", platform), With("token", platformPushToken), WithError(err))
		return "", err
	}
	return *resp.EndpointArn, nil
}

var cbGetEndpointAttributes = func(snsClient SnsClient, input *sns.GetEndpointAttributesInput) (*sns.GetEndpointAttributesOutput, error) {
	return snsClient.GetEndpointAttributes(input)
}

func GetEndpointAttributes(endpointArn string) (map[string]*string, error) {
	input := &sns.GetEndpointAttributesInput{
		EndpointArn: aws.String(endpointArn),
	}
	snsClient, err := createSnsClient()
	if err != nil {
		return nil, Log.Error("createSnsClient failed", WithError(err))
	}
	resp, err := cbGetEndpointAttributes(snsClient, input)
	if err != nil {
		errMessage := fmt.Sprintf("%s", err)
		if strings.HasPrefix(errMessage, "NotFound") {
			Log.Warn("EndPoint NotFound", WithError(EndPointNotFound))
			return nil, EndPointNotFound
		}
		Log.Error("get endpoint attributes failed.", With("endpointArn", endpointArn), WithError(err))
		return nil, err
	}
	return resp.Attributes, nil
}

var cbSetEndpointAttributes = func(snsClient SnsClient, input *sns.SetEndpointAttributesInput) (*sns.SetEndpointAttributesOutput, error) {
	return snsClient.SetEndpointAttributes(input)
}

func SetEndpointAttributes(attributes map[string]string, endpointArn string) error {
	inputAttributes := make(map[string]*string)
	for key, value := range attributes {
		inputAttributes[key] = aws.String(value)
	}
	input := &sns.SetEndpointAttributesInput{Attributes: inputAttributes, EndpointArn: aws.String(endpointArn)}
	snsClient, err := createSnsClient()
	if err != nil {
		return Log.Error("createSnsClient failed", WithError(err))
	}
	_, err = cbSetEndpointAttributes(snsClient, input)
	if err != nil {
		Log.Error("set endpoint attributes failed.", With("input", input), WithError(err))
		return err
	}
	return nil
}

func createPlatformEndpoint(input *sns.CreatePlatformEndpointInput) (*sns.CreatePlatformEndpointOutput, error) {
	client, err := createSnsClient()
	if err != nil {
		return nil, err
	}
	return client.CreatePlatformEndpoint(input)
}

func DeletePlatformEndpoint(endpointArn string) (*sns.DeleteEndpointOutput, error) {
	input := &sns.DeleteEndpointInput{EndpointArn: aws.String(endpointArn)}
	client, err := createSnsClient()
	if err != nil {
		return nil, err
	}
	return client.DeleteEndpoint(input)
}

func createSnsClient() (SnsClient, error) {
	isMockSns := beego.AppConfig.String("isMockSns")
	if isMockSns == "true" {
		Log.Info("use mock sns client")
		return &MockSnsClient{}, nil
	}
	accessKey := beego.AppConfig.String(_SNS_ACCESS_KEY)
	secretKey := beego.AppConfig.String(_SNS_SECRET_KEY)
	region := beego.AppConfig.String(_SNS_REGION)
	if accessKey == "" || secretKey == "" || region == "" {
		return nil, errors.New(fmt.Sprintf("missing one or more conf: %s", []string{_SNS_SECRET_KEY, _SNS_ACCESS_KEY, _SNS_REGION}))
	}
	snsConfig := aws.NewConfig().WithCredentials(credentials.NewStaticCredentials(accessKey, secretKey, "")).WithRegion(region)
	return sns.New(session.New(), snsConfig), nil
}
