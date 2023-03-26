package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/o-my-god/observer/api/server/router"
	"github.com/sirupsen/logrus"
)

// versionMatcher定义了一个变量匹配器，当一个请求即将被服务时由路由器解析。
const versionMatcher = "/v{version:[0-9.]+}"

// Config提供API服务器的配置
type Config struct {
	TLSConfig *tls.Config
	Hosts     []string
}

// Server包含服务器的实例细节
type Server struct {
	cfg     *Config
	servers []*HTTPServer
	routers []router.Router
	//middlewares []middleware.Middleware
}

// New返回一个基于指定配置的服务器的新实例。
// 它为servapi分配所需的资源(ports, unix-sockets)。
func New(cfg *Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

// Accept设置服务器接受连接的监听器。
func (s *Server) Accept(addr string, listeners ...net.Listener) {
	for _, listener := range listeners {
		httpServer := &HTTPServer{
			srv: &http.Server{
				Addr:              addr,
				ReadHeaderTimeout: 5 * time.Minute,
			},
			l: listener,
		}
		s.servers = append(s.servers, httpServer)
	}
}

// Close关闭HTTPServer对入站请求的监听。
func (s *Server) Close() {
	for _, srv := range s.servers {
		if err := srv.Close(); err != nil {
			logrus.Error(err)
		}
	}
}

// serveAPI循环遍历所有初始化的服务器并生成一个例程
// 使用Serve方法。它还将createMux()设置为Handler。
func (s *Server) serveAPI() error {
	var chErrors = make(chan error, len(s.servers))

	for _, srv := range s.servers {
		srv.srv.Handler = s.createMux()
		go func(srv *HTTPServer) {
			var err error

			logrus.Infof("API listen on %s", srv.l.Addr())

			if err = srv.Serve(); err != nil && strings.Contains(err.Error(), "use of cloased network connections") {
				err = nil
			}
			chErrors <- err
		}(srv)
	}

	for range s.servers {
		err := <-chErrors
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) makeHTTPHandler(handler router.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/*
			ctx := context.WithValue(r.Context(), dockerversion.UAStringKey{}, r.Header.Get("User-Agent"))
			r = r.WithContext(ctx)
			handlerFunc := s.handlerWithGlobalMiddlewares(handler)
		*/

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		//if err := handlerFunc(r.Context(), w, r, vars); err != nil {
		if err := handler(r.Context(), w, r, vars); err != nil {
			/*
				if statusCode >= 500 {
					logrus.Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
				}
			*/
			//makeErrorHandler(err)(w, r)
		}
	}
}

// InitRouter初始化服务器的路由器列表。
func (s *Server) InitRouter(routers ...router.Router) {
	s.routers = append(s.routers, routers...)
}

// createMux初始化服务器使用的主路由器。
func (s *Server) createMux() *mux.Router {
	m := mux.NewRouter()

	logrus.Info("Registering routes...")
	for _, apiRouter := range s.routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())

			logrus.Infof("Registering %s, %s", r.Method(), r.Path())
			m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
		}
	}
	/*
		debugRouter := debug.NewRouter()
		s.routers = append(s.routers, debugRouter)
		for _, r := range debugRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler())
			m.Path("/debug" + r.Path()).Handler(f)
		}

		notFoundHandler := makeErrorHandler(pageNotFoundError{})
		m.HandleFunc(versionMatcher+"/{path:.*}", notFoundHandler)
		m.NotFoundHandler = notFoundHandler
		m.MethodNotAllowedHandler = notFoundHandler
	*/
	return m
}

// HTTPServer包含http服务器和监听器的实例。
// SRV *http.Server，包含创建http服务器和带有所有api endpoints的mux路由器的配置。
// l net.listener，是一个TCP或Socket监听器，它将传入的请求分派到路由器。
type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

// Serve开始监听入站请求。
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close关闭HTTPServer对入站请求的监听。
func (s *HTTPServer) Close() error {
	return s.l.Close()
}

// Wait阻塞服务器goroutine直到它退出。
// 如果API执行过程中有任何错误，它会发送一个错误消息。
func (s *Server) Wait(waitChan chan error) {
	if err := s.serveAPI(); err != nil {
		logrus.Errorf("ServeAPI error: %v", err)
		waitChan <- err
		return
	}
	waitChan <- nil
}
