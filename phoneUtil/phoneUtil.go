package phoneUtil

import (
	"regexp"
	"strings"

	. "github.com/gyf841010/pz-infra-new/logging"
)

const (
	MOBILE_VALID_REGEX  = "^(\\+)?(\\d){10,}$"
	DEFAULT_MOBILE_CODE = "1"
	MOBILE_PLUS_SYMBOL  = "+"
)

var mobileRegex = regexp.MustCompile(MOBILE_VALID_REGEX)

func IsValidMobile(mobile string) bool {
	isValid := mobileRegex.MatchString(mobile)
	if !isValid {
		Log.Warn("Invalid Email Format: ", With("mobile", mobile))
	}
	return isValid
}

func FormalizedMobile(mobile, countryCode string) string {
	if HasCountryCode(mobile) {
		return mobile
	}
	countryCodeMap := getCountryCodeMap()
	if mobileCode, found := countryCodeMap[countryCode]; found {
		return MOBILE_PLUS_SYMBOL + mobileCode + mobile
	}
	return MOBILE_PLUS_SYMBOL + DEFAULT_MOBILE_CODE + mobile
}

func UnformalizedMobile(mobile, defaultCountryCode string) string {
	if strings.HasPrefix(mobile, defaultCountryCode) {
		return mobile[2:]
	} else if strings.HasPrefix(mobile, "00"+defaultCountryCode[1:]) {
		return mobile[2+len(defaultCountryCode)-1:]
	} else {
		return mobile
	}
}

func HasCountryCode(mobile string) bool {
	if strings.HasPrefix(mobile, MOBILE_PLUS_SYMBOL) {
		return true
	}
	return false
}

// Replace " " and "-" characters in mobile number
func StripMobile(mobile string) string {
	result := strings.Replace(mobile, " ", "", -1)
	result = strings.Replace(result, "-", "", -1)
	return result
}
