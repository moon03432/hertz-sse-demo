package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/hertz-contrib/sse"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Reference: https://www.cloudwego.io/zh/docs/hertz/tutorials/basic-feature/protocol/sse/

func main() {
	h := server.Default()

	h.GET("/sse", handler)
	h.GET("/progress", progressHandler)

	h.Spin()
}

func handler(ctx context.Context, c *app.RequestContext) {
	// 客户端可以通过 Last-Event-ID 告知服务器收到的最后一个事件
	lastEventID := sse.GetLastEventID(c)
	hlog.CtxInfof(ctx, "last event ID: %s", lastEventID)

	// 在第一次渲染调用之前必须先行设置状态代码和响应头文件
	c.SetStatusCode(http.StatusOK)
	s := sse.NewStream(c)
	for t := range time.NewTicker(1 * time.Second).C {
		event := &sse.Event{
			Event: "timestamp",
			Data:  []byte(t.Format(time.RFC3339)),
		}
		err := s.Publish(event)
		if err != nil {
			return
		}
	}
}

func progressHandler(ctx context.Context, c *app.RequestContext) {
	// 客户端可以通过 Last-Event-ID 告知服务器收到的最后一个事件
	lastEventID := sse.GetLastEventID(c)
	hlog.CtxInfof(ctx, "last event ID: %s", lastEventID)

	var progress int64 = 0

	// 在第一次渲染调用之前必须先行设置状态代码和响应头文件
	c.SetStatusCode(http.StatusOK)
	s := sse.NewStream(c)
	for range time.NewTicker(1 * time.Second).C {
		event := &sse.Event{
			Event: "progress",
			Data:  []byte(strconv.FormatInt(progress, 10)),
		}
		err := s.Publish(event)
		if err != nil {
			return
		}

		if progress >= 100 {
			return
		}

		progress += 1
	}
}
