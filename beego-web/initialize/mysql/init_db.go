package mysql

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func InitMYSQL() {
	host, _ := beego.AppConfig.String("mysql_host")
	port, _ := beego.AppConfig.String("mysql_port")
	database, _ := beego.AppConfig.String("mysql_database")
	username, _ := beego.AppConfig.String("mysql_username")
	password, _ := beego.AppConfig.String("mysql_password")
	charset, _ := beego.AppConfig.String("mysql_charset")
	dsn := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=" + charset + "&parseTime=true"}, "")
	databaseInit(dsn)
}

func databaseInit(DNS string) {
	orm.Debug = true //sql调试打印
	DriveName := "mysql"
	orm.RegisterDriver(DriveName, orm.DRMySQL) //设置数据库类型
	logs.Debug("mysql conn params: %s", DNS)
	orm.RegisterDataBase("default", DriveName, DNS) //register DB
	orm.RunSyncdb("default", false, false)          //migration
}
