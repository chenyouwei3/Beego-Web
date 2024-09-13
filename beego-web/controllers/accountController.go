package controllers

import (
	"beego-web/models"
	"beego-web/utils/extendController"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
)

type AccountController struct {
	extendController.Controller
}

func (a *AccountController) Post() {
	var newAccount models.Account
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &newAccount)
	if err != nil {
		a.RaiseBodyError()
	}
	if newAccount.Name == "" {
		logs.Warning("用户名为空")
	}
	id, err := newAccount.Add(&newAccount)
	if err != nil {
		a.RaiseParamsError()
	}
	a.Data["json"] = id
	a.ServeJSON()
}
