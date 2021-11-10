package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go_template/dao/mysql"
	"go_template/dao/redis"
	"go_template/logger"
	"go_template/routes"
	"go_template/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//项目模板

func main() {
	//1.加载配置文件
	if err := settings.Init(); err != nil {
		log.Printf("init settings failed, err :%s", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(); err != nil {
		log.Printf("init logger failed, err :%s", err)
		return
	}
	defer zap.L().Sync()
	//3.初始化 MySQL 连接
	if err := mysql.Init(); err != nil {
		log.Printf("init mysql failed, err :%s", err)
		return
	}
	defer mysql.Close()
	//4.初始化 Redis
	if err := redis.Init(); err != nil {
		log.Printf("init redis failed, err :%s", err)
		return
	}
	defer redis.Close()
	//	5.注册路由
	r := routes.Setup()
	//6.启动服务，优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Error("Listen And Serve.err %s", zap.Error(err))
		}
	}()
	//等待终端信号来优雅关闭服务器，为关闭服务器设置一个 5 秒的超时
	//创建一个接收信号的通道
	quit := make(chan os.Signal, 1)
	//kill 会默认发送 system.SIGTERM
	//	kill -2 system.SIGTERM 信号，我们常用的 Ctrl+C 就是触发这个信号
	//kill -9 发送system.SIGKILL 信号，但是不能被捕获，所以不需要添加
	//signal.Notify 把收到的 syscall.SIGINT 或者 syscall.SIGTERM 型号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞在此，当接收到上述两个信号的时候，才会向下运行
	zap.L().Info("Shutdown Server ...")
	//创建一个 5s 的超时 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5 秒内优雅关闭服务（将未处理完的请求处理完再去关闭服务），超过 5 秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}
	zap.L().Info("Server exiting..")
}
