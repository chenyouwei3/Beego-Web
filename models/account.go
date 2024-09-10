package models

import (
	"github.com/beego/beego/v2/adapter/orm"
	"time"
)

func NewAccount() *Account {
	return &Account{}
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

func (a *Account) Add(o orm.Ormer) {
	accout:=NewAccount()
	err:=
}

func (a *Account) Deleted() {

}

func (a *Account) Update() {

}

func (a *Account) Get() {

}
