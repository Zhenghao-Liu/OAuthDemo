package utils

import (
	"github.com/Zhenghao-Liu/OAuth_demo/common"
	"math"
	"math/rand"
	"time"
)

func GenInt64() int64 {
	return rand.Int63n(math.MaxInt64)
}

func GenString() string {
	bytes := []byte(common.StringAll)
	var ans []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < common.StringUpper; i++ {
		ans = append(ans, bytes[r.Intn(len(bytes))])
	}
	return string(ans)
}
