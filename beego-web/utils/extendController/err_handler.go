package extendController

import (
	"github.com/beego/beego/v2/core/validation"
	beego "github.com/beego/beego/v2/server/web"
	"strings"
)

type ResponseMsg struct {
	ZhCn string `json:"zh-CN"`
	EnUs string `json:"en-US"`
}

type ErrMsg struct {
	ErrCode int         `json:"errCode"`
	Message ResponseMsg `json:"message"`
}

var MessageResponse = ResponseMsg{ZhCn: "成功", EnUs: "success"}
var SuccessResponse = ErrMsg{ErrCode: 200, Message: MessageResponse}

func (c *Controller) RaiseError(httpResponseCode int, errMsg ErrMsg) {
	c.Ctx.Output.Status = httpResponseCode // 1. 设置 HTTP 响应状态码
	c.Data["json"] = errMsg                // 2. 设置响应数据为传入的错误信息 (JSON 格式)
	c.ServeJSON()                          // 3. 将错误信息以 JSON 格式发送给客户端
	panic(beego.ErrAbort)                  // 4. 停止后续的请求处理
}

func (c *Controller) RaiseParamsError() {
	data := &ErrMsg{
		ErrCode: E_PARAMS,
		Message: ResponseMsg{
			ZhCn: "传参错误",
			EnUs: "params error",
		},
	}
	c.RaiseError(E_HTTP_PARAMS, *data)
}

func (c *Controller) RaiseBodyError() {
	data := &ErrMsg{
		ErrCode: E_BODY,
		Message: ResponseMsg{
			ZhCn: "Body传参错误",
			EnUs: "body params error",
		},
	}
	c.RaiseError(E_HTTP_PARAMS, *data)
}

func (c *Controller) RaiseDBError() {
	data := &ErrMsg{
		ErrCode: E_DB,
		Message: ResponseMsg{
			ZhCn: "数据库错误",
			EnUs: "database error",
		},
	}
	c.RaiseError(E_DB, *data)
}

func (c *Controller) RaiseParamsValidError(errs []*validation.Error) {
	messageList := strings.Split(errs[0].Message, ".")
	zh_cn, en_us := "", ""
	if len(messageList) > 1 {
		zh_cn = messageList[0]
		en_us = messageList[1]
	} else {
		zh_cn = messageList[0]
		en_us = messageList[0]
	}
	data := &ErrMsg{
		ErrCode: E_PARAMS,
		Message: ResponseMsg{
			ZhCn: zh_cn,
			EnUs: en_us,
		},
	}
	c.RaiseError(E_HTTP_PARAMS, *data)
}
