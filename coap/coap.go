package coap

import (
	"bytes"
	"context"
	"io"

	"github.com/plgd-dev/go-coap/v3/message"
	"github.com/plgd-dev/go-coap/v3/message/codes"
	"github.com/plgd-dev/go-coap/v3/message/pool"
	"github.com/plgd-dev/go-coap/v3/udp/client"
)

type Coap struct {
	Host string
	Port string
	Co   *client.Conn
}

type CoapClient interface {
	Connection(ctx context.Context, co *client.Conn, c string, auth bool, args ...string) (token string, err error)
}

func NewCoap() Coap {
	return Coap{}
}

/*
* github.com/plgd-dev/go-coap/v3/udp
*	co, err := udp.Dial(fmt.Sprintf("%s:%s", host, port))
*	if err != nil {
*		return
*	}
 */
func (c Coap) Connection(ctx context.Context, co *client.Conn, cliID string, auth bool, args ...string) (ctoken string, err error) {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)

	req.SetCode(codes.POST)
	token, err := message.GetToken()
	if err != nil {
		return
	}
	req.SetToken(token)

	path := "/mqtt/connection"

	req.SetPath(path)
	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cliID))
	if auth {
		req.AddOptionBytes(message.URIQuery, []byte("username="+args[0]))
		req.AddOptionBytes(message.URIQuery, []byte("password="+args[1]))
	}

	resp, err := co.Do(req)
	if err != nil {
		return
	}
	// fmt.Println("connect rs :", resp)
	buf, err := bodyToString(resp.Body())
	return string(buf), err
}

func (c Coap) DelConnect(ctx context.Context, co *client.Conn, ctoken, cid string) error {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)

	req.SetCode(codes.DELETE)

	// 构建请求路径，包含查询参数
	path := "mqtt/connection"

	req.SetPath(path)
	token, err := message.GetToken()
	if err != nil {
		return err
	}
	req.SetToken(token)

	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cid))
	req.AddOptionBytes(message.URIQuery, []byte("token="+ctoken))

	// resp, err := co.Do(req)
	_, err = co.Do(req)
	if err != nil {
		return err
	}
	// fmt.Println("delConnect rs :", resp)
	return nil
}

func (c Coap) Heartbeat(ctx context.Context, co *client.Conn, ctoken, cid string) error {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)

	req.SetCode(codes.PUT)
	token, err := message.GetToken()
	if err != nil {
		return err
	}

	req.SetToken(token)
	req.SetPath("mqtt/connection")

	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cid))
	req.AddOptionBytes(message.URIQuery, []byte("token="+ctoken))

	// resp, err := co.Do(req)
	_, err = co.Do(req)
	if err != nil {
		return err
	}
	// fmt.Println("heartbeat rs:", resp)
	return nil
}

func (c Coap) Push(ctx context.Context, co *client.Conn, ctoken, cliID, topic string, data []byte) (err error) {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)

	req.SetCode(codes.POST)
	if topic[0] == '/' {
		topic = topic[1:]
	}
	path := "ps/" + topic
	req.SetPath(path)
	token, err := message.GetToken()
	if err != nil {
		return
	}
	req.SetToken(token)
	req.SetContentFormat(message.TextPlain)

	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cliID))
	req.AddOptionBytes(message.URIQuery, []byte("token="+ctoken))
	// req.AddOptionBytes(message.URIQuery, []byte("retain=0"))
	// req.AddOptionBytes(message.URIQuery, []byte("qos=0"))
	// req.AddOptionBytes(message.URIQuery, []byte("expiry=5"))

	req.SetBody(bytes.NewReader(data))

	// 记录完整请求选项
	// reqOpts := req.Options()
	// optLog := make([]string, 0, len(reqOpts))
	// for _, o := range reqOpts {
	// 	fmt.Println("请求选项: key:", o.ID.String(), "--------------val:", string(o.Value))
	// 	optLog = append(optLog, fmt.Sprintf("%s=%v", o.ID.String(), string(o.Value)))
	// }
	// resp, err := co.Do(req)
	_, err = co.Do(req)
	if err != nil {
		return
	}
	// fmt.Println("Push rs :", resp)
	return nil

}
func (c Coap) Sub(ctx context.Context, co *client.Conn, ctoken, cliID, topic string) (data []byte, err error) {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)
	if topic[0] == '/' {
		topic = topic[1:]
	}
	path := "ps/" + topic

	token, err := message.GetToken()
	if err != nil {
		return
	}
	req.SetToken(token)
	req.SetContentFormat(message.TextPlain)

	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cliID))
	req.AddOptionBytes(message.URIQuery, []byte("token="+ctoken))
	req.AddOptionBytes(message.URIQuery, []byte("observe=0"))

	_, err = co.Observe(ctx, path, func(req *pool.Message) {
		if req == nil {
			return
		}
		content, err := bodyToString(req.Body())
		if err != nil {
			return
		}
		// fmt.Println("sub rs:", string(content))
		data = content
	}, req.Options()...)
	return
}

func (c Coap) DelTopic(ctx context.Context, co *client.Conn, ctoken, cliID, topic string) (err error) {
	req := co.AcquireMessage(ctx)
	defer co.ReleaseMessage(req)

	req.SetCode(codes.POST) //fuck docs;this is post
	if topic[0] == '/' {
		topic = topic[1:]
	}
	path := "ps/" + topic

	req.SetPath(path)
	token, err := message.GetToken()
	if err != nil {
		return
	}
	req.SetToken(token)

	req.AddOptionBytes(message.URIQuery, []byte("clientid="+cliID))
	req.AddOptionBytes(message.URIQuery, []byte("token="+ctoken))

	// resp, err := co.Do(req)
	_, err = co.Do(req)
	if err != nil {
		return
	}
	// fmt.Println("DelTopic rs :", resp)
	return err
}

func bodyToString(body io.Reader) ([]byte, error) {
	if body == nil {
		return []byte{}, nil
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(body)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}
