package main

import (
	"flag"
	"os"
	"os/signal"
	config "self/configs"
	"self/internal/dao"
	job2 "self/internal/job"
	dlog2 "self/internal/pkg/dlog"
	exit2 "self/internal/pkg/exit"
	runtime2 "self/internal/pkg/runtime"
	cron2 "self/internal/task"
	"syscall"
)

var configFile = flag.String("config", "./cfg/config.dev.toml", "config file path")

func main() {
	defer exit2.OnExit()
	// 解析命令参数
	flag.Parse()
	// 读取全局配置, 这一步必须放在最前面执行, 否则后续的代码无法正常执行
	config.MustInit(*configFile)
	// 初始化日志组件
	_ = dlog2.Init(config.Config.Dlog)
	// 设置runtime参数
	runtime2.MustInit()
	// 初始化MySQLDb
	dao.DbStart()
	// 初始化全局缓存（维表的数据缓存到本地）
	cron2.Init()
	// 日志流实时处理任务 (读取日志、处理日志、分发日志)
	job2.Run()
	// 维表数据定时更新任务（更新Mysql、本地缓存，日志流处理需要关联维表数据）
	cron2.Run()
	// 监听退出信号
	waitSignals()
}

func waitSignals() {
	pid := os.Getpid()

	// SIGTERM: 由kill或killall命令发送到进程默认的信号, 可以被捕获或忽略
	// SIGQUIT: 控制终端发送到进程的信号, 通常可以ctrl+\
	// SIGINT:  控制终端发送到进程的信号, 通常可以ctrl+C
	// SIGKILL: 发送到进程、以使进程立即终止, 该信号不能被捕获或忽略
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		s := <-sigs
		dlog2.Warningf("pid %v recv sig %v", pid, s)
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			dlog2.Warningf("pid %v exit", pid)
			//清理任务
			os.Exit(0)
		}
	}
}
