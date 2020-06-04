package router

import (
	"net/http"

	"github.com/liangjfblue/gpusher/web/controllers"

	"github.com/gin-gonic/gin"
)

type Router struct {
	G *gin.Engine
}

func NewRouter() *Router {
	return &Router{
		G: gin.Default(),
	}
}

func (r *Router) Init() {
	r.G.Use(gin.Recovery())
	r.G.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route")
	})

	r.initRouter()
}

func (r *Router) initRouter() {
	g := r.G.Group("/v1")
	{
		g.POST("/push", controllers.PushMsg)
	}
}
