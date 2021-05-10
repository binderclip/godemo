package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"golang.org/x/sync/errgroup"
)

func main() {
	log.Printf("main: starting")

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	hello := func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "Hello World!")
		if err != nil {
			log.Printf("write resp failed, err: %v", err)
		}
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

	g.Go(func() error {
		defer func() {
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

	g.Go(func() error {
		defer func() {
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
		if err != nil {
			log.Printf("shutdown server2 err: %v", err)
		}
		log.Printf("shutdown all servers done")
	}

	// os.sig & errgroup done
	g.Go(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		// TODO: syscall.SIGUSR1, syscall.SIGUSR2 各自控制一个 server 做 test
		select {
		case <-ctx.Done():
			log.Printf("errgroup meet error, will shutdown servers")
		case s := <-c:
			log.Printf("got signal: %v, will shutdown servers", s)
		}
		shutdownAllServers()
		return nil
	})

	// wait server done
	err := g.Wait()
	if err != nil {
		log.Printf("errgroup err: %v", err)
	}
	log.Printf("main: done. existing")
}
