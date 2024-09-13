package authCenterModels

type User struct {
	Id        int64  `bson:"_id,omitempty" json:"_id,omitempty"`
	Account   string `bson:"account" json:"account"`
	Password  string `bson:"password,omitempty" json:"password,omitempty"`
	Name      string `bson:"name" json:"name"`
	AvatarUrl string `bson:"AvatarUrl" json:"AvatarUrl"` //头像地址
	Sex       string `bson:"sex" json:"sex"`
	Phone     string `bson:"phone" json:"phone"`
	Salt      string `bson:"salt,omitempty" json:"salt,omitempty"`
	RoleId    int64  `bson:"roleId" json:"roleId"`
	Desc      string `bson:"desc" json:"desc"`
}
