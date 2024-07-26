package encrypt

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"
	"unsafe"

	"golang.org/x/crypto/pbkdf2"
)

func Salt() string {
	hash := md5.Sum([]byte(RandomStr(32)))
	return hex.EncodeToString(hash[:])
}

func Key(salt, password string) string {
	key := pbkdf2.Key([]byte(password), []byte(salt), 128, 32, sha256.New)
	return hex.EncodeToString(key)
}

const letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	// 6 bits to represent a letter index
	letterIdBits = 6
	// All 1-bits as many as letterIdBits
	letterIdMask = 1<<letterIdBits - 1
	letterIdMax  = 63 / letterIdBits
)

func RandomStr(n int) string {
	b := make([]byte, n)
	src := rand.NewSource(time.Now().UnixNano())

	// A rand.Int63() generates 63 random bits, enough for letterIdMax letters!
	for i, cache, remain := n-1, src.Int63(), letterIdMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdMax
		}
		if idx := int(cache & letterIdMask); idx < len(letters) {
			b[i] = letters[idx]
			i--
		}
		cache >>= letterIdBits
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
