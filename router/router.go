package router

import (
	c "github.com/Betchika99/tp_db_project/controllers"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
)

type Routes []Route

type Route struct {
	Method      string
	Path     string
	HandlerFunc fasthttp.RequestHandler
}

var routes = Routes{
	Route{
		"GET",
		"/api",
		c.Index,
	},
	Route{
		"POST",
		"/api/user/:nickname/create",
		c.UserCreate,
	},
	Route{
		"GET",
		"/api/user/:nickname/profile",
		c.UserGetOne,
	},
	Route{
		"POST",
		"/api/user/:nickname/profile",
		c.UserUpdate,
	},
	Route{
		"POST",
		"/api/forum/:slug",
		c.ForumCreate,
	},
	Route{
		"GET",
		"/api/forum/:slug/details",
		c.ForumGetOne,
	},
	Route{
		"POST",
		"/api/forum/:slug/create",
		c.ThreadCreate,
	},
	Route{
		"GET",
		"/api/forum/:slug/threads",
		c.ForumGetThreads,
	},
	Route{
		"POST",
		"/api/thread/:slug_or_id/create",
		c.PostsCreate,
	},
	Route{
		"POST",
		"/api/thread/:slug_or_id/vote",
		c.ThreadVote,
	},
	Route{
		"GET",
		"/api/thread/:slug_or_id/details",
		c.ThreadGetOne,
	},
	Route{
		"GET",
		"/api/thread/:slug_or_id/posts",
		c.ThreadGetPosts,
	},
	Route{
		"POST",
		"/api/thread/:slug_or_id/details",
		c.ThreadUpdate,
	},
	Route{
		"GET",
		"/api/forum/:slug/users",
		c.ForumGetUsers,
	},
	Route{
		"GET",
		"/api/post/:id/details",
		c.PostGetOne,
	},
	Route{
		"POST",
		"/api/post/:id/details",
		c.PostUpdate,
	},
	Route{
		"GET",
		"/api/service/status",
		c.Status,
	},
	Route{
		"POST",
		"/api/service/clear",
		c.Clear,
	},
}

func InitRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	for _, route := range routes {
		router.Handle(route.Method, route.Path, route.HandlerFunc)
	}
	log.Printf("Router inited")
	return router
}