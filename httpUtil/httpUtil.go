package httpUtil

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	. "github.com/yiGmMk/pz-infra-new/logging"
)

func PostJson(url string, header map[string]string, body interface{}) ([]byte, error) {
	bodyByte, err := json.Marshal(body)
	if err != nil {
		Log.Error("occur error when marshal object to json", WithError(err))
		return nil, err
	}

	if bodyByte != nil {
		Log.Debug("request body", With("requestBody", string(bodyByte)))
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyByte))
	if err != nil {
		Log.Error("occur error when new http request, ", WithError(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("occur error when get response, ", WithError(err))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Log.Error("Failed to read response body from HTTP Request", With("url", url), WithError(err))
	}

	return respBody, nil
}

func PostXmlWithCert(url string, body string, cacrtFile, crtFile, keyFile string) ([]byte, error) {
	pool := x509.NewCertPool()
	caCrt, err := ioutil.ReadFile(cacrtFile)
	if err != nil {
		Log.Error("Failed to Read Cert File, ", WithError(err))
		return nil, err
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		Log.Error("Failed to Load x509 Key Pair ", WithError(err))
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		Log.Error("occur error when new http request, ", WithError(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "text/xml:charset=UTF-8")

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("occur error when get response, ", WithError(err))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Log.Error("Failed to read response body from HTTP Request", With("url", url), WithError(err))
	}

	return respBody, nil
}

func GetJson(url string, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		Log.Error("occur error when new http request, ", WithError(err))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		Log.Error("occur error when get response, ", WithError(err))
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Log.Error("Failed to read response body from HTTP Request", With("url", url), WithError(err))
	}

	return respBody, nil
}

func IPAddrAcl(ip string) (RetStr string) {
	//本地环回地址属于白名单，允许访问
	if match0, _ := regexp.MatchString(`127\.0\.0\.1`, ip); match0 {
		RetStr = "white"
		return
	}
	//局域网地址：10.*.*.*需要鉴权，其中网关地址10.10.30.1，直接拒绝
	if match2, _ := regexp.MatchString(`10\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match2 {
		if match3, _ := regexp.MatchString(`10\.10\.30\.1`, ip); match3 {
			RetStr = "black"
			return
		}
		RetStr = "auth"
		return
	}
	//局域网地址：172.16.*.* - 172.31.*.* 需要鉴权
	if match4, _ := regexp.MatchString(`172\.((1[6-9])|(2[0-9])|(3[0-1]))\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match4 {
		RetStr = "auth"
		return
	}
	//局域网地址：192.168.*.* 需要鉴权
	if match5, _ := regexp.MatchString(
		`192\.168\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match5 {
		RetStr = "auth"
		return
	}
	//其余地址或非法字符串或入参为空，均为非法，直接拒绝。如果对上层调用者不信任，这里可以再细化区别处理。
	RetStr = "black"
	return RetStr
}

func IPAddrIsLan(ip string) bool {
	//本地环回地址属于白名单，允许访问
	if match0, _ := regexp.MatchString(`127\.0\.0\.1`, ip); match0 {
		return true
	}
	//局域网地址：10.*.*.*需要鉴权，其中网关地址10.10.30.1，直接拒绝
	if match2, _ := regexp.MatchString(`10\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match2 {
		if match3, _ := regexp.MatchString(`10\.10\.30\.1`, ip); match3 {
			return false
		}
		return true
	}
	//局域网地址：172.16.*.* - 172.31.*.* 需要鉴权
	if match4, _ := regexp.MatchString(`172\.((1[6-9])|(2[0-9])|(3[0-1]))\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match4 {
		return true
	}
	//局域网地址：192.168.*.* 需要鉴权
	if match5, _ := regexp.MatchString(
		`192\.168\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])\.(25[0-5]|2[0-4][0-9]|[0-1]?[0-9]?[0-9])`, ip); match5 {
		return true
	}
	//其余地址或非法字符串或入参为空，均为非法，直接拒绝。如果对上层调用者不信任，这里可以再细化区别处理。
	return false
}
