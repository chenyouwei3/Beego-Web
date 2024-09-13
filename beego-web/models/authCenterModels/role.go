package authCenterModels

type Role struct {
	Id         int64   `bson:"_id,omitempty" json:"_id,omitempty"`
	Name       string  `bson:"name" json:"name"`
	Code       string  `bson:"code" json:"code"` //标识
	Apis       []int64 `bson:"apis" json:"apis"`
	RoleRoutes string  `bson:"roleRoutes" json:"roleRoutes"` //角色所拥有的路由
	FirstPage  string  `bson:"firstPage" json:"firstPage"`   //角色首页
	Desc       string  `bson:"desc" json:"desc"`
}
