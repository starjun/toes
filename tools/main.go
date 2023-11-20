package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"toes/internal/utils"
)

func Md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getRealKey(_key, _tp string) string {
	if _tp == "mysql" {
		return Md5Sum(_key + "1" + _tp)
	} else if _tp == "redis" {
		return Md5Sum(_key + "2" + _tp)
	} else if _tp == "jwt" {
		return Md5Sum(_key + "3" + _tp)
	} else {
		return _key
	}
}

func EncryptInternalValue(_key, _value, _tp string) string {
	diykey := getRealKey(_key, _tp)
	return utils.EncryptString(_value, diykey)
}

func DecryptInternalValue(_key, _value, _tp string) string {
	diykey := getRealKey(_key, _tp)
	return utils.DecryptString(_value, diykey)
}

func main() {
	// basekey 加密
	bk := "x8dsafasdf98asdfjasdfi90"
	b64bk := base64.StdEncoding.EncodeToString([]byte(bk))
	fmt.Println("seckey:basekey is", bk, " 加密后：", b64bk) // 暂时就直接 base64 了

	// 加密 mysql
	mysqlpsd := "mysqlpasswdisxxxx"
	aesmysqlpsd := EncryptInternalValue(bk, mysqlpsd, "mysql")
	fmt.Println("mysql 加密后: ", aesmysqlpsd)

	// 加密 redis
	redispsd := "myredispassxxxxx"
	aesredispsd := EncryptInternalValue(bk, redispsd, "redis")
	fmt.Println("redis 加密后: ", aesredispsd)

	// 加密 jwt
	jwtpsd := "myheaderxxxx"
	aesjwtpsd := EncryptInternalValue(bk, jwtpsd, "jwt")
	fmt.Println("jwt 加密后: ", aesjwtpsd)

	//	解密/生成密码
	_bk, _ := base64.StdEncoding.DecodeString(b64bk)
	fmt.Println("seckey:basekey is", _bk) // 暂时就直接 base64 了

	// mysql 解密
	_mysqlpsd := DecryptInternalValue(bk, aesmysqlpsd, "mysql")
	fmt.Println("mysql 解密后: ", _mysqlpsd)

	// redis 解密
	_redispsd := DecryptInternalValue(bk, aesredispsd, "redis")
	fmt.Println("redis 解密后: ", _redispsd)

	// mysql 解密
	_jwtpsd := DecryptInternalValue(bk, aesjwtpsd, "mysql")
	fmt.Println("mysql 解密后: ", _jwtpsd)

	// 防重放使用的 key
	_ntd := utils.Md5Sum(bk + "CheckHeaderReq")
	fmt.Println("防重放使用的 解密后: ", _ntd)

}
