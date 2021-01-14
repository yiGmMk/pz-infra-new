package gopay

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	. "github.com/gyf841010/pz-infra-new/logging"
	"github.com/gyf841010/pz-infra-new/payUtil/gopay/client"
	"github.com/gyf841010/pz-infra-new/payUtil/gopay/common"
	"github.com/gyf841010/pz-infra-new/payUtil/gopay/util"
)

func AliWebCallback(w http.ResponseWriter, r *http.Request) (*common.AliWebPayResult, error) {
	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	for k, v := range r.Form {
		// k不会有多个值的情况
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}

	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA" {
		//错误日志
		return nil, errors.New("签名类型未知")
	}

	err := client.DefaultAliWebClient().CheckSign(signData, m["sign"])
	if err != nil {
		//log.Error("签名验证失败：", err, signData, m)
		return nil, err
	}

	var aliPay common.AliWebPayResult
	err = util.MapStringToStruct(m, &aliPay)
	if err != nil {
		//log.Error(err)
		w.Write([]byte("error"))
		return nil, err
	}

	// err = biz.AliWebCallBack(&aliPay)
	// if err != nil {
	// 	//log.Error(err)
	// 	w.Write([]byte("error"))
	// 	return nil, err
	// }
	w.Write([]byte("success"))
	return &aliPay, nil
}

// 支付宝app支付回调
func AliAppCallback(w http.ResponseWriter, r *http.Request) (*common.AliWebPayResult, error) {
	var m = make(map[string]string)
	var signSlice []string
	r.ParseForm()
	for k, v := range r.Form {
		m[k] = v[0]
		if k == "sign" || k == "sign_type" {
			continue
		}
		signSlice = append(signSlice, fmt.Sprintf("%s=%s", k, v[0]))
	}

	Log.Debug("Original Response", With("m", m))

	sort.Strings(signSlice)
	signData := strings.Join(signSlice, "&")
	if m["sign_type"] != "RSA" && m["sign_type"] != "RSA2" {
		//log.Error(m)
		w.Write([]byte("error"))
		return nil, errors.New("签名类型未知")
	}

	err := client.DefaultAliAppClient().CheckSign(signData, m["sign"])
	if err != nil {
		//log.Error(err, m, signData)
		w.Write([]byte("error"))
		return nil, err
	}

	var aliPay common.AliWebPayResult
	var aliPayMap common.AliWebPayMap
	err = util.MapStringToStruct(m, &aliPayMap)
	if err != nil {
		//log.Error(err)
		w.Write([]byte("error"))
		return nil, err
	}

	if aliPayMap.FundBillList != "" {
		err = json.Unmarshal([]byte(aliPayMap.FundBillList), &aliPay.FundBillList)
		if err != nil {
			Log.Error("Failed to Unmarshal Fund Bill List", WithError(err))
			w.Write([]byte("error"))
			return nil, err
		}
	}
	aliPay.NotifyTime = aliPayMap.NotifyTime
	aliPay.NotifyType = aliPayMap.NotifyType
	aliPay.NotifyID = aliPayMap.NotifyID
	aliPay.SignType = aliPayMap.SignType
	aliPay.Sign = aliPayMap.Sign
	aliPay.OutTradeNo = aliPayMap.OutTradeNo
	aliPay.Subject = aliPayMap.Subject
	aliPay.TradeNo = aliPayMap.TradeNo
	aliPay.TradeStatus = aliPayMap.TradeStatus
	aliPay.GmtPayMent = aliPayMap.GmtPayMent
	aliPay.SellerEmail = aliPayMap.SellerEmail
	aliPay.BuyerLogonId = aliPayMap.BuyerLogonId
	aliPay.SellerID = aliPayMap.SellerID
	aliPay.BuyerID = aliPayMap.BuyerID
	aliPay.TotalAmount = aliPayMap.TotalAmount
	aliPay.MDiscountAmount = aliPayMap.MDiscountAmount
	aliPay.DiscountAmount = aliPayMap.DiscountAmount
	aliPay.PassbackParams = aliPayMap.PassbackParams

	w.Write([]byte("success"))
	return &aliPay, nil
}

// WeChatCallback 微信app支付
func WeChatAppCallback(w http.ResponseWriter, r *http.Request) (*common.WeChatPayResult, error) {
	var returnCode = "FAIL"
	var returnMsg = ""
	defer func() {
		formatStr := `<xml><return_code><![CDATA[%s]]></return_code>
                  <return_msg>![CDATA[%s]]</return_msg></xml>`
		returnBody := fmt.Sprintf(formatStr, returnCode, returnMsg)
		w.Write([]byte(returnBody))
	}()
	var reXML common.WeChatPayResult
	//body := cb.Ctx.Input.RequestBody
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		//log.Error(string(body))
		returnCode = "FAIL"
		returnMsg = "Bodyerror"
		return nil, err
	}
	err = xml.Unmarshal(body, &reXML)
	if err != nil {
		//log.Error(err, string(body))
		returnMsg = "参数错误"
		returnCode = "FAIL"
		return nil, err
	}

	if reXML.ReturnCode != "SUCCESS" {
		//log.Error(reXML)
		returnCode = "FAIL"
		return &reXML, errors.New(reXML.ReturnCode)
	}
	m, err := util.XmlToMap(body)
	if err != nil {
		//log.Error(err, body)
		returnMsg = "参数错误"
		returnCode = "FAIL"
		return nil, err
	}
	//log.Info(m)
	var signData []string
	for k, v := range m {
		if k == "sign" {
			continue
		}
		signData = append(signData, fmt.Sprintf("%v=%v", k, v))
	}
	sort.Strings(signData)
	signData2 := strings.Join(signData, "&")
	err = client.DefaultWechatAppClient().CheckSign(signData2, m["sign"])
	if err != nil {
		returnCode = "FAIL"
		return nil, err
	}

	returnCode = "SUCCESS"
	return &reXML, nil
}
