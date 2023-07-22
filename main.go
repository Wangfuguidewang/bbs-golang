package main

import (
	"bbs-go/model"
	"bbs-go/routes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//创建监听退出chan
	c := make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("退出", s)
				ExitFunc()
			default:
				fmt.Println("other", s)
			}
		}
	}()

	model.Initdb()
	routes.InitRouter()

	fmt.Println("进程启动...")
	sum := 0
	for {
		sum++
		fmt.Println("sum:", sum)
		time.Sleep(time.Second)
	}
	fmt.Println("进程启动结束...")
}

func ExitFunc() {
	fmt.Println("开始退出...")
	fmt.Println("执行清理...")
	time.Sleep(3 * time.Second)
	fmt.Println("结束退出...")
	os.Exit(0)
}
