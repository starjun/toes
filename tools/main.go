package main

import (
	"encoding/base64"
	"fmt"
	"github.com/starjun/jobrunner"
	"time"
	"toes/internal/job"
	"toes/internal/utils"
)

func InitJob() {
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 10s", job.Job01{Test: "xxxxx1"}, "xxxxx")
}

func main() {
	// test jobrunner
	InitJob()

	time.Sleep(10 * time.Second)

	// basekey 加密
	bk := "x8dsafasdf98asdfjasdfi90"
	b64bk := base64.StdEncoding.EncodeToString([]byte(bk))
	fmt.Println("seckey:basekey is", bk, " 加密后：", b64bk) // 暂时就直接 base64 了

	// 加密 mysql
	mysqlpsd := "jMy30OC*IFwnDL(JgL"
	aesmysqlpsd := utils.EncryptInternalValue(b64bk, mysqlpsd, "mysql")
	fmt.Println("mysql 加密后: ", aesmysqlpsd)

	// 加密 redis
	redispsd := "passzj123"
	aesredispsd := utils.EncryptInternalValue(b64bk, redispsd, "redis")
	fmt.Println("redis 加密后: ", aesredispsd)

	// 加密 jwt
	jwtpsd := "myheaderxxxx"
	aesjwtpsd := utils.EncryptInternalValue(b64bk, jwtpsd, "jwt")
	fmt.Println("jwt 加密后: ", aesjwtpsd)

	//	解密/生成密码
	_bk, _ := base64.StdEncoding.DecodeString(b64bk)
	fmt.Println("seckey:basekey is", string(_bk)) // 暂时就直接 base64 了

	// mysql 解密
	_mysqlpsd := utils.DecryptInternalValue(b64bk, aesmysqlpsd, "mysql")
	fmt.Println("mysql 解密后: ", _mysqlpsd)

	// redis 解密
	_redispsd := utils.DecryptInternalValue(b64bk, aesredispsd, "redis")
	fmt.Println("redis 解密后: ", _redispsd)

	// jwt 解密
	_jwtpsd := utils.GetRealKey(bk, "jwt")
	fmt.Println("jwt 解密后（计算）: ", _jwtpsd)

	// 防重放使用的 key
	_ntd := utils.GetRealKey(bk, "CheckHeaderReq")
	fmt.Println("防重放使用的 解密后（计算）: ", _ntd)

}
