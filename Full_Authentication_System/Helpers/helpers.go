package Helpers

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	crand "crypto/rand"
	"encoding/hex"
	"html/template"
	"io"
	"io/ioutil"
	"log"
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
		log.Println("helpers.go Log1", err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(crand.Reader, nonce); err != nil {
		log.Println("helpers.go Log2", err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println("helpers.go Log3", err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println("helpers.go Log4", err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Println("helpers.go Log5", err.Error())
	}
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

func ParseTemplate(fileName string, templateData interface{}) (string, error) {
	var str string
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return str, err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, templateData); err != nil {
		return str, err
	}
	str = buffer.String()
	return str, nil
}