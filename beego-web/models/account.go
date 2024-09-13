package models

import (
	"errors"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"time"
)

func init() {
	orm.RegisterModel(new(Account))
}

type Account struct {
	ID           int64
	Name         string
	Password     string
	Birthday     *time.Time
	Telephone    string
	Email        string
	Addr         string
	Status       int8
	RoleId       int64
	DepartmentId int64
	CreatedAt    *time.Time `orm:"auto_now_add"`
	UpdatedAt    *time.Time `orm:"auto_now"`
	DeletedAt    *time.Time
	Description  string
	Sex          bool
}

func (a *Account) Add(account *Account) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(account)
	if err != nil {
		logs.Warning("Create Account Fail: ", err)
		return 0, errors.New("插入失败")
	} else {
		logs.Debug("Create Account success")
		return id, nil
	}
}

func (a *Account) Deleted() {

}

func (a *Account) Update() {

}

func (a *Account) Get() {

}
