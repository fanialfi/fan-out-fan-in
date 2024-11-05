package lib

import (
	"math/rand"
	"time"
)

const (
	TotalFile     = 3000
	ContentLength = 5000
)

var TempPath = "/dev/shm/worker-pool"

type FileInfo struct {
	Index       int
	FileName    string
	WorkerIndex int // sebagai penanda worker mana yang mengerjakan operasi pembuatan file ini
	Err         error
}

func RandomString(length int) string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, length)

	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}

	return string(b)
}
