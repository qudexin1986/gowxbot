package wxbot

import (
	"strconv"
	"time"
	"math/rand"
)

func get_unix_time(n uint8) string {
	unix_time := time.Now().UnixNano()
	return strconv.Itoa(int(unix_time))[:n]
}
func GetRandomStringFromNum(length int) string {
	bytes := []byte("0123456789")
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "e"+string(result)
}

func getMilSecond() int64 {
	return time.Now().UnixNano()/1000000
}