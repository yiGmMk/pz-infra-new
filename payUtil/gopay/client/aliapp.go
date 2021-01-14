package client

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	. "github.com/gyf841010/pz-infra-new/logging"
	"github.com/gyf841010/pz-infra-new/payUtil/gopay/common"

	"github.com/astaxie/beego"
)

var defaultAliAppClient *AliAppClient

type AliAppClient struct {
	AppID      string // 应用ID
	SignType   string // 验签模式
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func InitAliAppClient(c *AliAppClient) {
	defaultAliAppClient = c
}

// DefaultAliAppClient 得到默认支付宝app客户端
func DefaultAliAppClient() *AliAppClient {
	return defaultAliAppClient
}

// 获取退款接口
func (aa *AliAppClient) Refund(refund *common.Refund) error {
	err := checkRefund(refund)
	if err != nil {
		return err
	}

	ct := DefaultAliAppClient()
	re, err := ct.BuildRefundMap(refund)
	if err != nil {
		return err
	}
	reqStr := ct.BuildRefundRequestString(re)
	refundUrl := fmt.Sprintf("%s?%s", beego.AppConfig.String("ali.gatewayUrl"), reqStr)
	Log.Debug("refundUrl", With("refundUrl", refundUrl))

	client := &http.Client{}
	newReq, err := http.NewRequest("GET", refundUrl, nil)
	if err != nil {
		Log.Error("Failed to New Ali Refund Request", WithError(err))
		return err
	}
	Log.Info("Ali Refund Request", With("req", newReq))
	resp, err := client.Do(newReq)
	if err != nil {
		Log.Error("Failed to Do Ali Refund Request", WithError(err))
		return err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Log.Error("Failed to read Resp body for Ali Refund Request", WithError(err))
		return err
	}
	Log.Debug("Response for Ali Refund", With("result", string(result)))
	if resp.StatusCode != 200 {
		Log.Debug("Response for Ali Refund failed", With("statusCode", resp.StatusCode), With("result", result))
		return errors.New("ali refund failed")
	}
	var aliRefundResult common.AliRefundResult
	if err := json.Unmarshal([]byte(result), &aliRefundResult); err != nil {
		Log.Error("Failed to Unmarshal Resp body for Ali Refund", WithError(err))
		return err
	}
	if aliRefundResult.AlipayTradeRefundResp.Code != "10000" {
		Log.Warn("Response Error Msg for Ali Refund", With("Msg", aliRefundResult.AlipayTradeRefundResp.Msg))
		return errors.New(fmt.Sprintf("%s, %s", "response error message", aliRefundResult.AlipayTradeRefundResp.Msg))
	}
	Log.Debug("Response for Ali Refund success", With("aliRefundResult", aliRefundResult))
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	return nil
}

// 验证内容
func checkRefund(refund *common.Refund) error {
	if refund.PayMethod < 0 {
		return errors.New("payMethod less than 0")
	}
	if refund.RefundAmount < 0 {
		return errors.New("totalFee less than 0")
	}
	return nil
}

func (this *AliAppClient) BuildRefundMap(refund *common.Refund) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.trade.refund"
	m["charset"] = "utf-8"
	m["format"] = "JSON"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["sign_type"] = this.SignType
	m["biz_content"] = buildRefundBizContent(refund)

	sign, err := this.GenSign(m)
	if err != nil {
		return nil, err
	}
	m["sign"] = sign
	return m, nil
}

func buildRefundBizContent(refund *common.Refund) string {
	var mc = make(map[string]string)
	mc["trade_no"] = refund.TradeNo
	mc["refund_amount"] = fmt.Sprintf("%.2f", float64(refund.RefundAmount)/float64(100))
	byteBiz, _ := json.Marshal(mc)
	biz := string(byteBiz)
	return biz
}

func (this *AliAppClient) BuildRefundRequestString(m map[string]string) string {
	reqStr := ""
	reqStr += "app_id=" + url.QueryEscape(m["app_id"])
	reqStr += "&biz_content=" + url.QueryEscape(m["biz_content"])
	reqStr += "&charset=" + url.QueryEscape(m["charset"])
	reqStr += "&format=" + url.QueryEscape(m["format"])
	reqStr += "&method=" + url.QueryEscape(m["method"])
	reqStr += "&sign_type=" + url.QueryEscape(m["sign_type"])
	reqStr += "&timestamp=" + url.QueryEscape(m["timestamp"])
	reqStr += "&version=" + url.QueryEscape(m["version"])
	reqStr += "&sign=" + url.QueryEscape(m["sign"])

	return reqStr
}

// 获取支付接口
func (aa *AliAppClient) Pay(charge *common.Charge) (string, error) {
	err := checkCharge(charge)
	if err != nil {
		return "", err
	}

	ct := DefaultAliAppClient()
	re, err := ct.BuildPayMap(charge)
	if err != nil {
		return "", err
	}
	reqStr := ct.BuildPayRequestString(re)
	return reqStr, nil
}

// 验证内容
func checkCharge(charge *common.Charge) error {
	if charge.PayMethod < 0 {
		return errors.New("payMethod less than 0")
	}
	if charge.MoneyFee < 0 {
		return errors.New("totalFee less than 0")
	}
	return nil
}

func buildPayBizContent(charge *common.Charge) string {
	var mc = make(map[string]string)
	mc["subject"] = charge.Describe
	mc["out_trade_no"] = charge.TradeNum
	mc["product_code"] = "QUICK_MSECURITY_PAY"
	mc["total_amount"] = fmt.Sprintf("%.2f", float64(charge.MoneyFee)/float64(100))
	mc["passback_params"] = charge.PassbackParams
	byteBiz, _ := json.Marshal(mc)
	biz := string(byteBiz)
	return biz
}

func (this *AliAppClient) BuildPayRequestString(m map[string]string) string {
	reqStr := ""
	reqStr += "app_id=" + url.QueryEscape(m["app_id"])
	reqStr += "&biz_content=" + url.QueryEscape(m["biz_content"])
	reqStr += "&charset=" + url.QueryEscape(m["charset"])
	reqStr += "&format=" + url.QueryEscape(m["format"])
	reqStr += "&method=" + url.QueryEscape(m["method"])
	reqStr += "&notify_url=" + url.QueryEscape(m["notify_url"])
	reqStr += "&sign_type=" + url.QueryEscape(m["sign_type"])
	reqStr += "&timestamp=" + url.QueryEscape(m["timestamp"])
	reqStr += "&version=" + url.QueryEscape(m["version"])
	reqStr += "&sign=" + url.QueryEscape(m["sign"])

	return reqStr
}

func (this *AliAppClient) BuildPayMap(charge *common.Charge) (map[string]string, error) {
	var m = make(map[string]string)
	m["app_id"] = this.AppID
	m["method"] = "alipay.trade.app.pay"
	m["charset"] = "utf-8"
	m["format"] = "JSON"
	m["timestamp"] = time.Now().Format("2006-01-02 15:04:05")
	m["version"] = "1.0"
	m["notify_url"] = charge.CallbackURL
	m["sign_type"] = this.SignType
	m["biz_content"] = buildPayBizContent(charge)

	sign, err := this.GenSign(m)
	if err != nil {
		return nil, err
	}
	m["sign"] = sign
	return m, nil
}

// GenSign 产生签名
func (this *AliAppClient) GenSign(m map[string]string) (string, error) {
	delete(m, "sign")
	var data []string
	for k, v := range m {
		if v == "" {
			continue
		}
		data = append(data, fmt.Sprintf(`%s=%s`, k, v))
	}
	sort.Strings(data)
	signData := strings.Join(data, "&")
	if m["sign_type"] == "RSA" {
		s := sha1.New()
		_, err := s.Write([]byte(signData))
		if err != nil {
			return "", err
		}
		hashByte := s.Sum(nil)
		signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA1)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(signByte), nil
	} else if m["sign_type"] == "RSA2" {
		s := sha256.New()
		_, err := s.Write([]byte(signData))
		if err != nil {
			return "", err
		}
		hashByte := s.Sum(nil)
		signByte, err := this.PrivateKey.Sign(rand.Reader, hashByte, crypto.SHA256)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(signByte), nil
	} else {
		return "", errors.New("sign type err")
	}
}

// CheckSign 检测签名
func (this *AliAppClient) CheckSign(signData, sign string) error {
	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return err
	}
	if this.SignType == "RSA" {
		s := sha1.New()
		if _, err = s.Write([]byte(signData)); err != nil {
			return err
		}
		hash := s.Sum(nil)
		return rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA1, hash, signByte)
	} else if this.SignType == "RSA2" {
		s := sha256.New()
		if _, err = s.Write([]byte(signData)); err != nil {
			return err
		}
		hash := s.Sum(nil)
		return rsa.VerifyPKCS1v15(this.PublicKey, crypto.SHA256, hash, signByte)
	} else {
		return errors.New("签名类型未知")
	}
}
