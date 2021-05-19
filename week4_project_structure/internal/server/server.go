package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	v1 "github.com/binderclip/godemo/week4_project_structure/api/chatbox/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct{}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Printf("write resp failed, err: %v", err)
		}
	}()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, err = w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}
	req := &v1.SendMessageRequest{}
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	err = protojson.Unmarshal(bs, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("marshal request failed"))
		return
	}

	bs, err = protojson.Marshal(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}
	_, err = w.Write(bs)
}

func (srv *Server) Run() error {
	fmt.Println("start...")

	handler := http.NewServeMux()
	handler.HandleFunc("/send_message", sendMessage)
	server := http.Server{
		Addr:    "localhost:8000",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		log.Printf("ListenAndServe 1 closed")
		return nil
	}

	return err
}

func (srv *Server) Stop() error {
	fmt.Println("stop...")
	return nil
}

func NewServer() *Server {
	return &Server{}
}
