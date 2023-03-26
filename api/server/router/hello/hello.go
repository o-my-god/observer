package hello

import (
	"context"
	"fmt"
	"net/http"

	"github.com/o-my-god/observer/api/server/router"
)

type helloRouter struct {
	routes []router.Route
}

func NewRouter() router.Router {
	r := &helloRouter{}
	r.initRoutes()
	return r
}

func (r *helloRouter) Routes() []router.Route {
	return r.routes
}

func (r *helloRouter) initRoutes() {
	r.routes = []router.Route {
		router.NewPostRoute("/hello", postHello),
	}
}

func postHello(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	fmt.Fprintf(w, "welcome")
	return nil
}