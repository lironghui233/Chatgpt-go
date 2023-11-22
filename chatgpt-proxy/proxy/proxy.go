package proxy

import (
	"chatgpt-proxy/pkg/config"
	"chatgpt-proxy/pkg/log"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	upStreamUrl *url.URL
	upStream    *httputil.ReverseProxy
}

func NewProxy() *Proxy {
	cnf := config.GetConf()
	upStreamUrl, err := url.Parse(strings.TrimSuffix(cnf.Chat.BaseURL, "/v1"))
	if err != nil {
		log.Fatal(err)
	}
	upStream := httputil.NewSingleHostReverseProxy(upStreamUrl)
	return &Proxy{
		upStreamUrl: upStreamUrl,
		upStream:    upStream,
	}
}

func (proxy *Proxy) ChatProxy(ctx *gin.Context) {
	ctx.Request.Host = proxy.upStreamUrl.Host         //将请求的主机头更改为上游服务器的地址
	proxy.upStream.ServeHTTP(ctx.Writer, ctx.Request) //将请求转发到上游服务器
}
