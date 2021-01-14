package commonUtil

import (
	"bytes"
	"crypto/md5"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gyf841010/pz-infra-new/log"

	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"
)

func UUID() string {
	return uuid.New()
}

func SetField(obj interface{}, name string, value interface{}) error {
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType != val.Type() {
		return errors.New("Provided value type didn't match obj field type")
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SToI(s string) (int, error) {
	return strconv.Atoi(s)
}

func ToJSON(object interface{}) (string, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// @return slice with the element at idx removed.
// panics if slice is not of Kind reflect.Slice
func DeleteSliceIndex(slice interface{}, idx int) interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {

		log.Errorf("Cannot call DeleteSliceIndex on a non-slice  of kind, slice: %+v, kind:%s", slice, v.Kind().String())
		return slice
	}

	if idx >= v.Len() {
		log.Errorf("Index for DeleteSliceIndex is out of bounds, idx: %d, slice: %+v, len: %d", idx, slice, v.Len())
		return slice
	}

	return reflect.AppendSlice(v.Slice(0, idx), v.Slice(idx+1, v.Len())).Interface()
}

func FloatValue(str string) float64 {
	rs, _ := strconv.ParseFloat(str, 64)
	return rs
}

func GetRandom(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

func GetMd5(v interface{}) string {
	content, _ := json.Marshal(v)
	md5Value := md5.New()
	io.WriteString(md5Value, string(content))
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(buffer, "%x", md5Value.Sum(nil))
	return buffer.String()
}

func LoadRSAPrivateKeyFromDisk(location string) *rsa.PrivateKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseRSAPrivateKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func LoadRSAPublicKeyFromDisk(location string) *rsa.PublicKey {
	keyData, e := ioutil.ReadFile(location)
	if e != nil {
		panic(e.Error())
	}
	key, e := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if e != nil {
		panic(e.Error())
	}
	return key
}

func ParsePublicKey(raw string) (result []byte) {
	return parseKey(raw, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----")
}

func ParsePrivateKey(raw string) (result []byte) {
	return parseKey(raw, "-----BEGIN RSA PRIVATE KEY-----", "-----END RSA PRIVATE KEY-----")
}

func parseKey(raw, prefix, suffix string) (result []byte) {
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var ll = 64
	var sl = len(raw)
	var c = sl / ll
	if sl%ll > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * ll
		var e = b + ll
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}

func ExternalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}

func GetMd5ForStr(signStr string) string {
	md5Value := md5.New()
	io.WriteString(md5Value, string(signStr))
	buffer := bytes.NewBuffer(nil)
	fmt.Fprintf(buffer, "%x", md5Value.Sum(nil))
	return buffer.String()
}

func StringRemoveDuplicates(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

func IntRemoveDuplicates(slc []int) []int {
	result := []int{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
