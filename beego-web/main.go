package main

import (
	"beego-web/initialize/mysql"
	"context"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	public_middleware "public/middleware"
	"syscall"
)

func init() {
	//if len(os.Args) > 1 && (os.Args[1] == "v" || os.Args[1] == "version") {
	//	var vs = strings.Split("default", "+")
	//	fmt.Println(vs[0], vs[1])
	//	os.Exit(0)
	//}
	logs.Reset()                                           //重置日志
	logs.SetLogger(logs.AdapterConsole, `{"color":false}`) //配置输出在控制台,并且不需要颜色
	//读取环境变量加载配置文件
	if os.Getenv("CONFIG_PATH") != "" {
		logs.Debug("load configure file from: %s", os.Getenv("CONFIG_PATH"))
		beego.LoadAppConfig("ini", os.Getenv("CONFIG_PATH"))   //指定了配置文件格式为 ini 格式，并从 configPath 指定的路径加载文件
		logs.Reset()                                           //重置日志
		logs.SetLogger(logs.AdapterConsole, `{"color":false}`) //配置输出在控制台,并且不需要颜色
	}
	//初始化mysql
	mysql.InitMYSQL()

	var logLevel = logs.LevelInfo
	//设置进入开发模式门槛
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
		logLevel = logs.LevelDebug
	}
	var logFormat = fmt.Sprintf(`{"level":%d,"color":false}`, logLevel)
	logs.SetLogger(logs.AdapterConsole, logFormat)
	logs.SetLogFuncCallDepth(3)
}

// 定义操作系统信号chan
var sign = make(chan os.Signal, 1)

func main() {
	//pprof检测程序性能
	go func() {
		log.Println(http.ListenAndServe("127.0.0.1:6066", nil))
	}()
	//中间件
	//幂等性中间件
	beego.InsertFilter("/*", beego.BeforeExec, public_middleware.IdempotencyMiddlewareBefore) //beego.BeforeExec 是 Beego 中定义的一个过滤器执行时机，表示这个过滤器会在执行主业务逻辑之前触发（即在控制器方法被执行之前
	beego.InsertFilter("/*", beego.AfterExec, public_middleware.IdempotencyMiddlewareAfter)   //这表示过滤器会在 控制器的主业务逻辑执行之后 触发。和 beego.BeforeExec 不同，beego.AfterExec 允许过滤器在主逻辑执行完毕（也就是说，控制器的处理方法已经完成）后执行一些额外的操作。
	//当请求进入controller层抛出请求中断错误[来设置全局的自定义错误恢复函数]
	beego.BConfig.RecoverFunc = public_middleware.ProcessRequest

	go beego.Run()
	//优雅关闭程序,中断处理,调试
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	stopService()
}

// 退出函数
func stopService() {
	<-sign
	beego.BeeApp.Server.Shutdown(context.Background())
	log.Println("the beego-service is shutdown")
}
