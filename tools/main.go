package main

import (
	"encoding/base64"
	"fmt"
	"github.com/starjun/jobrunner"
	"time"
	"toes/internal/job"
)

func InitJob() {
	jobrunner.Start() // optional: jobrunner.Start(pool int, concurrent int) (10, 1)
	jobrunner.Schedule("@every 1s", job.Job01{Test: "xxxxx1"}, "xxxxx")
}

func main() {
	// test jobrunner
	InitJob()

	time.Sleep(3 * time.Second)

	// basekey 加密
	bk := "x8dsafasdf98asdfjasdfi90"
	b64bk := base64.StdEncoding.EncodeToString([]byte(bk))
	fmt.Println("seckey:basekey is", bk, " 加密后：", b64bk) // 暂时就直接 base64 了

	re_tmp, _ := base64.StdEncoding.DecodeString(b64bk)
	fmt.Println("seckey:DecodeToString is", string(re_tmp))
}
