package utils

import (
	"math/rand"
	"time"
)

func Init() {
	rand.Seed(time.Now().Unix())
}
