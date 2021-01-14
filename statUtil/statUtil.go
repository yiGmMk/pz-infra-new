package statUtil

import (
	. "github.com/yiGmMk/pz-infra-new/logging"
	"github.com/yiGmMk/pz-infra-new/redisUtil"
)

func PushStatData(jsonStr string) {
	if err := redisUtil.LpushString("statistics", jsonStr); err != nil {
		Log.Error("Failed to Send Stat Data while error occurs", WithError(err))
	}
}
