package router

import (
	"net/http"
)

//在创建新路由时传入
type RouteWrapper func(r Route) Route

//localRoutelocalRoute定义了一个单独的API路由来连接observer守护进程。它实现了Route。
type localRoute struct {
	method  string
	path    string
	handler APIFunc
}

//handler返回APIFunc，让服务器用中间件包装它。
func (l localRoute) Handler() APIFunc {
	return l.handler
}

func (l localRoute) Method() string {
	return l.method
}

func (l localRoute) Path() string {
	return l.path
}

//NewRoute为Router初始化一个新的本地路由。
func NewRoute(method string, path string, handler APIFunc, opts ...RouteWrapper) Route {
	var r Route = localRoute{method, path, handler}

	for _, o := range opts {
		r = o(r)
	}

	return r
}

func NewGetRoute(path string, handler APIFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodGet, path, handler)
}

func NewPostRoute(path string, handler APIFunc, opts ...RouteWrapper) Route {
	return NewRoute(http.MethodPost, path, handler)
}

func NewPutRoute(path string, handler APIFunc) Route {
	return NewRoute(http.MethodPut, path, handler)
}

func NewDeleteRoute(path string, handler APIFunc) Route {
	return NewRoute(http.MethodDelete, path, handler)
}

func NewOptionsRoute(path string, handler APIFunc) Route {
	return NewRoute(http.MethodOptions, path, handler)
}

func NewHeadRoute(path string, handler APIFunc) Route {
	return NewRoute(http.MethodHead, path, handler)
}
