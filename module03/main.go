package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
G1 启动Http Server，调用/touch_err接口会往chan errOut中写入数据触发G2
G2 承上启下，当读到errOut有数据时，关掉Http Server，G1退出，触发G3的ctx.done，返回err退出
G3 监听退出信号，当监听到退出信号时，返回Err，触发G2的ctx.done,不再阻塞，执行下面的关掉Http Server，G1退出
*/

// 调用就写chan
var errOut = make(chan struct{})

func TouchErr(w http.ResponseWriter, r *http.Request) {
	errOut <- struct{}{}
}

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ping ok"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/touch_err", TouchErr)
	mux.HandleFunc("/ping", Ping)

	server := &http.Server{
		Addr:         ":1210",
		WriteTimeout: time.Second * 3, //设置3秒的写超时
		Handler:      mux,
	}
	g, ctx := errgroup.WithContext(context.Background())
	// G1
	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println("g1 http server one")
			fmt.Printf("%+v\n", err)
		}
		return err
	})
	// G2
	g.Go(func() error {
		select {
		case <-ctx.Done():
			fmt.Println("ctx cancel,g2 done")
		case <-errOut:
			fmt.Println("touch err,http server shutdown")
			fmt.Println("g2 done")
		}
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := server.Shutdown(timeoutCtx)
		return err
	})
	// G3
	g.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		select {
		case <-ctx.Done():
			fmt.Println("g3 done")
			return ctx.Err()
		case s := <-sigs:
			err := fmt.Sprintf("g3 recv sig %v", s)
			fmt.Println("s", err)
			return errors.New(err)
		}
	})

	err := g.Wait()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Println(ctx.Err())
}
