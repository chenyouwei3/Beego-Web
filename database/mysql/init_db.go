package mysql

import (
	"github.com/beego/beego/v2/adapter/orm"
	"log"
)

func mysqlInit() {
	orm.Debug = true      // 打印 sql 。可能影响性能。生产不建议使用
	driverName := "mysql" // 驱动名

	databaseName := "gophertok"

	dsn := "golang:123456@tcp(192.168.100.20:33060)/cmdb?charset=utf8mb4&parseTime=true"

	// 注册数据库驱动到 orm
	// 参数：自定义的数据库类型名，驱动类型（orm 中提供的）
	orm.RegisterDriver(driverName, orm.DRMySQL)

	// 参数：beego 必须指定默认的数据库名称，使用的驱动名称（orm 驱动类型名），数据库的配置信息，数据库（连接池），连接（池）名称
	err := orm.RegisterDataBase(databaseName, driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}
}
