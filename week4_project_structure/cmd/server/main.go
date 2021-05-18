package main

import (
	"fmt"

	v1 "github.com/binderclip/godemo/week4_project_structure/api/chatbox/v1"
)

func main() {
	fmt.Println("hello world")
	req := v1.SendMessageRequest{
		Name: "Tony",
		Mail: "tony@mybarbershop.com",
		Text: "Hi!",
	}
	fmt.Printf("req: %v\n", req)
}
