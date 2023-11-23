package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"runtime"
)

func IsZero(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func GetFileAndLine(v ...int) (string, int) {
	skip := 1
	if len(v) == 1 {
		skip = v[0]
	}
	_, file, line, _ := runtime.Caller(skip)

	return file, line
}

func JsonDecode(data string, v interface{}) error {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		file, line := GetFileAndLine(2)
		subject := fmt.Sprintf("%s:%d", file, line)
		// fmt.Println("json decode err", "object", data, "subject", subject, "result", err)
		log.Println("json decode err", "object", data, "subject", subject, "result", err)
	}

	return err
}

func JsonEncode(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("json marshal err", "err", err)
	}

	return string(data)
}

func JsonEncodeIndent(v interface{}) string {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println("json marshal indent err", "err", err)
	}

	return string(data)
}

func Md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unfading := int(origData[length-1])

	return origData[:(length - unfading)]
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padText...)
}

// AesEncrypt AES加密,CBC
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)

	return encrypted, nil
}

// AesDecrypt AES解密
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("AesDecrypt err=", err)
			// ErrorNotify(err)
		}
	}()
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)

	return origData, nil
}

func EncryptString(originData, _aeskey string) string {
	encryptedData, _ := AesEncrypt([]byte(originData), []byte(_aeskey))
	return base64.StdEncoding.EncodeToString(encryptedData)
}

func DecryptString(encryptedData, _aeskey string) string {
	encrypted, _ := base64.StdEncoding.DecodeString(encryptedData)
	originData, err := AesDecrypt(encrypted, []byte(_aeskey))
	if err != nil {
		log.Println("aes decrypt", "err", err.Error(), "encryptedData", encryptedData)
		return ""
	}
	return string(originData)
}

func GetRealKey(_key, _tp string) string {
	_bk, _ := base64.StdEncoding.DecodeString(_key)
	_sbk := string(_bk)
	if _tp == "mysql" {
		return Md5Sum(_sbk + "1" + _tp)
	} else if _tp == "redis" {
		return Md5Sum(_sbk + "2" + _tp)
	} else if _tp == "jwt" {
		return Md5Sum(_sbk + "3" + _tp)
	} else if _tp == "CheckHeaderReq" {
		return Md5Sum(_sbk + "CheckHeaderReq")
	} else {
		return _sbk
	}
}

func EncryptInternalValue(_key, _value, _tp string) string {
	diykey := GetRealKey(string(_key), _tp)
	return EncryptString(_value, diykey)
}

func DecryptInternalValue(_key, _value, _tp string) string {
	diykey := GetRealKey(string(_key), _tp)
	return DecryptString(_value, diykey)
}
