package gopay

import (
	"errors"

	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/client"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/common"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/constant"

	//"github.com/guidao/gopay/util"
	"strconv"
)

func Pay(charge *common.Charge) (string, error) {
	err := checkCharge(charge)
	if err != nil {
		//log.Error(err, charge)
		return "", err
	}

	ct := getPayClient(charge.PayMethod)
	re, err := ct.Pay(charge)
	if err != nil {
		//log.Error("支付失败:", err, charge)
		return "", err
	}
	return re, err
}

func Refund(refund *common.Refund) error {
	err := checkRefund(refund)
	if err != nil {
		//log.Error(err, refund)
		return err
	}

	ct := getPayClient(refund.PayMethod)
	err = ct.Refund(refund)
	if err != nil {
		//log.Error("退款失败:", err, refund)
		return err
	}
	return err
}

func checkCharge(charge *common.Charge) error {
	var id uint64
	var err error
	if charge.UserID == "" {
		id = 0
	} else {
		id, err = strconv.ParseUint(charge.UserID, 10, -1)
		if err != nil {
			return err
		}
	}
	if id < 0 {
		return errors.New("userID less than 0")
	}
	if charge.PayMethod < 0 {
		return errors.New("payMethod less than 0")
	}
	if charge.MoneyFee < 0 {
		return errors.New("totalFee less than 0")
	}

	if charge.CallbackURL == "" {
		return errors.New("callbackURL is NULL")
	}
	return nil
}

func checkRefund(refund *common.Refund) error {
	if refund.PayMethod < 0 {
		return errors.New("payMethod less than 0")
	}
	if refund.RefundAmount < 0 {
		return errors.New("totalFee less than 0")
	}
	return nil
}

// getPayClient 得到需要支付的客户端
func getPayClient(payMethod int64) common.PayClient {
	//如果使用余额支付
	switch payMethod {
	case constant.ALI_WEB:
		return client.DefaultAliWebClient()
	case constant.ALI_APP:
		//return client.DefaultAliAppOldClient()
		return client.DefaultAliAppClient()
	case constant.WECHAT:
		return client.DefaultWechatAppClient()
	}
	return nil
}
