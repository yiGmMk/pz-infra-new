package uuidUtil

import (
	"github.com/pborman/uuid"
	"math/rand"
	"sync"
	"time"
	"strconv"
)

const UPPER_RANDOM_NUMBER = 1000

var lock *sync.RWMutex = new(sync.RWMutex)

func GetUUID() string {
	return uuid.New()
}

//Get Unique Number , it's generated based on timestamp,
//16 Digits Number, Unix()+ MilliSecond + RandomNumber
func GetUniqueNumber() int64 {
	lock.Lock()
	defer lock.Unlock()
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := int64(r1.Intn(UPPER_RANDOM_NUMBER))
	time.Sleep(time.Millisecond * 1)
	return randomNumber + (getCurrentMilliSecond() * 1000)
}

func getCurrentMilliSecond() int64 {
	utcSecond := time.Now().Unix()
	millisSecond := time.Now().Nanosecond() / 1000000
	return utcSecond*1000 + int64(millisSecond)
}

//Get Unique Number String
func GetUniqueNumberStr() string {
	uniqueNumber := GetUniqueNumber()
	return strconv.FormatInt(uniqueNumber, 10)
}
