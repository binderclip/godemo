package main

import "fmt"

func main() {
	checkProduct(1)
	checkProduct(123)
}

func checkProduct(pid int) {
	p, err := GetProduct(pid)
	if err != nil {
		fmt.Printf("get product failed, err: %+v\n", err)
		return
	}
	if p == nil {
		fmt.Printf("no product(id=%v) found\n", pid)
		return
	}
	fmt.Printf("get product(id=%v) succeed: %+v\n", pid, *p)
}
