package registry

import "github.com/gin-gonic/gin"

type Router interface {
	RegisterRoutes(base *gin.RouterGroup)
}

var routersFactory []Router

func RegisterRouter(r Router) {
	routersFactory = append(routersFactory, r)
}

func GetRouters() []Router {
	return routersFactory
}
