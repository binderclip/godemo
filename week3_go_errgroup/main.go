package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.Printf("main: starting")

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	serverDoneWg := sync.WaitGroup{}

	hello := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	}
	handler1 := http.NewServeMux()
	handler1.HandleFunc("/hello", hello)
	server1 := http.Server{
		Addr:    "localhost:8001",
		Handler: handler1,
	}
	handler2 := http.NewServeMux()
	handler2.HandleFunc("/hello", hello)
	server2 := http.Server{
		Addr:    "localhost:8002",
		Handler: handler2,
	}

	// run server

	serverDoneWg.Add(1)
	g.Go(func() error {
		defer func() {
			serverDoneWg.Done()
			log.Printf("server 1 done")
		}()
		err := server1.ListenAndServe()
		if err == http.ErrServerClosed {
			log.Printf("ListenAndServe 1 closed")
			return nil
		}
		log.Printf("ListenAndServe 1 err: %v", err)
		return err
	})

	serverDoneWg.Add(1)
	g.Go(func() error {
		defer func() {
			serverDoneWg.Done()
			log.Printf("server 2 done")
		}()
		err := server2.ListenAndServe()
		if err == http.ErrServerClosed {
			log.Printf("ListenAndServe 2 closed")
			return nil
		}
		log.Printf("ListenAndServe 2 err: %v", err)
		return err
	})

	shutdownAllServers := func() {
		err := server1.Shutdown(ctx)
		if err != nil {
			log.Printf("shutdown server1 err: %v", err)
		}
		err = server2.Shutdown(ctx)
		// TODO: 实际执行不到下面就已经结束了，是不是增加相关 wg
		if err != nil {
			log.Printf("shutdown server2 err: %v", err)
		}
		log.Printf("shutdown all servers done")
	}

	// os.sig
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		// TODO: syscall.SIGUSR1, syscall.SIGUSR2 各自控制一个 server 做 test
		s := <-c

		log.Printf("got signal: %v, will shutdown servers", s)
		shutdownAllServers()
	}()

	// errgroup
	go func() {
		<-ctx.Done()
		log.Printf("errgroup meet error, will shutdown servers")
		shutdownAllServers()
	}()

	// wait server done
	serverDoneWg.Wait()
	log.Printf("main: done. existing")
}
