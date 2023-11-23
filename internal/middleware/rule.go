package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/starjun/gotools"
	"net/http"
	"strings"
	"toes/internal/request"
	"toes/internal/utils"
)

var TestRules = `
[
  {
    "opt": "suffix",
    "rev": false,
    "lcon": "and",
    "restrlist": [
      "volcgslb.com"
    ],
    "malocation": "header_hostname"
  }
]
`

type GotoolsRule struct {
	Opt        string   // 匹配方式
	ReStrList  []string // 匹配字符串
	MaLocation string   // 匹配位置
	Des        string   // 规则描述
	Rev        bool     // 是否取反
	Lcon       string   // 规则连接符
	MaValue    string   // malocation是header args post时 校验header名
}

func CheckPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var rules []GotoolsRule
		utils.JsonDecode(TestRules, &rules)
		if !CheckRule(c, c.Request, rules) {
			request.WriteResponseErr(c, "1000", nil, "CheckRule error")
			c.Abort()

			return
		}

		c.Next()
	}
}

func CheckRule(c *gin.Context, _req *http.Request, listCu []GotoolsRule) bool {
	mapstr := make(map[string]string)
	mapstr["uri"] = _req.RequestURI
	mapstr["method"] = _req.Method
	mapstr["proto"] = _req.Proto

	mapstr["remoteaddr"] = _req.RemoteAddr
	// mapstr["contentlength"] = string(_req.ContentLength)
	mapstr["host"] = _req.Host
	mapstr["useragent"] = _req.UserAgent()
	mapstr["referer"] = _req.Referer()

	for s, _ := range _req.Header {
		mapstr["header_"+strings.ToLower(s)] = _req.Header.Get(s)
	}

	_ = _req.ParseForm()
	args_values := _req.URL.Query()
	cuLen := len(listCu)
	for i := 0; i < cuLen; i++ {
		tmpc := listCu[i]
		if strings.HasPrefix(tmpc.MaLocation, "args_") {
			mapstr[tmpc.MaLocation] = args_values.Get(tmpc.MaLocation[5:])

			continue
		}
		if strings.HasPrefix(tmpc.MaLocation, "post_") {
			mapstr[tmpc.MaLocation] = _req.PostFormValue(tmpc.MaLocation[5:])

			continue
		}
		if tmpc.MaLocation == "header" {
			listCu[i].MaLocation = tmpc.MaLocation + "_" + strings.ToLower(tmpc.MaValue)

			continue
		}
		if tmpc.MaLocation == "args" {
			mapstr[tmpc.MaLocation] = args_values.Get(tmpc.MaValue)

			continue
		}
		// if tmpc.MaLocation == "post" {
		// 把request的内容读取出来
		// var bodyBytes []byte
		// if _req.Body != nil {
		//	bodyBytes, _ = ioutil.ReadAll(_req.Body)
		// }
		// var reqMapBody map[string]interface{}
		// libs.JsonDecode(string(bodyBytes), &reqMapBody)
		// if value, ok := reqMapBody[tmpc.MaValue]; ok {
		//	mapstr[tmpc.MaLocation] = libs.Strval(value)
		// }
		// }
	}
	// json body 内容检查待确定
	var gotoolsListCu []gotools.CRule
	err := utils.JsonDecode(utils.JsonEncode(listCu), &gotoolsListCu)
	if err != nil {

		return false
	}
	re := gotools.MapCrulesListMatch(mapstr, gotoolsListCu)

	return re
}
