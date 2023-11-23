package middleware

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
	"toes/global"
	"toes/internal/request"
	"toes/internal/utils"

	"github.com/gin-gonic/gin"
)

type CheckHeaderReq struct {
	XMyNonce     string
	XMyTime      time.Time
	XMyTimeStamp string
	XMySign      string
	Uri          string
	RawTimeStamp string
}

func (chr *CheckHeaderReq) NonceKey() string {
	return fmt.Sprintf("XMyNonce_%s", chr.XMyNonce)
}

func (chr *CheckHeaderReq) SignKey() string {
	md5key := utils.GetRealKey(global.Cfg.Seckey.Basekey, "CheckHeaderReq")

	return md5key
}

func (chr *CheckHeaderReq) Check() (bool, string, string) {
	checkTime := global.Cfg.CheckHeader.Time
	seconds := global.Cfg.CheckHeader.Seconds
	global.LogDebugw("CheckHeader", "时间差", time.Since(chr.XMyTime).Seconds(), " 配置时间(s)", seconds)
	if checkTime && time.Since(chr.XMyTime).Seconds() > seconds {
		return false, "1000", "时间值和服务时间差值异常 [X-My-Time]"
	}

	nonce, _ := global.Cache.Get(chr.NonceKey())
	checkNonce := global.Cfg.CheckHeader.Nonce
	if checkNonce && nonce != nil {
		return false, "1001", "随机串重复 [X-My-Notice]"
	}

	u, tmpErr := url.ParseRequestURI(chr.Uri)
	if tmpErr != nil {
		global.LogErrorw("CheckHeader", "url.ParseRequestURI err", tmpErr)
	}

	global.LogDebugw("CheckHeader", "XMyTimeStamp",
		chr.XMyTimeStamp, "XMyNonce", chr.XMyNonce,
		"Path", u.Path, "signKey", chr.SignKey())
	v := fmt.Sprintf("%s%s%s%s", chr.XMyTimeStamp, chr.XMyNonce, u.Path, chr.SignKey())
	sign := utils.Md5Sum(v)

	global.LogDebugw("CheckHeader", "sign: ", sign, "XMySign: ", chr.XMySign)

	checkSign := global.Cfg.CheckHeader.Sign
	if checkSign && sign != chr.XMySign {
		return false, "1002", "签名错误 [X-My-Sign]"
	}

	return true, "0000", "签名通过..."
}

func CheckHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkAll := global.Cfg.CheckHeader.All
		if !checkAll {
			c.Next()

			return
		}

		var chr CheckHeaderReq
		ts := c.GetHeader("X-My-Time")
		chr.XMyTimeStamp = ts
		if len(ts) > 10 {
			ts = ts[:10]
		}

		timeStamp, _ := strconv.Atoi(ts)
		chr.XMyTime = time.Unix(int64(timeStamp), 0)
		chr.XMyNonce = c.GetHeader("X-My-Nonce")
		chr.XMySign = c.GetHeader("X-My-Sign")
		requestUri, tmpErr := url.QueryUnescape(c.Request.RequestURI)
		if tmpErr != nil {
			global.LogGin(c).Sugar().Errorw("重放检测异常", "err", tmpErr.Error())
		}

		chr.Uri = strings.ReplaceAll(requestUri, "\t", "")
		re, code, msg := chr.Check()
		if re == false {
			c.JSON(200, request.Response{
				Code:    code,
				Message: msg,
				Data:    nil,
				Meta:    nil,
			})
			c.Abort()
			return
		}

		n := global.Cfg.CheckHeader.NonceCacheSeconds
		global.Cache.Set(chr.NonceKey(), 1, time.Second*time.Duration(n))

		c.Next()
	}
}
