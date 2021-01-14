package gopay

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/client"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/common"
	"github.com/yiGmMk/pz-infra-new/payUtil/gopay/constant"
)

func TestPay(t *testing.T) {
	initClient()
	initHandle()
	charge := new(common.Charge)
	charge.PayMethod = constant.WECHAT
	charge.MoneyFee = 1
	charge.Describe = "test pay"
	charge.TradeNum = "1111111111"

	fdata, err := Pay(charge)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fdata)
}

func initClient() {
	client.InitAliAppClient(&client.AliAppClient{
		AppID:      "xxx",
		SignType:   "RSA",
		PrivateKey: nil,
		PublicKey:  nil,
	})
}

func initHandle() {
	http.HandleFunc("callback/aliappcallback", func(w http.ResponseWriter, r *http.Request) {
		aliResult, err := AliAppCallback(w, r)
		if err != nil {
			fmt.Println(err)
			//log.xxx
			return
		}
		selfHandler(aliResult)
	})
}

func selfHandler(i interface{}) {
}
