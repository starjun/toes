package request

import "github.com/gin-gonic/gin"

var (
	RespErrdata = map[string]string{
		"0":    "success", // 请求成功
		"1000": "默认服务器端错误",
		"1001": "参数错误",
		"1002": "[X-My-Time] 时间异常",
		"1003": "[X-My-Notice] 随机数异常",
		"1004": "[X-My-Sign] 签名错误",
		"1005": "访问权限不足",
	}
)

type ListOBJResponse struct {
	TotalCount int64         `json:"totalCount"`
	Objs       []interface{} `json:"data"`
}

type Response struct {
	// Code defines the business code.
	Code string `json:"code"`

	// Message contains the detail of this message.
	// This message is suitable to be exposed to external
	Message string `json:"message"`

	Data interface{} `json:"data"`

	Meta interface{} `json:"meta"`
}

func WriteResponseErr(c *gin.Context, code string, data interface{}, msgExt string) {
	if code == "" {
		code = "1000"
	}
	c.JSON(200, Response{
		Code:    code,
		Message: RespErrdata[code] + " " + msgExt,
		Data:    data,
		Meta:    nil,
	})
}

func WriteResponseOk(c *gin.Context, code string, data interface{}, msgExt string) {
	if code == "" {
		code = "0"
	}
	c.JSON(200, Response{
		Code:    code,
		Message: RespErrdata[code] + " " + msgExt,
		Data:    data,
		Meta:    nil,
	})
}

func WriteResponseList(c *gin.Context, code string, obj ListOBJResponse, meta interface{}) {
	if code == "" {
		code = "0"
	}
	c.JSON(200, Response{
		Code:    code,
		Message: RespErrdata[code],
		Data:    obj,
		Meta:    meta,
	})

}
