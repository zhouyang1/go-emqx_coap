package main

import (
	"context"
	"fmt"
	"time"

	"github.com/plgd-dev/go-coap/v3/udp"
	"github.com/zhouyang1/go-emqx_coap/coap"
)

func main() {
	co, err := udp.Dial("127.0.0.1:5683")
	if err != nil {
		panic(err)
	}
	// token, err := new(coap.Coap).Connection(host, port, clientID)
	obj := new(coap.Coap)
	// 设置上下文和超时时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	cliID := "1234"
	token, err := obj.Connection(ctx, co, cliID, true, "cqset", "cqset_coap")
	fmt.Println(token, "-------", err)

	num := 0
	for {
		payload := []byte("  Hello, CoAP" + fmt.Sprintf("%d", num))

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		obj.Push(ctx, co, cliID, token, "/test/topic", payload)
		cancel()

		num++
		time.Sleep(time.Second * 5)
	}

}
