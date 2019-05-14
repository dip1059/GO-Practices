package Helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	mrand "math/rand"
	"os"
	"time"
)

var Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
var lenLetters = len(Letters)

func RandomString(n int) string {
	mrand.Seed(time.Now().UnixNano())
	//sec := time.Now().Second()
	b := make([]byte, n)
	for i := range b {
		//b[i] = Letters[rand.Intn(sec)]
		b[i] = Letters[mrand.Intn(lenLetters)]
	}
	return string(b)
}

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(CreateHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("helpers.go line:41", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		fmt.Println("helpers.go line:45", err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	//fmt.Println("helpers.go line:48", string(ciphertext))
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("helpers.go line:55", err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("helpers.go line:59", err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		fmt.Println("helpers.go line:65", err.Error())
	}
	//fmt.Println("helpers.go line:68", string(plaintext))
	return plaintext
}

func EncryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(Encrypt(data, passphrase))
}

func DecryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}