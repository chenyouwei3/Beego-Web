package authCenterModels

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

func init() {
	orm.RegisterModel(new(Api))
}

type Api struct {
	Id       int64     `json:"id"`                                      // 主键，自增长
	Name     string    `json:"name" valid:"Required"`                   // API 名称
	Url      string    `json:"url"`                                     // API 路径
	Method   string    `json:"method"`                                  // HTTP 方法（GET, POST, etc.）
	Desc     string    `json:"desc"`                                    // 描述
	CreateAt time.Time `json:"create_at" orm:"type(datetime);auto_now"` // 创建时间
	UpdateAt time.Time `json:"update_at" orm:"type(datetime);auto_now"` // 更新时间
}

func (a *Api) TableName() string {
	return "apis"
}

func (a *Api) IsExistSame(o orm.Ormer, api *Api) bool {
	exists := o.QueryTable(a.TableName()).Filter("name", api.Name).Filter("url", api.Url).Filter("method", api.Method).Exist()
	return exists
}

func (a *Api) IsExist(o orm.Ormer, api *Api) bool {
	exists := o.QueryTable(a.TableName()).Filter("id", api.Id).Exist()
	return exists
}

func (a *Api) CreateApi(api *Api) (*Api, error) {
	api.CreateAt = time.Now()
	if api.Method != "POST" && api.Method != "PUT" && api.Method != "GET" && api.Method != "DELETED" {
		return api, errors.New("method error")
	}
	o := orm.NewOrm()
	isExist := a.IsExistSame(o, api)
	if isExist {
		return api, errors.New("api is exist")
	}
	_, err := o.Insert(api)
	if err != nil {
		logs.Warning("create api fail: %v", err)
		return api, err
	} else {
		return api, nil
	}
}

func (a *Api) DeletedApi(api *Api) error {
	o := orm.NewOrm()
	isExist := a.IsExistSame(o, api)
	if !isExist {
		return errors.New("api is not exist")
	}
	_, err := o.Delete(api)
	if err != nil {
		logs.Warning("deleted api fail: %v", err)
		return err
	}
	return nil
}

func (a *Api) UpdateApi(api *Api) error {
	o := orm.NewOrm()
	isExist := a.IsExistSame(o, api)
	if isExist {
		return errors.New("api is same")
	}
	api.UpdateAt = time.Now()
	_, err := o.Update(api)
	if err != nil {
		logs.Warning("update api fail: %v", err)
		return err
	}
	return nil
}

func (a *Api) GetApi(offset, limit int) ([]Api, error) {
	var apiDB []Api
	o := orm.NewOrm()
	_, err := o.Raw(`SELECT * FROM apis ORDER BY create_at DESC LIMIT ? OFFSET ?`, limit, offset).QueryRows(apiDB)
	if err != nil {
		return nil, err
	}
	return apiDB, nil
}

//func (a *Api) GetAll(id int64) Api {
//	var apiDB Api
//
//	return apiDB
//}
