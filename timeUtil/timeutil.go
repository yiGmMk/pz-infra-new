package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	. "github.com/yiGmMk/pz-infra-new/errorUtil"
	. "github.com/yiGmMk/pz-infra-new/logging"
)

const (
	DATE_FORMAT       = "2006-01-02"
	TIME_FORMAT       = "2006-01-02 15:04:05"
	SPORT_UPDATE_TIME = 5
	WEEK_LEN          = 7 // 一周长度
	HOUR_SECOND       = 3600
	MIN5_MIRCO        = 5 * 60 * 1000
	MIN1_MIRCO        = 1 * 60 * 1000
	MIN_SECOND        = 1 * 60
	DAY_SECOND        = 24 * 3600
	HOUR_MIRCO        = 60 * 60 * 1000
	DAY_MICRO         = 24 * 60 * 60 * 1000
	FIFTEEN_MINS      = 15 * 60 * 1000
	TEN_SECOND        = 10
)

func CurrentUnix() int64 {
	return time.Now().Unix()
}

func CurrentUnixInt() int {
	return int(time.Now().Unix())
}

func CurrentUnixBigInt() int64 {
	return time.Now().UnixNano() / 1e6
}

func ConvertTimeToUnixBigInt(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func CurrentDateString(locs ...*time.Location) string {
	if len(locs) > 0 {
		time.Now().In(locs[0]).Format(DATE_FORMAT)
	}
	return ConvertTimeInLocation(time.Now()).Format(DATE_FORMAT)
}

func CurrentDateStringBeiJing() string {
	return ConvertTimeInLocation(time.Now()).Format(DATE_FORMAT)
}

func CurrentTomorrowDateStringBeiJing() string {
	return ConvertTimeInLocation(time.Now().Add(24 * time.Hour)).Format(DATE_FORMAT)
}

func TimeToString(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func TimeToStringDate(t time.Time) string {
	return t.Format(DATE_FORMAT)
}

func UnixToDate(unix int64) string {
	return time.Unix(unix, 0).Format(DATE_FORMAT)
}

func MicroToDate(unix int64) string {
	return ConvertInt64ToTimeInLocation(unix).Format(DATE_FORMAT)
}

func MicroToDateTime(unix int64) string {
	return ConvertInt64ToTimeInLocation(unix).Format(TIME_FORMAT)
}

func UnixToDateTime(unix int64, loc *time.Location) string {
	time.Local = loc
	return time.Unix(unix, 0).Format(TIME_FORMAT)
}

// Get Display String for Time
func DisplayForTime(seconds int, locale string) string {
	if locale == "" {
		locale = "US"
	}
	result := ""
	if seconds >= 3600 {
		result = strconv.Itoa(seconds/3600) + "h "
	}
	if seconds%3600 >= 60 {
		result = result + strconv.Itoa(seconds%3600/60) + "min"
	}
	if seconds < 60 {
		result = strconv.Itoa(seconds) + "sec"
	}
	return strings.TrimSpace(result)
}

// statDate格式 20160131
func GetTimestampByStatDate(statDate string, loc *time.Location) (int, error) {
	if len(statDate) != 8 {
		return 0, NewHErrorCustom(ERROR_CODE_TIME_FORMAT_ERROR)
	}
	y := statDate[0:4]
	m := statDate[4:6]
	d := statDate[6:8]

	dateStr := fmt.Sprintf("%s-%s-%s 00:00:00", y, m, d)
	t, err := time.ParseInLocation(TIME_FORMAT, dateStr, loc)
	if err != nil {
		return 0, err
	}
	return int(t.Unix()), nil
}

// statDate格式 2016-01-31
func GetOneDayStartTimeByDateStr(statDate string) time.Time {
	y := statDate[0:4]
	m := statDate[5:7]
	d := statDate[8:10]

	dateStr := fmt.Sprintf("%s-%s-%s 00:00:00", y, m, d)
	t, _ := time.ParseInLocation(TIME_FORMAT, dateStr, GMT8())
	return t
}

// statDate格式 20160131
func GetStatDateByTimestamp(t int, loc *time.Location) string {
	tm := time.Unix(int64(t), 0).In(loc)
	str := tm.Format(time.RFC3339)
	str = str[0:10]
	return strings.Replace(str, "-", "", -1)
}

func GetSeatlleLocation() *time.Location {
	loc, _ := time.LoadLocation("US/Pacific")
	return loc
}

func MicroSecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func GMT8() *time.Location {
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	return GMT8
}

func ConvertDateToToday(t time.Time) int64 {
	year, month, day := time.Now().Date()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	retTime := time.Date(year, month, day, h, m, s, 0, GMT8).Unix() * 1000
	return retTime
}

func ConvertDateToTomorrow(t time.Time) int64 {
	year, month, day := time.Now().Add(24 * time.Hour).Date()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	retTime := time.Date(year, month, day, h, m, s, 0, GMT8).Unix() * 1000
	return retTime
}

func DateToTomorrow(t time.Time) time.Time {
	year, month, day := time.Now().Add(24 * time.Hour).Date()
	h := t.Hour()
	m := t.Minute()
	s := t.Second()
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	retTime := time.Date(year, month, day, h, m, s, 0, GMT8)
	return retTime
}

func GetTodayStartTimeStamp() int64 {
	return GetTodayStartTimeInBeijing(time.Now()).UnixNano() / 1e6
}

func GetTomorrowStartTimeStamp() int64 {
	return GetTodayStartTimeInBeijing(time.Now()).UnixNano()/1e6 + 86400000
}

func GetTodayStartTime() time.Time {
	year, month, day := time.Now().Date()
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	retTime := time.Date(year, month, day, 0, 0, 0, 0, GMT8)
	return retTime
}

func GetTodayStartTimeByTime(t int64, hour int) time.Time {
	year, month, day := ConvertInt64ToTimeInLocation(t).Date()
	GMT8, _ := time.LoadLocation("Asia/Shanghai")
	retTime := time.Date(year, month, day, hour, 0, 0, 0, GMT8)
	return retTime
}

func FormatTimeToStringWithZone(t time.Time) string {
	timeStamp := t.Unix()
	timestr := time.Unix(timeStamp, 0).In(GMT8()).String()
	timeArrary := strings.Split(timestr, " +0800 ")

	tmpStr := timestr
	if len(timeArrary) > 0 {
		tmpStr = timeArrary[0]
	}

	return tmpStr
}

func ConvertTimeInLocation(t time.Time, locs ...*time.Location) time.Time {
	if len(locs) > 0 {
		return t.In(locs[0])
	}
	return t.In(GMT8())
}

func ConvertInt64ToTimeInLocation(micro int64) time.Time {
	return ConvertTimeInLocation(time.Unix(0, micro*1e6))
}

func GetTodayStartTimeInBeijing(t time.Time) time.Time {
	y, m, d := ConvertTimeInLocation(t).Date()
	newDay := time.Date(y, m, d, 0, 0, 0, 0, GMT8())
	return newDay
}

var addTime = 24 * time.Hour

func DoSthTomorrowNOclock(hour int, f func()) {
	tomorrowZero := GetTodayStartTimeByTime(CurrentUnixBigInt(), hour).Add(addTime)
	dur := tomorrowZero.Sub(time.Now())
	time.AfterFunc(dur, func() {
		worker(f)
	})
}

func worker(f func()) {
	defer func() {
		if r := recover(); r != nil {
			Log.Error("######catch timmer worker gorutine panic:", With("recover", r))
			time.Sleep(5 * time.Second)
			worker(f)
		}
	}()
	f()
	for range time.Tick(addTime) {
		f()
	}
}

func GetTodayStartTimeMicro(t int64) int64 {
	return GetTodayStartTimeInBeijing(ConvertInt64ToTimeInLocation(t)).Unix() * 1000
}

// 获取开始和结束中间的天
func GetDaysBetween(beg, end int64) []int64 {
	beg = GetTodayStartTimeMicro(beg)
	days := make([]int64, 0)
	tmpDay := beg
	for tmpDay < end {
		days = append(days, tmpDay)
		tmpDay += DAY_MICRO
	}
	return days
}

func IsSamedate(first, second int64) bool {
	return GetTodayStartTimeByTime(first, 0).Unix() == GetTodayStartTimeByTime(second, 0).Unix()
}

func ConvertDateToOneDay(t string) string {
	return "2016-10-14 " + t[11:]
}

func Max(a, b time.Time) time.Time {
	if a.Sub(b).Nanoseconds() > 0 {
		return a
	}
	return b
}

func Min(a, b time.Time) time.Time {
	if a.Sub(b).Nanoseconds() < 0 {
		return a
	}
	return b
}

func FromNowToTomorrowBeijingSecond(t time.Time) int {
	tmp := ConvertTimeInLocation(t)
	tomorrow := DateToTomorrow(t)
	return int(tomorrow.Unix() - tmp.Unix())
}

func BeginOfThisMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, GMT8())
}

func EndOfThisMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, -1, GMT8())
}

func Between(now, from, to int64) bool {
	if now >= from && now <= to {
		return true
	}
	return false
}
