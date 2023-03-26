package router

import (
	"context"
	"net/http"
)

type APIFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error

//Router是一个接口，指定了一组可以注册到oberserver服务器中的路由
type Router interface {
	Routes() []Route
}

//Route定义了一个observer服务器独立的API route
type Route interface {
	Path() string
	Method() string
	Handler() APIFunc
}
