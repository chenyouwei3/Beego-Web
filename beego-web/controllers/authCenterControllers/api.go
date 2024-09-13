package authCenterControllers

import (
	"beego-web/models/authCenterModels"
	"beego-web/utils/extendController"
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
)

type ApiController struct {
	extendController.Controller
}

//func (a *ApiController) GetAll() {
//	offset, _ := a.GetInt("offset", 0)
//	limit, _ := a.GetInt("limit", 10)
//
//}

func (a *ApiController) Post() {
	valid := validation.Validation{}
	var api authCenterModels.Api
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &api)
	if err != nil {
		logs.Warn("post data is valid json.")
		a.RaiseBodyError()
	}
	is_valid, _ := valid.Valid(&api)
	if !is_valid {
		for _, err := range valid.Errors {
			logs.Debug(err.Key, err.Message)
		}
		a.RaiseParamsValidError(valid.Errors)
	}
	_, err = api.CreateApi(&api)
	if err != nil {
		a.RaiseParamsError()
	}
	a.Data["json"] = api
	a.ServeJSON()
}
