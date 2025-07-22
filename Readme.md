## 简介
基于 https://github.com/plgd-dev/go-coap仓库封装，开箱即用


## 功能

### 创建连接
```
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
```
### 断开连接
```
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
    err = obj.DelConnect(ctx, co, token, cliID)
	fmt.Println("err-------", err)
```
### 心跳
```
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
	err = obj.Heartbeat(ctx, co, token, cliID)
	fmt.Println("err-------", err)
```

### 消息发布
```
    payload := []byte("  Hello, CoAP")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	obj.Push(ctx, co, cliID, token, "/test/topic", payload)
```
### 订阅
```
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    obj.Sub(ctx, co, cliID, token, "/test/topic")
```
### 取消订阅（未实现）
```
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
    defer cancel()
    obj.DelTopic(ctx, co, token, cliID, "/test/topic")
```