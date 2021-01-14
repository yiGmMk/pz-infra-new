package client

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/yiGmMk/pz-infra-new/httpUtil"
	. "github.com/yiGmMk/pz-infra-new/logging"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/common"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/util"
)

var defaultWechatAppClient *WechatAppClient

// DefaultWechatAppClient 默认微信app客户端
func DefaultWechatAppClient() *WechatAppClient {
	return defaultWechatAppClient
}

// WechatAppClient 微信app支付
type WechatAppClient struct {
	AppID       string // AppID
	MchID       string // 商户号ID
	CallbackURL string // 回调地址
	Key         string // 密钥
	PayURL      string // 支付地址
	RefundURL   string // 退款地址
	CaCertFile  string // 双向证书根证书目录
	CertFile    string // 双向证书文件目录
	KeyFile     string // 双向证书Key目录
}

func InitWechatClient(c *WechatAppClient) {
	defaultWechatAppClient = c
}

// Pay 支付
func (wechat *WechatAppClient) Pay(charge *common.Charge) (string, error) {
	var m = make(map[string]string)
	m["appid"] = wechat.AppID
	m["mch_id"] = wechat.MchID
	m["nonce_str"] = util.RandomStr()
	m["body"] = charge.Describe
	m["out_trade_no"] = charge.TradeNum
	m["total_fee"] = fmt.Sprintf("%d", charge.MoneyFee)
	m["spbill_create_ip"] = util.LocalIP()
	m["notify_url"] = wechat.CallbackURL
	m["trade_type"] = "APP"

	sign, err := wechat.GenSign(m)
	if err != nil {
		return "", err
	}
	m["sign"] = sign
	// 转出xml结构
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s><![CDATA[%s]]></%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	re, err := HTTPSC.PostData(wechat.PayURL, "text/xml; charset=utf-8", xmlStr)
	if err != nil {
		return "", err
	}

	var xmlRe common.WeChatReResult
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return "", err
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return "", errors.New(xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 支付失败
		return "", errors.New(xmlRe.ErrCodeDes)
	}

	var c = make(map[string]string)
	c["appid"] = wechat.AppID
	c["partnerid"] = wechat.MchID
	c["prepayid"] = xmlRe.PrepayID
	c["package"] = "Sign=WXPay"
	c["noncestr"] = util.RandomStr()
	c["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())

	sign2, err := wechat.GenSign(c)
	if err != nil {
		return "", err
	}
	//c["signType"] = "MD5"
	c["paySign"] = strings.ToUpper(sign2)

	jsonC, err := json.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(jsonC), nil
}

// 获取退款接口
func (wechat *WechatAppClient) Refund(refund *common.Refund) error {
	var m = make(map[string]string)
	m["appid"] = wechat.AppID
	m["mch_id"] = wechat.MchID
	m["nonce_str"] = util.RandomStr()
	m["transaction_id"] = refund.TradeNo
	m["out_trade_no"] = refund.OutTradeNo
	m["out_refund_no"] = refund.OutRefundNo
	m["total_fee"] = fmt.Sprintf("%d", refund.RefundAmount)
	m["refund_fee"] = fmt.Sprintf("%d", refund.RefundAmount)
	//m["op_user_id"] = wechat.MchID

	sign, err := wechat.GenSign(m)
	if err != nil {
		return err
	}
	m["sign"] = sign
	// 转出xml结构
	buf := bytes.NewBufferString("")
	for k, v := range m {
		buf.WriteString(fmt.Sprintf("<%s>%s</%s>", k, v, k))
	}
	xmlStr := fmt.Sprintf("<xml>%s</xml>", buf.String())

	Log.Debug("Refund Request", With("m", m), With("xmlStr", xmlStr))

	re, err := httpUtil.PostXmlWithCert(wechat.RefundURL, xmlStr, wechat.CaCertFile, wechat.CertFile, wechat.KeyFile)
	if err != nil {
		Log.Error("Failed to Post Xml With Cert", WithError(err))
		return err
	}

	var xmlRe common.WeChatReResult
	err = xml.Unmarshal(re, &xmlRe)
	if err != nil {
		return err
	}

	if xmlRe.ReturnCode != "SUCCESS" {
		// 通信失败
		return errors.New(xmlRe.ReturnMsg)
	}

	if xmlRe.ResultCode != "SUCCESS" {
		// 退款失败
		return errors.New(xmlRe.ErrCodeDes)
	}

	return nil
}

// GenSign 产生签名
func (wechat *WechatAppClient) GenSign(m map[string]string) (string, error) {
	delete(m, "sign")
	delete(m, "key")
	var signData []string
	for k, v := range m {
		if v != "" {
			signData = append(signData, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(signData)
	signStr := strings.Join(signData, "&")
	signStr = signStr + "&key=" + wechat.Key
	c := md5.New()
	_, err := c.Write([]byte(signStr))
	if err != nil {
		return "", err
	}
	signByte := c.Sum(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", signByte), nil
}

// CheckSign 检查签名
func (wechat *WechatAppClient) CheckSign(data string, sign string) error {
	signData := data + "&key=" + wechat.Key
	c := md5.New()
	_, err := c.Write([]byte(signData))
	if err != nil {
		return err
	}
	signOut := fmt.Sprintf("%x", c.Sum(nil))
	if strings.ToUpper(sign) == strings.ToUpper(signOut) {
		return nil
	}
	return errors.New("签名交易错误")
}
